package cmd

import (
	"os/exec"
	"strings"
)

func gitBranchListShort() ([]string, error) {
	output, err := exec.Command("git", "branch", "--list", "--format=\"%(refname:short)\"").CombinedOutput()

	if err != nil {
		cmdError := &Error{Message: string(output)}
		return nil, cmdError
	}

	// format and filter the list of branches
	branches := strings.Split(string(output), "\n")
	filtered_branches := []string{}

	for _, v := range branches {
		if v != "" {
			branch := strings.ReplaceAll(v, "\"", "")
			filtered_branches = append(filtered_branches, branch)
		}
	}

	return filtered_branches, nil
}

func gitBranchDelete(branch string, force bool) (string, error) {
	deleteFlag := "-d"
	if force {
		deleteFlag = "-D"
	}
	out, err := exec.Command("git", "branch", deleteFlag, branch).CombinedOutput()
	outputString := string(out)
	if err != nil {
		gitError := &Error{Message: outputString}
		return "", gitError
	}
	return outputString, nil
}

func gitCheckout(item string) (string, error) {
	out, err := exec.Command("git", "checkout", item).CombinedOutput()
	if err != nil {
		return "", &Error{Message: string(out)}
	}
	return string(out), nil
}
