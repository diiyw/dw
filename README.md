# DW

[![Go Version](https://img.shields.io/badge/go-1.24.1+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A command-line tool to convert Defold HTML5 bundles to WeChat Mini Games (WeApp Mini Games).

[中文文档](README_zh.md) | [English](README.md)

这是一个可以将 Defold 打包好的 H5 包转换到微信小游戏平台的工具。

## Features

- Converts Defold HTML5 WebAssembly bundles to WeChat Mini Games format
- Automatically modifies JavaScript files for WeChat Mini Game compatibility
- Compresses WASM files using Brotli compression
- Processes and renames archive files with proper extensions
- Generates necessary WeChat Mini Game configuration files
- Includes WeChat Mini Game adapter for full compatibility

## Prerequisites

- Go 1.24.1 or later
- A Defold HTML5 export (WebAssembly build)
- WeChat Mini Game AppID

## Installation

### Download Pre-built Binaries

You can download pre-built binaries for your platform from the [Releases](https://github.com/diiyw/dw/releases) page.

### From Source

```bash
git clone https://github.com/diiyw/dw.git
cd dw
go build -o dw .
```

### Using Go Install

```bash
go install github.com/diiyw/dw@latest
```

## Usage

```bash
dw -source <source_directory> -target <target_directory> -appid <wechat_appid>
```

### Parameters

- `-source`: Path to the Defold HTML5 export directory (containing `_wasm.js` files)
- `-target`: Output directory where WeChat Mini Game files will be generated
- `-appid`: Your WeChat Mini Game AppID

### Example

```bash
# Convert a Defold game to WeChat Mini Game
dw -source /path/to/defold/html5/export -target /path/to/wechat/output -appid wx1234567890abcdef

# The tool will create a subdirectory in the target directory named after your project
# Output structure:
# /path/to/wechat/output/
# └── your_project_name/
#     ├── game.js
#     ├── game.json
#     ├── project.config.json
#     ├── project.private.config.json
#     ├── dmloader.js
#     ├── your_game_wasm.js
#     ├── your_game.wasm.br
#     ├── archive/
#     │   ├── archive_files.json
#     │   └── *.bin
#     └── weapp-adapter/
#         └── ...
```

## What It Does

1. **JavaScript Modification**: Updates `_wasm.js` files to use `GameGlobal` instead of global variables for WeChat compatibility
2. **WASM Compression**: Compresses `.wasm` files using Brotli compression (`.wasm.br`)
3. **Archive Processing**: Renames archive files with `.bin` extensions and updates references
4. **Configuration Generation**: Creates necessary WeChat Mini Game configuration files
5. **Adapter Integration**: Copies WeChat Mini Game adapter files for full compatibility

## Requirements

- Defold HTML5 export must be a WebAssembly build
- Source directory must contain `_wasm.js` files
- Valid WeChat Mini Game AppID
- Write permissions to target directory

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please open an issue on GitHub.
