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
	output, err := gitCommand("branch", "--list", "--format=\"%(refname:short)\"")

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
	out, err := gitCommand("branch", deleteFlag, branch)
	outputString := string(out)
	if err != nil {
		gitError := &Error{Message: outputString}
		return "", gitError
	}
	return outputString, nil
}

func gitCheckout(item string) (string, error) {
	out, err := gitCommand("checkout", item)
	if err != nil {
		return "", &Error{Message: string(out)}
	}
	return string(out), nil
}

func gitAddFiles(files []string) (string, error) {
	commandArgs := append([]string{"add"}, files...)
	out, err := gitCommand(commandArgs...)
	if err != nil {
		return "", &Error{Message: string(out)}
	}
	return string(out), nil
}

func gitRemoveFiles(files []string) (string, error) {
	commandArgs := append([]string{"restore", "--staged"}, files...)
	out, err := gitCommand(commandArgs...)
	if err != nil {
		return "", &Error{Message: string(out)}
	}
	return string(out), nil
}

func gitGetModifiedTrackedFiles() ([]string, error) {
	return gitOutput("diff", "--name-only")
}

func gitGetModifiedFiles() ([]string, error) {
	return gitOutput("ls-files", "--others", "--modified", "--exclude-standard")
}

func gitGetUntrackedFiles() ([]string, error) {
	return gitOutput("ls-files", "--others", "--exclude-standard")
}

func gitGetStagedFiles() ([]string, error) {
	return gitOutput("diff", "--staged", "--name-only")
}

func gitCommand(commands ...string) ([]byte, error) {
	return exec.Command("git", commands...).CombinedOutput()
}

func gitOutput(args ...string) ([]string, error) {
	var o GitOutput
	o.Content, o.Err = gitCommand(args...)
	if o.Err != nil {
		return nil, &Error{Message: o.ToString()}
	}

	// Remove lines with empty strings
	var files []string
	for _, line := range o.Lines() {
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
}
