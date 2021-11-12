# gbl
A cli tool for switching between local git branches. 

# How to use
## Installation
Prerequisite: You need Go to use this tool. Go to https://golang.org/doc/install to install.

Run this in your terminal to install:
```
go install github.com/boksasp/gbl@latest
```

## Using gbl
In a directory containing a git repository, run `gbl` in the terminal to get a list of all the branches you have locally, and press 'enter' to select the one you want to checkout.

```
/code/myrepo$ gbl
Search: █
? Select branch to checkout:
  > jira-1911
    develop
    feat-a
    feat-b
    master
```

Start typing to search for a branch

```
/code/myrepo$ gbl
Search: fea█
? Select branch to checkout:
  > feat-a
    feat-b
```

# Why
I found myself switching between branches a lot, and I enjoy using the terminal. So I wanted to make my life just a bit easier.

This tool started as a bash script which used Whiptail for creating the select menu.

I wanted to learn Go, so I chose this tool to be my first project.  
It turned out to be (way) better for cross-platform use as well, so it's nice to be able to use this on Windows and not just Mac/Linux.
