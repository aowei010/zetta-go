/*
Copyright © 2024 Spring Zhang <spring.zhang@zettablock.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3"
)

type Payload struct {
	Organization   string            `json:"organization"`
	Project        string            `json:"project"`
	ApiKey         string            `json:"api_key"`
	GithubUrl      string            `json:"github_url"`
	Pat            string            `json:"pat"`
	Pipelines      []PipelinePayload `json:"pipelines"`
	ZsourceVersion string            `json:"zsource_version"`
	Version        string            `json:"version"`
}

type PipelinePayload struct {
	Name string `json:"name"`
}

type Config struct {
	ConfigFile     string `yaml:"config_file"`
	Org            string `yaml:"org"`
	Kind           string `yaml:"kind"`
	Network        string `yaml:"network"`
	Version        string `yaml:"version"`
	Name           string `yaml:"name"`
	ApiKey         string `yaml:"api_key"`
	GithubUrl      string `yaml:"github_url"`
	Pat            string `yaml:"pat"`
	ZsourceVersion string `yaml:"zsource_version"`
}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy current project to hosted zrunner service container.",
	Long: `The deploy command will collect all information of your current zrunner 
plugin project, generate a payload that includes project info and details for all 
pipelines, then submit a request to trigger the server side deployment.`,
	Run: func(_ *cobra.Command, args []string) {
		err := deployProject(args)
		// cobra.CheckErr(err)
		if err != nil {
			fmt.Println("Error occured:", err)
		} else {
			fmt.Println("Your project is submitted for deployment.")
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deployProject(args []string) error {
	serviceURL := "https://qugate-dev.prod-czff.zettablock.dev/api/v1/zrunner/pipeline"
	apiKey := "004fb30a-334b-47b0-b0d4-ec3e2e55977f"

	payload, err := generatePayload(args)
	if err != nil {
		return err
	}

	if payload == nil {
		return errors.New("Invalid payload returned. Payload is blank.")
	}

	fmt.Println("payload:", string(payload))

	req, err := http.NewRequest("POST", serviceURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Request failed with status: %s", resp.Status))
	}

	fmt.Println("Request successful")

	return err
}

func generatePayload(args []string) ([]byte, error) {
	var err error
	var pipelines []PipelinePayload

	payload := Payload{}
	configs, err := collectProjectInfo(args)
	if err != nil {
		return nil, err
	}

	for i, config := range configs {
		if i == 0 {
			parts := strings.Split(filepath.Dir(config.ConfigFile), "/")
			payload.Project = parts[len(parts)-2]
			payload.Organization = config.Org
			payload.Version = config.Version
			payload.ApiKey = config.ApiKey
			payload.GithubUrl = config.GithubUrl
			payload.Pat = config.Pat
			payload.ZsourceVersion = zsourceVersion()
		}
		name := config.Name
		pipelines = append(pipelines, PipelinePayload{name})
	}

	payload.Pipelines = pipelines

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func collectProjectInfo(args []string) ([]Config, error) {
	configs := []Config{}
	config := Config{}

	configFiles, err := findConfig()

	if len(configFiles) == 0 || err != nil {
		return configs, err
	}

	for _, configFile := range configFiles {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return configs, err
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return configs, err
		}

		config.ConfigFile = configFile
		configs = append(configs, config)
	}

	return configs, nil
}

// TODO: For now we support run in project folder only. Can extend to support path parameter.
func findConfig() ([]string, error) {
	var configFiles []string
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "config.yml" {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			configFiles = append(configFiles, absPath)
			return filepath.SkipDir // Stop walking after finding the file
		}
		return nil
	})
	if err != nil {
		return configFiles, err
	}

	return configFiles, nil
}

func zsourceVersion() string {
	modFile := "./go.mod"
	moduleName := "github.com/Zettablock/zsource"

	data, err := os.ReadFile(modFile)
	if err != nil {
		return ""
	}

	f, err := modfile.Parse(modFile, data, nil)
	if err != nil {
		return ""
	}

	moduleVersion := ""
	for _, req := range f.Require {
		if req.Mod.Path == moduleName {
			moduleVersion = req.Mod.Version
			break
		}
	}

	if moduleVersion == "" {
		return ""
	}

	version := &semver.Version{}
	version, err = semver.NewVersion(moduleVersion)
	if err != nil {
		return ""
	}

	return (*version).String()
}
