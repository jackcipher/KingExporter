# KingExporter

[English](https://github.com/jackcipher/KingExporter/blob/master/README-en.md) | [中文](https://github.com/jackcipher/KingExporter/blob/master/README.md)


KingExporter is a high-performance, multi-threaded tool designed for batch exporting documents from WPS/金山文档. It facilitates efficient document migration by preserving folder structures and supporting enterprise-scale operations.Built with Go, it leverages concurrent pattern to optimize file exports.



## Features

### Performance Optimizations
- Multi-threaded architecture utilizing Go's goroutines
- Intelligent channel management with fan-in/fan-out patterns
- Parallel file processing and downloads
- Optimized resource utilization

### File Handling Capabilities
- Direct downloads for standard files
- Automatic format conversion for KDocs-specific formats
- Preservation of original folder structure
- Real-time progress tracking and status updates

### Flexible Export Options
- Personal workspace export
- Team workspace export
- Bulk export across all accessible spaces
- Silent mode for automated operations

## Getting Started

### Installation
```bash
go get github.com/jackcipher/KingExporter
```

### Prerequisites
- Go 1.16+
- Valid KDocs account
- Sufficient storage space

### Authentication Setup

To use KingExporter, you'll need your KDocs "wps_sid" cookie. Here's how to find it:

![](https://raw.githubusercontent.com/jackcipher/static_resource/refs/heads/master/imgs/kingexporter01.png)

- Visit "https://kdocs.cn/latest" in your browser
- Open Developer Tools (F12 or Ctrl/Cmd + Shift + I)
- Navigate to Application > Cookies > https://www.kdocs.cn
- Locate and copy the "wps_sid" value

### Usage Examples

**Export Personal Files**
```bash
KingExporter --sid=YOUR_SID --download_dir=PATH_TO_DOWNLOAD
```

**Export Team Files**
```bash
KingExporter --sid=YOUR_SID --download_dir=PATH_TO_DOWNLOAD --group_id=YOUR_GROUP_ID
```

**Export All Accessible Files**
```bash
KingExporter --sid=YOUR_SID --download_dir=PATH_TO_DOWNLOAD -A
```

### Command Line Options

| Option | Description | Required |
|--------|-------------|----------|
| --sid | KDocs session ID | Yes |
| --download_dir | Download destination | Yes |
| --group_id | Team ID for export | No |
| -A | Export all accessible files | No |
| -s | Enable silent mode | No |

## Technical Details

### File Processing
- Direct download for Office formats
- Automatic conversion of KDocs formats (.otl, .ksheet) to standard Office formats
- Maintains original directory structure

### Performance Features
- Concurrent processing of conversion and download tasks
- Independent channels for different task types
- Parallel download optimization
- [ ] TODO: Real-time progress monitoring

## Development

### Contributing
We welcome contributions! Please follow these steps:
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to your branch
5. Create a Pull Request

### Support

- **Issues**: Submit bug reports via GitHub Issues
- **Discussions**: Join our GitHub Discussions
- **Documentation**: Check our Wiki for detailed guides

## License
This project is licensed under the MIT License - see the LICENSE file for details.