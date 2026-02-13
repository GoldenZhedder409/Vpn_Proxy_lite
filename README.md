# 🔥 Vpn_Proxy_lite 🔥

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go" />
  <img src="https://img.shields.io/badge/Platform-Android%20|%20Linux%20|%20macOS-blue?style=for-the-badge" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" />
</p>

<p align="center">
  <b>🌍 Rotate proxies across 13+ countries automatically with mystical patterns! 🌍</b>
</p>

---

## 📖 **Table of Contents**
- [What is This?](#-what-is-this)
- [Why This Project?](#-why-this-project)
- [Features](#-features)
- [How It Works](#-how-it-works)
- [Installation](#-installation)
- [Usage](#-usage)
- [Proxy List](#-proxy-list)
- [Mystical Patterns](#-mystical-patterns)
- [Stealth Mode](#-stealth-mode)
- [Configuration](#-configuration)
- [Requirements](#-requirements)
- [Contributing](#-contributing)
- [Disclaimer](#-disclaimer)
- [License](#-license)

---

## 🤔 **What is This?**

**Ultimate Proxy Rotator** is a powerful, educational Go program that **automatically rotates your IP address** through a pool of **30+ public proxies from 13+ countries** including USA, Germany, Japan, Brazil, India, Singapore, Netherlands, Canada, and more! 

It creates a **local SOCKS5 proxy server** on your device that cycles through different countries at your chosen interval, making it appear as if you're browsing from different locations around the world every few seconds!

---

## 🎯 **Why This Project?**

### **Educational Purpose** 📚
This project was created to demonstrate:
- **Network programming** in Go (sockets, proxies, protocols)
- **Concurrent programming** with goroutines and channels
- **SOCKS5 protocol** implementation and testing
- **Proxy rotation** techniques used in real-world applications
- **Clean UI** development in terminal applications

### **Practical Applications** 💡
- 🕵️ **Privacy enhancement** - Hide your real IP address
- 🌐 **Bypass geo-restrictions** - Access content from different countries
- 📊 **Web scraping** - Rotate IPs to avoid rate limiting
- 🎓 **Learning networking** - Understand how proxies work
- 😄 **Fun** - Watch your IP dance around the world with mystical patterns!

---

## ✨ **Features**

| Feature | Description |
|---------|-------------|
| 🌍 **30+ Proxies** | From 13+ countries including USA, Germany, Japan, Brazil, India, Singapore, Netherlands, Canada, UK, Australia |
| 🔄 **Auto-Rotation** | Automatically switches proxy every X seconds (configurable 5-60s) |
| 🎭 **Mystical Patterns** | Each rotation gets a fun pattern name (Bermuda Triangle, Area 51 Raid, etc.) |
| 🧠 **Health Checking** | Automatically tests proxies and only uses alive ones |
| ⚡ **Latency Display** | Shows ping time for each proxy |
| 📊 **Live Status** | Displays healthy proxy count in real-time |
| 🎨 **Clean UI** | Beautiful terminal interface with emojis and formatting |
| ⚙️ **Configurable** | Change interval, port, and selected countries via menu |
| 💾 **Save Settings** | Configuration saved to JSON file |
| 🚀 **Lightweight** | Single binary, no dependencies |
| 📱 **Cross-Platform** | Works on Android (Termux), Linux, macOS |

---

## 🔧 **How It Works**

1. **Your browser** connects to `127.0.0.1:1080` (local proxy)
2. **The program** forwards traffic to a random healthy proxy
3. **Every X seconds**, it switches to a different proxy
4. **Your IP appears** to change locations automatically!

---

## ⚠️ **Disclaimer** (Made by me "GolDer409")
THIS SOFTWARE IS PROVIDED FOR EDUCATIONAL PURPOSES ONLY!

- The proxies included are PUBLIC proxies found on the internet
- We do not own or operate any of these proxies
- Use at your own risk
- Do not use for illegal activities
- Some websites may block proxy traffic
- Respect websites' terms of service
- The author is not responsible for misuse

By using this software, you agree to these terms.


## 📦 **Installation**

### **On Android (Termux)**
```bash
# Update packages
pkg update && pkg upgrade

# Install Go
pkg install golang

# Clone or create the file
mkdir ~/proxy-rotator
cd ~/proxy-rotator
nano proxy-rotator.go
# (Copy and paste the code)

# Run it!
go run proxy-rotator.go
