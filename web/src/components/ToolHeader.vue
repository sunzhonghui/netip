<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  title: string
  description: string
  icon?: Component
}>()
</script>

<template>
  <div class="mb-6 flex items-start gap-4">
    <div
      v-if="icon"
      class="w-12 h-12 rounded-xl bg-brand-50 text-brand-600 dark:bg-brand-950/60 dark:text-brand-400 border border-brand-200/60 dark:border-brand-800/60 flex items-center justify-center flex-shrink-0"
    >
      <component :is="icon" class="w-6 h-6" />
    </div>
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">{{ title }}</h1>
      <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ description }}</p>
    </div>
  </div>
</template>
