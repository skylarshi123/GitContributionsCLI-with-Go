# GitContributionsCLI 🚀

<div align="center">
  <img src="https://miro.medium.com/v2/resize:fit:1400/format:webp/1*WY7ELhXIVxbGlUwmhA1PSw.jpeg" alt="Git Contribution Graph Demo" />
  
  [![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
  [![Git](https://img.shields.io/badge/Git-Required-F05032?style=for-the-badge&logo=git)](https://git-scm.com/)
  [![Docker](https://img.shields.io/badge/Docker-Optional-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
  
  <p align="center">
    <strong>🔍 Local Git Analytics</strong> • <strong>📊 Contribution Visualization</strong> • <strong>📈 File Statistics</strong>
  </p>
</div>

A powerful command-line tool that brings GitHub-style contribution visualization to your local Git repositories with advanced analytics! 📊

## Features ✨

### Core Functionality 🎯
- **Local Repository Scanning** 🔍
  - Automatically discovers and tracks Git repositories
  - Maintains a clean configuration in ~/.gogitlocalstats
  - Smart directory traversal with vendor/node_modules exclusion

- **GitHub-Style Contribution Graph** 📅
  - Beautiful calendar heatmap visualization
  - 6-month contribution history
  - Day-of-week based layout
  - Real-time contribution tracking

- **File Type Analytics** 📈
  - Comprehensive file extension statistics
  - Most modified file types tracking
  - Detailed file count per extension
  - Sorted by frequency

### Technical Highlights 🛠️
- Written in pure Go
- Efficient repository processing
- Docker support
- Minimal dependencies

## Installation 📦

### Prerequisites
- Go 1.23 or higher
- Git
- (Optional) Docker

### Standard Installation
```bash
# Clone the repository
git clone https://github.com/skylarshi123/GitContributionsCLI-with-Go
cd GitContributionsCLI-with-Go

# Install dependencies
go mod download

# Build the project
go build
```

### Docker Installation 🐳
```bash
# Build the Docker image
docker build -t gitcontrib .
```

## Usage 💡

### Viewing Contributions
To view your contribution statistics:

```bash
# Standard usage
go run main.go scan.go stats.go --email "your@email.com"

# Docker usage
docker run -it \
  -v $HOME:/root \
  -v /path/to/your/repos:/repos \
  gitcontrib --email "your@email.com"
```

## Output Example 🎨

```
             Jul                 Aug             Sep                 Oct             Nov             Dec             
       -   -   -   -   -   -   1   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -  11   - 
 Fri   -   -   -   -   -   2   -   1   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   2   - 
       -   -   -   -   -   -   -   -   1   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   9 
 Wed   -   -   -   -   -   3   -   1   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   - 
       -   -   -   -   -   1   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   - 
 Mon   -   -   -   -   -   -   2   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   - 
       -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   -   - 

File Type Statistics:
===================
.map              364 files
.css              299 files
.jpg              287 files
.js               273 files
.go                55 files
.html              52 files
.png               51 files
.txt               39 files
.md                28 files
.gitignore         13 files
```

## How It Works 🔧
The tool works by:
1. Scanning your local Git repositories
2. Processing commit history with your email
3. Generating a visual contribution graph
4. Analyzing file types across repositories
5. Presenting statistics in a clear, formatted output

## Contributing 🤝
Contributions are welcome! Feel free to:
- Open issues for bugs or feature requests
- Submit pull requests
- Improve documentation
- Share feedback

## License 📄
MIT License - feel free to use and modify as you wish!

## Acknowledgments 🙏
- Inspired by GitHub's contribution graph
- Built with Go's powerful standard library
- Uses [go-git](https://github.com/go-git/go-git) for Git operations

## To-Do 📝
Future enhancements planned:
- [ ] Interactive CLI interface
- [ ] Contribution streak tracking
- [ ] Multiple email support
- [ ] Custom date range selection
- [ ] JSON/CSV export options

---
⭐ If you find this tool useful, please consider giving it a star!

