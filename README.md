# NetIP - 现代化自托管网络诊断与 IP 工具箱

<p align="center">
  <b>快速、干净、轻量、自托管、IPv6 友好的网络诊断工具箱与 IPW.cn 替代方案</b>
</p>

---

## ⚡ 30 秒快速部署 (Quick Start)

仅需两步即可在拥有公网 IPv4 / IPv6 的 Linux 服务器上完整启动：

```bash
# 1. 克隆代码并进入项目目录
git clone https://github.com/netip-org/netip.git
cd netip

# 2. 复制并调整环境变量
cp .env.example .env

# 3. 一键构建并启动 NetIP + Caddy 反向代理
docker compose up -d --build
```

启动完成后：
* 访问 `https://ip.yourdomain.com` 即可打开 Web 诊断控制台。
* 终端执行 `curl https://4.yourdomain.com` 即可快速获取当前 IPv4 地址。
* 终端执行 `curl https://6.yourdomain.com` 即可快速获取当前 IPv6 地址。

---

## 🌐 域名与 DNS 解析配置规范

为了让双栈网络检测与纯文本脚本接口正常工作，请在您的域名 DNS 解析提供商处完成以下配置（以 `ipw.3x.cx` 为例）：

| 主机记录 | 类型 | 记录值 | 说明 |
| :--- | :--- | :--- | :--- |
| `ip.ipw.3x.cx` | **A** | 服务器公网 IPv4 | 主站 Web UI 与 REST API (双栈) |
| `ip.ipw.3x.cx` | **AAAA** | 服务器公网 IPv6 | 主站 Web UI 与 REST API (双栈) |
| `4.ipw.3x.cx` | **A** | 服务器公网 IPv4 | **仅配置 A 记录**（不得配置 AAAA），供 IPv4 探测与脚本使用 |
| `6.ipw.3x.cx` | **AAAA** | 服务器公网 IPv6 | **仅配置 AAAA 记录**（不得配置 A），供 IPv6 探测与脚本使用 |
| `test.ipw.3x.cx` | **A** | 服务器公网 IPv4 | 双栈优先级检测（同时配置 A 与 AAAA） |
| `test.ipw.3x.cx` | **AAAA** | 服务器公网 IPv6 | 双栈优先级检测（同时配置 A 与 AAAA） |

> **⚠️ 注意**：`4.` 子域名**切勿添加 AAAA 记录**，`6.` 子域名**切勿添加 A 记录**。请直接将域名 A/AAAA 指向服务器 IP，不要启用 CDN 代理（如 Cloudflare 橙色云朵），避免客户端真实握手协议与真实 IP 被隐藏。

---

## ✨ 核心特性

* 🚀 **即开即测**：首页并发秒级检测当前网络公网 IPv4、IPv6 及协议访问优先级（IPv6 优先 / IPv4 优先）。
* 📜 **极简 IP API**：专为 Shell、DDNS、路由固件打造的纯文本 IP 输出（`curl https://4.yourdomain.com` / `curl https://6.yourdomain.com` / `/ip` / `/json`）。
* 🧭 **多 DNS 解析器查询**：并发查询系统 DNS、Cloudflare (1.1.1.1)、Google (8.8.8.8)、AliDNS (223.5.5.5)、DNSPod (119.29.29.29)。支持 A、AAAA、CNAME、MX、TXT、NS、CAA、SRV、PTR、SOA。
* 📡 **双模 Ping 诊断**：支持原生底层 ICMP Echo 丢包率测试与基于传输层端口握手的 TCPing。
* 🌐 **网站 IPv6 双栈全面体检**：自动验证 DNS A/AAAA、HTTP/HTTPS 双栈连通性并输出就绪结论。
* 🔒 **SSL / TLS 证书深度查询**：解析 X.509 证书链、剩余天数、SANs 多域名、TLS 协议版本与密码套件。
* 🔍 **WHOIS / RDAP 域名与 IP 查询**：现代化 RDAP 优先，权威 WHOIS 兜底，结构化展示注册信息与生命周期。
* ⚡ **网站 HTTP 测速与瀑布分析**：精确拆解 DNS、TCP、TLS、TTFB 耗时，并在 5MB 安全限额内测算下载速率 (Mbps)。
* 🛡️ **军工级 SSRF & DNS Rebinding 防护**：严格阻断所有内网保留地址、Link-Local、CGNAT、私有网段与本地主机名；采用预先解析与安全绑定套接字，根治 DNS 重绑定风险。
* 📦 **单二进制极简部署**：Go 1.25+ 原生全栈，Vue 3 前端静态资源通过 `embed.FS` 内嵌于单一二进制文件中；无需 MySQL、PostgreSQL 或 Redis。
* 🛰️ **预留分布式探测架构 (Probe Mode)**：同一可执行文件支持一键切换为探测节点，通过 HMAC-SHA256 安全鉴权与主站协同工作。

---

## 🛠️ 技术栈

* **后端**：Go 1.25+ / Gin / `miekg/dns` / `golang.org/x/net/icmp` / `log/slog` / Prometheus Metrics
* **前端**：Vue 3 / TypeScript / Vite / Tailwind CSS / Pinia / Lucide Icons
* **代理与容器**：Caddy 2 / Docker Multi-Stage / Docker Compose

---

## 💻 命令行与脚本调用指南

### 1. 终端一行命令获取本机 IP

```bash
# 获取公网 IPv4
curl https://4.ipw.3x.cx

# 获取公网 IPv6
curl https://6.ipw.3x.cx

# 测试当前网络首选协议 IP
curl https://test.ipw.3x.cx

# 获取 JSON 格式响应
curl https://4.ipw.3x.cx/json
```

### 2. RESTful API 示例

```bash
# 查询当前客户端 IP 与归属地
curl https://ip.ipw.3x.cx/api/v1/me

# 查询指定 IP
curl https://ip.ipw.3x.cx/api/v1/ip/8.8.8.8

# 查询 ASN 信息
curl https://ip.ipw.3x.cx/api/v1/asn/AS4134

# DNS 解析查询
curl -X POST https://ip.ipw.3x.cx/api/v1/dns \
  -H "Content-Type: application/json" \
  -d '{"name":"ipw.3x.cx","type":"A"}'

# ICMP Ping 测试
curl -X POST https://ip.ipw.3x.cx/api/v1/ping \
  -H "Content-Type: application/json" \
  -d '{"target":"1.1.1.1","count":4}'

# TCPing 端口握手测试
curl -X POST https://ip.ipw.3x.cx/api/v1/tcping \
  -H "Content-Type: application/json" \
  -d '{"target":"ipw.3x.cx","port":443,"count":4}'

# 网站 IPv6 双栈体检
curl -X POST https://ip.ipw.3x.cx/api/v1/ipv6-check \
  -H "Content-Type: application/json" \
  -d '{"target":"ipw.3x.cx"}'

# SSL 证书查询
curl -X POST https://ip.ipw.3x.cx/api/v1/ssl \
  -H "Content-Type: application/json" \
  -d '{"hostname":"ipw.3x.cx","port":443}'
```

---

## 🗄️ 离线 IP 归属地与 ASN 数据库 (可选)

NetIP 在没有额外数据库文件时即可正常启动并使用内置/在线数据源。如需高性能本地离线定位，可将以下格式文件放入 `./data/ipdb/` 目录：

1. `ip2region.xdb`（用于精确查询中国国内省市及运营商）
2. `GeoLite2-City.mmdb`（MaxMind 全球地理位置数据库）
3. `GeoLite2-ASN.mmdb`（MaxMind 全球 ASN 数据库）

挂载至容器后 NetIP 将自动热加载并提升解析精度。

---

## 🛰️ 分布式探测节点部署 (Probe Node)

如需在异地（例如北京、上海、广州或海外 VPS）部署探测节点：

在探测节点服务器上运行相同二进制：

```bash
# 环境变量配置
MODE=probe
HTTP_ADDR=:8080
PROBE_ID=cn-beijing-01
PROBE_NAME=北京
PROBE_ISP=中国电信
PROBE_SECRET=your-shared-secret

# 启动节点
./netip probe
```

在主站 `config/probes.yaml` 中配置节点列表：

```yaml
probes:
  - id: cn-beijing-01
    name: 北京
    isp: 中国电信
    url: https://bj-probe.yourdomain.com
    secret: your-shared-secret
```

主站与节点间通信通过 **HMAC-SHA256** 与时间戳窗口（±30 秒）进行双向签名鉴权，杜绝非法调用与未授权网络扫描。

---

## 🛡️ 安全基准

* **SSRF 防御**：彻底封禁 `127.0.0.0/8`、`10.0.0.0/8`、`172.16.0.0/12`、`192.168.0.0/16`、`169.254.0.0/16`、`::1`、`fc00::/7`、`fe80::/10` 及元数据地址 `169.254.169.254`。
* **重定向防御**：HTTP 请求跳转最多追溯 3 跳，每次重定向均强制重新执行完整的 SSRF 目标校验。
* **端口白名单**：TCPing 与 SSL 检测受 `ALLOWED_TCP_PORTS` 环境变量严格管控。
* **限流与并发控制**：基于内存令牌桶与全局信号量控制，有效抵御滥用与资源耗尽。

---

## 🧪 本地开发与测试

```bash
# 运行后端单元测试
go test -v ./...

# 运行前端测试与构建
cd web
npm install
npm run typecheck
npm run test
npm run build

# 本地运行服务
make dev
```

---

## 📄 开源许可证

本项目基于 [MIT License](LICENSE) 开源。
