# filter-videos 📹
Moves all videos in the current directory to a new directory.  
Looks for the following file extensions:  
`".mp4", ".avi", ".webm", ".mkv", ".flv", ".wmv"`

## Usage
`$ > filter-videos -quiet -dir="./pics"`  
`$ > filter-videos`

## Arguments

- `-quiet` disables logging output
- `-verbose` enables per-file logging output
- `-dir (default: "./videos")` sets the target directory. Will create it if it does not exist.