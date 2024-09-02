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
	"fmt"
	"os"

	"github.com/Zettablock/zetta-go/internal"

	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:     "init",
		Aliases: []string{"initialize", "initialise", "create"},
		Short:   "Initialize a zrunner project",
		Long: `Init will create a new zrunner project, with a default
config file and the appropriate structure for a zrunner plugin.

zetta-cli must be run inside of a github repository.)
`,

		Run: func(_ *cobra.Command, args []string) {
			path, err := initializeProject()
			cobra.CheckErr(err)
			fmt.Println("Your zrunner project is ready at: ", path)
			fmt.Println()
			fmt.Println("Please edit project.yml, pipeline.yml, block_handlers.go and event_handlers.go files to add your business logic.")
		},
	}
)

func init() {

}

func initializeProject() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	project := &internal.Project{
		WorkingDir: wd,
	}

	if err = project.Create(); err != nil {
		return "", err
	}

	return project.WorkingDir, nil
}
