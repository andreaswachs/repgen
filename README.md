# repgen

`repgen` is a test report genenerator tool. It works by consuming the JSON formatted test output from running `go test`, and builds a cohesive report with the test results.

The test reults within the report is viewed in a table that supports column ordering and search along with the ability to read output from each test.

## Installation

Using the Go CLI tool with version minimum `1.24.5` run the following command:

```sh
go install github.com/andreaswachs/repgen@v0.6.0
```

## Precompiled Binaries

Precompiled binaries for the most popular platforms and architectures are available on the [releases page](https://github.com/andreaswachs/repgen/releases). You can download the appropriate binary for your system without needing to install Go or build from source.

## Usage

By default, `repgen` reads the JSON data from stdin and writes the report to `report.html` in the current directory.

`repgen` can instead read from a specific file by using the `-i <filename>` flag. `repgen` can also write to a specific file using the `-o <filename>` flag

**Examples**:

```sh  
go test -json ./some/packages/... | repgen
```

```sh 
go test -json ./some/packages/... > file.json
repgen -i file.json -o my_report.html
```

