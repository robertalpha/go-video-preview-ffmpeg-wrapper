# Go-video-preview-ffmpeg-wrapper
A simple helper wrapper to generate small webm video previews using ffmpeg, useful for web previews. 

## Getting Started
- use CheckSystem() once during you application initialization to make sure `ffmpeg` and `ffprobe` are installed and found by `go-video-preview-ffmpeg-wrapper`
- use CreatePreviewDefaults(sourcePath string, destinationPath string, parts int, partDurationSeconds float64), for example t create a 3 second preview with 3 parts of 1 second run:
```
CreatePreviewDefaults("~/BigBuckBunny.mp4", "~/BigBuckBunny-preview.webm", 3, 1)
```

### Prerequisites
You need to have `ffmpeg` and `ffprobe` installed.

## Running the tests
The tests are dockerized, to ensure they run with `ffmpeg` and `ffprobe` available.
```
docker build -t go-video-preview-ffmpeg-wrapper-tests . && docker run go-video-preview-ffmpeg-wrapper-tests
```
coverage: `85.4% of statements`

### Test description
These tests make sure the commands to ffmpeg and ffprobe are working. Additionally, a clip preview consisting of 3 parts of 1second each should be 3 seconds long, this is included in the test. Visually the output is assumed to be correct.

## Versioning
`v0` for now as I'll probably only be using it. 

## License
Distributed under the MIT License. See `LICENSE` for more information.

## Authors
* **Robert van Alphen** - *Initial work* - [robertalpha](https://github.com/robertalpha)

## Acknowledgments
* Thanks to the creators of the excellent open source project `ffmpeg`
