<script setup lang="ts">
import { Code2, Terminal } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import CopyButton from '@/components/CopyButton.vue'

const host = window.location.hostname
const parts = host.split('.')
let baseDomain = host
if (parts.length > 2) {
  baseDomain = parts.slice(1).join('.')
}

const apis = [
  {
    title: '纯文本 IPv4 地址获取 (脚本 / DDNS / cURL)',
    method: 'GET',
    url: `https://4.${baseDomain}/`,
    desc: '仅配置 A 记录，返回调用者的公网 IPv4 地址与换行符。',
    example: `curl https://4.${baseDomain}`,
    response: `123.123.123.123`,
  },
  {
    title: '纯文本 IPv6 地址获取',
    method: 'GET',
    url: `https://6.${baseDomain}/`,
    desc: '仅配置 AAAA 记录，返回调用者的公网 IPv6 地址与换行符。',
    example: `curl https://6.${baseDomain}`,
    response: `240e:390:1234::1`,
  },
  {
    title: '网络双栈优先级连通 IP',
    method: 'GET',
    url: `https://test.${baseDomain}/`,
    desc: '同时配置 A 与 AAAA 记录，根据客户端首选连接协议返回地址。',
    example: `curl https://test.${baseDomain}`,
    response: `240e:390:1234::1`,
  },
  {
    title: 'JSON 格式 IP 地址',
    method: 'GET',
    url: `https://4.${baseDomain}/json`,
    desc: '返回当前客户端 IP 及 IP 协议版本号。',
    example: `curl https://4.${baseDomain}/json`,
    response: `{\n  "ip": "123.123.123.123",\n  "version": 4\n}`,
  },
  {
    title: '当前公网 IP 及归属地详细信息',
    method: 'GET',
    url: `/api/v1/me`,
    desc: '返回当前客户端的公网 IP、地理位置、运营商及 ASN。',
    example: `curl https://${host}/api/v1/me`,
    response: `{\n  "success": true,\n  "data": {\n    "ip": "1.1.1.1",\n    "version": 4,\n    "country": "Australia",\n    "country_code": "AU",\n    "isp": "Cloudflare",\n    "asn": 13335,\n    "as_name": "CLOUDFLARENET"\n  },\n  "request_id": "req-123"\n}`,
  },
  {
    title: '指定 IP 地址归属地查询',
    method: 'GET',
    url: `/api/v1/ip/{ip}`,
    desc: '查询指定 IPv4 / IPv6 地址的地理位置与网络信息。',
    example: `curl https://${host}/api/v1/ip/8.8.8.8`,
    response: `{\n  "success": true,\n  "data": {\n    "ip": "8.8.8.8",\n    "version": 4,\n    "country": "United States",\n    "country_code": "US",\n    "isp": "Google",\n    "asn": 15169,\n    "as_name": "GOOGLE"\n  },\n  "request_id": "req-456"\n}`,
  },
  {
    title: '多 DNS 解析器查询',
    method: 'POST',
    url: `/api/v1/dns`,
    desc: '向系统 DNS、1.1.1.1、8.8.8.8、223.5.5.5 等并发查询记录。',
    example: `curl -X POST https://${host}/api/v1/dns \\\n  -H "Content-Type: application/json" \\\n  -d '{"name":"ipw.3x.cx","type":"A"}'`,
    response: `{\n  "success": true,\n  "data": {\n    "name": "ipw.3x.cx",\n    "type": "A",\n    "results": [\n      {\n        "resolver": "1.1.1.1",\n        "latency_ms": 18,\n        "answers": [{"type":"A","value":"1.2.3.4","ttl":300}]\n      }\n    ]\n  }\n}`,
  },
  {
    title: 'ICMP Ping 连通性测试',
    method: 'POST',
    url: `/api/v1/ping`,
    desc: '向目标主机发送 ICMP Echo 报文并统计往返时延与丢包率。',
    example: `curl -X POST https://${host}/api/v1/ping \\\n  -H "Content-Type: application/json" \\\n  -d '{"target":"1.1.1.1","count":4}'`,
    response: `{\n  "success": true,\n  "data": {\n    "target": "1.1.1.1",\n    "resolved_ip": "1.1.1.1",\n    "sent": 4,\n    "received": 4,\n    "loss_percent": 0,\n    "avg_ms": 14.2\n  }\n}`,
  },
  {
    title: 'TCPing 端口握手延迟',
    method: 'POST',
    url: `/api/v1/tcping`,
    desc: '针对指定 TCP 端口测量三次握手连接建立耗时。',
    example: `curl -X POST https://${host}/api/v1/tcping \\\n  -H "Content-Type: application/json" \\\n  -d '{"target":"ipw.3x.cx","port":443,"count":4}'`,
    response: `{\n  "success": true,\n  "data": {\n    "target": "ipw.3x.cx",\n    "port": 443,\n    "success": 4,\n    "avg_ms": 28.5\n  }\n}`,
  },
  {
    title: '网站 IPv6 双栈支持检测',
    method: 'POST',
    url: `/api/v1/ipv6-check`,
    desc: '检测目标网站的 DNS A/AAAA、HTTP 与 HTTPS IPv6 可用性。',
    example: `curl -X POST https://${host}/api/v1/ipv6-check \\\n  -H "Content-Type: application/json" \\\n  -d '{"target":"ipw.3x.cx"}'`,
    response: `{\n  "success": true,\n  "data": {\n    "domain": "ipw.3x.cx",\n    "supported": true,\n    "conclusion": "该网站完整支持 IPv6"\n  }\n}`,
  },
  {
    title: 'SSL / TLS 证书查询',
    method: 'POST',
    url: `/api/v1/ssl`,
    desc: '连接目标 HTTPS 端口并提取 X.509 证书链、SANs、加密套件。',
    example: `curl -X POST https://${host}/api/v1/ssl \\\n  -H "Content-Type: application/json" \\\n  -d '{"hostname":"ipw.3x.cx","port":443}'`,
    response: `{\n  "success": true,\n  "data": {\n    "hostname": "ipw.3x.cx",\n    "valid": true,\n    "days_remaining": 82,\n    "tls_version": "TLS 1.3"\n  }\n}`,
  },
  {
    title: '网站 HTTP 响应与耗时分析',
    method: 'POST',
    url: `/api/v1/http`,
    desc: '记录 DNS、TCP、TLS、TTFB 耗时与响应头信息。',
    example: `curl -X POST https://${host}/api/v1/http \\\n  -H "Content-Type: application/json" \\\n  -d '{"target":"https://ipw.3x.cx"}'`,
    response: `{\n  "success": true,\n  "data": {\n    "dns_ms": 12,\n    "tcp_ms": 24,\n    "tls_ms": 38,\n    "ttfb_ms": 65,\n    "status_code": 200\n  }\n}`,
  },
  {
    title: '网站 HTTP 测速',
    method: 'POST',
    url: `/api/v1/speed`,
    desc: '下载最高 5MB 流量计算下载速率 (Mbps)。',
    example: `curl -X POST https://${host}/api/v1/speed \\\n  -H "Content-Type: application/json" \\\n  -d '{"target":"https://ipw.3x.cx"}'`,
    response: `{\n  "success": true,\n  "data": {\n    "speed_mbps": 48.2,\n    "download_bytes": 1048576,\n    "download_ms": 174\n  }\n}`,
  },
]
</script>

<template>
  <div class="space-y-6">
    <ToolHeader
      title="API 开发接口文档"
      description="NetIP 提供面向 Shell / DDNS 脚本的轻量纯文本接口，以及适合程序调用的标准化 RESTful JSON API"
      :icon="Code2"
    />

    <!-- Quick cURL Showcase -->
    <div class="custom-card p-6 bg-slate-900 text-slate-100 font-mono space-y-4">
      <div class="flex items-center justify-between text-xs text-slate-400 border-b border-slate-800 pb-3 font-sans">
        <div class="flex items-center gap-2 font-semibold">
          <Terminal class="w-4 h-4 text-brand-400" />
          <span>终端快捷使用示例 (CLI Quick Start)</span>
        </div>
      </div>

      <div class="space-y-3 text-xs sm:text-sm">
        <div class="flex items-center justify-between gap-4 p-2.5 rounded bg-slate-950/60">
          <span class="text-emerald-400">curl https://4.{{ baseDomain }}</span>
          <CopyButton :text="`curl https://4.${baseDomain}`" label="复制" />
        </div>
        <div class="flex items-center justify-between gap-4 p-2.5 rounded bg-slate-950/60">
          <span class="text-purple-400">curl https://6.{{ baseDomain }}</span>
          <CopyButton :text="`curl https://6.${baseDomain}`" label="复制" />
        </div>
        <div class="flex items-center justify-between gap-4 p-2.5 rounded bg-slate-950/60">
          <span class="text-blue-400">curl https://test.{{ baseDomain }}</span>
          <CopyButton :text="`curl https://test.${baseDomain}`" label="复制" />
        </div>
      </div>
    </div>

    <!-- API List -->
    <div class="space-y-6">
      <div
        v-for="(apiItem, idx) in apis"
        :key="idx"
        class="custom-card p-6 space-y-4"
      >
        <div class="flex flex-wrap items-center justify-between gap-3 border-b border-slate-100 dark:border-slate-800 pb-3">
          <div class="flex items-center gap-3">
            <span
              class="px-2 py-0.5 rounded text-xs font-mono font-bold uppercase"
              :class="apiItem.method === 'GET' ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/60 dark:text-blue-200' : 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/60 dark:text-emerald-200'"
            >
              {{ apiItem.method }}
            </span>
            <span class="font-mono text-sm font-bold text-slate-800 dark:text-slate-200">{{ apiItem.url }}</span>
          </div>
          <span class="text-xs text-slate-500 font-medium">{{ apiItem.title }}</span>
        </div>

        <p class="text-xs sm:text-sm text-slate-600 dark:text-slate-400 leading-relaxed">
          {{ apiItem.desc }}
        </p>

        <!-- Example & Response tabs -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 text-xs font-mono">
          <div class="space-y-1.5">
            <div class="flex items-center justify-between text-slate-400 font-sans">
              <span>调用示例:</span>
              <CopyButton :text="apiItem.example" :icon-only="true" />
            </div>
            <pre class="p-3 rounded-lg bg-slate-900 text-emerald-400 overflow-x-auto whitespace-pre-wrap">{{ apiItem.example }}</pre>
          </div>

          <div class="space-y-1.5">
            <div class="flex items-center justify-between text-slate-400 font-sans">
              <span>响应示例:</span>
            </div>
            <pre class="p-3 rounded-lg bg-slate-900 text-slate-200 overflow-x-auto max-h-40">{{ apiItem.response }}</pre>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
