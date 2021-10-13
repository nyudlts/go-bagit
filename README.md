# go-bagit
early Go implementation of the bagit specification - currently only able to validate bags

## Build
`go build -o go-bagit main/main.go`


## Usage:
`go-bagit [command] [flags]`

## Available Commands:
    help        Help about any command
    validate    validate a bag

## Validate Command

### Usage
`go-bagit validate [flags]`

### Flags
    --bag string   bag to be validated
    -h, --help     help for validate


