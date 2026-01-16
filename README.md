
# secret-scan

A fast, multi-threaded CLI tool for scanning files and directories to detect exposed secrets, credentials, and sensitive information.

## Features

- 🚀 **Multi-threaded scanning** - defaults to CPU cores - 1
- 📄 **Multiple file formats** - txt, log, yaml, json, csv, etc, xlsx
- 🔌 **Extensible patterns** - Pattern definitions via Lua scripts
- 🔌 **Pattern validation** - Pattern validation support via calculated string entropy or paired line contents 
- 📊 **JSON output** - Structured findings for easy parsing and integration
- 🎯 **Configurable** - Control verbosity, threading, patterns, ignores, and output location

## Installation
### Clone the repository
git clone <your-repo-url> cd secret-scan

### Build
go build -o secret-scan ./cmd

## Usage

### Scan current directory using embedded patterns
./secret-scan
### Scan specific path
./secret-scan /path/to/scan
### Scan with verbose output
./secret-scan -verbose /path/to/scan

## Output Options
### Output to file
./secret-scan -out findings.json /path/to/scan
### Pipe to jq for filtering
./secret-scan /path/to/scan | jq '.severity == "critical"'
### Silent mode (errors only)
./secret-scan -silent /path/to/scan > findings.json

## Performance Tuning
./secret-scan -threads 4 /path/to/scan

## Command-Line Flags

| Flag                   | Description                              | Default |
|------------------------|------------------------------------------|-----|
| `-verbose`             | Enable verbose debug output              | `false` |
| `-silent`              | Suppress all output except errors        | `false` |
| `-out `                | Write findings to file instead of stdout | stdout |
| `-threads`             | Number of worker threads                 | CPU cores - 1 |
| `-patterns`            | define a patterns directory              | ""|
| `-no-default-patterns` | excludes embedded patterns               | `false`|
| `-ignore`              | define a git-style ignore file           | ""|

## Output Format

Findings are output as JSON, one per line:

Patterns are defined in `patterns/patterns.lua`. See the file for examples of how to add custom patterns.

## Supported File Types

- **Text files**: `.txt`, `.log`, `.yaml`, `.yml`, `.json`, `.md`, `.conf`, `.cfg`, `.csv`
- **Excel files**: `.xlsx`
- **default ignores**: ".git/", ".svn/", ".hg/", "node_modules/", "vendor/", "*.exe", "*.dll", "*.so", "*.dylib", "*.zip", "*.tar", "*.tar.gz", "*.rar", "*.7z", "*.bin", "*.class", "*.jar", "*.war", "*.ear", "*.pyc", "*.pyo", "*.o", "*.a", "*.lib", "*.obj", "*.iso", "*.img", "*.dmg"


## Examples
### Scan a web application directory
./secret-scan -verbose /var/www/myapp
### Scan and save high-severity findings
./secret-scan /project | jq 'select(.severity == "critical" or .severity == "high")' > critical-findings.json
### Quick scan with minimal output
./secret-scan -silent -out report.json ~/Documents
### Use custom patterns
./secret-scan -patterns ./my-patterns/ /path/to/scan
./secret-scan -patterns ~/.secret-scan/patterns/ /path/to/scan

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
## Roadmap
 - plugin support for additional extractors
 - plugin support for additional validators
