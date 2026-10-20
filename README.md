# 🚀 Luminious Auto-Boost (Discord Server Boost Automation)

<div align="center">

```text
  ██╗     ██╗   ██╗███╗   ███╗██╗███╗   ██╗██╗ ██████╗ ██╗   ██╗███████╗
  ██║     ██║   ██║████╗ ████║██║████╗  ██║██║██╔═══██╗██║   ██║██╔════╝
  ██║     ██║   ██║██╔████╔██║██║██╔██╗ ██║██║██║   ██║██║   ██║███████╗
  ██║     ██║   ██║██║╚██╔╝██║██║██║╚██╗██║██║██║   ██║██║   ██║╚════██║
  ███████╗╚██████╔╝██║ ╚═╝ ██║██║██║ ╚████║██║╚██████╔╝╚██████╔╝███████║
  ╚══════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═╝ ╚═════╝  ╚═════╝ ╚══════╝
```

**High-Performance • Multi-Account • Transfer-Aware • Safe Rate-Limit**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue)](https://github.com/Luminiou4/Auto-boost-)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

[Tiếng Việt](#-tiếng-việt) | [English](#-english)

</div>

---

## 🇻🇳 Tiếng Việt

**Luminious Auto-Boost** là công cụ dòng lệnh (CLI) bằng ngôn ngữ **Golang** hỗ trợ tự động hóa việc tăng cấp (Boost) máy chủ Discord từ danh sách nhiều tài khoản (Tokens). Công cụ được thiết kế chuyên nghiệp, xử lý thông minh tình trạng token, slot đang bận ở server khác, cooldown và xuất báo cáo đa định dạng.

### ✨ Tính Năng Nổi Bật

- ⚡ **Kiểm tra hàng loạt (Multi-threaded Checker):** Quét nhanh toàn bộ tài khoản với số luồng song song có thể tùy chỉnh trong config.
- 🔍 **Nhận diện Nitro & Trạng thái Slot chuyên sâu:**
  - Phát hiện gói Nitro (Nitro, Classic, Basic).
  - Phân loại chính xác từng slot: **SẴN SÀNG (Ready)**, **ĐANG DÙNG (Busy)** ở server khác, hoặc **ĐANG COOLDOWN (Chờ thời gian hồi)**.
  - Hiển thị đồng hồ đếm ngược cooldown chi tiết (ngày, giờ, phút).
- 🔄 **Hỗ trợ chuyển Boost thông minh (Smart Transfer):** Tự động rút boost đang gán ở server cũ để đập sang server mục tiêu theo yêu cầu.
- 🛡️ **Bảo vệ tài khoản & Chống Rate-Limit:**
  - Tự động ngẫu nhiên hóa độ trễ giữa các lần boost (`min_boost_delay` - `max_boost_delay`).
  - Xử lý mã lỗi `HTTP 429` (Too Many Requests), tự động đọc `X-RateLimit-Reset-After` để nghỉ và thử lại.
  - Hỗ trợ xoay vòng Proxy (HTTP & SOCKS5) tránh bị chặn dải IP.
- 📊 **Xuất báo cáo 3 định dạng tự động:**
  - `JSON Lines`: Chứa dữ liệu gốc chi tiết của từng lượt boost.
  - `TXT Report`: Bảng báo cáo trực quan, thống kê tỷ lệ thành công và nguyên nhân lỗi.
  - `CSV (Excel)`: Có sẵn UTF-8 BOM, mở bằng Microsoft Excel không bao giờ bị lỗi font.

---

### ⚙️ Cấu Hình (`config.json`)

```json
{
  "tokens_file": "tokens.txt",
  "proxy_file": "",
  "log_dir": "logs",
  "max_concurrent": 10,
  "min_boost_delay": "1.5s",
  "max_boost_delay": "3.5s",
  "check_delay": "800ms",
  "max_retries": 3,
  "base_url": "https://discord.com/api/v10",
  "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
}
```

- `tokens_file`: Đường dẫn file token (mỗi dòng một token).
- `proxy_file`: File proxy dạng `ip:port` hoặc `user:pass@ip:port` (để trống nếu không dùng).
- `log_dir`: Thư mục lưu file báo cáo kết quả sau khi chạy.
- `max_concurrent`: Số luồng check tài khoản song song.
- `min_boost_delay` / `max_boost_delay`: Khoảng thời gian nghỉ ngẫu nhiên giữa 2 lượt boost (ví dụ: `"1.5s"`, `"3.5s"`).
- `check_delay`: Thời gian nghỉ giữa các lượt kiểm tra token (`"800ms"`).

---

### 🚀 Hướng Dẫn Cài Đặt & Sử Dụng

1. **Cài đặt Golang:** Tải và cài đặt Go (phiên bản 1.21 trở lên) từ [golang.org](https://go.dev/dl/).
2. **Tải mã nguồn:**
   ```bash
   git clone https://github.com/Luminiou4/Auto-boost-.git
   cd Auto-boost-
   ```
3. **Cài đặt các thư viện phụ thuộc:**
   ```bash
   go mod tidy
   ```
4. **Chuẩn bị Token:**
   - Sao chép file `tokens.example.txt` thành `tokens.txt`.
   - Dán danh sách token Discord vào file `tokens.txt` (mỗi dòng 1 token).
5. **Chạy trực tiếp:**
   ```bash
   go run main.go
   ```
   Hoặc biên dịch thành file thực thi `.exe`:
   ```bash
   go build -ldflags="-s -w" -o luminious.exe .
   .\luminious.exe
   ```

---

## 🇬🇧 English

**Luminious Auto-Boost** is a high-performance Discord Server Boost automation tool built with **Golang**. It allows you to check token statuses, detect Nitro tiers, track cooldown timers, transfer existing boosts from other servers, and safely execute boosting tasks with rate-limit protections and proxy support.

### ✨ Key Features

- ⚡ **Concurrent Checker:** Check hundreds of tokens rapidly with configurable goroutine workers.
- 🔍 **Accurate Status Breakdown:**
  - Nitro subscription detection (Nitro, Classic, Basic).
  - Slot status classification: **Ready**, **Busy** (boosting other servers), and **Cooldown**.
  - Detailed countdown timer for slots currently on cooldown.
- 🔄 **Smart Boost Transfer:** Safely transfer active boosts from previous guilds to the target server.
- 🛡️ **Anti-Detection & Rate-Limit Safe:**
  - Randomized delay interval between boosts to simulate natural human activity.
  - Automatic `HTTP 429` rate-limit backoff using Discord's `Retry-After` header.
  - Proxy rotation (HTTP/SOCKS5) to protect your home IP address.
- 📊 **Multi-Format Export:**
  - `JSON Lines`: Machine-readable raw execution logs.
  - `TXT Report`: Human-readable summary table and failure details.
  - `CSV (Excel)`: Pre-configured with UTF-8 BOM for clean spreadsheet import.

---

## ⚠️ CẢNH BÁO / DISCLAIMER

> [!WARNING]
> **VI PHẠM ĐIỀU KHOẢN DỊCH VỤ DISCORD (ToS):**
> 
> Công cụ này được thiết kế và chia sẻ **CHỈ NHẰM MỤC ĐÍCH HỌC TẬP, NGHIÊN CỨU KỸ THUẬT VÀ KIỂM THỬ BẢO MẬT (Educational & Security Research Only)**.
> 
> Việc tự động hóa tài khoản người dùng cá nhân (Self-bot) là **hành vi vi phạm nghiêm trọng Điều khoản Dịch vụ (Terms of Service - ToS) và Chính sách Nhà phát triển của Discord**. Hành động này có thể dẫn đến việc tài khoản bị khóa vĩnh viễn (Terminated) hoặc bị cấm số điện thoại/IP.
> 
> **Tác giả hoàn toàn KHÔNG chịu trách nhiệm** cho bất kỳ thiệt hại nào, bao gồm nhưng không giới hạn ở việc mất tài khoản, khóa thẻ thanh toán, hay bất kỳ hình thức xử lý kỷ luật nào từ Discord. Bạn hoàn toàn tự chịu trách nhiệm khi quyết định chạy công cụ này.

---

## 💼 NHẬN CODE BOT & PHẦN MỀM THEO YÊU CẦU / CUSTOM DEV SERVICES

Bạn cần phát triển bot, phần mềm tự động hóa hoặc công cụ MMO tùy chỉnh?
- 🤖 **Discord Bot / Telegram Bot** (quản lý cộng đồng, thông báo biến động, auto tasks, mini apps).
- 🌐 **Web Scraper & Crawler** (cào dữ liệu tự động, vượt Cloudflare, reCAPTCHA).
- ⚡ **Auto-click, Automation Tools** bằng Golang, Python, C#, Node.js.
- 🔗 **Blockchain / Web3 Bots** (Sniper bot, DEX bot, volume maker, wallet tracking).
- 🛠️ **Reverse Engineering & Custom API Integration**.

📩 **Liên hệ đặt code:**
- **GitHub:** [@Luminiou4](https://github.com/Luminiou4)
- Mở một **Issue** hoặc liên hệ trực tiếp qua thông tin trên hồ sơ GitHub.

---

## ☕ ỦNG HỘ TÁC GIẢ / DONATION & SPONSORSHIP

Nếu công cụ này hữu ích cho bạn, hãy cân nhắc mời tác giả một ly cà phê qua các địa chỉ ví tiền mã hóa dưới đây. Sự ủng hộ của bạn là động lực lớn để tiếp tục nâng cấp và bảo trì dự án! ❤️

| Blockchain | Địa Chỉ Ví (Wallet Address) |
|---|---|
| **🟣 Solana (SOL / SPL)** | `HfMx9usAgnrTfstFCHU9kEWG9wBA2PV4SZdRwPT56zMP` |
| **🔷 Ethereum (ETH / ERC-20)** | `0x200a7e99946b32f06874a8c55e05e923b2fc6423` |
| **🟡 BNB Chain (BNB / BEP-20)** | `0x200a7e99946b32f06874a8c55e05e923b2fc6423` |
| **🟠 Bitcoin (Taproot - P2TR)** | `bc1pu52zpwzeauluxl7fs7f3zrrlmalamd0d06mg8z4qde87hwemswds7xts56` |
| **🟠 Bitcoin (Native SegWit - bech32)** | `bc1q7kn2pwwg7k9n5fnulzs0yrqcg5szgputvsevs9` |

---

<div align="center">
  <sub>Phát triển với sự tận tâm bởi <b>Luminious</b>. Cảm ơn các bạn đã ghé thăm!</sub>
</div>
