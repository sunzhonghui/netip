<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Layers, Search, ArrowRight, Building, Globe, Database } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import { api } from '@/utils/api'
import type { ASNResult } from '@/types'

const route = useRoute()
const router = useRouter()

const query = ref('AS4134')
const loading = ref(false)
const error = ref('')
const result = ref<ASNResult | null>(null)

async function executeQuery() {
  const target = query.value.trim()
  if (!target) return

  loading.value = true
  error.value = ''

  try {
    result.value = await api.getASN(target)
    if (route.params.query !== target) {
      router.push(`/asn/${encodeURIComponent(target)}`)
    }
  } catch (e: any) {
    error.value = e.message || 'ASN 查询失败'
    result.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  const q = route.params.query as string
  if (q) {
    query.value = q
  }
  if (query.value) {
    executeQuery()
  }
})

watch(
  () => route.params.query,
  (newQ) => {
    if (newQ && newQ !== query.value) {
      query.value = newQ as string
      executeQuery()
    }
  }
)
</script>

<template>
  <div class="space-y-6">
    <ToolHeader
      title="ASN 自治系统查询"
      description="查询 Autonomous System 自治系统编号、BGP 广播运营商名称、注册局与所属国家"
      :icon="Layers"
    />

    <!-- Input Form -->
    <div class="custom-card p-5 sm:p-7 shadow-card">
      <form @submit.prevent="executeQuery" class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="query"
            type="text"
            placeholder="输入 AS 号或 IP，例如 AS4134 或 1.1.1.1"
            aria-label="输入 ASN 或 IP"
            class="w-full pl-11 pr-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700/80 bg-slate-50/60 dark:bg-slate-950/60 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition-all shadow-inner"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-4 top-3.5" />
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="btn-primary"
        >
          <span>{{ loading ? '查询中...' : '查询 ASN' }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </form>
    </div>

    <!-- Error Alert -->
    <ErrorAlert v-if="error" :message="error" />

    <!-- Result View -->
    <div v-if="result" class="custom-card p-6 space-y-6">
      <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-4">
        <div>
          <div class="font-mono text-2xl font-bold text-slate-900 dark:text-white">
            AS{{ result.asn }}
          </div>
          <div class="text-sm font-medium text-brand-600 dark:text-brand-400 mt-1">
            {{ result.as_name || '未知机构' }}
          </div>
        </div>
        <span class="text-xs font-mono font-medium px-2.5 py-1 rounded bg-orange-50 text-orange-700 dark:bg-orange-950/60 dark:text-orange-300">
          自治系统
        </span>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 text-sm">
        <div v-if="result.country" class="space-y-1">
          <div class="text-xs text-slate-400 font-medium flex items-center gap-1.5">
            <Globe class="w-3.5 h-3.5" />
            <span>注册国家/地区</span>
          </div>
          <div class="font-semibold text-slate-800 dark:text-slate-200">{{ result.country }}</div>
        </div>

        <div v-if="result.registry" class="space-y-1">
          <div class="text-xs text-slate-400 font-medium flex items-center gap-1.5">
            <Building class="w-3.5 h-3.5" />
            <span>区域互联网注册局 (RIR)</span>
          </div>
          <div class="font-mono font-semibold text-slate-800 dark:text-slate-200">{{ result.registry }}</div>
        </div>

        <div v-if="result.network" class="space-y-1">
          <div class="text-xs text-slate-400 font-medium">
            广播网段
          </div>
          <div class="font-mono text-slate-800 dark:text-slate-200">{{ result.network }}</div>
        </div>

        <div v-if="result.source" class="space-y-1">
          <div class="text-xs text-slate-400 font-medium flex items-center gap-1.5">
            <Database class="w-3.5 h-3.5" />
            <span>查询来源</span>
          </div>
          <div class="font-mono text-xs text-slate-500">{{ result.source }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
