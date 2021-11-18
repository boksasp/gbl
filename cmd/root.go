/*
Copyright © 2021 Trond Boksasp <trond@hey.com>

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
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gbl",
	Short: "cli prompt for switching between local git branches",
	Long: `gbl lists all branches in the current repo which you have locally.
	Select which branch you want to check out with the arrow keys.	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)

		// verify that gbl is run from a git repository
		_, fileErr := os.Stat(".git")
		if fileErr != nil {
			log.Print("Not a git repository")
			return
		}

		// get local branch names
		out, err := exec.Command("git", "branch", "--list", "--format=\"%(refname)\"").CombinedOutput()

		if err != nil {
			log.Fatal(string(out))
		}

		// format and filter the list of branches
		output := string(out)
		branches := strings.Split(output, "\n")
		filtered_branches := []string{}

		for _, v := range branches {
			if v != "" {
				branch_name := strings.Replace(v, "refs/heads/", "", 1)
				filtered_branches = append(filtered_branches, strings.ReplaceAll(branch_name, "\"", ""))
			}
		}

		// define and run the prompt
		searcher := func(input string, index int) bool {
			branch := filtered_branches[index]
			return strings.Contains(branch, input)
		}
		prompt := promptui.Select{
			Label:             "Select branch to checkout",
			Items:             filtered_branches,
			Size:              20,
			Searcher:          searcher,
			StartInSearchMode: true,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		// attempt to checkout the selected branch
		if result != "" {
			out, err := exec.Command("git", "checkout", result).CombinedOutput()
			if err != nil {
				log.Fatal(string(out))
			}
			fmt.Print(string(out))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gbl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gbl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gbl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
