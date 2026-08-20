<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  Globe2,
  Network,
  Activity,
  CheckCircle2,
  Lock,
  Search,
  Gauge,
  Layers,
  ArrowUpDown,
  RefreshCw,
} from 'lucide-vue-next'
import IPCard from '@/components/IPCard.vue'
import { api, fetchWithTimeout } from '@/utils/api'
import type { IPDetails } from '@/types'

const loading = ref(true)
const ipv4 = ref<string | undefined>()
const ipv4Details = ref<IPDetails | undefined>()
const ipv6 = ref<string | undefined>()
const ipv6Details = ref<IPDetails | undefined>()
const preferred = ref<'ipv4' | 'ipv6' | 'unknown'>('unknown')

const quickTools = [
  {
    name: 'IP 查询',
    desc: '查询指定 IP 的地理归属地、运营商与 ASN',
    path: '/ip',
    icon: Globe2,
    color: 'text-blue-600 dark:text-blue-400',
    bg: 'bg-blue-50 dark:bg-blue-950/50',
  },
  {
    name: 'DNS 查询',
    desc: '多 DNS 解析器并发查询 A、AAAA、CNAME、MX 等',
    path: '/dns',
    icon: Network,
    color: 'text-indigo-600 dark:text-indigo-400',
    bg: 'bg-indigo-50 dark:bg-indigo-950/50',
  },
  {
    name: 'Ping / TCPing',
    desc: '测试目标主机的 ICMP 丢包率与 TCP 端口握手延迟',
    path: '/ping',
    icon: Activity,
    color: 'text-emerald-600 dark:text-emerald-400',
    bg: 'bg-emerald-50 dark:bg-emerald-950/50',
  },
  {
    name: 'IPv6 检测',
    desc: '全面检测目标网站的 IPv6 双栈支持与可访问性',
    path: '/ipv6',
    icon: CheckCircle2,
    color: 'text-purple-600 dark:text-purple-400',
    bg: 'bg-purple-50 dark:bg-purple-950/50',
  },
  {
    name: 'SSL 证书',
    desc: '查询 HTTPS 证书有效期、颁发者、SANs 与加密套件',
    path: '/ssl',
    icon: Lock,
    color: 'text-amber-600 dark:text-amber-400',
    bg: 'bg-amber-50 dark:bg-amber-950/50',
  },
  {
    name: 'WHOIS / RDAP',
    desc: '查询域名或 IP 地址的注册人、注册商及生命周期',
    path: '/whois',
    icon: Search,
    color: 'text-cyan-600 dark:text-cyan-400',
    bg: 'bg-cyan-50 dark:bg-cyan-950/50',
  },
  {
    name: '网站测速',
    desc: '测量 DNS、TCP、TLS、TTFB 及真实 HTTP 下载速度',
    path: '/speed',
    icon: Gauge,
    color: 'text-rose-600 dark:text-rose-400',
    bg: 'bg-rose-50 dark:bg-rose-950/50',
  },
  {
    name: 'ASN 查询',
    desc: '查询自治系统编号、路由前缀与网络运营商',
    path: '/asn',
    icon: Layers,
    color: 'text-orange-600 dark:text-orange-400',
    bg: 'bg-orange-50 dark:bg-orange-950/50',
  },
]

async function detectNetwork() {
  loading.value = true
  ipv4.value = undefined
  ipv4Details.value = undefined
  ipv6.value = undefined
  ipv6Details.value = undefined
  preferred.value = 'unknown'

  // Dynamic host determination
  const host = window.location.hostname
  const parts = host.split('.')
  let baseDomain = host
  if (parts.length > 2) {
    baseDomain = parts.slice(1).join('.')
  }
  const protocol = window.location.protocol

  const v4Url = `${protocol}//4.${baseDomain}/json`
  const v6Url = `${protocol}//6.${baseDomain}/json`
  const testUrl = `${protocol}//test.${baseDomain}/json`

  // Concurrent detection with 4000ms timeout
  const [resV4, resV6, resTest, resMe] = await Promise.allSettled([
    fetchWithTimeout<{ ip: string; version: number }>(v4Url, 4000),
    fetchWithTimeout<{ ip: string; version: number }>(v6Url, 4000),
    fetchWithTimeout<{ ip: string; version: number }>(testUrl, 4000),
    api.getMe(),
  ])

  // Handle IPv4 result
  if (resV4.status === 'fulfilled' && resV4.value?.ip) {
    ipv4.value = resV4.value.ip
  }

  // Handle IPv6 result
  if (resV6.status === 'fulfilled' && resV6.value?.ip) {
    ipv6.value = resV6.value.ip
  }

  // Handle Priority result
  if (resTest.status === 'fulfilled' && resTest.value?.version) {
    preferred.value = resTest.value.version === 6 ? 'ipv6' : 'ipv4'
  }

  // Fallback: If 4. / 6. subdomains didn't respond (e.g. dev localhost), use /api/v1/me
  if (!ipv4.value && !ipv6.value && resMe.status === 'fulfilled' && resMe.value?.ip) {
    const me = resMe.value
    if (me.version === 6) {
      ipv6.value = me.ip
      ipv6Details.value = me
      preferred.value = 'ipv6'
    } else {
      ipv4.value = me.ip
      ipv4Details.value = me
      preferred.value = 'ipv4'
    }
  }

  // Enrich details asynchronously
  if (ipv4.value && !ipv4Details.value) {
    api.getIP(ipv4.value).then((d) => {
      ipv4Details.value = d
    }).catch(() => {})
  }

  if (ipv6.value && !ipv6Details.value) {
    api.getIP(ipv6.value).then((d) => {
      ipv6Details.value = d
    }).catch(() => {})
  }

  loading.value = false
}

onMounted(() => {
  detectNetwork()
})
</script>

<template>
  <div class="space-y-8">
    <!-- Header Section -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3.5">
        <img src="/favicon.svg" alt="NetIP" class="w-10 h-10 sm:w-12 sm:h-12 rounded-2xl shadow-md flex-shrink-0" />
        <div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white">
            您的网络
          </h1>
          <p class="mt-0.5 text-xs sm:text-sm text-slate-500 dark:text-slate-400">
            实时公网 IP 地址与 IPv4 / IPv6 双栈支持状态
          </p>
        </div>
      </div>

      <button
        type="button"
        @click="detectNetwork"
        :disabled="loading"
        class="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold rounded-lg bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-200 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700 shadow-sm transition-all"
      >
        <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
        <span>重新检测</span>
      </button>
    </div>

    <!-- Main IP Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <IPCard
        version="IPv4"
        :ip="ipv4"
        :details="ipv4Details"
        :available="!!ipv4"
        :loading="loading"
      />
      <IPCard
        version="IPv6"
        :ip="ipv6"
        :details="ipv6Details"
        :available="!!ipv6"
        :loading="loading"
      />
    </div>

    <!-- Network Priority & Status Card -->
    <div class="custom-card p-6 bg-gradient-to-r from-brand-500/5 via-transparent to-brand-500/5">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-brand-500/10 text-brand-600 dark:text-brand-400 flex items-center justify-center flex-shrink-0">
            <ArrowUpDown class="w-5 h-5" />
          </div>
          <div>
            <div class="text-xs uppercase tracking-wider font-semibold text-slate-400 dark:text-slate-500">
              网络访问优先级
            </div>
            <div class="text-lg font-bold text-slate-900 dark:text-white mt-0.5">
              <span v-if="loading" class="text-slate-400">检测中...</span>
              <span v-else-if="preferred === 'ipv6'" class="text-purple-600 dark:text-purple-400">IPv6 优先</span>
              <span v-else-if="preferred === 'ipv4'" class="text-blue-600 dark:text-blue-400">IPv4 优先</span>
              <span v-else class="text-slate-500">常规模式</span>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-4 text-xs font-mono font-medium">
          <div class="flex items-center gap-1.5 px-3 py-1 rounded-md border" :class="ipv4 ? 'bg-emerald-50 border-emerald-200 text-emerald-700 dark:bg-emerald-950/40 dark:border-emerald-800 dark:text-emerald-300' : 'bg-slate-100 border-slate-200 text-slate-400 dark:bg-slate-800 dark:border-slate-700'">
            <span>IPv4</span>
            <span>{{ ipv4 ? '✓' : '✕' }}</span>
          </div>
          <div class="flex items-center gap-1.5 px-3 py-1 rounded-md border" :class="ipv6 ? 'bg-emerald-50 border-emerald-200 text-emerald-700 dark:bg-emerald-950/40 dark:border-emerald-800 dark:text-emerald-300' : 'bg-slate-100 border-slate-200 text-slate-400 dark:bg-slate-800 dark:border-slate-700'">
            <span>IPv6</span>
            <span>{{ ipv6 ? '✓' : '✕' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Tools Grid -->
    <div>
      <h2 class="text-lg font-bold text-slate-900 dark:text-white mb-4">
        网络诊断工具箱
      </h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <router-link
          v-for="tool in quickTools"
          :key="tool.path"
          :to="tool.path"
          class="custom-card custom-card-hover p-4 flex flex-col justify-between group"
        >
          <div>
            <div class="w-10 h-10 rounded-lg flex items-center justify-center mb-3 transition-transform group-hover:scale-105" :class="tool.bg">
              <component :is="tool.icon" class="w-5 h-5" :class="tool.color" />
            </div>
            <h3 class="font-bold text-sm text-slate-900 dark:text-white group-hover:text-brand-600 dark:group-hover:text-brand-400 transition-colors">
              {{ tool.name }}
            </h3>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 leading-relaxed">
              {{ tool.desc }}
            </p>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>
