# GitRepo Manifest Branch Protection Updater

This tool aims to ease the Github branch protection process when working on multiple repositories listed inside a GitRepo manifest file.

## Usage

Launch grmbpu executable. A config.json file should be available next to the executable, and correctly filled.

See the attached example file.

## Configuration

Based on the config.json file, fill the different information regarding the Github server, and the branches to configure protection for.

- ExcludeList: put here some repositories listed inside your manifest that you don't want to configure a protection for.

Define inside "Branches" the different branches to configure. All the keywords are based on Github API and https://github.com/google/go-github

The parts:
- RequiredStatusChecks
- RequiredPullRequestReviews

may be removed from the configuration to disable these checks. 

## TODO

- [ ] Manage multiple remote servers
- [ ] Force a specific port for ssh remotes
- [ ] Manage multiple remote credentials
- [ ] Manage included manifest files
- [ ] Manage CLI arguments
- [ ] Config file fancy generator

