<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Gauge, Search, ArrowRight, Download, Zap } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import { api } from '@/utils/api'
import { formatBytes } from '@/utils/format'
import type { SpeedTestResult } from '@/types'

const route = useRoute()
const router = useRouter()

const target = ref('https://ipw.3x.cx')
const loading = ref(false)
const error = ref('')
const result = ref<SpeedTestResult | null>(null)

async function executeSpeedTest() {
  const url = target.value.trim()
  if (!url) return

  loading.value = true
  error.value = ''

  try {
    result.value = await api.speedTest(url)
    router.replace({ query: { url } })
  } catch (e: any) {
    error.value = e.message || '测速失败'
    result.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (route.query.url) {
    target.value = route.query.url as string
  }
  if (target.value) {
    executeSpeedTest()
  }
})
</script>

<template>
  <div class="space-y-6">
    <ToolHeader
      title="网站 HTTP 测速"
      description="连接目标网站并下载最高 5MB 静态数据包，计算实际下载吞吐速率 (Mbps) 与分段延迟"
      :icon="Gauge"
    />

    <!-- Input Form -->
    <div class="custom-card p-5 sm:p-7 shadow-card">
      <form @submit.prevent="executeSpeedTest" class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="target"
            type="text"
            placeholder="输入目标网址，例如 https://ipw.3x.cx"
            aria-label="输入待测速网址"
            class="w-full pl-11 pr-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700/80 bg-slate-50/60 dark:bg-slate-950/60 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition-all shadow-inner"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-4 top-3.5" />
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="btn-primary"
        >
          <span>{{ loading ? '测速中...' : '开始测速' }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </form>
    </div>

    <!-- Error Alert -->
    <ErrorAlert v-if="error" :message="error" />

    <!-- Result View -->
    <div v-if="result" class="space-y-6">
      <!-- Big Speed Highlight Card -->
      <div class="custom-card p-8 bg-gradient-to-br from-brand-500/10 via-transparent to-purple-500/10 border-brand-200 dark:border-brand-900/60 text-center space-y-3">
        <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-brand-50 text-brand-700 dark:bg-brand-950/60 dark:text-brand-300 border border-brand-200/60 dark:border-brand-800">
          <Zap class="w-3.5 h-3.5" />
          <span>下载速率</span>
        </div>

        <div class="font-mono text-4xl sm:text-6xl font-extrabold text-slate-900 dark:text-white tracking-tight">
          {{ result.speed_mbps }} <span class="text-xl sm:text-2xl font-bold text-brand-600 dark:text-brand-400">Mbps</span>
        </div>

        <div class="flex flex-wrap items-center justify-center gap-4 text-xs font-mono text-slate-500 dark:text-slate-400 pt-2">
          <span>下载流量: {{ formatBytes(result.download_bytes) }}</span>
          <span>·</span>
          <span>下载耗时: {{ result.download_ms }} ms</span>
          <span v-if="result.resolved_ip">· 解析 IP: {{ result.resolved_ip }}</span>
        </div>
      </div>

      <!-- Breakdown Grid -->
      <div class="custom-card p-6 space-y-4">
        <h4 class="text-base font-bold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3 flex items-center gap-2">
          <Download class="w-4 h-4 text-brand-600" />
          <span>耗时分段明细</span>
        </h4>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 font-mono text-xs sm:text-sm text-center">
          <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/50">
            <div class="text-xs text-slate-400 font-sans">DNS 解析</div>
            <div class="mt-1 font-bold text-slate-800 dark:text-slate-200">{{ result.dns_ms }} ms</div>
          </div>

          <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/50">
            <div class="text-xs text-slate-400 font-sans">TCP 连接</div>
            <div class="mt-1 font-bold text-slate-800 dark:text-slate-200">{{ result.connect_ms }} ms</div>
          </div>

          <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/50">
            <div class="text-xs text-slate-400 font-sans">TLS 握手</div>
            <div class="mt-1 font-bold text-slate-800 dark:text-slate-200">{{ result.tls_ms }} ms</div>
          </div>

          <div class="p-3 rounded-lg bg-slate-50 dark:bg-slate-800/50">
            <div class="text-xs text-slate-400 font-sans">首字节 (TTFB)</div>
            <div class="mt-1 font-bold text-slate-800 dark:text-slate-200">{{ result.ttfb_ms }} ms</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
