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
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	initCmd = &cobra.Command{
		Use:     "init pipeline",
		Aliases: []string{"initialize", "initialise", "create"},
		Short:   "Initialize a zrunner plugin pipeline",
		Long: `Initialize (zfaas init) will create a new pipeline, with a default
config file and the appropriate structure for a zrunner plugin.

Zfaas init must be run inside of a go module (please run "go mod init <MODNAME>" first)
`,
		Args: cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			pipelinePath, err := initializePipeline(args)
			cobra.CheckErr(err)
			cobra.CheckErr(goGet("github.com/spf13/cobra"))
			if viper.GetBool("useViper") {
				cobra.CheckErr(goGet("github.com/spf13/viper"))
			}
			cobra.CheckErr(goGet("github.com/Zettablock/zsource"))
			fmt.Printf("Your zrunner plugin pipeline is ready at\n%s.\n", pipelinePath)
			fmt.Println("Please edit config.yml, block_handlers.go or event_handlers.go files to add your business logic.")
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func initializePipeline(args []string) (string, error) {
	var modName string
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(args) > 0 {
		if args[0] != "." {
			modName = args[0]
			wd = fmt.Sprintf("%s/%s", wd, args[0])
		}
	}

	pipeline := &Pipeline{
		AbsolutePath: wd,
		PkgName:      modName,
		Legal:        getLicense(),
		Copyright:    copyrightLine(),
		Viper:        viper.GetBool("useViper"),
		AppName:      path.Base(modName),
	}

	if err := pipeline.Create(); err != nil {
		return "", err
	}

	return pipeline.AbsolutePath, nil
}

func goGet(mod string) error {
	return exec.Command("go", "get", mod).Run()
}
