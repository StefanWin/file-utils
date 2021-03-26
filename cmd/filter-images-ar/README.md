# filter-aspect-ratio 📷
Copy all images in the current directory to a new folder based on their aspect ratio,  
i.e. widescreen (`width >= height`) and portrait (`width < height`).


## Usage
`$ > filter-aspect-ratio -quiet`  
`$ > filter-aspect-ratio -portraitDir="phone" -wideScreenDir="desktop"`

## Arguments
- `-quiet` disables logging output
- `-verbose` enables per-file logging output
- `-portraitDir (default : "./portrait")` sets the target directory for portrait images.
- `-wideScreenDir (default : "./widescreen")` sets the target directory for widescreen images.