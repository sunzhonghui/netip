<script setup lang="ts">
import { ShieldCheck, Globe, HeartHandshake, Code2, Sparkles } from 'lucide-vue-next'
</script>

<template>
  <footer class="border-t border-slate-200/80 dark:border-slate-800 bg-white/60 dark:bg-slate-900/60 py-10 mt-16 transition-colors">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 space-y-6 text-xs text-slate-500 dark:text-slate-400">
      <!-- Main Highlights Bar -->
      <div class="flex flex-wrap items-center justify-between gap-4 border-b border-slate-100 dark:border-slate-800/80 pb-6">
        <div class="flex flex-wrap items-center gap-3 font-medium">
          <!-- IPv6 Badge -->
          <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-purple-50 text-purple-700 dark:bg-purple-950/60 dark:text-purple-300 font-semibold border border-purple-200/60 dark:border-purple-800/60">
            <Globe class="w-3.5 h-3.5" />
            <span>本站完整支持 IPv6</span>
          </span>

          <!-- Free & Permanent Badge -->
          <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-50 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-300 font-semibold border border-emerald-200/60 dark:border-emerald-800/60">
            <HeartHandshake class="w-3.5 h-3.5" />
            <span>开源 · 永久开放免费 · 承诺持续稳定运行</span>
          </span>

          <!-- No Logging Badge -->
          <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-blue-50 text-blue-700 dark:bg-blue-950/60 dark:text-blue-300 border border-blue-200/60 dark:border-blue-800/60">
            <ShieldCheck class="w-3.5 h-3.5" />
            <span>隐私保护 · 不记录 IP 查询历史</span>
          </span>
        </div>

        <div class="flex items-center gap-4 text-xs">
          <router-link to="/docs/api" class="text-brand-600 dark:text-brand-400 hover:underline font-semibold flex items-center gap-1">
            <Code2 class="w-3.5 h-3.5" />
            <span>API 开发文档</span>
          </router-link>
          <span>·</span>
          <router-link to="/ipv6" class="hover:text-slate-900 dark:hover:text-white transition-colors">
            网站 IPv6 检测
          </router-link>
          <span>·</span>
          <router-link to="/dns" class="hover:text-slate-900 dark:hover:text-white transition-colors">
            DNS 查询
          </router-link>
          <span>·</span>
          <router-link to="/ping" class="hover:text-slate-900 dark:hover:text-white transition-colors">
            Ping / TCPing
          </router-link>
        </div>
      </div>

      <!-- Bottom Credits & Description -->
      <div class="flex flex-col sm:flex-row items-center justify-between gap-3 text-center sm:text-left">
        <div>
          <span>NetIP — 现代化自托管网络诊断与 IP 工具箱</span>
          <span class="hidden sm:inline mx-2">·</span>
          <span class="block sm:inline mt-1 sm:mt-0 text-slate-400">高效、轻量、纯净、无广告，专注纯粹的网络探测体验</span>
        </div>
        <div class="flex items-center gap-1 text-slate-400 text-[11px]">
          <Sparkles class="w-3 h-3 text-amber-500" />
          <span>Powered by Go & Vue 3</span>
        </div>
      </div>
    </div>
  </footer>
</template>
