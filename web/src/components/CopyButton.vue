<script setup lang="ts">
import { ref } from 'vue'
import { Copy, Check } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    text: string
    label?: string
    iconOnly?: boolean
  }>(),
  {
    label: '复制',
    iconOnly: false,
  }
)

const copied = ref(false)

async function handleCopy() {
  if (!props.text) return
  try {
    await navigator.clipboard.writeText(props.text)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch {
    // Fallback
    const input = document.createElement('textarea')
    input.value = props.text
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  }
}
</script>

<template>
  <button
    type="button"
    @click.stop="handleCopy"
    :title="copied ? '已复制到剪贴板' : '点击复制'"
    :aria-label="copied ? '已复制' : '复制'"
    class="inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-lg transition-all border select-none shadow-2xs active:scale-95"
    :class="
      copied
        ? 'bg-emerald-50 text-emerald-700 border-emerald-300 dark:bg-emerald-950/60 dark:text-emerald-300 dark:border-emerald-700'
        : 'bg-white text-slate-700 border-slate-200/90 hover:bg-slate-50 hover:text-brand-600 hover:border-brand-300 dark:bg-slate-800/90 dark:text-slate-300 dark:border-slate-700/80 dark:hover:bg-slate-700/90 dark:hover:text-brand-400 dark:hover:border-brand-600'
    "
  >
    <Check v-if="copied" class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400" />
    <Copy v-else class="w-3.5 h-3.5 opacity-70" />
    <span v-if="!iconOnly">{{ copied ? '已复制' : label }}</span>
  </button>
</template>
