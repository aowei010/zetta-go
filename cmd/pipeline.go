package cmd

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed templates/config.yml.tmpl
var configYmlTemplate string

//go:embed templates/block_handlers.go.tmpl
var blockHandlersTemplate string

//go:embed templates/event_handlers.go.tmpl
var eventHandlersTemplate string

//go:embed templates/go.mod.tmpl
var goModTemplate string

// Project contains name, license and paths to projects.
type Project struct {
	PkgName      string
	Copyright    string
	AbsolutePath string
	Legal        License
	Viper        bool
	AppName      string
}

func (p *Project) Create() error {
	var err error
	var mode os.FileMode = 0755
	// check if AbsolutePath exists
	if _, err = os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err = os.Mkdir(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	// create example-pipeline dir
	examplePipelineDir := fmt.Sprintf("%s/example-pipeline", p.AbsolutePath)
	if _, err = os.Stat(examplePipelineDir); os.IsNotExist(err) {
		// create directory
		if err = os.Mkdir(examplePipelineDir, 0755); err != nil {
			return err
		}
	}

	// create config.yml
	configFileName := fmt.Sprintf("%s/config.yml", examplePipelineDir)
	err = os.WriteFile(configFileName, []byte(configYmlTemplate), mode)
	if err != nil {
		return err
	}

	// create block_handlers.go
	blockHandlersFileName := fmt.Sprintf("%s/block_handlers.go", examplePipelineDir)
	err = os.WriteFile(blockHandlersFileName, []byte(blockHandlersTemplate), mode)
	if err != nil {
		return err
	}

	// create event_handlers.go
	eventHandlersFileName := fmt.Sprintf("%s/event_handlers.go", examplePipelineDir)
	err = os.WriteFile(eventHandlersFileName, []byte(eventHandlersTemplate), mode)
	if err != nil {
		return err
	}

	// create go.mod
	goModFileName := fmt.Sprintf("%s/go.mod", p.AbsolutePath)
	err = os.WriteFile(goModFileName, []byte(goModTemplate), mode)
	if err != nil {
		return err
	}

	return nil
}
