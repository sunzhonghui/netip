<script setup lang="ts">
import { ShieldCheck } from 'lucide-vue-next'
</script>

<template>
  <footer class="border-t border-slate-200/80 dark:border-slate-800 bg-white/50 dark:bg-slate-900/50 py-8 mt-16 transition-colors">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-slate-500 dark:text-slate-400">
      <div class="flex items-center gap-2">
        <ShieldCheck class="w-4 h-4 text-brand-600 dark:text-brand-400" />
        <span>网络检测需要读取您的公网 IP。本站默认不建立个人 IP 查询历史。</span>
      </div>

      <div class="flex items-center gap-4">
        <router-link to="/docs/api" class="hover:text-brand-600 dark:hover:text-brand-400">
          API 接口
        </router-link>
        <span>·</span>
        <span>NetIP 开源自托管网络诊断工具</span>
      </div>
    </div>
  </footer>
</template>
