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

	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"gorm.io/rawsql"
	"gorm.io/gen"
)

// ormgenCmd represents the ormgen command
var ormgenCmd = &cobra.Command{
	Use:   "ormgen",
	Short: "Generate GORM DAO models from provided schema sql file",
	Long: `Command ormgen generates DAO models from SQL files. 
	
	All schema files should contain "create table" script for your tables and be stored into folder "schema".`,
	Run: func(_ *cobra.Command, args []string) {
		daoPath, err := generateOrm(args)
		cobra.CheckErr(err)
		fmt.Printf("Models are generated at\n%s.\n", daoPath)
	},
}

func init() {
	rootCmd.AddCommand(ormgenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ormgenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ormgenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generateOrm(args []string) (string, error) {
	packagePath := "dao"
	g := gen.NewGenerator(gen.Config{
		// OutPath: "./query",
		ModelPkgPath: fmt.Sprintf("./%s", packagePath),
	    Mode: gen.WithoutContext|gen.WithDefaultQuery|gen.WithQueryInterface,
	  })

	gormdb, err := gorm.Open(rawsql.New(rawsql.Config{
		DriverName: "postgres",
		FilePath: []string{
			// "./schema/base_mints.sql", // create table sql file
	    	"./schema", // create table sql file directory
		},
	}))
	if err != nil {
		return "", err
	}

	fmt.Println(gormdb.Migrator().GetTables())

	g.UseDB(gormdb)

	// g.GenerateModel("base_mints")
	g.GenerateAllTable()

	g.Execute()

	return packagePath, nil
}
