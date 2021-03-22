# batch-download
Asynchronously downloads all files in a text file (new line separated).  
Automatically checks for `CRLF` or `LF`.  

## Usage
`$ > batch-download -quiet -file="urls.txt" -output="out" -threads=5`

## Arguments
- `-quiet` disables logging output
- `-file (default: "./urls.txt")` specifies the file with the urls
- `-output (default: "./")` specifies an output directory
- `-threads (default: 3)` number of threads to use