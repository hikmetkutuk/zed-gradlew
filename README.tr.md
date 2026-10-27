![GitHub repo size](https://img.shields.io/github/repo-size/hikmetkutuk/zed-gradlew?color=inactive&logo=github&style=for-the-badge)
![Go](https://img.shields.io/static/v1?logo=go&label=Go&message=1.26&color=00ADD8&style=for-the-badge)
![Rust](https://img.shields.io/static/v1?logo=rust&label=Rust&message=1.99&color=000000&style=for-the-badge)

> 🇬🇧 [English README](README.md)

# zed-gradlew

Zed editor için Gradle LSP extension'ı. `build.gradle` dosyalarında code lens butonları (`▶ bootRun`, `▶ test`) gösterir ve görevi doğrudan editörden çalıştırır.

## Nasıl çalışır

```
Zed → LSP (zed-gradlew binary) → ./gradlew tasks --all → code lens
                                → ./gradlew <task>     → görev çalıştırma
```

Extension, Zed'in LSP desteğini kullanır. Özel UI paneli veya sidebar yoktur — code lens Zed'in kendi arayüzüne entegre olur.

## Kurulum

### 1. Binary'yi derle ve PATH'e ekle

```bash
cd server
go build -o zed-gradlew .
sudo mv zed-gradlew /usr/local/bin/
```

### 2. Extension'ı Zed'e yükle

```
Zed → Extensions → Install Dev Extension → extension/ klasörünü seç
```

## Gereksinimler

- Go 1.21+
- Rust (extension build için)
- Zed 0.130+
- Projenin kök dizininde `gradlew` betiği

## Proje yapısı

```
zed-gradlew/
├── server/          # Go LSP sunucusu
│   ├── main.go
│   ├── handler/     # LSP method implementasyonları
│   └── gradle/      # Görev keşfi ve çalıştırma
└── extension/       # Zed extension Rust shim'i
    ├── Cargo.toml
    ├── extension.toml
    ├── rust-toolchain.toml
    └── src/lib.rs
```

## Desteklenen dosyalar

`build.gradle`, `build.gradle.kts`
