# baidupan-cli

百度网盘命令行工具（交互式 shell，基于 `grumble`）。

## 功能

- **授权登录**：`auth`（默认自动打开浏览器扫码）
- **用户信息**：`userinfo`
- **容量查询**：`cap`
- **文件列表**：`fs` / `files`（支持排序/递归/筛选/表格输出）
- **单个重命名**：`rename`（文件/目录都支持）
- **批量重命名**：`rename-batch`（别名 `rb`，支持 sed 替换模式 + 正则模式 + 进度条 + 出错继续）

## 快速开始

### 配置文件

项目根目录准备一个配置文件（示例：`config.yaml`）：

```yaml
baidu-pan:
  app-id:
  app-key:
  secret-key:
  sign-key:
```

### 运行

- **运行二进制**（推荐，token 默认保存到二进制同目录）：

```bash
go build -o baidupan-cli .
./baidupan-cli -c config.yaml
```

- **go run**（会生成临时可执行文件，建议指定 token 固定目录）：

```bash
export BAIDUPAN_CLI_TOKEN_DIR="$PWD/.debug"
go run . -c config.yaml
```

### 授权

在交互 shell 里执行：

```bash
auth
```

- 默认会 **打开浏览器** 进行扫码授权；如果你不想打开浏览器：

```bash
auth --open-browser=false
```

## Token 持久化（免重复授权）

授权成功后会写入 `token.json`，并在启动时自动读取/过期自动刷新。

- **默认位置**：可执行文件同目录的 `token.json`
- **覆盖位置**（两者选一）：
  - **`BAIDUPAN_CLI_TOKEN_FILE`**：直接指定 token 文件路径
  - **`BAIDUPAN_CLI_TOKEN_DIR`**：指定目录，文件名固定 `token.json`

## 常用命令

### 文件列表（fs / files）

```bash
fs --dir "/我的文档" --only-folder -v
fs --dir "/我的文档" --recurse --limit 200 -v
```

### 单个重命名（rename）

> **注意**：文件/目录路径必须以 `/` 开头（绝对路径）。

```bash
rename --path "/我的文档/uml设计图" --newname "UML设计图"
```

### 批量重命名（rename-batch / rb）

`rb` 默认是 **sed 替换模式**：给两个位置参数 `FIND TO`，对名称做“包含替换”。

#### 示例 1：替换模式（默认）

把目录名 `UML设计图` 中的 `设计` 替换为 `分析`：

```bash
rb --dir "/我的文档" --target dirs 设计 分析
rb --dir "/我的文档" --target dirs 设计 分析 --apply
```

如需把 `FIND` 当作正则（而不是纯文本），加 `--find-regex`：

```bash
rb --dir "/我的文档" --target dirs --find-regex '设.' 分析 --apply
```

#### 示例 2：正则模式（--pattern/--replace）

把 `xxx.mp4` 改名为 `xxx_1080p.mp4`：

```bash
# 只是预览，不会执行
rb --dir "/video" --pattern '^(.*)\.mp4$' --replace '${1}_1080p.mp4'
# 真正执行
rb --dir "/video" --pattern '^(.*)\.mp4$' --replace '${1}_1080p.mp4' --apply
```

#### 示例 3：去掉中文【】里的内容

把 `xxx【xxx】xxx` 变成 `xxxxxx`：

```bash
rb --dir "/xxx" --target dirs --pattern '^(.*)【[^】]*】(.*)$' --replace '${1}${2}' --apply
```

#### 示例 5：去掉前导 0（012xxx -> 12xxx）

```bash
rb --dir "/xxx" --target dirs --pattern '^0+([0-9]+.*)$' --replace '${1}' --apply
```

## 批量执行与大数据量建议

当目录下条目很多时，单次请求可能慢/超时：

- 调小每次请求包含的条目数：`-s/--size`
- 或使用异步任务：`--async`
- 默认会显示进度与转圈：`-p/--progress`

示例：

```bash
rb --dir "/xxx" --target dirs 设计 分析 --apply -s 50
rb --dir "/xxx" --target dirs 设计 分析 --apply --async
```

## 出错继续（跳过失败项）

批量执行时如果希望“出错也继续处理后续”：

```bash
rb --dir "/xxx" --target dirs 设计 分析 --apply -c
```

如果希望即使有失败也返回成功退出码（适合脚本流水线）：

```bash
rb --dir "/xxx" --target dirs 设计 分析 --apply -c -i
```

## 常见错误排查

- **`errno=-6`**：常见是参数/路径问题（路径必须以 `/` 开头、目录不存在等）
- **`errno=12` + `info[].errno=-8`**：批量部分失败，通常是**重名冲突**（替换后目标名称已存在）。建议先 dry-run 看计划，或保留部分信息避免撞名。

## 批量重命名执行顺序说明（包含目录时）

当 `--target dirs` 或 `--target all` 且在递归目录下批量重命名时，为避免“先改父目录导致子路径失效”，程序会按以下顺序提交：

- 先重命名 **文件**
- 再重命名 **目录**（按路径深度从深到浅）
