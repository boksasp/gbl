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

	"github.com/spf13/cobra"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/viper"
)

var cfgFile string
var deleteBranch bool
var forceDeleteBranch bool

func delete(branch string) error {
	_, err := gitBranchDelete(branch, forceDeleteBranch)
	if err != nil {
		log.Printf("❌ %s", branch)
		return err
	}
	log.Printf("✅ %s", branch)
	return nil
}

func deletePrompt(branches []string) {
	prompt := &survey.MultiSelect{
		Message: "Select branch(es) to delete?",
		Options: branches,
	}
	selected := []string{}
	survey.AskOne(prompt, &selected)

	if len(selected) > 0 {
		errorMessages := []error{}
		for _, branch := range selected {
			err := delete(branch)
			if err != nil {
				errorMessages = append(errorMessages, err)
			}
		}
		if len(errorMessages) > 0 {
			for _, v := range errorMessages {
				log.Print(v)
			}
		}
	} else {
		log.Println("No branch selected")
	}
}

func checkoutPrompt(branches []string) {
	prompt := &survey.Select{
		Message: "Select branch to checkout?",
		Options: branches,
	}
	var selected string
	survey.AskOne(prompt, &selected)

	// attempt to checkout the selected branch
	if selected != "" {
		out, err := gitCheckout(selected)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(out)
	} else {
		log.Println("No branch selected")
	}
}

var rootCmd = &cobra.Command{
	Use:   "gbl",
	Short: "cli prompt for switching between local git branches",
	Long: `gbl lists all branches in the current repo which you have locally.
	Select which branch you want to check out with the arrow keys.	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)

		// verify that gbl is run from a git repository
		_, err := os.Stat(".git")
		if err != nil {
			log.Print("Not a git repository")
			return
		}

		branches, err := gitBranchListShort()

		if err != nil {
			log.Fatal(err)
		}

		if deleteBranch == true || forceDeleteBranch == true {
			deletePrompt(branches)
		} else {
			checkoutPrompt(branches)
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

	rootCmd.Flags().BoolVarP(&deleteBranch, "delete", "d", false, "Delete local branches (git branch -d <branch>)")
	rootCmd.Flags().BoolVarP(&forceDeleteBranch, "force-delete", "D", false, "Force delete local branches (git branch -D <branch>)")
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
