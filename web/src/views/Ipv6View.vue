<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CheckCircle2, XCircle, Search, ArrowRight } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import LatencyBadge from '@/components/LatencyBadge.vue'
import { api } from '@/utils/api'
import type { IPv6CheckResponse } from '@/types'

const route = useRoute()
const router = useRouter()

const domain = ref('ipw.3x.cx')
const loading = ref(false)
const error = ref('')
const result = ref<IPv6CheckResponse | null>(null)

async function executeCheck() {
  const target = domain.value.trim()
  if (!target) return

  loading.value = true
  error.value = ''

  try {
    result.value = await api.checkIPv6(target)
    router.replace({ query: { domain: target } })
  } catch (e: any) {
    error.value = e.message || '检测失败'
    result.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (route.query.domain) {
    domain.value = route.query.domain as string
  }
  if (domain.value) {
    executeCheck()
  }
})
</script>

<template>
  <div class="space-y-6">
    <ToolHeader
      title="网站 IPv6 支持检测"
      description="全面测试目标网站是否已具备 IPv6 DNS (AAAA) 解析、IPv6 HTTP 以及 IPv6 HTTPS 双栈连通性"
      :icon="CheckCircle2"
    />

    <!-- Input Form -->
    <div class="custom-card p-4 sm:p-6">
      <form @submit.prevent="executeCheck" class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="domain"
            type="text"
            placeholder="输入域名，例如 ipw.3x.cx 或 baidu.com"
            aria-label="输入待检测网站域名"
            class="w-full pl-10 pr-4 py-2.5 rounded-lg border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/50"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-3.5 top-3.5" />
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="inline-flex items-center justify-center gap-2 px-6 py-2.5 rounded-lg bg-brand-600 hover:bg-brand-700 text-white text-sm font-semibold shadow-sm transition-all disabled:opacity-50"
        >
          <span>{{ loading ? '检测中...' : '开始检测' }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </form>
    </div>

    <!-- Error Alert -->
    <ErrorAlert v-if="error" :message="error" />

    <!-- Result View -->
    <div v-if="result" class="space-y-6">
      <!-- Conclusion Banner -->
      <div
        class="custom-card p-6 border-l-4"
        :class="
          result.supported
            ? 'border-l-emerald-500 bg-emerald-50/40 dark:bg-emerald-950/20'
            : 'border-l-amber-500 bg-amber-50/40 dark:bg-amber-950/20'
        "
      >
        <div class="flex items-start gap-4">
          <CheckCircle2 v-if="result.supported" class="w-8 h-8 text-emerald-600 dark:text-emerald-400 flex-shrink-0" />
          <XCircle v-else class="w-8 h-8 text-amber-500 flex-shrink-0" />
          <div>
            <h3 class="text-lg font-bold text-slate-900 dark:text-white">
              {{ result.conclusion }}
            </h3>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 font-mono">
              目标域名: {{ result.domain }}
            </p>
          </div>
        </div>
      </div>

      <!-- Test Breakdown Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- IPv4 Breakdown -->
        <div class="custom-card p-6 space-y-4">
          <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-3">
            <span class="font-bold text-base text-blue-600 dark:text-blue-400">IPv4 协议栈检测</span>
            <span class="text-xs px-2 py-0.5 rounded bg-blue-50 text-blue-700 dark:bg-blue-950/60 dark:text-blue-300 font-mono">
              A 记录
            </span>
          </div>

          <div class="space-y-3 text-sm">
            <div class="flex items-center justify-between py-1">
              <span class="text-slate-600 dark:text-slate-300">DNS 解析 (A)</span>
              <span v-if="result.ipv4.dns" class="text-emerald-600 dark:text-emerald-400 font-mono text-xs font-bold">
                ✓ 已配置 ({{ result.ipv4.addresses.length }} 个 IP)
              </span>
              <span v-else class="text-rose-500 text-xs font-mono">✕ 未配置</span>
            </div>

            <div v-if="result.ipv4.addresses.length > 0" class="bg-slate-50 dark:bg-slate-800/50 p-2.5 rounded-lg text-xs font-mono space-y-1 text-slate-600 dark:text-slate-300">
              <div v-for="ip in result.ipv4.addresses" :key="ip">{{ ip }}</div>
            </div>

            <div class="flex items-center justify-between py-1 border-t border-slate-100 dark:border-slate-800">
              <span class="text-slate-600 dark:text-slate-300">HTTP 连接 (Port 80)</span>
              <div v-if="result.ipv4.http.supported" class="flex items-center gap-2">
                <span class="text-xs text-emerald-600 font-mono font-bold">{{ result.ipv4.http.status_code }} OK</span>
                <LatencyBadge :latency-ms="result.ipv4.http.latency_ms || 0" />
              </div>
              <span v-else class="text-xs text-slate-400">{{ result.ipv4.http.error || '未启用' }}</span>
            </div>

            <div class="flex items-center justify-between py-1 border-t border-slate-100 dark:border-slate-800">
              <span class="text-slate-600 dark:text-slate-300">HTTPS 连接 (Port 443)</span>
              <div v-if="result.ipv4.https.supported" class="flex items-center gap-2">
                <span class="text-xs text-emerald-600 font-mono font-bold">{{ result.ipv4.https.status_code }} OK</span>
                <LatencyBadge :latency-ms="result.ipv4.https.latency_ms || 0" />
              </div>
              <span v-else class="text-xs text-slate-400">{{ result.ipv4.https.error || '未启用' }}</span>
            </div>
          </div>
        </div>

        <!-- IPv6 Breakdown -->
        <div class="custom-card p-6 space-y-4">
          <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-3">
            <span class="font-bold text-base text-purple-600 dark:text-purple-400">IPv6 协议栈检测</span>
            <span class="text-xs px-2 py-0.5 rounded bg-purple-50 text-purple-700 dark:bg-purple-950/60 dark:text-purple-300 font-mono">
              AAAA 记录
            </span>
          </div>

          <div class="space-y-3 text-sm">
            <div class="flex items-center justify-between py-1">
              <span class="text-slate-600 dark:text-slate-300">DNS 解析 (AAAA)</span>
              <span v-if="result.ipv6.dns" class="text-emerald-600 dark:text-emerald-400 font-mono text-xs font-bold">
                ✓ 已配置 ({{ result.ipv6.addresses.length }} 个 IP)
              </span>
              <span v-else class="text-rose-500 text-xs font-mono">✕ 未配置</span>
            </div>

            <div v-if="result.ipv6.addresses.length > 0" class="bg-slate-50 dark:bg-slate-800/50 p-2.5 rounded-lg text-xs font-mono space-y-1 text-slate-600 dark:text-slate-300 break-all">
              <div v-for="ip in result.ipv6.addresses" :key="ip">{{ ip }}</div>
            </div>

            <div class="flex items-center justify-between py-1 border-t border-slate-100 dark:border-slate-800">
              <span class="text-slate-600 dark:text-slate-300">HTTP 连接 (Port 80)</span>
              <div v-if="result.ipv6.http.supported" class="flex items-center gap-2">
                <span class="text-xs text-emerald-600 font-mono font-bold">{{ result.ipv6.http.status_code }} OK</span>
                <LatencyBadge :latency-ms="result.ipv6.http.latency_ms || 0" />
              </div>
              <span v-else class="text-xs text-slate-400">{{ result.ipv6.http.error || '未启用' }}</span>
            </div>

            <div class="flex items-center justify-between py-1 border-t border-slate-100 dark:border-slate-800">
              <span class="text-slate-600 dark:text-slate-300">HTTPS 连接 (Port 443)</span>
              <div v-if="result.ipv6.https.supported" class="flex items-center gap-2">
                <span class="text-xs text-emerald-600 font-mono font-bold">{{ result.ipv6.https.status_code }} OK</span>
                <LatencyBadge :latency-ms="result.ipv6.https.latency_ms || 0" />
              </div>
              <span v-else class="text-xs text-slate-400">{{ result.ipv6.https.error || '未启用' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
