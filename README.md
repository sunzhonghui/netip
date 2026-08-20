# NetIP - 现代化自托管网络诊断与 IP 工具箱

<p align="center">
  <img src="web/public/favicon.svg" alt="NetIP Logo" width="80" height="80" />
</p>

<p align="center">
  <b>快速、高颜值、轻量、自托管、IPv6 友好的网络诊断工具箱与 IPW.cn 替代方案</b>
</p>

<p align="center">
  <a href="https://github.com/sunzhonghui/netip/releases"><img src="https://img.shields.io/github/v/release/sunzhonghui/netip?color=2563eb&label=Release" alt="Release"></a>
  <a href="https://hub.docker.com/r/sunzhonghui/netip"><img src="https://img.shields.io/docker/pulls/sunzhonghui/netip?color=2563eb&label=Docker%20Pulls" alt="Docker Pulls"></a>
  <img src="https://img.shields.io/badge/Architecture-amd64%20%7C%20arm64-blue" alt="Multi-Arch">
  <img src="https://img.shields.io/badge/IPv6-Supported-purple" alt="IPv6 Ready">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-emerald" alt="License"></a>
</p>

---

## ⚡ 一键 Docker 极速启动 (Quick Start)

### 方式一：Docker CLI 一行命令启动（最推荐）

```bash
# 使用 host 网络模式启动（支持原生 IPv4/IPv6 双栈与低延迟 ICMP Ping）
docker run -d \
  --name netip \
  --restart unless-stopped \
  --network host \
  -v $(pwd)/data/ipdb:/data/ipdb \
  sunzhonghui/netip:latest
```

启动完成后，打开浏览器访问 `http://你的服务器IP:8080` 即可！

---

### 方式二：Docker Compose 启动

创建 `docker-compose.yml`：

```yaml
services:
  netip:
    image: sunzhonghui/netip:latest
    container_name: netip
    restart: unless-stopped
    network_mode: host
    volumes:
      # 持久化离线 IP 库（容器会自动定时拉取最新 ip2region / MaxMind 数据库）
      - ./data/ipdb:/data/ipdb
    environment:
      - HTTP_ADDR=:8080
      - TRUSTED_PROXIES=127.0.0.1,::1
      - IPDB_AUTO_UPDATE=true
      - IPDB_UPDATE_INTERVAL=168h
```

执行启动：
```bash
docker compose up -d
```

---

### 方式三：源码一键部署（包含 Caddy 自动配置 HTTPS 证书）

```bash
# 1. 克隆代码并进入项目目录
git clone https://github.com/sunzhonghui/netip.git
cd netip

# 2. 复制并调整环境变量
cp .env.example .env

# 3. 一键构建并启动 NetIP + Caddy 反向代理
docker compose up -d --build
```

---

## 🌐 域名与 DNS 解析配置规范

为了让双栈网络检测与纯文本脚本接口正常工作，建议在您的域名 DNS 解析服务商处完成以下配置（以 `ipw.3x.cx` 为例）：

| 主机记录 | 类型 | 记录值 | 说明 |
| :--- | :--- | :--- | :--- |
| `ip.ipw.3x.cx` | **A** | 服务器公网 IPv4 | 主站 Web UI 与 REST API (双栈) |
| `ip.ipw.3x.cx` | **AAAA** | 服务器公网 IPv6 | 主站 Web UI 与 REST API (双栈) |
| `4.ipw.3x.cx` | **A** | 服务器公网 IPv4 | **仅配置 A 记录**（不得配置 AAAA），供 IPv4 探测与脚本使用 |
| `6.ipw.3x.cx` | **AAAA** | 服务器公网 IPv6 | **仅配置 AAAA 记录**（不得配置 A），供 IPv6 探测与脚本使用 |
| `test.ipw.3x.cx` | **A** | 服务器公网 IPv4 | 双栈优先级检测（同时配置 A 与 AAAA） |
| `test.ipw.3x.cx` | **AAAA** | 服务器公网 IPv6 | 双栈优先级检测（同时配置 A 与 AAAA） |

> **⚠️ 注意**：`4.` 子域名**切勿添加 AAAA 记录**，`6.` 子域名**切勿添加 A 记录**。请直接将域名 A/AAAA 指向服务器 IP，不要开启 CDN 代理（如 Cloudflare 橙色云朵），避免客户端真实握手协议与真实 IP 被隐藏。

---

## ✨ 核心特性

* 🚀 **即开即测**：首页并发秒级检测当前网络公网 IPv4、IPv6 及协议访问优先级（IPv6 优先 / IPv4 优先）。
* 🎨 **高质感极光蓝 UI**：现代化磨砂玻璃卡片、微交互触感反馈、实时雷达脉冲绿灯，支持深色/浅色自适应主题。
* 📜 **极简 IP API**：专为 Shell、DDNS、路由固件打造的纯文本 IP 输出（`curl https://4.yourdomain.com` / `curl https://6.yourdomain.com` / `/ip` / `/json`）。
* 🧭 **多 DNS 并发对比解析**：并发向系统 DNS、Cloudflare (1.1.1.1)、Google (8.8.8.8)、AliDNS (223.5.5.5)、DNSPod (119.29.29.29) 查询 A、AAAA、CNAME、MX、TXT、NS 等。
* 📡 **双模 Ping 诊断**：支持原生底层 ICMP Echo 丢包率与时延测试，以及基于传输层端口握手的 TCPing。
* 🌐 **网站 IPv6 双栈全面体检**：自动验证 DNS A/AAAA、HTTP/HTTPS 双栈连通性并输出就绪评估。
* 🔒 **SSL / TLS 证书深度查询**：解析 X.509 证书链、剩余天数、SANs 多域名、TLS 协议版本与密码套件。
* 🔍 **WHOIS / RDAP 域名与 IP 查询**：现代化 RDAP 协议优先，权威 WHOIS 兜底，结构化展示注册信息。
* ⚡ **网站 HTTP 测速与瀑布流分析**：精确拆解 DNS、TCP、TLS、TTFB 首包耗时，并在 5MB 安全限额内测算实际下载速率 (Mbps)。
* 🗄️ **持久库自动定时更新与热加载**：内置后台安全拉取 `ip2region.xdb`、`GeoLite2-ASN.mmdb` 等，0ms 宕机无感热重载。
* 🛡️ **军工级 SSRF & DNS Rebinding 防护**：严格阻断所有内网保留地址、Link-Local、CGNAT 与私有网段。
* 📦 **单二进制极简部署**：Go 1.25+ 原生全栈，Vue 3 前端静态资源内嵌于单一二进制文件中，无需任何外部数据库。

---

## 💻 命令行与终端快捷调用

```bash
# 获取公网 IPv4
curl https://4.ipw.3x.cx

# 获取公网 IPv6
curl https://6.ipw.3x.cx

# 测试当前网络首选协议 IP
curl https://test.ipw.3x.cx

# 获取 JSON 格式响应
curl https://4.ipw.3x.cx/json

# 查询当前客户端 IP 与归属地
curl https://ip.ipw.3x.cx/api/v1/me

# 查询指定 IP
curl https://ip.ipw.3x.cx/api/v1/ip/8.8.8.8
```

---

## 🛠️ 技术栈

* **后端**：Go 1.25+ / Gin / `miekg/dns` / `golang.org/x/net/icmp` / `log/slog`
* **前端**：Vue 3 / TypeScript / Vite / Tailwind CSS / Pinia / Lucide Icons
* **容器**：Docker Hub (`sunzhonghui/netip`) / Multi-Arch (`linux/amd64`, `linux/arm64`)

---

## 📄 开源许可证

本项目基于 [MIT License](LICENSE) 开源。欢迎 Star 与 Fork 支持！
