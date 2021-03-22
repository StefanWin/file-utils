# flatten
Finds all files of a certain type in subdirectories and copies them to a new directory.  
### ❗ Warning: not fully stable yet, only supports depth of 1, output not configurable ❗

## Usage
`$ > flatten -quiet -images`  
`$ > flatten -videos`

## Arguments
- `-quiet` disables logging output
- `-images (default)` only copies images. Can't be used with `-videos`.
- `-videes` only copies videos. Can't be used with `-images`.