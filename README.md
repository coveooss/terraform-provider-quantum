# terraform-provider-quantum
A custom provider for terraform 


![Coveo](https://img.shields.io/badge/Coveo-awesome-f58020.svg)
[![Build Status](https://travis-ci.org/coveo/terraform-provider-quantum.svg?branch=master)](https://travis-ci.org/coveo/terraform-provider-quantum)
[![Go Report Card](https://goreportcard.com/badge/github.com/coveo/terraform-provider-quantum)](https://goreportcard.com/report/github.com/coveo/terraform-provider-quantum)

## Usage

### quantum_list_files

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

#### Configuration variables

- `folders` - (Optional) - The source list for folders
- `patterns` - (Optional) - The patterns to match files, uses [golang's filepath.Match](http://godoc.org/path/filepath#Match)
- `recursive` - (Optional) - Default `false`, walk directory recursively
- `files` - (Optional) - A static list of files to match


## Develop

```sh
go get github.com/coveo/terraform-provider-quantum
cd $GOPATH/src/github.com/coveo/terraform-provider-quantum
$EDITOR .

# Get dependency tool
go get -u github.com/golang/dep/cmd/dep
# Get the dependencies
dep ensure
```