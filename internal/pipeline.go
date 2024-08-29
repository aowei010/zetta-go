package internal

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
)

type Pipeline struct {
	WorkingDir string
	Name       string
}

func (p *Pipeline) Create() error {
	var err error
	var mode os.FileMode = 0755
	// check if WorkingDir exists
	if _, err = os.Stat(p.WorkingDir); os.IsNotExist(err) {
		// create directory
		if err = os.Mkdir(p.WorkingDir, 0754); err != nil {
			return err
		}
	}

	// create pipeline dir
	pipelineDir := fmt.Sprintf("%s/%s", p.WorkingDir, p.Name)
	if _, err = os.Stat(pipelineDir); os.IsNotExist(err) {
		// create directory
		if err = os.Mkdir(pipelineDir, mode); err != nil {
			return err
		}
	}

	pipelineYmlTemplate = strings.Replace(pipelineYmlTemplate, "[pipeline-name]", p.Name, -1)

	// create pipeline.yml
	configFileName := fmt.Sprintf("%s/%s", pipelineDir, pipelineYml)
	err = os.WriteFile(configFileName, []byte(pipelineYmlTemplate), mode)
	if err != nil {
		return err
	}

	// create block_handlers.go
	blockHandlersFileName := fmt.Sprintf("%s/%s", pipelineDir, blockHandlersFile)
	err = os.WriteFile(blockHandlersFileName, []byte(blockHandlersTemplate), mode)
	if err != nil {
		return err
	}

	// create event_handlers.go
	eventHandlersFileName := fmt.Sprintf("%s/%s", pipelineDir, eventHandlersFile)
	err = os.WriteFile(eventHandlersFileName, []byte(eventHandlersTemplate), mode)
	if err != nil {
		return err
	}

	return nil
}
