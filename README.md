# DW

Defold for wechat.

Export Defold bundle to WeApp Mini Games.

这是一个可以将Defold打包好的H5包转换到微信小游戏平台的工具

# Get start

- Download `dw` from release page.
- Exculate the shell from your terminal.

```bash
dw -source defold_wasm_built_dir -target the_weeapp_dir -appid weappid
```
Example:
```bash
dw -source /e/Game/defold/clgy/wasm-web/clgy -target /e/Game/wx -appid 12312fsa123ss
```
