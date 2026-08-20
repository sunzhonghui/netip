<script setup lang="ts">
import { computed } from 'vue'
import { XCircle, Globe, ShieldCheck, Layers, ArrowRight } from 'lucide-vue-next'
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
    class="custom-card p-6 sm:p-7 relative overflow-hidden transition-all duration-300 group"
    :class="[
      available
        ? 'border-slate-200/90 dark:border-slate-800 shadow-card hover:shadow-card-hover'
        : 'border-slate-200/60 dark:border-slate-800/60 bg-white/60 dark:bg-slate-900/60'
    ]"
  >
    <!-- Top Gradient Accent Line -->
    <div
      class="absolute top-0 left-0 right-0 h-1 transition-all duration-300"
      :class="[
        available
          ? version === 'IPv6'
            ? 'bg-gradient-to-r from-purple-500 to-indigo-500'
            : 'bg-gradient-to-r from-brand-600 via-brand-500 to-sky-400'
          : 'bg-slate-200 dark:bg-slate-800'
      ]"
    ></div>

    <!-- Header / Status Badge -->
    <div class="flex items-center justify-between mb-5">
      <div class="flex items-center gap-2.5">
        <span
          class="px-2.5 py-1 rounded-lg text-xs font-mono font-bold tracking-wider uppercase border shadow-2xs"
          :class="[
            version === 'IPv6'
              ? 'bg-purple-50 text-purple-700 border-purple-200/80 dark:bg-purple-950/60 dark:text-purple-300 dark:border-purple-800/80'
              : 'bg-brand-50 text-brand-700 border-brand-200/80 dark:bg-brand-950/60 dark:text-brand-300 dark:border-brand-800/80'
          ]"
        >
          {{ version }}
        </span>

        <!-- Live Pulse Radar Status -->
        <span
          v-if="available"
          class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-300 border border-emerald-200/80 dark:border-emerald-800/80"
        >
          <span class="relative flex h-2 w-2">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
          <span>已连接</span>
        </span>

        <span
          v-else-if="!loading"
          class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium text-slate-400 dark:text-slate-500 bg-slate-100 dark:bg-slate-800/80 border border-slate-200/60 dark:border-slate-800"
        >
          <XCircle class="w-3.5 h-3.5 text-slate-400" />
          <span>未检测到</span>
        </span>
      </div>

      <CopyButton v-if="ip && available" :text="ip" label="复制" />
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="animate-pulse space-y-4 py-2">
      <div class="h-8 bg-slate-200/80 dark:bg-slate-800 rounded-xl w-3/4"></div>
      <div class="space-y-2 pt-2">
        <div class="h-4 bg-slate-200/70 dark:bg-slate-800/70 rounded-md w-1/2"></div>
        <div class="h-4 bg-slate-200/70 dark:bg-slate-800/70 rounded-md w-2/3"></div>
      </div>
    </div>

    <!-- Available State -->
    <div v-else-if="available && ip" class="space-y-5">
      <div>
        <div class="font-mono text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white break-all select-all leading-tight">
          {{ ip }}
        </div>
      </div>

      <!-- Info Meta Rows -->
      <div class="space-y-2.5 text-sm border-t border-slate-100 dark:border-slate-800/80 pt-4">
        <!-- Location -->
        <div v-if="locationText" class="flex items-center gap-2.5 text-slate-600 dark:text-slate-300">
          <div class="w-6 h-6 rounded-md bg-slate-100 dark:bg-slate-800 flex items-center justify-center flex-shrink-0 text-slate-500 dark:text-slate-400">
            <Globe class="w-3.5 h-3.5" />
          </div>
          <span class="font-medium text-slate-800 dark:text-slate-200">{{ locationText }}</span>
        </div>

        <!-- ISP -->
        <div v-if="details?.isp" class="flex items-center gap-2.5 text-slate-700 dark:text-slate-200">
          <div class="w-6 h-6 rounded-md bg-brand-50 dark:bg-brand-950/60 flex items-center justify-center flex-shrink-0 text-brand-600 dark:text-brand-400">
            <ShieldCheck class="w-3.5 h-3.5" />
          </div>
          <span class="font-semibold text-slate-900 dark:text-slate-100">{{ details.isp }}</span>
        </div>

        <!-- ASN -->
        <div v-if="details?.asn" class="flex items-center gap-2.5 text-xs text-slate-500 dark:text-slate-400 pt-0.5">
          <div class="w-6 h-6 rounded-md bg-slate-100 dark:bg-slate-800 flex items-center justify-center flex-shrink-0 text-slate-500 dark:text-slate-400">
            <Layers class="w-3.5 h-3.5" />
          </div>
          <div class="flex items-center gap-1.5 font-mono">
            <router-link
              :to="`/asn/${details.asn}`"
              class="text-brand-600 dark:text-brand-400 hover:underline font-semibold flex items-center gap-1"
            >
              <span>AS{{ details.asn }}</span>
              <span v-if="details.as_name">({{ details.as_name }})</span>
              <ArrowRight class="w-3 h-3" />
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Unavailable State -->
    <div v-else class="py-6 text-center space-y-2">
      <div class="w-12 h-12 mx-auto rounded-2xl bg-slate-100 dark:bg-slate-800/80 flex items-center justify-center text-slate-400">
        <XCircle class="w-6 h-6" />
      </div>
      <p class="text-sm font-medium text-slate-600 dark:text-slate-400">
        当前网络环境未检测到 {{ version }} 地址
      </p>
      <p class="text-xs text-slate-400 dark:text-slate-500 max-w-xs mx-auto">
        可能是您的本地宽带未开启 {{ version }}，或当前路由器未配置分配。
      </p>
    </div>
  </div>
</template>
