package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	survey "github.com/AlecAivazis/survey/v2"
)

var deleteBranch bool
var forceDeleteBranch bool

type Error struct {
	Message string
}

func (r *Error) Error() string {
	return r.Message
}

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
	rootCmd.Flags().BoolVarP(&deleteBranch, "delete", "d", false, "Delete local branches (git branch -d <branch>)")
	rootCmd.Flags().BoolVarP(&forceDeleteBranch, "force-delete", "D", false, "Force delete local branches (git branch -D <branch>)")
}
