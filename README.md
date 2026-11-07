# file-utils
Utility scripts for media file management.  
Generally moves files with a specific mime-type to a directory.

## Scripts
The utility scripts do not delete the files but move or copy them to a directory.  

All scripts take the `-quiet` argument to disable logging output,  
and the `-verbose` argument to enable per-file logging output.  
Additionally, you can use the `-h` flag to show all options.  

### Implemented Scripts
- `filter-videos` moves all videos to a new directory
- `filter-images` moves all images to a new directory
- `flatten` copies all images or videos from subdirectories to a new directory (experimental, only supports `depth=1`)
- `batch-download` download files from a textfile of urls asynchronously.
- `download` downloads all urls given via cli args.
- `filter-images-ar` copies images into widescreen/portrait directories based on their aspect ratio
- `gpu-convert` converts all videos with Nvidia's `h264_nvenc` encoder.

### Installing
Requires Go 1.27 or newer.  

Install all scripts:
```bash
go install github.com/StefanWin/file-utils/cmd/...@latest
```

Or a single one, e.g. `flatten`:
```bash
go install github.com/StefanWin/file-utils/cmd/flatten@latest
```

The binaries are installed to `$(go env GOPATH)/bin`, so add that to your `$PATH`.

### Building from source
- `go build ./cmd/<script_name>` to create a binary in the current directory
- `go install ./...` from the project root to install all scripts
- `go test ./...` to run the tests

## Code 
The `pkg` directory contains general utility functions for directory filters.  
Look at the code for more info.
