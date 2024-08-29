/*
Copyright © 2024 Zettablock

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
package zrunner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3"
)

const (
	qugateEndpoint = "https://qugate-dev.prod-czff.zettablock.dev/api/v1/zrunner/pipeline"
	goModFile      = "go.mod"
	zsourceModule  = "github.com/Zettablock/zsource"
	pipelineYml    = "pipeline.yml"
	projectYml     = "project.yml"
)

type Payload struct {
	Org            string            `json:"org"`
	Project        string            `json:"project"`
	ApiKey         string            `json:"api_key"`
	GithubRepo     string            `json:"github_repo"`
	Pat            string            `json:"pat"`
	Pipelines      []PipelinePayload `json:"pipelines"`
	ZSourceVersion string            `json:"zsource_version"`
	Version        string            `json:"version"`
}

type PipelinePayload struct {
	Name string `json:"name"`
}

type ProjectConfig struct {
	Dir            string
	Org            string
	Kind           string
	Network        string
	Version        string
	Name           string
	ApiKey         string
	GithubRepo     string `yaml:"githubRepo"`
	Pat            string
	ZSourceVersion string
	Pipelines      []PipelineConfig
}

type PipelineConfig struct {
	Name string
	Dir  string
}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the project to the hosted zrunner service",
	Run: func(cmd *cobra.Command, args []string) {
		err := deployProject(cmd)
		cobra.CheckErr(err)
		fmt.Println("Deployment submitted.")
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	deployCmd.Flags().String("api-key", "", "Zettablock api key")
	deployCmd.Flags().String("pat", "", "github repo personal access token, necessary if the repo is private")
	deployCmd.MarkFlagRequired("api-key")
}

func deployProject(cmd *cobra.Command) error {
	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}
	if apiKey == "" {
		return errors.New("api-key is required")
	}
	pat, err := cmd.Flags().GetString("pat")
	if err != nil {
		return err
	}

	payload, err := generatePayload()
	if err != nil {
		return err
	}
	payload.ApiKey = apiKey
	payload.Pat = pat

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", qugateEndpoint, bytes.NewBuffer(b))
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
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("request failed: %s", string(body))
	}

	return err
}

func generatePayload() (*Payload, error) {
	var err error
	var pipelines []PipelinePayload

	payload := &Payload{}
	config, err := collectProjectInfo()
	if err != nil {
		return nil, err
	}

	zsourceVer, err := zsourceVersion()
	if err != nil {
		return nil, err
	}
	config.ZSourceVersion = zsourceVer

	err = validateConfig(&config)
	if err != nil {
		return nil, err
	}

	payload.Project = config.Name
	payload.Org = config.Org
	payload.Version = config.Version
	payload.GithubRepo = config.GithubRepo
	payload.ZSourceVersion = config.ZSourceVersion

	for _, pipelineCfg := range config.Pipelines {
		pipelines = append(pipelines, PipelinePayload{pipelineCfg.Name})
	}
	payload.Pipelines = pipelines

	return payload, nil
}

func validateConfig(config *ProjectConfig) error {
	if config.Name == "" {
		return errors.New("project name should not be empty")
	}
	if config.Name != config.Dir {
		return fmt.Errorf("project name: %s should be the same as the project folder name: %s", config.Name, config.Dir)
	}
	if config.Org == "" {
		return errors.New("org should not be empty")
	}
	// org can only contain underscore and alphanumeric characters
	pattern := `^[a-zA-Z0-9_]+$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the string matches the pattern
	if !re.MatchString(config.Org) {
		return errors.New("org should only contain alphanumeric characters and underscore")
	}

	if config.Kind == "" {
		return errors.New("kind should not be empty")
	}
	if config.Network == "" {
		return errors.New("network should not be empty")
	}
	if config.Version == "" {
		return errors.New("version should not be empty")
	}
	// validate version
	v, err := semver.NewVersion(config.Version)
	if err != nil {
		return errors.New("invalid version")
	}
	config.Version = v.String()
	if config.GithubRepo == "" {
		return errors.New("github repo should not be empty")
	}
	// validate github repo url
	if strings.HasPrefix(config.GithubRepo, "https://") {
		config.GithubRepo = strings.TrimPrefix(config.GithubRepo, "https://")
	}
	if strings.HasPrefix(config.GithubRepo, "http://") {
		config.GithubRepo = strings.TrimPrefix(config.GithubRepo, "http://")
	}
	if !strings.Contains(config.GithubRepo, "github.com") {
		return errors.New("invalid github repo url")
	}
	for _, pipeline := range config.Pipelines {
		if pipeline.Name == "" {
			return errors.New("pipeline name should not be empty")
		}
		if pipeline.Dir != pipeline.Name {
			return fmt.Errorf("pipeline name: %s should be the same as the pipeline folder name: %s", pipeline.Name, pipeline.Dir)
		}
	}
	return nil
}

func collectProjectInfo() (ProjectConfig, error) {
	projectCfg := ProjectConfig{}
	projectCfgLoc, err := findProjectConfig()
	if err != nil {
		return projectCfg, err
	}

	projectCfg.Dir = filepath.Base(filepath.Dir(projectCfgLoc))

	data, err := os.ReadFile(projectCfgLoc)
	if err != nil {
		return projectCfg, err
	}
	err = yaml.Unmarshal(data, &projectCfg)
	if err != nil {
		return ProjectConfig{}, err
	}

	pipelineCfgs, err := findPipelineConfig()

	if len(pipelineCfgs) == 0 || err != nil {
		return projectCfg, err
	}

	for _, cfgLoc := range pipelineCfgs {
		cfg := PipelineConfig{}
		data, err = os.ReadFile(cfgLoc)
		if err != nil {
			return projectCfg, err
		}

		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			return projectCfg, err
		}

		cfg.Dir = filepath.Base(filepath.Dir(cfgLoc))
		projectCfg.Pipelines = append(projectCfg.Pipelines, cfg)
	}

	return projectCfg, nil
}

func findProjectConfig() (string, error) {
	var projectConfigLoc string
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == projectYml {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			projectConfigLoc = absPath
			return filepath.SkipDir // Stop walking after finding the file
		}
		return nil
	})
	if err != nil {
		return projectConfigLoc, err
	}

	return projectConfigLoc, nil
}

// TODO: For now we support run in project folder only. Can extend to support path parameter.
func findPipelineConfig() ([]string, error) {
	var configFiles []string
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == pipelineYml {
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

func zsourceVersion() (string, error) {
	data, err := os.ReadFile(goModFile)
	if err != nil {
		return "", err
	}

	f, err := modfile.Parse(goModFile, data, nil)
	if err != nil {
		return "", err
	}

	moduleVersion := ""
	for _, req := range f.Require {
		if req.Mod.Path == zsourceModule {
			moduleVersion = req.Mod.Version
			break
		}
	}

	if moduleVersion == "" {
		return "", errors.New("zsource module not found in go.mod file")
	}

	v, err := semver.NewVersion(moduleVersion)
	if err != nil {
		return "", err
	}

	return v.String(), nil
}
