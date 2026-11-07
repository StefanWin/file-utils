# batch-download ⬇️ ⬇️ ⬇️
Asynchronously downloads all files in a text file (new line separated).  
Handles `CRLF` and `LF` line endings and skips blank lines.  
Exits with status 1 if any download failed.  

## Usage
`$ > batch-download -quiet -file="urls.txt" -output="out" -threads=5`

## Arguments
- `-quiet` disables logging output
- `-verbose` enables per-file logging output
- `-file (default: "./urls")` specifies the file with the urls
- `-output (default: "./")` specifies an output directory
- `-threads (default: 3)` number of threads to use