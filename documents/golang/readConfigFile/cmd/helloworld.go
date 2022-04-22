/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

// helloworldCmd represents the helloworld command
var helloworldCmd = &cobra.Command{
	Use:   "helloworld",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("helloworld called")
		toggle, _ := cmd.Flags().GetBool("toggle")
		e, _ := cmd.Flags().GetString("echo")
		if toggle {
			fmt.Println("Toggle was called")
		}
		if len(e) > 0 {
			fmt.Println(e)
		}
		/*
			test.yaml passed looks like this
				type: "yaml"
				hello: "world"
				language:
				  - "english"
				  - "spanish"
				  - "german"
			The below calls viper which reads the file in root.go
		*/
		fmt.Println(viper.Get("hello"))
		language := viper.GetStringSlice("language")
		// I have to range over them since it's an slice of string
		for _, l := range language {
			fmt.Println(l)
		}
	},
}

func init() {
	rootCmd.AddCommand(helloworldCmd)
	helloworldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	helloworldCmd.Flags().StringP("echo", "e", "", "A help for echo")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloworldCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloworldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
