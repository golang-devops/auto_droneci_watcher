package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"os"

	"os/exec"

	"sync"

	"github.com/golang-devops/auto_droneci_watcher/config"
	"github.com/golang-devops/auto_droneci_watcher/logging"
)

//Checker is responsible for continually checking modified stamps of drone yaml files and then calling the drone sign and drone secret commands
type Checker struct {
	Interval time.Duration
	Cfg      *config.Config

	timeLayout   string
	yamlModTimes map[string]string
}

type changedProjectsResult struct {
	ChangedProjects []*config.Project
	Errors          []string
}

func (c *Checker) getChangedProjects(logger logging.Logger) (result *changedProjectsResult) {
	result = &changedProjectsResult{}

	for _, proj := range c.Cfg.Projects {
		info, err := os.Stat(proj.YamlFile)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Cannot get yaml file '%s' stats, error: %s", proj.YamlFile, err.Error()))
			continue
		}

		modTimeString := info.ModTime().Format(c.timeLayout)
		if !strings.EqualFold(c.yamlModTimes[proj.YamlFile], modTimeString) {
			c.yamlModTimes[proj.YamlFile] = modTimeString
			result.ChangedProjects = append(result.ChangedProjects, proj)
		}
	}

	return
}

func (c *Checker) executeProjectDroneCmd(proj *config.Project, args ...string) error {
	cmd := exec.Command("drone", args...)
	cmd.Dir = filepath.Dir(proj.YamlFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Command failed. Error: %s. Output: %s", err.Error(), string(out))
	}
	return nil
}

func (c *Checker) executeProjectCommands(projectsWG *sync.WaitGroup, logger logging.Logger, project *config.Project) {
	defer projectsWG.Done()

	logger.WithField("repository", project.Repository).Info("Executing commands of project")

	logger.Info("Signing drone yaml")
	if err := c.executeProjectDroneCmd(project, "sign", project.Repository); err != nil {
		logger.WithError(err).Error("Unable to execute drone sign")
	}

	logger.Info("Adding secrets")
	for _, secretLine := range project.Secrets {
		secret, err := config.ParseSecretLine(secretLine)
		if err != nil {
			logger.WithError(err).Error("Unable to parse secret")
			continue
		}

		args := []string{"secret", "add"}
		for _, image := range secret.Images {
			args = append(args, "--image", image)
		}
		args = append(args, project.Repository)
		args = append(args, secret.Key)
		args = append(args, secret.Value)

		logger.WithField("secret-key", secret.Key).Info("Adding secret")
		if err := c.executeProjectDroneCmd(project, args...); err != nil {
			logger.WithError(err).Error("Unable to execute add secret")
			continue
		}
	}
}

//Run will run the checker
func (c *Checker) Run(logger logging.Logger) error {
	c.timeLayout = "2006-01-02 15:04:05"
	c.yamlModTimes = map[string]string{}

	for _, proj := range c.Cfg.Projects {
		c.yamlModTimes[proj.YamlFile] = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Format(c.timeLayout)
	}

	iteration := 0
	for {
		iteration++
		tmpLogger := logger.WithField("iteration", iteration)

		result := c.getChangedProjects(tmpLogger)
		if len(result.Errors) > 0 {
			tmpErr := fmt.Errorf("Error(s): %s", strings.Join(result.Errors, " & "))
			msg := fmt.Sprintf("Failure to check %d/%d projects", (len(c.Cfg.Projects) - len(result.ChangedProjects)), len(c.Cfg.Projects))
			tmpLogger.WithError(tmpErr).Error(msg)
		}

		if len(result.ChangedProjects) > 0 {
			changedRepositoryNamesCombined := strings.Join(config.ProjectSlice(result.ChangedProjects).RepositoryNames(), ",")
			tmpLogger.WithField("repository-names", changedRepositoryNamesCombined).Info(fmt.Sprintf("Detected %d changed projects", len(result.ChangedProjects)))
		} else {
			tmpLogger.Debug("No projects changed")
		}

		var projectsWG sync.WaitGroup
		for _, proj := range result.ChangedProjects {
			loopLogger := tmpLogger.WithField("repository", proj.Repository).WithField("drone-file", proj.YamlFile)

			projectsWG.Add(1)
			go c.executeProjectCommands(&projectsWG, loopLogger, proj)
		}
		projectsWG.Wait()

		tmpLogger.Debug(fmt.Sprintf("Sleeping for %s", c.Interval))
		time.Sleep(c.Interval)
	}
}
