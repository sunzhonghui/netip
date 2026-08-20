import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
      meta: {
        title: 'NetIP - 我的网络与公网 IPv4 / IPv6 状态',
        description: '实时检测当前网络公网 IPv4、IPv6 地址、访问优先级、网络运营商与 ASN 信息，支持一键复制。',
      },
    },
    {
      path: '/ip',
      name: 'ip',
      component: () => import('@/views/IpView.vue'),
      meta: {
        title: 'IP 地址与归属地查询 - NetIP',
        description: '精准查询 IPv4 / IPv6 公网地址的地理位置、国家省市、网络运营商、所属 ASN 自治系统及广播网段。',
      },
    },
    {
      path: '/ip/:address',
      name: 'ip-detail',
      component: () => import('@/views/IpView.vue'),
      meta: {
        title: 'IP 地址详情查询 - NetIP',
        description: '查询指定 IP 地址的地理归属地、运营商、ASN 自治系统和网络信息。',
      },
    },
    {
      path: '/dns',
      name: 'dns',
      component: () => import('@/views/DnsView.vue'),
      meta: {
        title: '多 DNS 解析器对比查询 - NetIP',
        description: '并发向系统 DNS、1.1.1.1、8.8.8.8、223.5.5.5 等公共 DNS 查询 A/AAAA/CNAME/MX/TXT 记录并对比解析延迟。',
      },
    },
    {
      path: '/ping',
      name: 'ping',
      component: () => import('@/views/PingView.vue'),
      meta: {
        title: 'Ping / TCPing 网络连通性测试 - NetIP',
        description: '支持 ICMP Echo 丢包率测试与 TCP 端口三次握手连接耗时测量，诊断主机与端口连通性。',
      },
    },
    {
      path: '/ipv6',
      name: 'ipv6',
      component: () => import('@/views/Ipv6View.vue'),
      meta: {
        title: '网站 IPv6 双栈支持体检 - NetIP',
        description: '全面检测目标网站的 DNS A/AAAA 记录解析、HTTP 与 HTTPS IPv6 双栈连通性并输出就绪结论。',
      },
    },
    {
      path: '/ssl',
      name: 'ssl',
      component: () => import('@/views/SslView.vue'),
      meta: {
        title: 'SSL / TLS 证书查询 - NetIP',
        description: '深度解析 HTTPS 目标服务器 X.509 证书链、颁发机构、剩余有效期天数、SANs 多域名及加密套件。',
      },
    },
    {
      path: '/whois',
      name: 'whois',
      component: () => import('@/views/WhoisView.vue'),
      meta: {
        title: 'WHOIS / RDAP 注册信息查询 - NetIP',
        description: '基于现代 RDAP 协议与权威 WHOIS 服务器，查询域名或 IP 地址的注册商、注册时间、过期时间及状态。',
      },
    },
    {
      path: '/asn',
      name: 'asn',
      component: () => import('@/views/AsnView.vue'),
      meta: {
        title: 'ASN 自治系统查询 - NetIP',
        description: '查询自治系统编号 (ASN)、BGP 广播运营商名称、注册局与所属国家。',
      },
    },
    {
      path: '/asn/:query',
      name: 'asn-detail',
      component: () => import('@/views/AsnView.vue'),
      meta: {
        title: 'ASN 详情查询 - NetIP',
        description: '查询指定自治系统 (ASN) 或 IP 对应的运营商与自治域信息。',
      },
    },
    {
      path: '/http',
      name: 'http',
      component: () => import('@/views/HttpView.vue'),
      meta: {
        title: '网站 HTTP 响应与瀑布耗时分析 - NetIP',
        description: '精确拆解 HTTP 请求的 DNS 解析、TCP 握手、TLS 协商、首字节时间 (TTFB) 及完整响应流程。',
      },
    },
    {
      path: '/speed',
      name: 'speed',
      component: () => import('@/views/SpeedView.vue'),
      meta: {
        title: '网站 HTTP 测速 - NetIP',
        description: '连接目标网站并下载最高 5MB 静态数据包，计算实际下载吞吐速率 (Mbps) 与分段延迟。',
      },
    },
    {
      path: '/docs/api',
      name: 'api-docs',
      component: () => import('@/views/ApiDocsView.vue'),
      meta: {
        title: 'API 开发接口文档 - NetIP',
        description: '面向 Shell / DDNS 脚本的轻量纯文本接口 (curl https://4.ipw.3x.cx) 与标准化 RESTful JSON API 文档。',
      },
    },
  ],
})

router.afterEach((to) => {
  // Update Title
  if (to.meta.title) {
    document.title = to.meta.title as string
  }

  // Update Meta Description
  if (to.meta.description) {
    let descMeta = document.querySelector('meta[name="description"]')
    if (!descMeta) {
      descMeta = document.createElement('meta')
      descMeta.setAttribute('name', 'description')
      document.head.appendChild(descMeta)
    }
    descMeta.setAttribute('content', to.meta.description as string)
  }

  // Update Canonical URL
  let canonicalLink = document.querySelector('link[rel="canonical"]')
  if (!canonicalLink) {
    canonicalLink = document.createElement('link')
    canonicalLink.setAttribute('rel', 'canonical')
    document.head.appendChild(canonicalLink)
  }
  canonicalLink.setAttribute('href', `https://ip.ipw.3x.cx${to.path}`)
})

export default router
