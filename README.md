# utrunner
UTRunner walks a directory tree, runs all go units tests in the tree, and generates a report.

UTRunner will run any Go unit tests found in the top level directory it is run from, and all sub-directories.

## Build

```shell
go build
```

## Run

UTRunner is meant to be run from the top level directory that you want to walk and run tests in (including sub-directories).

```
./utrunner
```
