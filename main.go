package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/urfave/cli"
)

const (
	// JSONDestinationFile Json destination file
	JSONDestinationFile = ".poguy"
)

// Config of this app
type Config struct {
	Projects []Project `json:"projects"`
}

// Project Configuration
type Project struct {
	Name      string    `json:"name"`
	Directory string    `json:"directory"`
	Programs  []Program `json:"programs"`
}

// Program Configuration
type Program struct {
	Name    string `json:"name"`
	Execute string `json:"execute"`
}

// GetConfig get config content
func GetConfig() Config {
	// Get user home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("fail to get current user")
		os.Exit(1)
	}

	// Get Config file path
	filepath := path.Join(usr.HomeDir, JSONDestinationFile)

	// Read Config file path
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("open file error: %v\n", err)
		os.Exit(1)
	}

	// Deserialize Config file
	var config Config
	json.Unmarshal(file, &config)
	return config
}

// ListProjects list all projects in config
func ListProjects() {
	// Get Config
	config := GetConfig()

	// Print out all projects name
	fmt.Println("Projects Available")
	fmt.Println("==================")
	for i, project := range config.Projects {
		fmt.Printf("%d. %s\n", i+1, project.Name)
	}
}

// OpenProject open project
func OpenProject(projectName string) {
	// Get config
	config := GetConfig()

	// Run the project programs
	for _, project := range config.Projects {
		if project.Name == projectName {
			for _, program := range project.Programs {
				s := strings.Split(program.Execute, " ")
				cmd := exec.Command(s[0], s[1:]...)
				cmd.Dir = project.Directory
				cmd.Run()
			}
			break
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "POGuy"
	app.Usage = "project opener"
	app.Commands = []cli.Command{
		{
			Name:  "list",
			Usage: "list all projects available in config",
			Action: func(c *cli.Context) error {
				ListProjects()
				return nil
			},
		},
		{
			Name:  "open",
			Usage: "open project with execute all programs",
			Action: func(c *cli.Context) error {
				OpenProject(c.Args().Get(0))
				return nil
			},
		},
	}
	app.Run(os.Args)
}
