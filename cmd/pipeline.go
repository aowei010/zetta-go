package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/Zettablock/zfaas/tpl"
)

// Pipeline contains name, license and paths to projects.
type Pipeline struct {
	PkgName      string
	Copyright    string
	AbsolutePath string
	Legal        License
	Viper        bool
	AppName      string
}

func (p *Pipeline) Create() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	// create config.yml
	configFile, err := os.Create(fmt.Sprintf("%s/config.yml", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer configFile.Close()

	configTemplate := template.Must(template.New("config").Parse(string(tpl.ConfigTemplate())))
	err = configTemplate.Execute(configFile, p)
	if err != nil {
		return err
	}

	// create block_handlers.go
	blockHandlersFile, err := os.Create(fmt.Sprintf("%s/block_handlers.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer blockHandlersFile.Close()

	blockHandlersTemplate := template.Must(template.New("blockHandlers").Parse(string(tpl.BlockHandlersTemplate())))
	err = blockHandlersTemplate.Execute(blockHandlersFile, p)
	if err != nil {
		return err
	}

	// create event_handlers.go
	eventHandlersFile, err := os.Create(fmt.Sprintf("%s/event_handlers.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer eventHandlersFile.Close()

	eventHandlersTemplate := template.Must(template.New("eventHandlers").Parse(string(tpl.EventHandlersTemplate())))
	err = eventHandlersTemplate.Execute(eventHandlersFile, p)
	if err != nil {
		return err
	}

	//Ensure required
	p.ensureRequiredMod()

	return nil
}

func (p *Pipeline) ensureRequiredMod() error {
	requiredModules := []string{"github.com/Zettablock/zsource"}
	for _, module := range requiredModules {
		fmt.Println("Module: ", module)
		cmd := exec.Command("go", "get", module)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(string(output))
			return err
		}
	}

	return nil
}
