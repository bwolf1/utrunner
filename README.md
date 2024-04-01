# utrunner

UTRunner walks a directory tree, runs all go units tests in the tree, and
generates a report.

UTRunner will run any Go unit tests found in the top level directory it is run
from, and all sub-directories.

## Configure

Edit the `config.json` file with the base path, search depth and a list of
directories to skip.

Note: Search depth is directly related to the base path. For example, a base
path of `.` would require a search depth of `0`, and a base path of `../repos`
would require a search depth of `2`.

## Build

```shell
go build
```

## Run

UTRunner is meant to be run from the top level directory that you want to walk
and run tests in (including sub-directories).

```shell
./utrunner
```
