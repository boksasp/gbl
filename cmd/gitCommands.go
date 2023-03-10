package cmd

import (
	"os/exec"
	"strings"
)

type GitOutput struct {
	Content []byte
	Err     error
}

func (g *GitOutput) ToString() string {
	return string(g.Content)
}

func (g *GitOutput) IsEmpty() bool {
	return len(g.Content) == 0
}

func (g *GitOutput) Lines() []string {
	if g.IsEmpty() {
		return nil
	}

	var lines []string

	for _, val := range strings.Split(g.ToString(), "\n") {
		if val != "" {
			lines = append(lines, val)
		}
	}

	return lines
}

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

func gitAddFiles(files []string) (string, error) {
	commandArgs := append([]string{"add"}, files...)
	out, err := exec.Command("git", commandArgs...).CombinedOutput()
	if err != nil {
		return "", &Error{Message: string(out)}
	}
	return string(out), nil
}

func gitGetModifiedFiles() ([]string, error) {
	var o GitOutput
	o.Content, o.Err = exec.Command("git", "diff", "--name-only").CombinedOutput()
	if o.Err != nil {
		return nil, &Error{Message: o.ToString()}
	}

	return o.Lines(), nil
}

func gitRemoveFiles(files []string) (string, error) {
	commandArgs := append([]string{"restore", "--staged"}, files...)
	out, err := exec.Command("git", commandArgs...).CombinedOutput()
	if err != nil {
		return "", &Error{Message: string(out)}
	}
	return string(out), nil
}

func gitGetStagedFiles() ([]string, error) {
	var o GitOutput
	o.Content, o.Err = exec.Command("git", "diff", "--staged", "--name-only").CombinedOutput()
	if o.Err != nil {
		return nil, &Error{Message: o.ToString()}
	}

	return o.Lines(), nil
}
