## MyClipboard
TUI-based clipboard text history manager.

## Key Features
* Records text data when you select and copy to clipboard
* Easily copy recorded text to the clipboard on your device.
* Preserves the existing history even after the application is closed and reopened.

## How To Use
### Installing Guide
To use this application, you will need a machine running Darwin OS with arm64 architecture.

```bash
# Clone this repository
$ git clone https://github.com/ddhyun93/mac-clipboard-history

# Go into the repository
$ cd mac-clipboard-history

# Build application
$ go build

# Run the app
$ ./myclipboard
```

> **Note**
> This application was written in Go 1.21, and building it may require a corresponding environment.

### Usage
![](https://github.com/ddhyun93/mac-clipboard-history/assets/58629967/5b11abdd-b8e8-4a43-a51f-3c935eea4855)
1. Just "CMD + C" for record your clipboard history
2. Your history will be saved and displayed in TUI application
3. Just select and "Enter" what you want to copy!

## Credits

This software uses the following open source packages:

- [Bubbles](github.com/charmbracelet/bubbles)
- [Bubbletea](github.com/charmbracelet/bubbletea)
- [golang-design/clipboard](golang.design/x/clipboard)

## Milestones
- [ ] Add tests and comments for repository forker
- [ ] Support image type clipboard
- [ ] Various data storage
- [ ] Share clipboard history within same network environments


## License

MIT

---
