<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Activity, Search, ArrowRight, Radio } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import LatencyBadge from '@/components/LatencyBadge.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import { api } from '@/utils/api'
import type { MultiNodePingResponse, MultiNodeTCPingResponse } from '@/types'

const route = useRoute()
const router = useRouter()

const mode = ref<'icmp' | 'tcp'>('icmp')
const target = ref('1.1.1.1')
const port = ref(443)
const count = ref(4)
const loading = ref(false)
const error = ref('')

const icmpResult = ref<MultiNodePingResponse | null>(null)
const tcpResult = ref<MultiNodeTCPingResponse | null>(null)

async function executeTest() {
  const host = target.value.trim()
  if (!host) return

  loading.value = true
  error.value = ''

  try {
    if (mode.value === 'icmp') {
      icmpResult.value = await api.ping(host, count.value)
      tcpResult.value = null
    } else {
      tcpResult.value = await api.tcping(host, port.value, count.value)
      icmpResult.value = null
    }
    router.replace({
      query: {
        mode: mode.value,
        target: host,
        port: mode.value === 'tcp' ? port.value : undefined,
      },
    })
  } catch (e: any) {
    error.value = e.message || '测试失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (route.query.target) {
    target.value = route.query.target as string
  }
  if (route.query.mode === 'tcp' || route.query.mode === 'icmp') {
    mode.value = route.query.mode
  }
  if (route.query.port) {
    port.value = parseInt(route.query.port as string) || 443
  }
  if (target.value) {
    executeTest()
  }
})
</script>

<template>
  <div class="space-y-6">
    <ToolHeader
      title="Ping / TCPing 网络连通性测试"
      description="支持原生 ICMP Echo 丢包率与往返时延测试，以及基于 TCP 端口三次握手建立速度的 TCPing 诊断"
      :icon="Activity"
    />

    <!-- Control Form -->
    <div class="custom-card p-5 sm:p-7 space-y-4 shadow-card">
      <!-- Mode toggle -->
      <div class="flex items-center gap-4 border-b border-slate-100 dark:border-slate-800/80 pb-3">
        <label class="flex items-center gap-2 text-xs sm:text-sm font-bold cursor-pointer select-none px-3 py-1.5 rounded-lg transition-all" :class="mode === 'icmp' ? 'bg-brand-50 text-brand-700 dark:bg-brand-950/60 dark:text-brand-300' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'">
          <input
            v-model="mode"
            type="radio"
            value="icmp"
            class="text-brand-600 focus:ring-brand-500"
          />
          <span>ICMP Ping (网络层)</span>
        </label>

        <label class="flex items-center gap-2 text-xs sm:text-sm font-bold cursor-pointer select-none px-3 py-1.5 rounded-lg transition-all" :class="mode === 'tcp' ? 'bg-brand-50 text-brand-700 dark:bg-brand-950/60 dark:text-brand-300' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'">
          <input
            v-model="mode"
            type="radio"
            value="tcp"
            class="text-brand-600 focus:ring-brand-500"
          />
          <span>TCPing (传输层端口)</span>
        </label>
      </div>

      <form @submit.prevent="executeTest" class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="target"
            type="text"
            placeholder="输入域名或 IP，例如 1.1.1.1 或 ipw.3x.cx"
            aria-label="输入目标地址"
            class="w-full pl-11 pr-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700/80 bg-slate-50/60 dark:bg-slate-950/60 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition-all shadow-inner"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-4 top-3.5" />
        </div>

        <div v-if="mode === 'tcp'" class="w-full sm:w-28">
          <input
            v-model.number="port"
            type="number"
            placeholder="端口"
            aria-label="目标端口"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700/80 bg-slate-50/60 dark:bg-slate-950/60 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition-all"
          />
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="btn-primary"
        >
          <span>{{ loading ? '测试中...' : '发起测试' }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </form>
    </div>

    <!-- Error Alert -->
    <ErrorAlert v-if="error" :message="error" />

    <!-- ICMP Results -->
    <div v-if="icmpResult" class="space-y-4">
      <div class="custom-card overflow-hidden">
        <div class="px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-2 font-mono text-sm font-bold text-slate-800 dark:text-slate-200">
            <span>{{ icmpResult.target }}</span>
            <span v-if="icmpResult.resolved_ip" class="text-slate-400 font-normal">({{ icmpResult.resolved_ip }})</span>
          </div>
          <span class="text-xs font-mono font-medium px-2.5 py-1 rounded bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300">
            ICMP Echo 协议
          </span>
        </div>

        <div class="divide-y divide-slate-100 dark:divide-slate-800">
          <div
            v-for="(node, idx) in icmpResult.nodes"
            :key="idx"
            class="p-6 space-y-4"
          >
            <div class="flex flex-wrap items-center justify-between gap-4">
              <div class="font-bold text-base text-slate-900 dark:text-white flex items-center gap-2">
                <Radio class="w-4 h-4 text-brand-600" />
                <span>{{ node.node }}</span>
              </div>

              <!-- Latency summary badges -->
              <div v-if="!node.error && node.received > 0" class="flex items-center gap-3">
                <div class="text-xs font-mono text-slate-500 dark:text-slate-400">
                  丢包: <span :class="node.loss_percent === 0 ? 'text-emerald-600' : 'text-rose-500 font-bold'">{{ node.loss_percent }}%</span>
                </div>
                <LatencyBadge :latency-ms="node.avg_ms" :show-label="true" />
              </div>
            </div>

            <div v-if="node.error" class="text-sm text-rose-500">
              {{ node.error }}
            </div>

            <!-- Stats Grid -->
            <div v-else class="grid grid-cols-2 sm:grid-cols-4 gap-4 bg-slate-50 dark:bg-slate-800/50 p-4 rounded-xl text-center font-mono text-xs sm:text-sm">
              <div>
                <div class="text-xs text-slate-400 dark:text-slate-500 font-sans">发 / 收</div>
                <div class="font-bold text-slate-800 dark:text-slate-200 mt-0.5">{{ node.sent }} / {{ node.received }}</div>
              </div>
              <div>
                <div class="text-xs text-slate-400 dark:text-slate-500 font-sans">最短耗时</div>
                <div class="font-bold text-emerald-600 dark:text-emerald-400 mt-0.5">{{ node.min_ms }} ms</div>
              </div>
              <div>
                <div class="text-xs text-slate-400 dark:text-slate-500 font-sans">平均耗时</div>
                <div class="font-bold text-blue-600 dark:text-blue-400 mt-0.5">{{ node.avg_ms }} ms</div>
              </div>
              <div>
                <div class="text-xs text-slate-400 dark:text-slate-500 font-sans">最长耗时</div>
                <div class="font-bold text-amber-600 dark:text-amber-400 mt-0.5">{{ node.max_ms }} ms</div>
              </div>
            </div>

            <!-- Samples stream -->
            <div v-if="node.samples && node.samples.length > 0" class="flex items-center gap-2 text-xs font-mono text-slate-500 dark:text-slate-400">
              <span>采样明细:</span>
              <span
                v-for="(s, sIdx) in node.samples"
                :key="sIdx"
                class="px-2 py-0.5 rounded bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700"
              >
                {{ s }}ms
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TCPing Results -->
    <div v-if="tcpResult" class="space-y-4">
      <div class="custom-card overflow-hidden">
        <div class="px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-2 font-mono text-sm font-bold text-slate-800 dark:text-slate-200">
            <span>{{ tcpResult.target }}:{{ tcpResult.port }}</span>
            <span v-if="tcpResult.resolved_ip" class="text-slate-400 font-normal">({{ tcpResult.resolved_ip }})</span>
          </div>
          <span class="text-xs font-mono font-medium px-2.5 py-1 rounded bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300">
            TCP 握手协议
          </span>
        </div>

        <div class="divide-y divide-slate-100 dark:divide-slate-800">
          <div
            v-for="(node, idx) in tcpResult.nodes"
            :key="idx"
            class="p-6 space-y-4"
          >
            <div class="flex flex-wrap items-center justify-between gap-4">
              <div class="font-bold text-base text-slate-900 dark:text-white flex items-center gap-2">
                <Radio class="w-4 h-4 text-brand-600" />
                <span>{{ node.node }}</span>
              </div>

              <div v-if="node.success > 0" class="flex items-center gap-3">
                <div class="text-xs font-mono text-slate-500 dark:text-slate-400">
                  成功: <span class="text-emerald-600 font-bold">{{ node.success }}</span> / 失败: <span :class="node.failed > 0 ? 'text-rose-500' : 'text-slate-400'">{{ node.failed }}</span>
                </div>
                <LatencyBadge :latency-ms="node.avg_ms" :show-label="true" />
              </div>
            </div>

            <!-- Stats Grid -->
            <div v-if="node.success > 0" class="grid grid-cols-3 gap-4 bg-slate-50 dark:bg-slate-800/50 p-4 rounded-xl text-center font-mono text-xs sm:text-sm">
              <div>
                <div class="text-xs text-slate-400 dark:text-slate-500 font-sans">最短握手</div>
                <div class="font-bold text-emerald-600 dark:text-emerald-400 mt-0.5">{{ node.min_ms }} ms</div>
              </div>
              <div>
                <div class="text-xs text-slate-400 dark:text-slate-500 font-sans">平均握手</div>
                <div class="font-bold text-blue-600 dark:text-blue-400 mt-0.5">{{ node.avg_ms }} ms</div>
              </div>
              <div>
                <div class="text-xs text-slate-400 dark:text-slate-500 font-sans">最长握手</div>
                <div class="font-bold text-amber-600 dark:text-amber-400 mt-0.5">{{ node.max_ms }} ms</div>
              </div>
            </div>

            <!-- Sample Attempts -->
            <div v-if="node.samples" class="space-y-1.5 font-mono text-xs">
              <div
                v-for="(s, sIdx) in node.samples"
                :key="sIdx"
                class="flex items-center justify-between p-2 rounded bg-slate-50 dark:bg-slate-800/60"
              >
                <span>第 {{ sIdx + 1 }} 次连接</span>
                <span v-if="s.success" class="text-emerald-600 dark:text-emerald-400 font-bold">{{ s.latency_ms }} ms</span>
                <span v-else class="text-rose-500 font-sans">{{ s.error || '连接失败' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
