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
  ArrowRight,
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
    name: 'IP 归属地查询',
    desc: '查询指定 IP 的高精度地理归属地、运营商与 ASN',
    path: '/ip',
    icon: Globe2,
    color: 'text-brand-600 dark:text-brand-400',
    bg: 'bg-brand-50 dark:bg-brand-950/60',
    border: 'group-hover:border-brand-500/40',
  },
  {
    name: '多 DNS 对比解析',
    desc: '并发向 1.1.1.1、8.8.8.8、223.5.5.5 等解析器查询记录',
    path: '/dns',
    icon: Network,
    color: 'text-indigo-600 dark:text-indigo-400',
    bg: 'bg-indigo-50 dark:bg-indigo-950/60',
    border: 'group-hover:border-indigo-500/40',
  },
  {
    name: 'Ping / TCPing 诊断',
    desc: '测试目标主机 ICMP 丢包率与 TCP 端口三次握手耗时',
    path: '/ping',
    icon: Activity,
    color: 'text-emerald-600 dark:text-emerald-400',
    bg: 'bg-emerald-50 dark:bg-emerald-950/60',
    border: 'group-hover:border-emerald-500/40',
  },
  {
    name: '网站 IPv6 体检',
    desc: '全面检测目标网站的 IPv6 双栈连通性与解析就绪度',
    path: '/ipv6',
    icon: CheckCircle2,
    color: 'text-purple-600 dark:text-purple-400',
    bg: 'bg-purple-50 dark:bg-purple-950/60',
    border: 'group-hover:border-purple-500/40',
  },
  {
    name: 'SSL / TLS 证书解析',
    desc: '深度解析目标服务器 HTTPS 证书链、颁发者与剩余有效期',
    path: '/ssl',
    icon: Lock,
    color: 'text-cyan-600 dark:text-cyan-400',
    bg: 'bg-cyan-50 dark:bg-cyan-950/60',
    border: 'group-hover:border-cyan-500/40',
  },
  {
    name: 'WHOIS / RDAP 查询',
    desc: '查询域名或 IP 地址的权威注册商、创建时间与过期日期',
    path: '/whois',
    icon: Search,
    color: 'text-blue-600 dark:text-blue-400',
    bg: 'bg-blue-50 dark:bg-blue-950/60',
    border: 'group-hover:border-blue-500/40',
  },
  {
    name: 'HTTP 测速与耗时分析',
    desc: '精确测量 TTFB 首包时延、吞吐速率 (Mbps) 与下载瀑布流',
    path: '/speed',
    icon: Gauge,
    color: 'text-rose-600 dark:text-rose-400',
    bg: 'bg-rose-50 dark:bg-rose-950/60',
    border: 'group-hover:border-rose-500/40',
  },
  {
    name: 'ASN 自治系统查询',
    desc: '查询自治系统编号 (ASN)、广播路由前缀与网络组织',
    path: '/asn',
    icon: Layers,
    color: 'text-amber-600 dark:text-amber-400',
    bg: 'bg-amber-50 dark:bg-amber-950/60',
    border: 'group-hover:border-amber-500/40',
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

  // Fallback: If 4. / 6. subdomains didn't respond, use /api/v1/me
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
    <!-- Hero Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div class="flex items-center gap-3.5 sm:gap-4">
        <div class="relative group flex-shrink-0">
          <div class="absolute -inset-0.5 bg-gradient-to-r from-brand-600 to-sky-400 rounded-2xl blur opacity-30 group-hover:opacity-60 transition duration-300"></div>
          <img src="/favicon.svg" alt="NetIP" class="relative w-12 h-12 sm:w-14 sm:h-14 rounded-2xl shadow-card" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white">
              您的网络与 IP 状态
            </h1>
          </div>
          <p class="mt-0.5 text-xs sm:text-sm text-slate-500 dark:text-slate-400 font-medium">
            实时公网 IPv4 / IPv6 双栈探测与访问优先级分析
          </p>
        </div>
      </div>

      <button
        type="button"
        @click="detectNetwork"
        :disabled="loading"
        class="btn-secondary self-start sm:self-auto"
      >
        <RefreshCw class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" :class="{ 'animate-spin': loading }" />
        <span>重新检测</span>
      </button>
    </div>

    <!-- Main IP Hero Cards Grid -->
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
    <div class="custom-card p-5 sm:p-6 bg-gradient-to-r from-brand-500/10 via-brand-500/3 to-sky-500/5 border-brand-200/60 dark:border-brand-900/50 shadow-card">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div class="flex items-center gap-3.5">
          <div class="w-11 h-11 rounded-xl bg-brand-600/10 dark:bg-brand-500/20 text-brand-600 dark:text-brand-400 flex items-center justify-center flex-shrink-0 shadow-2xs">
            <ArrowUpDown class="w-5 h-5" />
          </div>
          <div>
            <div class="text-[11px] uppercase tracking-wider font-bold text-slate-400 dark:text-slate-500">
              网络访问优先级分析
            </div>
            <div class="text-lg font-extrabold text-slate-900 dark:text-white mt-0.5 flex items-center gap-2">
              <span v-if="loading" class="text-slate-400 text-sm">正在测定双栈优先级...</span>
              <span v-else-if="preferred === 'ipv6'" class="text-purple-600 dark:text-purple-400 flex items-center gap-1.5">
                <span>IPv6 访问优先</span>
                <span class="text-xs px-2 py-0.5 rounded-full bg-purple-100 dark:bg-purple-950/80 font-normal">现代化网络</span>
              </span>
              <span v-else-if="preferred === 'ipv4'" class="text-brand-600 dark:text-brand-400 flex items-center gap-1.5">
                <span>IPv4 访问优先</span>
                <span class="text-xs px-2 py-0.5 rounded-full bg-brand-100 dark:bg-brand-950/80 font-normal">经典模式</span>
              </span>
              <span v-else class="text-slate-500">常规模式</span>
            </div>
          </div>
        </div>

        <!-- Dual Stack Status Badges -->
        <div class="flex items-center gap-2.5 text-xs font-mono font-semibold self-stretch sm:self-auto justify-end">
          <div
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl border transition-all"
            :class="ipv4 ? 'bg-emerald-50 border-emerald-200/90 text-emerald-700 dark:bg-emerald-950/60 dark:border-emerald-800 dark:text-emerald-300' : 'bg-slate-100 border-slate-200 text-slate-400 dark:bg-slate-800 dark:border-slate-700'"
          >
            <span class="w-1.5 h-1.5 rounded-full" :class="ipv4 ? 'bg-emerald-500' : 'bg-slate-400'"></span>
            <span>IPv4 {{ ipv4 ? '已就绪' : '未就绪' }}</span>
          </div>

          <div
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl border transition-all"
            :class="ipv6 ? 'bg-purple-50 border-purple-200/90 text-purple-700 dark:bg-purple-950/60 dark:border-purple-800 dark:text-purple-300' : 'bg-slate-100 border-slate-200 text-slate-400 dark:bg-slate-800 dark:border-slate-700'"
          >
            <span class="w-1.5 h-1.5 rounded-full" :class="ipv6 ? 'bg-purple-500' : 'bg-slate-400'"></span>
            <span>IPv6 {{ ipv6 ? '已就绪' : '未就绪' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Tools Grid Section -->
    <div class="space-y-4 pt-2">
      <div class="flex items-center justify-between">
        <h2 class="text-lg sm:text-xl font-extrabold text-slate-900 dark:text-white flex items-center gap-2">
          <span>网络诊断工具箱</span>
          <span class="text-xs font-normal text-slate-400 px-2 py-0.5 rounded-full bg-slate-100 dark:bg-slate-800">8 大全能工具</span>
        </h2>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-5">
        <router-link
          v-for="tool in quickTools"
          :key="tool.path"
          :to="tool.path"
          class="custom-card custom-card-hover p-5 flex flex-col justify-between group relative overflow-hidden"
          :class="tool.border"
        >
          <div>
            <div class="w-11 h-11 rounded-xl flex items-center justify-center mb-4 transition-transform duration-200 group-hover:scale-110 shadow-2xs" :class="tool.bg">
              <component :is="tool.icon" class="w-5 h-5" :class="tool.color" />
            </div>
            <h3 class="font-bold text-sm text-slate-900 dark:text-white group-hover:text-brand-600 dark:group-hover:text-brand-400 transition-colors flex items-center justify-between">
              <span>{{ tool.name }}</span>
              <ArrowRight class="w-4 h-4 text-slate-300 dark:text-slate-600 group-hover:text-brand-600 dark:group-hover:text-brand-400 group-hover:translate-x-1 transition-all duration-200" />
            </h3>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-1.5 leading-relaxed">
              {{ tool.desc }}
            </p>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>
