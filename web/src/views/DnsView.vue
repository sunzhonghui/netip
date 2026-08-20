<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Network, Search, ArrowRight } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import LatencyBadge from '@/components/LatencyBadge.vue'
import CopyButton from '@/components/CopyButton.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import { api } from '@/utils/api'
import type { DNSQueryResult } from '@/types'

const route = useRoute()
const router = useRouter()

const domain = ref('ipw.3x.cx')
const recordType = ref('A')
const loading = ref(false)
const error = ref('')
const result = ref<DNSQueryResult | null>(null)

const types = ['A', 'AAAA', 'CNAME', 'MX', 'TXT', 'NS', 'CAA', 'SRV', 'PTR', 'SOA']

async function executeQuery() {
  const target = domain.value.trim()
  if (!target) return

  loading.value = true
  error.value = ''

  try {
    result.value = await api.queryDNS(target, recordType.value)
    router.replace({
      query: { name: target, type: recordType.value },
    })
  } catch (e: any) {
    error.value = e.message || 'DNS 查询失败'
    result.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (route.query.name) {
    domain.value = route.query.name as string
  }
  if (route.query.type && types.includes((route.query.type as string).toUpperCase())) {
    recordType.value = (route.query.type as string).toUpperCase()
  }
  if (domain.value) {
    executeQuery()
  }
})
</script>

<template>
  <div class="space-y-6">
    <ToolHeader
      title="多 DNS 解析器查询"
      description="并发向系统 DNS、Cloudflare (1.1.1.1)、Google (8.8.8.8)、阿里公共 DNS (223.5.5.5)、DNSPod (119.29.29.29) 发起查询"
      :icon="Network"
    />

    <!-- Form -->
    <div class="custom-card p-5 sm:p-7 space-y-4 shadow-card">
      <form @submit.prevent="executeQuery" class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="domain"
            type="text"
            placeholder="输入域名，例如 ipw.3x.cx 或 8.8.8.8"
            aria-label="输入待查询域名"
            class="w-full pl-11 pr-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700/80 bg-slate-50/60 dark:bg-slate-950/60 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition-all shadow-inner"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-4 top-3.5" />
        </div>

        <select
          v-model="recordType"
          aria-label="选择 DNS 记录类型"
          class="px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700/80 bg-slate-50/60 dark:bg-slate-950/60 text-slate-900 dark:text-white font-mono text-sm font-semibold focus:outline-none focus:ring-2 focus:ring-brand-500/40"
        >
          <option v-for="t in types" :key="t" :value="t">{{ t }} 记录</option>
        </select>

        <button
          type="submit"
          :disabled="loading"
          class="btn-primary"
        >
          <span>{{ loading ? '查询中...' : '查询' }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </form>

      <!-- Record type quick pills -->
      <div class="flex flex-wrap items-center gap-2 border-t border-slate-100 dark:border-slate-800/80 pt-3">
        <span class="text-xs text-slate-400 dark:text-slate-500 mr-1 font-medium">常用类型:</span>
        <button
          v-for="t in types"
          :key="t"
          type="button"
          @click="recordType = t; executeQuery()"
          class="px-3 py-1 rounded-lg text-xs font-mono font-bold transition-all border"
          :class="
            recordType === t
              ? 'bg-brand-600 text-white border-brand-700 shadow-sm'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100 dark:bg-slate-800 dark:text-slate-300 dark:border-slate-700 hover:border-brand-500/40'
          "
        >
          {{ t }}
        </button>
      </div>
    </div>

    <!-- Error Alert -->
    <ErrorAlert v-if="error" :message="error" />

    <!-- Result Table -->
    <div v-if="result" class="custom-card overflow-hidden">
      <div class="px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex items-center justify-between">
        <div class="flex items-center gap-2 font-mono text-sm font-bold text-slate-800 dark:text-slate-200">
          <span>{{ result.name }}</span>
          <span class="text-slate-400">·</span>
          <span class="text-brand-600 dark:text-brand-400">{{ result.type }} 记录查询结果</span>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="text-xs text-slate-500 dark:text-slate-400 uppercase bg-slate-50/75 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-800 font-mono">
            <tr>
              <th scope="col" class="px-6 py-3">解析器 / 节点</th>
              <th scope="col" class="px-6 py-3">响应耗时</th>
              <th scope="col" class="px-6 py-3">TTL</th>
              <th scope="col" class="px-6 py-3">解析值</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-800 font-mono text-xs sm:text-sm">
            <tr
              v-for="(r, idx) in result.results"
              :key="idx"
              class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
            >
              <td class="px-6 py-4 font-semibold text-slate-800 dark:text-slate-200 whitespace-nowrap">
                <div>{{ r.node_name || r.resolver }}</div>
                <div v-if="r.isp" class="text-[11px] font-normal text-slate-400">{{ r.isp }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <LatencyBadge v-if="!r.error" :latency-ms="r.latency_ms" />
                <span v-else class="text-xs text-rose-500">超时 / 错误</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-slate-500 dark:text-slate-400">
                <span v-if="r.answers && r.answers.length > 0">{{ r.answers[0].ttl }}s</span>
                <span v-else>-</span>
              </td>
              <td class="px-6 py-4">
                <div v-if="r.answers && r.answers.length > 0" class="space-y-1.5">
                  <div
                    v-for="(ans, aIdx) in r.answers"
                    :key="aIdx"
                    class="flex items-center justify-between gap-4 font-mono font-medium text-slate-900 dark:text-slate-100"
                  >
                    <span>{{ ans.value }}</span>
                    <CopyButton :text="ans.value" :icon-only="true" />
                  </div>
                </div>
                <div v-else-if="r.error" class="text-xs text-rose-500 font-sans">
                  {{ r.error }}
                </div>
                <div v-else class="text-xs text-slate-400 font-sans">
                  无此类记录 (NXDOMAIN / Empty)
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
