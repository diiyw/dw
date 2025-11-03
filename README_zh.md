# DW

[![Go Version](https://img.shields.io/badge/go-1.24.1+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

一个将 Defold HTML5 包转换为微信小游戏（WeApp Mini Games）的命令行工具。

[中文文档](README_zh.md) | [English](README.md)

这是一个可以将 Defold 打包好的 H5 包转换到微信小游戏平台的工具。

## 功能特性

- 将 Defold HTML5 WebAssembly 包转换为微信小游戏格式
- 自动修改 JavaScript 文件以兼容微信小游戏
- 使用 Brotli 压缩 WASM 文件
- 处理并重命名存档文件，使用正确的扩展名
- 生成必要的微信小游戏配置文件
- 包含微信小游戏适配器以实现完全兼容

## 前置要求

- Go 1.24.1 或更高版本
- Defold HTML5 导出包（WebAssembly 构建）
- 微信小游戏 AppID

## 安装

### 下载预构建二进制文件

您可以从 [Releases](https://github.com/diiyw/dw/releases) 页面下载适用于您平台的预构建二进制文件。

### 从源码安装

```bash
git clone https://github.com/diiyw/dw.git
cd dw
go build -o dw .
```

### 使用 Go Install

```bash
go install github.com/diiyw/dw@latest
```

## 使用方法

```bash
dw -source <源目录> -target <目标目录> -appid <微信appid>
```

### 参数说明

- `-source`：Defold HTML5 导出目录的路径（包含 `_wasm.js` 文件）
- `-target`：生成微信小游戏文件的输出目录
- `-appid`：您的微信小游戏 AppID

### 示例

```bash
# 将 Defold 游戏转换为微信小游戏
dw -source /path/to/defold/html5/export -target /path/to/wechat/output -appid wx1234567890abcdef

# 工具将在目标目录中创建一个以项目命名的子目录
# 输出结构：
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

## 工作原理

1. **JavaScript 修改**：更新 `_wasm.js` 文件，使用 `GameGlobal` 替代全局变量以兼容微信
2. **WASM 压缩**：使用 Brotli 压缩将 `.wasm` 文件压缩为 `.wasm.br`
3. **存档处理**：为存档文件添加 `.bin` 扩展名并更新引用
4. **配置生成**：创建必要的微信小游戏配置文件
5. **适配器集成**：复制微信小游戏适配器文件以实现完全兼容

## 系统要求

- Defold HTML5 导出必须是 WebAssembly 构建
- 源目录必须包含 `_wasm.js` 文件
- 有效的微信小游戏 AppID
- 对目标目录的写权限

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 支持

如果您遇到任何问题或有疑问，请在 GitHub 上打开一个 issue。
