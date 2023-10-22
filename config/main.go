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

type Pane struct {
	ShellCommand []string `yaml:"shell_command"`
}

type Window struct {
	WindowName         string   `yaml:"window_name"`
	Layout             string   `yaml:"layout"`
	ShellCommandBefore []string `yaml:"shell_command_before"`
	Panes              []Pane   `yaml:"panes"`
}

type Configuration struct {
	SessionName   string   `yaml:"name"`
	EditorCommand string   `yaml:"editor"`
	WorkingDir    string   `yaml:"root"`
	Windows       []Window `yaml:"windows"`
	LastOpened    time.Time
	Attach        bool   `yaml:"attach"`
	StartupWindow string `yaml:"startup_window"`
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

func GetEditorCommand(projectName string) (string, error) {
	configFile := ProjectConfigFilePath(projectName)

	config, err := Load(configFile)
	if err != nil {
		return "", err
	}

	return config.EditorCommand, nil
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

func CreateProject(sessionName, workingDir string, windows []Window) (string, error) {
	configFile := ProjectConfigFilePath(sessionName)

	if _, err := os.Stat(configFile); err == nil {
		return "", fmt.Errorf("Project with the name '%s' already exists", sessionName)
	}

	newConfig := &Configuration{
		SessionName: sessionName,
		WorkingDir:  workingDir,
		Windows:     windows,
		Attach:      true,
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

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
