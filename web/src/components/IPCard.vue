<script setup lang="ts">
import { computed } from 'vue'
import { CheckCircle2, XCircle, Globe, ShieldCheck } from 'lucide-vue-next'
import CopyButton from './CopyButton.vue'
import type { IPDetails } from '@/types'

const props = defineProps<{
  version: 'IPv4' | 'IPv6'
  ip?: string
  details?: IPDetails
  available: boolean
  loading?: boolean
}>()

const locationText = computed(() => {
  if (!props.details) return ''
  const parts = [
    props.details.country,
    props.details.province,
    props.details.city,
  ].filter(Boolean)
  return parts.join(' · ')
})
</script>

<template>
  <div
    class="custom-card p-6 relative overflow-hidden transition-all"
    :class="
      available
        ? 'border-brand-200/80 dark:border-brand-900/60 shadow-sm'
        : 'border-slate-200 dark:border-slate-800 opacity-90'
    "
  >
    <!-- Background accent indicator -->
    <div
      class="absolute top-0 left-0 right-0 h-1.5"
      :class="available ? 'bg-brand-500' : 'bg-slate-300 dark:bg-slate-700'"
    ></div>

    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-2">
        <span
          class="px-2.5 py-0.5 rounded-md text-xs font-bold font-mono tracking-wide uppercase border"
          :class="
            version === 'IPv6'
              ? 'bg-purple-50 text-purple-700 border-purple-200 dark:bg-purple-950/60 dark:text-purple-300 dark:border-purple-800'
              : 'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950/60 dark:text-blue-300 dark:border-blue-800'
          "
        >
          {{ version }}
        </span>
        <span
          v-if="available"
          class="inline-flex items-center gap-1 text-xs font-medium text-emerald-600 dark:text-emerald-400"
        >
          <CheckCircle2 class="w-3.5 h-3.5" />
          <span>支持</span>
        </span>
        <span
          v-else-if="!loading"
          class="inline-flex items-center gap-1 text-xs font-medium text-slate-400 dark:text-slate-500"
        >
          <XCircle class="w-3.5 h-3.5" />
          <span>未检测到</span>
        </span>
      </div>

      <CopyButton v-if="ip && available" :text="ip" />
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="animate-pulse space-y-3 py-2">
      <div class="h-7 bg-slate-200 dark:bg-slate-800 rounded w-3/4"></div>
      <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-1/2"></div>
      <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-2/3"></div>
    </div>

    <!-- Available State -->
    <div v-else-if="available && ip" class="space-y-4">
      <div>
        <div class="font-mono text-xl md:text-2xl font-bold tracking-tight text-slate-900 dark:text-white break-all select-all">
          {{ ip }}
        </div>
      </div>

      <div class="space-y-2 text-sm text-slate-600 dark:text-slate-300 border-t border-slate-100 dark:border-slate-800/80 pt-3">
        <div v-if="locationText" class="flex items-center gap-2">
          <Globe class="w-4 h-4 text-slate-400 flex-shrink-0" />
          <span>{{ locationText }}</span>
        </div>

        <div v-if="details?.isp" class="flex items-center gap-2 font-medium text-slate-800 dark:text-slate-200">
          <ShieldCheck class="w-4 h-4 text-brand-500 flex-shrink-0" />
          <span>{{ details.isp }}</span>
        </div>

        <div v-if="details?.asn" class="flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
          <router-link
            :to="`/asn/${details.asn}`"
            class="font-mono hover:text-brand-600 dark:hover:text-brand-400 hover:underline"
          >
            AS{{ details.asn }}
          </router-link>
          <span v-if="details.as_name" class="truncate">· {{ details.as_name }}</span>
        </div>
      </div>

      <div class="pt-1 flex items-center gap-1.5 text-xs text-emerald-700 dark:text-emerald-400 font-medium">
        <CheckCircle2 class="w-4 h-4 text-emerald-500" />
        <span>✓ {{ version }} 可用</span>
      </div>
    </div>

    <!-- Unavailable State -->
    <div v-else class="space-y-4 py-2">
      <div class="text-lg font-medium text-slate-400 dark:text-slate-500">
        未检测到 {{ version }}
      </div>
      <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed border-t border-slate-100 dark:border-slate-800/80 pt-3">
        ✕ 当前网络可能未分配 {{ version }} 公网地址或路由器未开启对应协议支持。
      </p>
    </div>
  </div>
</template>
