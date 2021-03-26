# filter-images 📷
Moves all images in the current directory to a new directory.  
Looks for the following file extensions:  
`".jpg", ".jpeg", ".png", ".bmp", ".gif"`

## Usage
`$ > filter-images -quiet -dir="./pics"`  
`$ > filter-images`

## Arguments

- `-quiet` disables logging output
- `-verbose` enables per-file logging output
- `-dir (default: "./images")` sets the target directory. Will create it if it does not exist.