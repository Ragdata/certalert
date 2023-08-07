/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/print"
	"certalert/internal/utils"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var printAll bool
var outputFormat string
var supportedOutputFormats = []string{"text", "json", "yaml"}

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if !utils.IsInList(outputFormat, supportedOutputFormats) {
			fmt.Printf("Unsupported output format: %s. Supported formats are: %s\n", outputFormat, strings.Join(supportedOutputFormats, ", "))
			cmd.Help()
			os.Exit(1)
		}

		if printAll {
			// Handle --all flag
			output, err := print.EvaluateOutputFormat(outputFormat, config.App.Certs)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			fmt.Println(output)
			return
		}

		// Handle arguments
		if len(args) < 1 {
			fmt.Println("Please provide at least one argument or use the --all flag")
			cmd.Help()
			os.Exit(1)
		}

		var certs []certificates.Certificate
		// Create a list with all wanted certificates
		for _, arg := range args {
			certificate, err := certificates.GetCertificateByName(arg, config.App.Certs)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			certs = append(certs, *certificate)
		}

		// Print the certificates
		output, err := print.EvaluateOutputFormat(outputFormat, certs)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	printCmd.PersistentFlags().BoolVarP(&printAll, "all", "A", false, "Prints all certificates")

	printCmd.Flags().StringVarP(&outputFormat, "output", "o", "", fmt.Sprintf("Output format. One of: %s", strings.Join(supportedOutputFormats, "|")))
	printCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return supportedOutputFormats, cobra.ShellCompDirectiveDefault
	})
}