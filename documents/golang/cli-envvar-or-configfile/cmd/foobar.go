/*
Copyright © 2022 Christian Hernandez christian@email.com

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
	"github.com/spf13/viper"
)

// foobarCmd represents the foobar command
var foobarCmd = &cobra.Command{
	Use:   "foobar",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := viper.GetString("name")
		fmt.Println(name)
	},
}

func init() {
	rootCmd.AddCommand(foobarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// foobarCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	foobarCmd.Flags().StringP("name", "n", "", "Name of the person to greet")

	/*
		CHX -> Below is the line. Now look for "name" (which will turn into MUSTER_NAME because of what's in root)
		for the config or get it from the CLI or load it from the YAML file. What takes precedence is:
		1. CLI
		2. Env var
		3. Config file
	*/
	viper.BindPFlag("name", foobarCmd.Flags().Lookup("name"))
}
