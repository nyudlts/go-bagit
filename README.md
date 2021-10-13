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
    --bag string        location of the bag to be validated, Mandatory
    -h, --help          help for validate
    --fast              test whether the bag directory has the expected payload specified in the checksum manifests without performing checksum validation to detect corruption
    --completeness-only only test whether the bag directory has the number of files and total size specified in Payload-Oxum without performing checksum validation to detect corruption
