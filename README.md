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


## Code 
The `pkg` directory contains general utility functions for directory filters.  
Look at the code for more info.