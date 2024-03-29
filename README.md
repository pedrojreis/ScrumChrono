<div align="center">

![ScrumChrono](./assets/example.gif)

# ScrumChrono

  • <a href="#about">About</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usages</a> •
  <a href="#features">Features</a> •
  <a href="#built-with">Built With</a> •
</div>

---

[![Go Report Card](https://goreportcard.com/badge/github.com/pedrojreis/ScrumChrono)](https://goreportcard.com/report/github.com/pedrojreis/ScrumChrono) ![Golang Linter](https://github.com/pedrojreis/Scrumchrono/actions/workflows/linter.yml/badge.svg)
---

## About

ScrumChrono is a terminal UI to aid Scrum meetings. The goal of this project is to make easy to track time and provide some extra information.
Feel free to open an issue with any bug you encounter of any suggestion.

## Installation

Via [Homebrew](https://brew.sh)
```
brew install pedrojreis/tap/scrumchrono
```

Via [Winget](https//github.com/microsoft/winget-cli)
```
winget install ScrumChrono
```

## Usage

```shell
# Current version
ScrumChrono version

# Starts a wizard to init the configuration
ScrumChrono config wizard

# View Config - it will include path and content
ScrumChrono config view

# Allows the user to delete a team in his configuration
ScrumChrono config delete

# Run for firstteamname
ScrumChrono -t firstteamname
```

Please refer to the [wiki](https://github.com/pedrojreis/ScrumChrono/wiki) for more information and how-to's

## Features

* Config your teams by members, time and font.
* Countdown will change color when 33%, 66% and 100% of time has elapsed.
* Pause
* Atlassian Integration
* Soon: Statistics

<p align="center">
  <img width="420" height="640" src="./assets/jira_example.png">
</p>

## Built With

* [Go 1.21.7](https://go.dev/dl/) - Framework
    * [Cobra](https://github.com/spf13/cobra) - lib to create cli app
    * [Viper](https://github.com/spf13/viper) - configuration solution
    * [Termui](https://github.com/gizak/termui) - terminal dashboard and widget library
    * [Go-Figure](https://github.com/common-nighthawk/go-figure) - beautiful ASCII
    * [Go-jira](https://github.com/andygrunwald/go-jira) - client library for Atlassian Jira
* ❤️
