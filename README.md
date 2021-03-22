# file-utils
Utility scripts media file management.  
Generally moves files with a specific mime-type to a directory.

## Scripts
The utility scripts do not delete the files but move them to a directory.  

All scripts take two base commandline-arguments.  
Additionally, you can use the `-h` flag to show all options.
- `-dir` to set the target directory (surrounded with quotes).
- `-quiet` to disable logging output to `stdout`.

### Implemented Scripts
- `filter-videos`
- `filter-images`

### Compiling
Requires `go 1.16`.  
- `cd /cmd/<script_name>`
- `go build` to create a binary in the current directory
- `go install` to install the binary to `$GOPATH/bin`

Easiest solution is to add `$GOPATH/bin` to your `$PATH`.


## Code 
The `pkg` directory contains general utility functions for directory filters.  
Look at the code for more info.