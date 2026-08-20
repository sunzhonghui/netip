import{_ as x}from"./ToolHeader.vue_vue_type_script_setup_true_lang-YvSb8OyN.js";import{_ as p}from"./CopyButton.vue_vue_type_script_setup_true_lang-C_PcGsw8.js";import{c as _,d as h,a as i,e as r,u as l,C as b,b as e,t as o,F as v,k as f,o as d,n as T}from"./index-Cve9KufU.js";/**
 * @license lucide-vue-next v0.475.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const g=_("TerminalIcon",[["polyline",{points:"4 17 10 11 4 5",key:"akl6gq"}],["line",{x1:"12",x2:"20",y1:"19",y2:"19",key:"q2wloq"}]]),P={class:"space-y-6"},y={class:"rounded-xl p-6 bg-slate-950 text-slate-100 font-mono space-y-4 border border-slate-800/90 shadow-md"},S={class:"flex items-center justify-between text-xs text-slate-400 border-b border-slate-800 pb-3 font-sans"},w={class:"flex items-center gap-2 font-semibold"},A={class:"space-y-3 text-xs sm:text-sm"},k={class:"flex items-center justify-between gap-4 p-3 rounded-lg bg-black border border-slate-800"},$={class:"text-emerald-400 font-mono font-medium"},C={class:"flex items-center justify-between gap-4 p-3 rounded-lg bg-black border border-slate-800"},O={class:"text-purple-400 font-mono font-medium"},j={class:"flex items-center justify-between gap-4 p-3 rounded-lg bg-black border border-slate-800"},I={class:"text-blue-400 font-mono font-medium"},N={class:"space-y-6"},E={class:"flex flex-wrap items-center justify-between gap-3 border-b border-slate-100 dark:border-slate-800 pb-3"},D={class:"flex items-center gap-3"},H={class:"font-mono text-sm font-bold text-slate-800 dark:text-slate-200"},L={class:"text-xs text-slate-500 font-medium"},G={class:"text-xs sm:text-sm text-slate-600 dark:text-slate-400 leading-relaxed"},X={class:"grid grid-cols-1 lg:grid-cols-2 gap-4 text-xs font-mono"},q={class:"space-y-1.5"},B={class:"flex items-center justify-between text-slate-400 font-sans"},U={class:"p-3.5 rounded-lg bg-black border border-slate-800 text-emerald-400 overflow-x-auto whitespace-pre-wrap"},F={class:"space-y-1.5"},M={class:"p-3.5 rounded-lg bg-black border border-slate-800 text-slate-200 overflow-x-auto max-h-40"},K=h({__name:"ApiDocsView",setup(R){const s=window.location.hostname,c=s.split(".");let t=s;c.length>2&&(t=c.slice(1).join("."));const u=[{title:"纯文本 IPv4 地址获取 (脚本 / DDNS / cURL)",method:"GET",url:`https://4.${t}/`,desc:"仅配置 A 记录，返回调用者的公网 IPv4 地址与换行符。",example:`curl https://4.${t}`,response:"123.123.123.123"},{title:"纯文本 IPv6 地址获取",method:"GET",url:`https://6.${t}/`,desc:"仅配置 AAAA 记录，返回调用者的公网 IPv6 地址与换行符。",example:`curl https://6.${t}`,response:"240e:390:1234::1"},{title:"网络双栈优先级连通 IP",method:"GET",url:`https://test.${t}/`,desc:"同时配置 A 与 AAAA 记录，根据客户端首选连接协议返回地址。",example:`curl https://test.${t}`,response:"240e:390:1234::1"},{title:"JSON 格式 IP 地址",method:"GET",url:`https://4.${t}/json`,desc:"返回当前客户端 IP 及 IP 协议版本号。",example:`curl https://4.${t}/json`,response:`{
  "ip": "123.123.123.123",
  "version": 4
}`},{title:"当前公网 IP 及归属地详细信息",method:"GET",url:"/api/v1/me",desc:"返回当前客户端的公网 IP、地理位置、运营商及 ASN。",example:`curl https://${s}/api/v1/me`,response:`{
  "success": true,
  "data": {
    "ip": "1.1.1.1",
    "version": 4,
    "country": "Australia",
    "country_code": "AU",
    "isp": "Cloudflare",
    "asn": 13335,
    "as_name": "CLOUDFLARENET"
  },
  "request_id": "req-123"
}`},{title:"指定 IP 地址归属地查询",method:"GET",url:"/api/v1/ip/{ip}",desc:"查询指定 IPv4 / IPv6 地址的地理位置与网络信息。",example:`curl https://${s}/api/v1/ip/8.8.8.8`,response:`{
  "success": true,
  "data": {
    "ip": "8.8.8.8",
    "version": 4,
    "country": "United States",
    "country_code": "US",
    "isp": "Google",
    "asn": 15169,
    "as_name": "GOOGLE"
  },
  "request_id": "req-456"
}`},{title:"多 DNS 解析器查询",method:"POST",url:"/api/v1/dns",desc:"向系统 DNS、1.1.1.1、8.8.8.8、223.5.5.5 等并发查询记录。",example:`curl -X POST https://${s}/api/v1/dns \\
  -H "Content-Type: application/json" \\
  -d '{"name":"ipw.3x.cx","type":"A"}'`,response:`{
  "success": true,
  "data": {
    "name": "ipw.3x.cx",
    "type": "A",
    "results": [
      {
        "resolver": "1.1.1.1",
        "latency_ms": 18,
        "answers": [{"type":"A","value":"1.2.3.4","ttl":300}]
      }
    ]
  }
}`},{title:"ICMP Ping 连通性测试",method:"POST",url:"/api/v1/ping",desc:"向目标主机发送 ICMP Echo 报文并统计往返时延与丢包率。",example:`curl -X POST https://${s}/api/v1/ping \\
  -H "Content-Type: application/json" \\
  -d '{"target":"1.1.1.1","count":4}'`,response:`{
  "success": true,
  "data": {
    "target": "1.1.1.1",
    "resolved_ip": "1.1.1.1",
    "sent": 4,
    "received": 4,
    "loss_percent": 0,
    "avg_ms": 14.2
  }
}`},{title:"TCPing 端口握手延迟",method:"POST",url:"/api/v1/tcping",desc:"针对指定 TCP 端口测量三次握手连接建立耗时。",example:`curl -X POST https://${s}/api/v1/tcping \\
  -H "Content-Type: application/json" \\
  -d '{"target":"ipw.3x.cx","port":443,"count":4}'`,response:`{
  "success": true,
  "data": {
    "target": "ipw.3x.cx",
    "port": 443,
    "success": 4,
    "avg_ms": 28.5
  }
}`},{title:"网站 IPv6 双栈支持检测",method:"POST",url:"/api/v1/ipv6-check",desc:"检测目标网站的 DNS A/AAAA、HTTP 与 HTTPS IPv6 可用性。",example:`curl -X POST https://${s}/api/v1/ipv6-check \\
  -H "Content-Type: application/json" \\
  -d '{"target":"ipw.3x.cx"}'`,response:`{
  "success": true,
  "data": {
    "domain": "ipw.3x.cx",
    "supported": true,
    "conclusion": "该网站完整支持 IPv6"
  }
}`},{title:"SSL / TLS 证书查询",method:"POST",url:"/api/v1/ssl",desc:"连接目标 HTTPS 端口并提取 X.509 证书链、SANs、加密套件。",example:`curl -X POST https://${s}/api/v1/ssl \\
  -H "Content-Type: application/json" \\
  -d '{"hostname":"ipw.3x.cx","port":443}'`,response:`{
  "success": true,
  "data": {
    "hostname": "ipw.3x.cx",
    "valid": true,
    "days_remaining": 82,
    "tls_version": "TLS 1.3"
  }
}`},{title:"网站 HTTP 响应与耗时分析",method:"POST",url:"/api/v1/http",desc:"记录 DNS、TCP、TLS、TTFB 耗时与响应头信息。",example:`curl -X POST https://${s}/api/v1/http \\
  -H "Content-Type: application/json" \\
  -d '{"target":"https://ipw.3x.cx"}'`,response:`{
  "success": true,
  "data": {
    "dns_ms": 12,
    "tcp_ms": 24,
    "tls_ms": 38,
    "ttfb_ms": 65,
    "status_code": 200
  }
}`},{title:"网站 HTTP 测速",method:"POST",url:"/api/v1/speed",desc:"下载最高 5MB 流量计算下载速率 (Mbps)。",example:`curl -X POST https://${s}/api/v1/speed \\
  -H "Content-Type: application/json" \\
  -d '{"target":"https://ipw.3x.cx"}'`,response:`{
  "success": true,
  "data": {
    "speed_mbps": 48.2,
    "download_bytes": 1048576,
    "download_ms": 174
  }
}`}];return(V,n)=>(d(),i("div",P,[r(x,{title:"API 开发接口文档",description:"NetIP 提供面向 Shell / DDNS 脚本的轻量纯文本接口，以及适合程序调用的标准化 RESTful JSON API",icon:l(b)},null,8,["icon"]),e("div",y,[e("div",S,[e("div",w,[r(l(g),{class:"w-4 h-4 text-brand-400"}),n[0]||(n[0]=e("span",{class:"text-slate-200 font-bold"},"终端快捷使用示例 (CLI Quick Start)",-1))])]),e("div",A,[e("div",k,[e("span",$,"curl https://4."+o(l(t)),1),r(p,{text:`curl https://4.${l(t)}`,label:"复制"},null,8,["text"])]),e("div",C,[e("span",O,"curl https://6."+o(l(t)),1),r(p,{text:`curl https://6.${l(t)}`,label:"复制"},null,8,["text"])]),e("div",j,[e("span",I,"curl https://test."+o(l(t)),1),r(p,{text:`curl https://test.${l(t)}`,label:"复制"},null,8,["text"])])])]),e("div",N,[(d(),i(v,null,f(u,(a,m)=>e("div",{key:m,class:"custom-card p-6 space-y-4"},[e("div",E,[e("div",D,[e("span",{class:T(["px-2 py-0.5 rounded text-xs font-mono font-bold uppercase",a.method==="GET"?"bg-blue-100 text-blue-800 dark:bg-blue-900/60 dark:text-blue-200":"bg-emerald-100 text-emerald-800 dark:bg-emerald-900/60 dark:text-emerald-200"])},o(a.method),3),e("span",H,o(a.url),1)]),e("span",L,o(a.title),1)]),e("p",G,o(a.desc),1),e("div",X,[e("div",q,[e("div",B,[n[1]||(n[1]=e("span",null,"调用示例:",-1)),r(p,{text:a.example,"icon-only":!0},null,8,["text"])]),e("pre",U,o(a.example),1)]),e("div",F,[n[2]||(n[2]=e("div",{class:"flex items-center justify-between text-slate-400 font-sans"},[e("span",null,"响应示例:")],-1)),e("pre",M,o(a.response),1)])])])),64))])]))}});export{K as default};
