package projectconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Name        string `yaml:"name"`
	SessionName string `yaml:"session_name"`
	WorkingDir  string `yaml:"working_dir"` // New field for working directory
	Tabs        []struct {
		Name     string   `yaml:"name"`
		Commands []string `yaml:"commands"`
	} `yaml:"tabs"`
	LastOpened time.Time `yaml:"last_opened"`
}

var configDir string

func initConfigDir() string {
	var configDirPath string

	if runtime.GOOS == "windows" {
		usr, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configDirPath = filepath.Join(usr, "AppData", "Local", "pee")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configDirPath = filepath.Join(homeDir, ".config", "pee")
	}

	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(configDirPath, 0755); err != nil {
			panic(err)
		}
	}

	return configDirPath
}

func Init() {
	configDir = initConfigDir() // Initialize the configDir variable.
}

func ProjectConfigFilePath(projectName string) string {
	return filepath.Join(configDir, projectName+".yml")
}

func Load(filename string) (*Configuration, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// config.LastOpened = time.Now()

	err = WriteConfigToFile(filename, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func UpdateLastOpened(projectName string) error {
	configFile := ProjectConfigFilePath(projectName)

	config, err := Load(configFile)
	if err != nil {
		return err
	}

	config.LastOpened = time.Now()

	err = WriteConfigToFile(configFile, config)
	if err != nil {
		return err
	}

	return nil
}

func ListProjects() (map[string]*Configuration, error) {
	projectConfigs := make(map[string]*Configuration)

	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		projectName := strings.TrimSuffix(file.Name(), ".yml")

		projectConfigFile := filepath.Join(configDir, file.Name())
		config, err := Load(projectConfigFile)
		if err != nil {
			return nil, err
		}

		projectConfigs[projectName] = config
	}

	return projectConfigs, nil
}

func GetProjectConfig(projectName string) (*Configuration, error) {
	projectConfigFile := ProjectConfigFilePath(projectName)
	config, err := Load(projectConfigFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func UpdateProjectConfig(projectName string, updatedConfig *Configuration) error {
	configFile := ProjectConfigFilePath(projectName)

	err := WriteConfigToFile(configFile, updatedConfig)
	if err != nil {
		return err
	}

	return nil
}

func ProjectExists(projectName string) bool {
	configFile := ProjectConfigFilePath(projectName)

	if _, err := os.Stat(configFile); err != nil {
		return false
	}

	return true
}

func CreateProject(projectName, sessionName, workingDir string, tabs []struct {
	Name     string
	Commands []string
},
) (string, error) {
	configFile := ProjectConfigFilePath(projectName)

	if _, err := os.Stat(configFile); err == nil {
		return "", fmt.Errorf("Project with the name '%s' already exists", projectName)
	}

	var tabsWithYAMLTags []struct {
		Name     string   `yaml:"name"`
		Commands []string `yaml:"commands"`
	}

	for _, tab := range tabs {
		tabWithYAMLTags := struct {
			Name     string   `yaml:"name"`
			Commands []string `yaml:"commands"`
		}{
			Name:     tab.Name,
			Commands: tab.Commands,
		}
		tabsWithYAMLTags = append(tabsWithYAMLTags, tabWithYAMLTags)
	}

	newConfig := &Configuration{
		Name:        projectName,
		SessionName: sessionName,
		WorkingDir:  workingDir,
		Tabs:        tabsWithYAMLTags,
		LastOpened:  time.Now(),
	}

	err := WriteConfigToFile(configFile, newConfig)
	if err != nil {
		return "", err
	}

	return configFile, nil
}

func WriteConfigToFile(filename string, config *Configuration) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	indentedYAML := indentYAML(string(data), "") // Convert data to string

	err = os.WriteFile(filename, []byte(indentedYAML), 0644)
	if err != nil {
		return err
	}

	return nil
}

func indentYAML(yamlString, prefix string) string {
	lines := strings.Split(yamlString, "\n")
	indentedLines := make([]string, len(lines))

	for i, line := range lines {
		indentedLines[i] = prefix + line
	}

	return strings.Join(indentedLines, "\n")
}
