<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Globe2, Search, ArrowRight, ShieldCheck, Layers, MapPin, Database, Network } from 'lucide-vue-next'
import ToolHeader from '@/components/ToolHeader.vue'
import CopyButton from '@/components/CopyButton.vue'
import ErrorAlert from '@/components/ErrorAlert.vue'
import { api } from '@/utils/api'
import type { IPDetails } from '@/types'

const route = useRoute()
const router = useRouter()

const query = ref('')
const loading = ref(false)
const error = ref('')
const result = ref<IPDetails | null>(null)

async function executeQuery() {
  const target = query.value.trim()
  error.value = ''

  if (!target) {
    // Lookup current IP
    loading.value = true
    try {
      result.value = await api.getMe()
      query.value = result.value.ip
    } catch (e: any) {
      error.value = e.message || '获取当前 IP 失败'
    } finally {
      loading.value = false
    }
    return
  }

  loading.value = true
  try {
    result.value = await api.getIP(target)
    if (route.params.address !== target) {
      router.push(`/ip/${encodeURIComponent(target)}`)
    }
  } catch (e: any) {
    error.value = e.message || 'IP 查询失败'
    result.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  const addr = route.params.address as string
  if (addr) {
    query.value = addr
  }
  executeQuery()
})

watch(
  () => route.params.address,
  (newAddr) => {
    if (newAddr && newAddr !== query.value) {
      query.value = newAddr as string
      executeQuery()
    }
  }
)
</script>

<template>
  <div class="space-y-8">
    <ToolHeader
      title="IP 地址与归属地查询"
      description="精准查询 IPv4 / IPv6 公网地址的地理位置、运营商、所属 ASN 自治系统及广播网段"
      :icon="Globe2"
    />

    <!-- Input Form Card -->
    <div class="custom-card p-5 sm:p-7 shadow-card">
      <form @submit.prevent="executeQuery" class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="query"
            type="text"
            placeholder="输入 IPv4 或 IPv6 地址（留空查询本机公网 IP）"
            aria-label="输入 IP 地址"
            class="w-full pl-11 pr-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700/80 bg-slate-50/60 dark:bg-slate-950/60 text-slate-900 dark:text-white font-mono text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/40 focus:border-brand-500 transition-all shadow-inner"
          />
          <Search class="w-4 h-4 text-slate-400 absolute left-4 top-3.5" />
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="btn-primary"
        >
          <span>{{ loading ? '查询中...' : '查询' }}</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </form>
    </div>

    <!-- Error Alert -->
    <ErrorAlert v-if="error" :message="error" />

    <!-- Result Card -->
    <div v-if="result" class="custom-card p-6 sm:p-8 space-y-6 shadow-card">
      <div class="flex flex-wrap items-center justify-between gap-4 border-b border-slate-100 dark:border-slate-800 pb-5">
        <div class="flex items-center gap-3">
          <span
            class="px-2.5 py-1 rounded-lg text-xs font-mono font-bold uppercase border shadow-2xs"
            :class="
              result.version === 6
                ? 'bg-purple-50 text-purple-700 border-purple-200 dark:bg-purple-950/60 dark:text-purple-300 dark:border-purple-800'
                : 'bg-brand-50 text-brand-700 border-brand-200 dark:bg-brand-950/60 dark:text-brand-300 dark:border-brand-800'
            "
          >
            IPv{{ result.version }}
          </span>
          <span class="font-mono text-xl sm:text-3xl font-extrabold text-slate-900 dark:text-white select-all tracking-tight">
            {{ result.ip }}
          </span>
        </div>
        <CopyButton :text="result.ip" label="复制 IP" />
      </div>

      <!-- Info Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 pt-2">
        <!-- Geolocation -->
        <div class="space-y-1.5 p-4 rounded-xl bg-slate-50/60 dark:bg-slate-950/40 border border-slate-100 dark:border-slate-800/80">
          <div class="text-xs text-slate-400 dark:text-slate-500 font-semibold flex items-center gap-1.5">
            <MapPin class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
            <span>地理位置</span>
          </div>
          <div class="text-base font-bold text-slate-900 dark:text-slate-100">
            {{ [result.country, result.province, result.city].filter(Boolean).join(' · ') || '暂无归属地数据' }}
          </div>
        </div>

        <!-- ISP -->
        <div class="space-y-1.5 p-4 rounded-xl bg-slate-50/60 dark:bg-slate-950/40 border border-slate-100 dark:border-slate-800/80">
          <div class="text-xs text-slate-400 dark:text-slate-500 font-semibold flex items-center gap-1.5">
            <ShieldCheck class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
            <span>网络运营商 (ISP)</span>
          </div>
          <div class="text-base font-bold text-slate-900 dark:text-slate-100">
            {{ result.isp || '未知' }}
          </div>
        </div>

        <!-- ASN -->
        <div class="space-y-1.5 p-4 rounded-xl bg-slate-50/60 dark:bg-slate-950/40 border border-slate-100 dark:border-slate-800/80">
          <div class="text-xs text-slate-400 dark:text-slate-500 font-semibold flex items-center gap-1.5">
            <Layers class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
            <span>自治系统 (ASN)</span>
          </div>
          <div class="text-base font-bold text-slate-900 dark:text-slate-100">
            <router-link
              v-if="result.asn"
              :to="`/asn/${result.asn}`"
              class="text-brand-600 dark:text-brand-400 hover:underline font-mono inline-flex items-center gap-1"
            >
              <span>AS{{ result.asn }}</span>
              <span v-if="result.as_name" class="text-sm font-sans font-medium text-slate-600 dark:text-slate-300">({{ result.as_name }})</span>
            </router-link>
            <span v-else class="text-slate-400">未知</span>
          </div>
        </div>

        <!-- Network CIDR -->
        <div v-if="result.network" class="space-y-1.5 p-4 rounded-xl bg-slate-50/60 dark:bg-slate-950/40 border border-slate-100 dark:border-slate-800/80">
          <div class="text-xs text-slate-400 dark:text-slate-500 font-semibold flex items-center gap-1.5">
            <Network class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
            <span>广播网段 (CIDR)</span>
          </div>
          <div class="font-mono text-sm font-semibold text-slate-800 dark:text-slate-200">
            {{ result.network }}
          </div>
        </div>

        <!-- Sources -->
        <div v-if="result.sources" class="space-y-1.5 p-4 rounded-xl bg-slate-50/60 dark:bg-slate-950/40 border border-slate-100 dark:border-slate-800/80">
          <div class="text-xs text-slate-400 dark:text-slate-500 font-semibold flex items-center gap-1.5">
            <Database class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
            <span>解析数据源</span>
          </div>
          <div class="text-xs font-mono font-medium text-slate-600 dark:text-slate-400">
            {{ Object.entries(result.sources).map(([k, v]) => `${k}:${v}`).join(' · ') }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
