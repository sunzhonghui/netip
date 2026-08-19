<script setup lang="ts">
import { computed } from 'vue'
import { getLatencyLevel } from '@/utils/format'

const props = defineProps<{
  latencyMs: number
  showLabel?: boolean
}>()

const level = computed(() => getLatencyLevel(props.latencyMs))
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-mono font-medium border"
    :class="[level.bg, level.text, level.border]"
  >
    <span
      class="w-1.5 h-1.5 rounded-full"
      :class="{
        'bg-emerald-500': level.color === 'emerald',
        'bg-blue-500': level.color === 'blue',
        'bg-amber-500': level.color === 'amber',
        'bg-rose-500': level.color === 'rose',
        'bg-gray-400': level.color === 'gray',
      }"
    ></span>
    <span>{{ latencyMs }} ms</span>
    <span v-if="showLabel" class="text-[10px] opacity-80">({{ level.label }})</span>
  </span>
</template>
