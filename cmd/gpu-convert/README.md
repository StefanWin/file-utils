# gpu-convert
Convert all video files in the current directory with 
Nvidia's ```h264_nvenc``` encoder.  
Writes `<name>_conv.mp4` next to each input and skips `*_conv` files from earlier runs.  

### ❗ Warning: Requires a Nvidia GPU ❗

## Usage
`$ > gpu-convert -quiet`  
`$ > gpu-convert`

## Arguments
- `-quiet` disables logging output
- `-verbose` enables per-file logging output prints the output from the ```ffmpeg``` command
