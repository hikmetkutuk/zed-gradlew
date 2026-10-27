![GitHub repo size](https://img.shields.io/github/repo-size/hikmetkutuk/zed-gradlew?color=inactive&logo=github&style=for-the-badge)
![Go](https://img.shields.io/static/v1?logo=go&label=Go&message=1.26&color=00ADD8&style=for-the-badge)
![Rust](https://img.shields.io/static/v1?logo=rust&label=Rust&message=1.99&color=000000&style=for-the-badge)

> 🇹🇷 [Türkçe README](README.tr.md)

# zed-gradlew

A Gradle LSP extension for the Zed editor. Shows code lens buttons (`▶ bootRun`, `▶ test`) on `build.gradle` files and runs tasks directly from the editor.

## How it works

```
Zed → LSP (zed-gradlew binary) → ./gradlew tasks --all → code lens
                                → ./gradlew <task>     → task execution
```

The extension uses Zed's built-in LSP support. There is no custom UI panel or sidebar — code lens integrates into Zed's native interface.

## Installation

### 1. Build the binary and add it to PATH

```bash
cd server
go build -o zed-gradlew .
sudo mv zed-gradlew /usr/local/bin/
```

### 2. Install the extension in Zed

```
Zed → Extensions → Install Dev Extension → select the extension/ folder
```

## Requirements

- Go 1.21+
- Rust (for extension build)
- Zed 0.130+
- `gradlew` script in the project root

## Project structure

```
zed-gradlew/
├── server/          # Go LSP server
│   ├── main.go
│   ├── handler/     # LSP method implementations
│   └── gradle/      # Task discovery and execution
└── extension/       # Zed extension Rust shim
    ├── Cargo.toml
    ├── extension.toml
    ├── rust-toolchain.toml
    └── src/lib.rs
```

## Supported files

`build.gradle`, `build.gradle.kts`
