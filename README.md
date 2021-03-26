# file-utils
Utility scripts media file management.  
Generally moves files with a specific mime-type to a directory.

## Scripts
The utility scripts do not delete the files but move them to a directory.  

All scripts take the `-quiet` argument to disable logging output,  
and the `-verbose` argument to enable per-file logging output.  
Additionally, you can use the `-h` flag to show all options.  

### Implemented Scripts
- `filter-videos` moves all videos to a new directory
- `filter-images` moves all images to a new directory
- `flatten` moves all images or videos from subdirectories to a new directory (experimental, only supports `depth=1`)
- `batch-download` download files from a textfile of urls asynchronously.
- `download` downloads all urls given via cli args.
- `filter-images-ar` filters images based on their aspect ratio

### Compiling
Requires `go 1.16`.  
- `cd /cmd/<script_name>`
- `go build` to create a binary in the current directory
- `go install` to install the binary to `$GOPATH/bin`

Easiest solution is to add `$GOPATH/bin` to your `$PATH`.


## Code 
The `pkg` directory contains general utility functions for directory filters.  
Look at the code for more info.