## go-chrome-cookies [![Go Reference][1]][2]
`go-chrome-cookies` is an open-source tool that could help you decrypt `data ( password|bookmark|cookie|history|credit-card|downloads link )` from the browser. It supports the most popular browsers on the market and runs on Windows, macOS and Linux.

## Disclaimer
> This tool is limited to security research only, and the user assumes all legal and related responsibilities arising from its use! The author assumes no legal responsibility!

## Supported Browser

### Windows
| Browser                             | Password | Cookie | Bookmark | History |
| :---------------------------------- | :------: | :----: | :------: | :-----: |
| Google Chrome |    ✓    |   ✓   |    ✓    |    ✓    |
| Google Chrome Beta |    ✓    |   ✓   |    ✓    |    ✓    |
| Chromium |    ✓    |    ✓    |    ✓    |    ✓    |
| Microsoft Edge |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox Beta |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox Dev |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox ESR |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox Nightly |    ✓    |   ✓   |    ✓    |    ✓    |


### MacOS

Based on Apple's security policy, some browsers **require a current user password** to decrypt.

| Browser                             | Password | Cookie | Bookmark | History |
| :---------------------------------- | :------: | :----: | :------: | :-----: |
| Google Chrome |    ✓    |   ✓   |    ✓    |    ✓    |
| Google Chrome Beta |    ✓    |   ✓   |    ✓    |    ✓    |
| Chromium |    ✓    |    ✓    |    ✓    |    ✓    |
| Microsoft Edge |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox Beta |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox Dev |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox ESR |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox Nightly |    ✓    |   ✓   |    ✓    |    ✓    |

### Linux

| Browser                             | Password | Cookie | Bookmark | History |
| :---------------------------------- | :------: | :----: | :------: | :-----: |
| Google Chrome |    ✓    |   ✓   |    ✓    |    ✓    |
| Google Chrome Beta |    ✓    |   ✓   |    ✓    |    ✓    |
| Chromium |    ✓    |    ✓    |    ✓    |    ✓    |
| Microsoft Edge Dev |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox Beta |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox Dev |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox ESR |    ✓    |   ✓   |    ✓    |    ✓    |
| Firefox Nightly |    ✓    |   ✓   |    ✓    |    ✓    |


## Getting started

### Install

Installation of `go-chrome-cookies` is dead-simple, just download [the release][3] and build it.

> In some situations, this security tool will be treated as a virus by Windows Defender or other antivirus software and can not be executed. The code is all open source, you can modify and compile by yourself.

### Building from source

support `go 1.14+`

```bash
git clone https://github.com/teocci/go-chrome-cookies
cd go-chrome-cookies
go get -v -t -d ./...
go build
```

[1]: https://pkg.go.dev/badge/github.com/teocci/go-chrome-cookies.svg
[2]: https://pkg.go.dev/github.com/teocci/go-chrome-cookies
[3]: https://github.com/teocci/go-chrome-cookies/releases/tag/v1.0.0
