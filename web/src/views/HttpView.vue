<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Activity, Search, ArrowRight, Clock, Server, FileText } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import LatencyBadge from '@/components/LatencyBadge.vue'
import { api } from '@/utils/api'
import { formatBytes } from '@/utils/format'
import type { HTTPCheckResult } from '@/types'

const route = useRoute()
const router = useRouter()

const targetUrl = ref('https://ipw.3x.cx')
const loading = ref(false)
const error = ref('')
const result = ref<HTTPCheckResult | null>(null)

async function executeCheck() {
  const url = targetUrl.value.trim()
  if (!url) return

  loading.value = true
  error.value = ''

  try {
    result.value = await api.checkHTTP(url)
    router.replace({ query: { url } })
  } catch (e: any) {
    error.value = e.message || 'HTTP 检测失败'
    result.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (route.query.url) {
    targetUrl.value = route.query.url as string
  }
  if (targetUrl.value) {
    executeCheck()
  }
})
</script>

<template>
  <div class="space-y-6">
    <ToolHeader
      title="网站 HTTP 响应与耗时分析"
      description="发起真实的 HTTP 请求并记录 DNS 解析、TCP 握手、TLS 协商、首字节到达时间 (TTFB) 及完整响应流程"
      :icon="Activity"
    />

    <!-- Input Form -->
    <div class="custom-card p-4 sm:p-6">
      <form @submit.prevent="executeCheck" class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="targetUrl"
            type="text"
            placeholder="输入目标网址，例如 https://ipw.3x.cx"
            aria-label="输入待检测网址"
            class="w-full pl-10 pr-4 py-2.5 rounded-lg border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/50"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-3.5 top-3.5" />
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="inline-flex items-center justify-center gap-2 px-6 py-2.5 rounded-lg bg-brand-600 hover:bg-brand-700 text-white text-sm font-semibold shadow-sm transition-all disabled:opacity-50"
        >
          <span>{{ loading ? '诊断中...' : '发起诊断' }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </form>
    </div>

    <!-- Error Alert -->
    <ErrorAlert v-if="error" :message="error" />

    <!-- Result View -->
    <div v-if="result" class="space-y-6">
      <!-- Status Card -->
      <div class="custom-card p-6 flex flex-wrap items-center justify-between gap-4 border-l-4 border-l-brand-500">
        <div>
          <div class="flex items-center gap-3">
            <span
              class="px-2.5 py-0.5 rounded text-xs font-mono font-bold uppercase border"
              :class="
                result.status_code >= 200 && result.status_code < 400
                  ? 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/60 dark:text-emerald-300'
                  : 'bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-950/60 dark:text-rose-300'
              "
            >
              {{ result.status_text }}
            </span>
            <span class="text-xs font-mono px-2 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 font-semibold">
              {{ result.protocol }}
            </span>
          </div>
          <div class="font-mono text-xs text-slate-500 mt-2">
            目标: {{ result.url }}
            <span v-if="result.resolved_ip">({{ result.resolved_ip }})</span>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <Clock class="w-5 h-5 text-brand-600" />
          <div class="text-right">
            <div class="text-xs text-slate-400">总耗时</div>
            <div class="text-xl font-bold font-mono text-slate-900 dark:text-white">{{ result.total_ms }} ms</div>
          </div>
        </div>
      </div>

      <!-- Waterfall Timeline Breakdown -->
      <div class="custom-card p-6 space-y-4">
        <h4 class="text-base font-bold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3">
          网络耗时瀑布图 (Waterfall Breakdown)
        </h4>

        <div class="space-y-3 font-mono text-xs sm:text-sm">
          <!-- DNS -->
          <div class="flex items-center justify-between p-3 rounded-lg bg-slate-50 dark:bg-slate-800/40">
            <span class="font-medium text-slate-700 dark:text-slate-300 font-sans">1. DNS 域名解析</span>
            <LatencyBadge :latency-ms="result.dns_ms" />
          </div>

          <!-- TCP Connect -->
          <div class="flex items-center justify-between p-3 rounded-lg bg-slate-50 dark:bg-slate-800/40">
            <span class="font-medium text-slate-700 dark:text-slate-300 font-sans">2. TCP 三次握手</span>
            <LatencyBadge :latency-ms="result.tcp_ms" />
          </div>

          <!-- TLS Handshake -->
          <div v-if="result.tls_ms > 0" class="flex items-center justify-between p-3 rounded-lg bg-slate-50 dark:bg-slate-800/40">
            <span class="font-medium text-slate-700 dark:text-slate-300 font-sans">3. TLS 安全握手</span>
            <LatencyBadge :latency-ms="result.tls_ms" />
          </div>

          <!-- TTFB -->
          <div class="flex items-center justify-between p-3 rounded-lg bg-slate-50 dark:bg-slate-800/40">
            <span class="font-medium text-slate-700 dark:text-slate-300 font-sans">4. 首字节时间 (TTFB)</span>
            <LatencyBadge :latency-ms="result.ttfb_ms" />
          </div>
        </div>
      </div>

      <!-- HTTP Headers Grid -->
      <div class="custom-card p-6 space-y-4">
        <h4 class="text-base font-bold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3">
          响应头与服务器信息
        </h4>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-sm font-mono">
          <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/50">
            <div class="text-xs text-slate-400 font-sans flex items-center gap-1">
              <Server class="w-3.5 h-3.5" />
              <span>Server</span>
            </div>
            <div class="font-bold text-slate-800 dark:text-slate-200 mt-1 truncate">{{ result.server || '-' }}</div>
          </div>

          <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/50">
            <div class="text-xs text-slate-400 font-sans flex items-center gap-1">
              <FileText class="w-3.5 h-3.5" />
              <span>Content-Type</span>
            </div>
            <div class="font-bold text-slate-800 dark:text-slate-200 mt-1 truncate">{{ result.content_type || '-' }}</div>
          </div>

          <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/50">
            <div class="text-xs text-slate-400 font-sans">Content-Length</div>
            <div class="font-bold text-slate-800 dark:text-slate-200 mt-1">
              {{ result.content_length && result.content_length > 0 ? formatBytes(result.content_length) : '-' }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
