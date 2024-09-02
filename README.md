# Overview

zetta-go is a Golang client for ZettaBlock AI Network. Currently only support zrunner pipelines.

## Pre-requisites
Golang 1.21 or higher

## Setup
```bash
go install github.com/zetta-block/zetta-go@latest
```

## Usage
```bash
❯ zetta-go -h
See the website at https://zettablock.com/ for documentation and more information about running code   
on ZettaBlock AI network.

Usage:
  zetta-go [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  zrunner     Manage your zrunner project

Flags:
  -h, --help      help for zetta-go
  -v, --version   version for zetta-go

Use "zetta-go [command] --help" for more information about a command.
```

### ZRunner
```bash
❯ zetta-go zrunner -h
Manage your zrunner project

Usage:
  zetta-go zrunner [command]

Available Commands:
  deploy      Deploy the project to the hosted zrunner service
  init        Initialize a zrunner project
  ormgen      Generate GORM DAO files from the provided .sql files
  pipeline    Manage your local zrunner pipeline

Flags:
  -h, --help   help for zrunner

Use "zetta-go zrunner [command] --help" for more information about a command.
```

#### Initialize a zrunner project
You must initialize a zrunner project inside a GitHub repo.
```bash
❯ zetta-go zrunner init
```
`zetta-go` will generate a scaffold for a zrunner project.

#### Generate GORM DAO files
`zetta-go` will scan /schemas folder and generate GORM DAO files for each .sql file.
```bash
❯ zetta-go zrunner ormgen
```

#### Create a pipeline template
`zetta-go` will generate a pipeline template in /your-pipeline-name folder.
```bash
❯ zetta-go zrunner pipeline create your-pipeline-name
```

#### Deploy the project
`zetta-go` will deploy the pipeline to the hosted zrunner service. `--pat` is required for private GitHub repo.
```bash
❯ zetta-go zrunner deploy --api-key zettablock-api-key [--pat your-github-pat] 
```