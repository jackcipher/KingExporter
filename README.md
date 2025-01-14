# KingExporter

[English](https://github.com/jackcipher/KingExporter/blob/master/README-en.md) | [中文](https://github.com/jackcipher/KingExporter/blob/master/README.md)

KingExporter 是一款为 WPS/金山文档设计的高性能文档导出工具。基于 Go 语言开发，通过多线程架构实现高效的批量文档迁移，支持企业级操作的同时保持文件夹结构完整性。

## 功能特性

### 性能优化
- 基于 Go goroutines 的多线程架构
- 智能通道管理（fan-in/fan-out 模式）
- 并行文件处理和下载
- 优化的资源利用

### 文件处理能力
- 标准文件直接下载
- 金山文档特有格式自动转换
- 保持原始文件夹结构
- 实时进度跟踪和状态更新

### 灵活的导出选项
- 个人空间导出
- 团队空间导出
- 全空间批量导出
- 静默模式支持

## 使用指南

### 安装
```bash
go get github.com/jackcipher/KingExporter
```

### 运行环境要求
- Go 1.16 或更高版本
- 有效的金山文档账号
- 足够的存储空间

### 身份认证设置

![](https://raw.githubusercontent.com/jackcipher/static_resource/refs/heads/master/imgs/kingexporter01.png)

使用 KingExporter 需要获取金山文档的 "wps_sid" cookie。获取步骤：

- 访问 "https://kdocs.cn/latest" 
- 打开开发者工具（F12 或 Ctrl/Cmd + Shift + I） 
- 导航至 应用程序 > Cookies > https://www.kdocs.cn
- 找到并复制 "wps_sid" 的值

### 使用示例

**导出个人文件**
```bash
KingExporter --sid=您的SID --download_dir=下载路径
```

**导出团队文件**
```bash
KingExporter --sid=您的SID --download_dir=下载路径 --group_id=团队ID
```

**导出所有可访问文件**
```bash
KingExporter --sid=您的SID --download_dir=下载路径 -A
```

### 命令行选项

| 选项 | 说明 | 是否必需 |
|--------|-------------|----------|
| --sid | 金山文档会话ID | 是 |
| --download_dir | 下载目标路径 | 是 |
| --group_id | 团队ID | 否 |
| -A | 导出所有可访问文件 | 否 |
| -s | 启用静默模式 | 否 |

## 技术细节

### 文件处理
- Office 格式文件直接下载
- 金山文档格式（.otl、.ksheet）自动转换为标准 Office 格式
- 保持原始目录结构

### 性能特性
- 转换和下载任务并发处理
- 不同任务类型独立通道
- 并行下载优化
- [ ] TODO: 实时进度监控

## 开发相关

### 贡献指南
我们欢迎各种形式的贡献！请遵循以下步骤：
1. Fork 项目仓库
2. 创建特性分支
3. 提交更改
4. 推送到您的分支
5. 创建 Pull Request

### 支持

- 问题反馈：通过 GitHub Issues 提交 bug 报告
- 讨论：加入我们的 GitHub Discussions
- 文档：查看我们的 Wiki 获取详细指南

## 许可证
本项目采用 MIT 许可证 - 详见 LICENSE 文件。