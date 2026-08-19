<script setup lang="ts">
import { AlertCircle } from 'lucide-vue-next'

defineProps<{
  title?: string
  message: string
}>()
</script>

<template>
  <div class="rounded-xl border border-rose-200 bg-rose-50/70 p-4 text-rose-800 dark:border-rose-900/60 dark:bg-rose-950/40 dark:text-rose-300">
    <div class="flex items-start gap-3">
      <AlertCircle class="w-5 h-5 text-rose-600 dark:text-rose-400 mt-0.5 flex-shrink-0" />
      <div>
        <h4 v-if="title" class="text-sm font-semibold mb-0.5">{{ title }}</h4>
        <p class="text-sm leading-relaxed">{{ message }}</p>
      </div>
    </div>
  </div>
</template>
