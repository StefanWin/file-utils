# download ⬇️
Download all URLS to the current directory.

## Usage
`$ > download <url_1> <url_2> <url_2> ...`  

Uses `filepath.Base` to determine the filename. Examples:  
- `https://foo.bar/file.ext -> file.ext`
- `http://domain.com/videos/video.mp4 -> video.mp4`