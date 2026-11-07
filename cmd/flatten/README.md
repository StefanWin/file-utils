# flatten ➖➖
Finds all files of a certain type in subdirectories and copies them to a new directory.  
### ❗ Warning: not fully stable yet, only supports depth of 1, output not configurable ❗

## Usage
`$ > flatten -quiet -images`  
`$ > flatten -videos`

## Arguments
- `-quiet` disables logging output
- `-verbose` enables per-file logging output
- `-images (default)` only copies images. Can't be explicitly combined with `-videos`.
- `-videos` only copies videos. Implies `-images=false`.
- `-images=false` without `-videos` copies both images and videos to `./all`.