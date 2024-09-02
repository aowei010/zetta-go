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
package pipeline

import (
	"errors"
	"os"
	"regexp"

	"github.com/Zettablock/zetta-go/internal"

	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create [pipeline-name]",
		Short: "Create a zrunner pipeline",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			err := createPipeline(args)
			cobra.CheckErr(err)
		},
	}
)

func init() {
}

func createPipeline(args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	pipelineName := args[0]

	pattern := `^[a-zA-Z0-9_-]+$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the string matches the pattern
	if !re.MatchString(pipelineName) {
		return errors.New("pipeline name should only contain alphanumeric characters, underscore and hyphen. No spaces or special characters allowed except hyphen and")
	}

	pipeline := &internal.Pipeline{
		WorkingDir: wd,
		Name:       pipelineName,
	}
	err = pipeline.Create()
	if err != nil {
		return err
	}

	return nil
}
