<script setup lang="ts">
import { ShieldCheck, Globe, HeartHandshake, Code2, Sparkles, ExternalLink } from 'lucide-vue-next'
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

          <!-- GitHub Repo Badge -->
          <a
            href="https://github.com/sunzhonghui/netip"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-slate-900 text-white dark:bg-slate-800 dark:text-slate-100 hover:bg-slate-800 dark:hover:bg-slate-700 font-medium transition-all shadow-sm group"
          >
            <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
              <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
            </svg>
            <span>GitHub 开源</span>
            <ExternalLink class="w-3 h-3 text-slate-400 group-hover:text-white" />
          </a>
        </div>

        <div class="flex items-center gap-4 text-xs">
          <a
            href="https://github.com/sunzhonghui/netip"
            target="_blank"
            rel="noopener noreferrer"
            class="text-slate-600 dark:text-slate-300 hover:text-brand-600 dark:hover:text-brand-400 transition-colors font-medium flex items-center gap-1"
          >
            <span>源码仓库</span>
          </a>
          <span>·</span>
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
          <span class="block sm:inline mt-1 sm:mt-0 text-slate-400">
            开源地址:
            <a
              href="https://github.com/sunzhonghui/netip"
              target="_blank"
              rel="noopener noreferrer"
              class="text-brand-600 dark:text-brand-400 hover:underline ml-1 font-mono"
            >
              github.com/sunzhonghui/netip
            </a>
          </span>
        </div>
        <div class="flex items-center gap-1 text-slate-400 text-[11px]">
          <Sparkles class="w-3 h-3 text-amber-500" />
          <span>Powered by Go & Vue 3</span>
        </div>
      </div>
    </div>
  </footer>
</template>
