package config

//Project describes config around a single .drone.yml file
type Project struct {
	YamlFile   string `yaml:"yaml_file"`
	Repository string
	Secrets    []string // Parsed with ParseSecretLine
}

type ProjectSlice []*Project

func (p ProjectSlice) RepositoryNames() (names []string) {
	for _, proj := range p {
		names = append(names, proj.Repository)
	}
	return
}
