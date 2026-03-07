# Port Manager - 端口管理工具

一个高性能的跨平台端口管理桌面应用，使用 Go + Wails + Svelte 构建。

## ✨ 功能特性

- 🔍 **端口扫描** - 快速扫描系统中所有 TCP/UDP 端口
- 💀 **进程管理** - 查看占用端口的进程并一键杀死
- 📊 **使用统计** - 追踪端口使用历史和时长统计
- 🎨 **丝滑动画** - 高性能 60fps+ 流畅动画
- 💾 **轻量级** - 内存占用仅 30-50MB，体积 ~10-15MB
- 🌙 **暗色主题** - 现代化暗色 UI 设计

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- Wails CLI v2

### 安装 Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 开发模式

```bash
# 克隆仓库
git clone https://github.com/zarttic/port-manager.git
cd port-manager

# 安装前端依赖
cd frontend
npm install

# 返回根目录并启动开发模式
cd ..
wails dev
```

### 构建生产版本

```bash
wails build
```

构建产物位于 `build/bin/` 目录。

## 📦 下载安装

从 [Releases](https://github.com/zarttic/port-manager/releases) 页面下载适合您系统的版本：

- **Windows**: `port-manager.exe` (直接运行) 或 `port-manager-windows-amd64.zip`
- **macOS (Intel)**: `port-manager-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `port-manager-darwin-arm64.tar.gz`

## 🏗️ 项目结构

```
port-manager/
├── backend/                 # Go 后端代码
│   ├── main.go             # 应用入口
│   ├── app.go              # 应用逻辑
│   ├── internal/           # 内部模块
│   │   ├── api/           # API 层
│   │   ├── service/       # 业务逻辑层
│   │   ├── repository/    # 数据访问层
│   │   └── model/         # 数据模型
│   ├── pkg/               # 可导出包
│   │   └── sysinfo/       # 系统信息获取
│   └── database/          # 数据库文件
├── frontend/               # Svelte 前端代码
│   ├── src/
│   │   ├── components/    # UI 组件
│   │   ├── lib/          # 工具库
│   │   └── styles/       # 样式文件
│   └── package.json
├── wails.json             # Wails 配置
└── .github/
    └── workflows/
        └── release.yml    # 自动发布工作流
```

## 🔧 技术栈

### 后端
- **Go 1.21** - 核心语言
- **Wails v2** - 桌面应用框架
- **SQLite** - 数据持久化
- **gopsutil** - 系统信息库

### 前端
- **Svelte 4** - UI 框架
- **TypeScript** - 类型安全
- **TailwindCSS** - 样式系统
- **ECharts** - 数据可视化

## 🤝 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Wails](https://wails.io/) - 优秀的 Go 桌面应用框架
- [Svelte](https://svelte.dev/) - 简洁高效的前端框架
- [gopsutil](https://github.com/shirou/gopsutil) - 强大的跨平台系统信息库
