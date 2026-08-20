<script setup lang="ts">
import type { Component } from 'vue'
import { ChevronRight } from 'lucide-vue-next'

defineProps<{
  title: string
  description: string
  icon?: Component
}>()
</script>

<template>
  <div class="mb-8 space-y-3">
    <!-- Breadcrumb -->
    <div class="flex items-center gap-1.5 text-xs text-slate-400 dark:text-slate-500 font-medium">
      <router-link to="/" class="hover:text-brand-600 dark:hover:text-brand-400 transition-colors">
        首页
      </router-link>
      <ChevronRight class="w-3.5 h-3.5 text-slate-300 dark:text-slate-600" />
      <span class="text-slate-600 dark:text-slate-300">{{ title }}</span>
    </div>

    <!-- Main Header Box -->
    <div class="flex items-start gap-4">
      <div
        v-if="icon"
        class="w-12 h-12 sm:w-14 sm:h-14 rounded-2xl bg-gradient-to-br from-brand-500/15 via-brand-500/5 to-transparent text-brand-600 dark:text-brand-400 border border-brand-500/20 dark:border-brand-500/30 flex items-center justify-center flex-shrink-0 shadow-card"
      >
        <component :is="icon" class="w-6 h-6 sm:w-7 sm:h-7" />
      </div>
      <div>
        <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white">
          {{ title }}
        </h1>
        <p class="mt-1 text-xs sm:text-sm text-slate-500 dark:text-slate-400 font-normal leading-relaxed max-w-3xl">
          {{ description }}
        </p>
      </div>
    </div>
  </div>
</template>
