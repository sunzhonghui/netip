<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search, ArrowRight, FileText, Calendar, Building, Globe } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import { api } from '@/utils/api'
import { formatDate } from '@/utils/format'
import type { WHOISResult } from '@/types'

const route = useRoute()
const router = useRouter()

const query = ref('ipw.3x.cx')
const loading = ref(false)
const error = ref('')
const result = ref<WHOISResult | null>(null)
const showRaw = ref(false)

async function executeQuery() {
  const target = query.value.trim()
  if (!target) return

  loading.value = true
  error.value = ''

  try {
    result.value = await api.queryWHOIS(target)
    router.replace({ query: { q: target } })
  } catch (e: any) {
    error.value = e.message || 'WHOIS 查询失败'
    result.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (route.query.q) {
    query.value = route.query.q as string
  }
  if (query.value) {
    executeQuery()
  }
})
</script>

<template>
  <div class="space-y-6">
    <ToolHeader
      title="WHOIS / RDAP 注册信息查询"
      description="基于现代 RDAP 协议与权威 WHOIS 服务器，查询域名或 IP 地址的注册机构、状态、创建及过期生命周期"
      :icon="Search"
    />

    <!-- Input Form -->
    <div class="custom-card p-5 sm:p-7 shadow-card">
      <form @submit.prevent="executeQuery" class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="query"
            type="text"
            placeholder="输入域名或 IP，例如 ipw.3x.cx 或 8.8.8.8"
            aria-label="输入待查询域名或 IP"
            class="w-full pl-11 pr-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700/80 bg-slate-50/60 dark:bg-slate-950/60 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition-all shadow-inner"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-4 top-3.5" />
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="btn-primary"
        >
          <span>{{ loading ? '查询中...' : '查询 WHOIS' }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </form>
    </div>

    <!-- Error Alert -->
    <ErrorAlert v-if="error" :message="error" />

    <!-- Result View -->
    <div v-if="result" class="space-y-6">
      <!-- Domain Details -->
      <div v-if="result.type === 'domain'" class="custom-card p-6 space-y-6">
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-4">
          <div>
            <div class="font-mono text-xl font-bold text-slate-900 dark:text-white">
              {{ result.domain || result.query }}
            </div>
            <div class="text-xs text-slate-400 mt-1 font-mono">
              协议源: {{ result.source }}
            </div>
          </div>
          <span class="text-xs font-mono font-medium px-2.5 py-1 rounded bg-brand-50 text-brand-700 dark:bg-brand-950/60 dark:text-brand-300">
            域名
          </span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
          <div v-if="result.registrar">
            <span class="text-xs text-slate-400 font-medium flex items-center gap-1">
              <Building class="w-3.5 h-3.5" />
              <span>注册商 (Registrar)</span>
            </span>
            <div class="font-mono text-slate-800 dark:text-slate-200 mt-1 font-semibold">{{ result.registrar }}</div>
          </div>

          <div v-if="result.dnssec">
            <span class="text-xs text-slate-400 font-medium">DNSSEC 状态</span>
            <div class="font-mono text-slate-800 dark:text-slate-200 mt-1">{{ result.dnssec }}</div>
          </div>

          <div v-if="result.created">
            <span class="text-xs text-slate-400 font-medium flex items-center gap-1">
              <Calendar class="w-3.5 h-3.5" />
              <span>注册时间 (Created)</span>
            </span>
            <div class="font-mono text-slate-700 dark:text-slate-300 mt-1">{{ formatDate(result.created) }}</div>
          </div>

          <div v-if="result.expires">
            <span class="text-xs text-slate-400 font-medium flex items-center gap-1">
              <Calendar class="w-3.5 h-3.5" />
              <span>到期时间 (Expires)</span>
            </span>
            <div class="font-mono text-slate-700 dark:text-slate-300 mt-1">{{ formatDate(result.expires) }}</div>
          </div>
        </div>

        <!-- Name Servers -->
        <div v-if="result.name_servers && result.name_servers.length > 0" class="border-t border-slate-100 dark:border-slate-800 pt-4">
          <span class="text-xs text-slate-400 font-medium block mb-2">权威 DNS 服务器 (Name Servers)</span>
          <div class="flex flex-wrap gap-2 font-mono text-xs">
            <span
              v-for="ns in result.name_servers"
              :key="ns"
              class="px-2.5 py-1 rounded bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300"
            >
              {{ ns }}
            </span>
          </div>
        </div>

        <!-- Status -->
        <div v-if="result.status && result.status.length > 0" class="border-t border-slate-100 dark:border-slate-800 pt-4">
          <span class="text-xs text-slate-400 font-medium block mb-2">域名状态 (Domain Status)</span>
          <div class="flex flex-wrap gap-1.5 font-mono text-xs">
            <span
              v-for="st in result.status"
              :key="st"
              class="px-2 py-0.5 rounded bg-slate-50 dark:bg-slate-800/80 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700"
            >
              {{ st }}
            </span>
          </div>
        </div>
      </div>

      <!-- IP Details -->
      <div v-else class="custom-card p-6 space-y-6">
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-4">
          <div>
            <div class="font-mono text-xl font-bold text-slate-900 dark:text-white">
              {{ result.query }}
            </div>
            <div class="text-xs text-slate-400 mt-1 font-mono">
              协议源: {{ result.source }}
            </div>
          </div>
          <span class="text-xs font-mono font-medium px-2.5 py-1 rounded bg-blue-50 text-blue-700 dark:bg-blue-950/60 dark:text-blue-300">
            IP 网段
          </span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
          <div v-if="result.network">
            <span class="text-xs text-slate-400 font-medium">广播范围 (Network)</span>
            <div class="font-mono text-slate-800 dark:text-slate-200 mt-1 font-semibold">{{ result.network }}</div>
          </div>

          <div v-if="result.organization">
            <span class="text-xs text-slate-400 font-medium flex items-center gap-1">
              <Building class="w-3.5 h-3.5" />
              <span>所属机构 (Organization)</span>
            </span>
            <div class="font-mono text-slate-800 dark:text-slate-200 mt-1 font-semibold">{{ result.organization }}</div>
          </div>

          <div v-if="result.country">
            <span class="text-xs text-slate-400 font-medium flex items-center gap-1">
              <Globe class="w-3.5 h-3.5" />
              <span>注册国家 (Country)</span>
            </span>
            <div class="font-mono text-slate-700 dark:text-slate-300 mt-1">{{ result.country }}</div>
          </div>
        </div>
      </div>

      <!-- Raw text toggle -->
      <div v-if="result.raw_text" class="custom-card p-4">
        <button
          type="button"
          @click="showRaw = !showRaw"
          class="text-xs font-semibold text-brand-600 dark:text-brand-400 hover:underline flex items-center gap-1.5"
        >
          <FileText class="w-4 h-4" />
          <span>{{ showRaw ? '折叠原始 WHOIS 输出' : '展开原始 WHOIS 纯文本输出' }}</span>
        </button>
        <pre
          v-if="showRaw"
          class="mt-4 p-4 rounded-lg bg-slate-900 text-slate-200 font-mono text-xs leading-relaxed overflow-x-auto max-h-96"
        >{{ result.raw_text }}</pre>
      </div>
    </div>
  </div>
</template>
