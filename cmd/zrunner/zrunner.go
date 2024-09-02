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
	"github.com/Zettablock/zetta-go-client/cmd/zrunner/pipeline"

	"github.com/spf13/cobra"
)

// Cmd represents the zrunner command
var Cmd = &cobra.Command{
	Use:   "zrunner [command]",
	Short: "Manage your zrunner project",
	Args:  cobra.ExactArgs(1),
}

func init() {
	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(deployCmd)
	Cmd.AddCommand(ormgenCmd)
	Cmd.AddCommand(pipeline.Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
