# terraform-provider-quantum

A custom provider for terraform.

![Coveo](https://img.shields.io/badge/Coveo-awesome-f58020.svg)
[![Build Status](https://travis-ci.org/coveo/terraform-provider-quantum.svg?branch=master)](https://travis-ci.org/coveo/terraform-provider-quantum)
[![Go Report Card](https://goreportcard.com/badge/github.com/coveo/terraform-provider-quantum)](https://goreportcard.com/report/github.com/coveo/terraform-provider-quantum)

## Usage

### quantum_list_files

#### Example Usage

Returns a list of files from a directory

```hcl
data "quantum_list_files" "data_files" {
    folders = ["./data"]
}
```

The output will look like this:

```sh
data.quantum_list_files.data_files.files = ["./data/file1.txt", "./data/file2.txt"]
```

#### Argument Reference

- `folders` - (Optional) - The source list for folders
- `patterns` - (Optional) - The patterns to match files, uses [golang's filepath.Match](http://godoc.org/path/filepath#Match)
- `recursive` - (Optional) - Default `false`, walk directory recursively

#### Attributes Reference

- `files` - The list of matched files

### quantum_elb

#### Example Usage

Find only ELB matching the given tags with at least one instance healthy

```hcl
data "quantum_elb" "elb1" {
    tags    = [{ "Key" : "key1" , "Value" : "value1"}]
    healthy = true
}
```

Find all ELB matching one of the given tags

```hcl
data "quantum_elb" "elb2" {
    tags = [
        { "Key" : "key1" , "Value" : "value1"},
        { "Key" : "key2" , "Value" : "value2"}
    ]
}
```

Find ELB matching all given tags with at least one healthy instance

```hcl
data "quantum_elb" "elb3" {
    tags = [
        { "Key" : "key1" , "Value" : "value1"},
        { "Key" : "key2" , "Value" : "value2"}
    ]
    healthy        = true
    match_all_tags = True
}
```

The output will look like this:

```sh
```

#### Argument Reference

#### Attributes Reference


## Installation

1. Download the latest [release](github.com/coveo/terraform-provider-quantum/releases) for your platform
2. rename the file to `terraform-provider-quantum`
3. Copy the file to the same directory as terraform `dirname $(which terraform)` is installed

## Develop

```sh
go get github.com/coveo/terraform-provider-quantum
cd $GOPATH/src/github.com/coveo/terraform-provider-quantum
go get ./...
$EDITOR .
```