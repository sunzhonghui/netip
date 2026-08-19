import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
      meta: { title: 'NetIP - 我的网络与 IP 状态' },
    },
    {
      path: '/ip',
      name: 'ip',
      component: () => import('@/views/IpView.vue'),
      meta: { title: 'IP 地址与归属地查询 - NetIP' },
    },
    {
      path: '/ip/:address',
      name: 'ip-detail',
      component: () => import('@/views/IpView.vue'),
      meta: { title: 'IP 地址查询 - NetIP' },
    },
    {
      path: '/dns',
      name: 'dns',
      component: () => import('@/views/DnsView.vue'),
      meta: { title: '多 DNS 解析器查询 - NetIP' },
    },
    {
      path: '/ping',
      name: 'ping',
      component: () => import('@/views/PingView.vue'),
      meta: { title: 'Ping / TCPing 网络连通性测试 - NetIP' },
    },
    {
      path: '/ipv6',
      name: 'ipv6',
      component: () => import('@/views/Ipv6View.vue'),
      meta: { title: '网站 IPv6 支持检测 - NetIP' },
    },
    {
      path: '/ssl',
      name: 'ssl',
      component: () => import('@/views/SslView.vue'),
      meta: { title: 'SSL / TLS 证书查询 - NetIP' },
    },
    {
      path: '/whois',
      name: 'whois',
      component: () => import('@/views/WhoisView.vue'),
      meta: { title: 'WHOIS / RDAP 注册信息查询 - NetIP' },
    },
    {
      path: '/asn',
      name: 'asn',
      component: () => import('@/views/AsnView.vue'),
      meta: { title: 'ASN 自治系统查询 - NetIP' },
    },
    {
      path: '/asn/:query',
      name: 'asn-detail',
      component: () => import('@/views/AsnView.vue'),
      meta: { title: 'ASN 查询 - NetIP' },
    },
    {
      path: '/http',
      name: 'http',
      component: () => import('@/views/HttpView.vue'),
      meta: { title: '网站 HTTP 响应与耗时分析 - NetIP' },
    },
    {
      path: '/speed',
      name: 'speed',
      component: () => import('@/views/SpeedView.vue'),
      meta: { title: '网站 HTTP 测速 - NetIP' },
    },
    {
      path: '/docs/api',
      name: 'api-docs',
      component: () => import('@/views/ApiDocsView.vue'),
      meta: { title: 'API 开发接口文档 - NetIP' },
    },
  ],
})

router.afterEach((to) => {
  if (to.meta.title) {
    document.title = to.meta.title as string
  }
})

export default router
