# nuxtblog CLI

NuxtBlog 官方命令行工具，用于插件开发、项目构建与发布管理。

## 安装

### 从 GitHub Releases 下载（推荐）

前往 [Releases](https://github.com/nuxtblog/nuxtblog-cli/releases) 页面，下载对应平台的预编译二进制：

| 平台 | 文件 |
|------|------|
| macOS (Apple Silicon) | `nuxtblog_*_darwin_arm64.tar.gz` |
| macOS (Intel) | `nuxtblog_*_darwin_amd64.tar.gz` |
| Linux (x86_64) | `nuxtblog_*_linux_amd64.tar.gz` |
| Linux (ARM64) | `nuxtblog_*_linux_arm64.tar.gz` |
| Windows (x86_64) | `nuxtblog_*_windows_amd64.zip` |

下载后解压，将 `nuxtblog`（Windows 下为 `nuxtblog.exe`）放到 `PATH` 目录中即可：

```bash
# macOS / Linux 示例
tar xzf nuxtblog_*_linux_amd64.tar.gz
sudo mv nuxtblog /usr/local/bin/

# 验证
nuxtblog --help
```

### 从源码安装

需要 Go 1.25+：

```bash
go install github.com/nuxtblog/nuxtblog-cli/cmd/nuxtblog@latest
```

或在项目目录内：

```bash
cd nuxtblog-cli
make install
```

## 使用

所有命令均在 **项目根目录**（包含 `go.work` 的目录）下执行。

### 构建项目

```bash
nuxtblog build
```

执行完整构建流水线：扫描 Go 插件 → 构建前端资源 → 同步到 server/builtin/ → 生成 plugins.go → 编译服务端。

### 插件管理

```bash
# 列出所有已安装的插件
nuxtblog plugin list

# 创建新插件（交互式向导）
nuxtblog plugin create

# 添加插件
nuxtblog plugin add <plugin-id>

# 发布插件到市场
nuxtblog plugin publish
```

#### `plugin list` 输出示例

```
ID                               Type     Runtime      Version  Status
--                               ----     -------      -------  ------
nuxtblog-plugin-ai-polish        builtin  compiled     2.0.0    bundled
nuxtblog-plugin-auto-excerpt     js       interpreted  2.0.0    installed
nuxtblog-plugin-comment-guard    builtin  compiled     2.0.0    bundled
nuxtblog-plugin-hello-js         js       interpreted  0.2.0    installed
nuxtblog-plugin-pinyin-slug      js       interpreted  2.0.0    installed
nuxtblog-plugin-reading-time     js       interpreted  2.0.0    installed
nuxtblog-plugin-telegram-notify  builtin  compiled     2.0.0    bundled
nuxtblog-plugin-view-counter     builtin  compiled     2.0.0    bundled
```

### Docker 部署

```bash
nuxtblog docker
```

## 本地开发

```bash
cd nuxtblog-cli

# 编译
make build

# 安装到 GOPATH/bin
make install

# 代码检查
make lint

# GoReleaser 本地预览（不发布）
make snapshot
```

## 发布新版本

推送 `v*` 标签后，GitHub Actions 会自动通过 GoReleaser 构建并发布所有平台的二进制：

```bash
git tag v1.0.0
git push origin v1.0.0
```

## 许可证

MIT
