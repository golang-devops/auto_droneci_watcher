package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	apex "github.com/francoishill/log"
	"github.com/go-zero-boilerplate/extended-apex-logger/logging"
	"github.com/go-zero-boilerplate/extended-apex-logger/logging/text_handler"

	"github.com/golang-devops/auto_droneci_watcher/config"
	local_logger "github.com/golang-devops/auto_droneci_watcher/logging"
)

var (
	GitSha1 = "NO_GIT_SHA1"
	Version = "1.0.0"
)

var (
	configFlag   = flag.String("config", "", "The path to the yaml config file")
	logLevel     = flag.String("loglevel", "debug", "The log level - github.com/francoishill/log")
	sampleConfig = flag.Bool("sampleconfig", false, "Prints a sample/default config to Stdout and exits")
)

func getLogger(logLevelString string) local_logger.Logger {
	defaultLevel := apex.DebugLevel
	level := defaultLevel
	if parsedLevel, err := apex.ParseLevel(logLevelString); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Unable to parse log level '%s', error: %s. Using default level '%s'", logLevelString, err.Error(), defaultLevel.String()))
	} else {
		level = parsedLevel
	}

	loggerFields := apex.Fields{}
	loggerFields["version"] = Version
	if len(strings.TrimSpace(GitSha1)) > 0 {
		//cater for scenario where git sha is not available
		loggerFields["git_sha1"] = GitSha1[:8]
	}
	apexEntry := apex.WithFields(loggerFields)

	logHandler := text_handler.New(os.Stdout, os.Stderr, text_handler.DefaultTimeStampFormat, text_handler.DefaultMessageWidth)
	exitOnEmergency := true
	return logging.NewApexLogger(level, logHandler, apexEntry, exitOnEmergency)
}

func printSampleConfig(writer io.Writer) error {
	sampleCfgBytes, err := config.NewSampleYamlBytes()
	if err != nil {
		return err
	}
	_, err = writer.Write(sampleCfgBytes)
	return err
}

func main() {
	flag.Parse()
	if *sampleConfig == true {
		if err := printSampleConfig(os.Stdout); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	logger := getLogger(*logLevel)
	logger.Info(fmt.Sprintf("Version %s, GitSha1 '%s'", Version, GitSha1))

	if *configFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.LoadConfigFile(*configFlag)
	if err != nil {
		logger.Emergency(err.Error())
	}

	checkInterval := 3 * time.Second
	logger.Info(fmt.Sprintf("Loaded config, num-projects=%d", len(cfg.Projects)))
	logger.Info(fmt.Sprintf("Using interval %s for checker", checkInterval))

	checker := &Checker{
		Interval: checkInterval,
		Cfg:      cfg,
	}
	err = checker.Run(logger)
	if err != nil {
		logger.Emergency(err.Error())
	}
}
