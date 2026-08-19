<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Lock, Search, ArrowRight, ShieldCheck, AlertTriangle, Calendar } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import { api } from '@/utils/api'
import { formatDate } from '@/utils/format'
import type { SSLCheckResult } from '@/types'

const route = useRoute()
const router = useRouter()

const hostname = ref('ipw.3x.cx')
const port = ref(443)
const loading = ref(false)
const error = ref('')
const result = ref<SSLCheckResult | null>(null)

async function executeCheck() {
  const host = hostname.value.trim()
  if (!host) return

  loading.value = true
  error.value = ''

  try {
    result.value = await api.checkSSL(host, port.value)
    router.replace({
      query: { host, port: port.value !== 443 ? port.value : undefined },
    })
  } catch (e: any) {
    error.value = e.message || 'SSL 证书检查失败'
    result.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (route.query.host) {
    hostname.value = route.query.host as string
  }
  if (route.query.port) {
    port.value = parseInt(route.query.port as string) || 443
  }
  if (hostname.value) {
    executeCheck()
  }
})
</script>

<template>
  <div class="space-y-6">
    <ToolHeader
      title="SSL / TLS 证书查询"
      description="连接目标 HTTPS 服务器并解析 X.509 证书链，查看证书颁发机构、剩余有效期、SANs 多域名及加密套件"
      :icon="Lock"
    />

    <!-- Input Form -->
    <div class="custom-card p-4 sm:p-6">
      <form @submit.prevent="executeCheck" class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="hostname"
            type="text"
            placeholder="输入域名或主机名，例如 ipw.3x.cx 或 github.com"
            aria-label="输入待检测主机名"
            class="w-full pl-10 pr-4 py-2.5 rounded-lg border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/50"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-3.5 top-3.5" />
        </div>

        <div class="w-full sm:w-28">
          <input
            v-model.number="port"
            type="number"
            placeholder="端口"
            aria-label="HTTPS 端口"
            class="w-full px-3 py-2.5 rounded-lg border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/50"
          />
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="inline-flex items-center justify-center gap-2 px-6 py-2.5 rounded-lg bg-brand-600 hover:bg-brand-700 text-white text-sm font-semibold shadow-sm transition-all disabled:opacity-50"
        >
          <span>{{ loading ? '查询中...' : '查询证书' }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </form>
    </div>

    <!-- Error Alert -->
    <ErrorAlert v-if="error" :message="error" />

    <!-- Result View -->
    <div v-if="result" class="space-y-6">
      <!-- Status Card -->
      <div
        class="custom-card p-6 border-l-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
        :class="
          result.valid && result.days_remaining > 0
            ? 'border-l-emerald-500 bg-emerald-50/40 dark:bg-emerald-950/20'
            : 'border-l-rose-500 bg-rose-50/40 dark:bg-rose-950/20'
        "
      >
        <div class="flex items-center gap-3">
          <ShieldCheck
            v-if="result.valid && result.days_remaining > 0"
            class="w-8 h-8 text-emerald-600 dark:text-emerald-400 flex-shrink-0"
          />
          <AlertTriangle v-else class="w-8 h-8 text-rose-500 flex-shrink-0" />
          <div>
            <h3 class="text-lg font-bold text-slate-900 dark:text-white">
              {{ result.valid && result.days_remaining > 0 ? '证书有效且信任' : '证书无效或已过期' }}
            </h3>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 font-mono">
              {{ result.hostname }}:{{ result.port }}
              <span v-if="result.resolved_ip">({{ result.resolved_ip }})</span>
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <span class="text-xs text-slate-500 dark:text-slate-400">剩余有效期:</span>
          <span
            class="px-3 py-1 rounded-full text-sm font-mono font-bold"
            :class="
              result.days_remaining > 30
                ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/60 dark:text-emerald-200'
                : result.days_remaining > 0
                ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/60 dark:text-amber-200'
                : 'bg-rose-100 text-rose-800 dark:bg-rose-900/60 dark:text-rose-200'
            "
          >
            {{ result.days_remaining }} 天
          </span>
        </div>
      </div>

      <!-- Details Grid -->
      <div class="custom-card p-6 space-y-6">
        <h4 class="text-base font-bold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3">
          证书核心信息
        </h4>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
          <div>
            <span class="text-xs text-slate-400 font-medium">使用者 (Subject)</span>
            <div class="font-mono text-slate-800 dark:text-slate-200 mt-1 font-semibold break-all">{{ result.subject }}</div>
          </div>

          <div>
            <span class="text-xs text-slate-400 font-medium">颁发机构 (Issuer)</span>
            <div class="font-mono text-slate-800 dark:text-slate-200 mt-1 font-semibold break-all">{{ result.issuer }}</div>
          </div>

          <div>
            <span class="text-xs text-slate-400 font-medium flex items-center gap-1">
              <Calendar class="w-3.5 h-3.5" />
              <span>生效日期 (Not Before)</span>
            </span>
            <div class="font-mono text-slate-700 dark:text-slate-300 mt-1">{{ formatDate(result.not_before) }}</div>
          </div>

          <div>
            <span class="text-xs text-slate-400 font-medium flex items-center gap-1">
              <Calendar class="w-3.5 h-3.5" />
              <span>到期日期 (Not After)</span>
            </span>
            <div class="font-mono text-slate-700 dark:text-slate-300 mt-1">{{ formatDate(result.not_after) }}</div>
          </div>

          <div>
            <span class="text-xs text-slate-400 font-medium">TLS 协议版本</span>
            <div class="font-mono text-slate-700 dark:text-slate-300 mt-1">{{ result.tls_version }}</div>
          </div>

          <div>
            <span class="text-xs text-slate-400 font-medium">加密套件 (Cipher Suite)</span>
            <div class="font-mono text-xs text-slate-700 dark:text-slate-300 mt-1 break-all">{{ result.cipher_suite }}</div>
          </div>
        </div>

        <!-- SANs List -->
        <div v-if="result.dns_names && result.dns_names.length > 0" class="border-t border-slate-100 dark:border-slate-800 pt-4">
          <span class="text-xs text-slate-400 font-medium block mb-2">主体备用名称 (SANs / DNS Names)</span>
          <div class="flex flex-wrap gap-1.5 font-mono text-xs">
            <span
              v-for="name in result.dns_names"
              :key="name"
              class="px-2 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-700"
            >
              {{ name }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
