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

	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

const (
	schemaPath  = "schemas"
	packagePath = "dao"
)

// ormgenCmd represents the ormgen command
var ormgenCmd = &cobra.Command{
	Use:   "ormgen",
	Short: "Generate GORM DAO files from the provided .sql files",
	Long: `ormgen generates DAO files from .sql files. 
	
	All schema files should contain "create table" script for your tables and be stored in /schemas.`,
	Run: func(_ *cobra.Command, args []string) {
		daoPath, err := generateOrm(args)
		cobra.CheckErr(err)
		fmt.Printf("Models are generated at\n%s.\n", daoPath)
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ormgenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ormgenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generateOrm(args []string) (string, error) {
	g := gen.NewGenerator(gen.Config{
		// OutPath: "./query",
		ModelPkgPath: fmt.Sprintf("./%s", packagePath),
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	gormdb, err := gorm.Open(rawsql.New(rawsql.Config{
		DriverName: "postgres",
		FilePath: []string{
			schemaPath, // create table sql file directory
		},
	}))
	if err != nil {
		return "", err
	}

	g.UseDB(gormdb)

	g.GenerateAllTable()

	g.Execute()

	return packagePath, nil
}
