# go-bagit
Early Go implementation of the bagit specification

## Build
`go build -o go-bagit main/main.go`

## Usage:
`go-bagit [command] [flags]`

## Available Commands:
    help        Help about any command
    validate    validate a bag
    create      create a bag

## Validate Command

### Usage
`go-bagit validate [flags]`

### Flags
    --bag string        location of the bag to be validated, Mandatory
    -h, --help          help for validate
    --fast              test whether the bag directory has the expected payload specified in the checksum manifests without performing checksum validation to detect corruption
    --completeness-only only test whether the bag directory has the number of files and total size specified in Payload-Oxum without performing checksum validation to detect corruption

## Create Commnd

### Usage
`go-bagit create [flags]`

### Flags

    --algorithm string   the algorithm used for checksums (default "md5")
    -h, --help           help for create
    --input-dir string   the directory to be bagged
    --processes int      Use multiple processes to calculate checksums faster (default 1)