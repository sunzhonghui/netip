<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  Activity,
  Sun,
  Moon,
  Menu,
  X,
  Globe2,
  Network,
  Lock,
  FileCode2,
} from 'lucide-vue-next'
import { useThemeStore } from '@/stores/theme'

const router = useRouter()
const themeStore = useThemeStore()
const mobileMenuOpen = ref(false)

function navigate(path: string) {
  router.push(path)
  mobileMenuOpen.value = false
}

const navLinks = [
  { name: '首页', path: '/' },
  { name: 'IP 查询', path: '/ip', icon: Globe2 },
  { name: 'DNS 查询', path: '/dns', icon: Network },
  { name: 'Ping / TCPing', path: '/ping', icon: Activity },
  { name: 'IPv6 检测', path: '/ipv6' },
  { name: 'SSL 证书', path: '/ssl', icon: Lock },
  { name: 'WHOIS', path: '/whois' },
  { name: '网站测速', path: '/speed' },
  { name: 'HTTP 诊断', path: '/http' },
  { name: 'ASN', path: '/asn' },
  { name: 'API 文档', path: '/docs/api', icon: FileCodeCodeWrapper },
]

function FileCodeCodeWrapper() {
  return FileCode2
}
</script>

<template>
  <header class="sticky top-0 z-50 bg-white/85 dark:bg-slate-900/85 backdrop-blur-md border-b border-slate-200/80 dark:border-slate-800 transition-colors">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between">
      <!-- Logo -->
      <router-link to="/" class="flex items-center gap-2.5 font-bold text-lg text-slate-900 dark:text-white tracking-tight group">
        <img src="/favicon.svg" alt="NetIP Logo" class="w-8 h-8 rounded-lg shadow-sm transition-transform group-hover:scale-105" />
        <span class="text-xl">Net<span class="text-brand-600 dark:text-brand-400">IP</span></span>
      </router-link>

      <!-- Desktop Navigation -->
      <nav class="hidden md:flex items-center gap-1 text-sm font-medium text-slate-600 dark:text-slate-300">
        <router-link
          to="/"
          class="px-3 py-1.5 rounded-xl hover:text-brand-600 hover:bg-slate-100/80 dark:hover:bg-slate-800/70 dark:hover:text-brand-400 transition-all text-xs lg:text-sm font-medium"
          active-class="!text-brand-600 dark:!text-brand-300 !bg-brand-500/10 dark:!bg-brand-500/15 !font-bold"
        >
          首页
        </router-link>

        <router-link
          to="/ip"
          class="px-3 py-1.5 rounded-xl hover:text-brand-600 hover:bg-slate-100/80 dark:hover:bg-slate-800/70 dark:hover:text-brand-400 transition-all text-xs lg:text-sm font-medium"
          active-class="!text-brand-600 dark:!text-brand-300 !bg-brand-500/10 dark:!bg-brand-500/15 !font-bold"
        >
          IP 归属地
        </router-link>

        <router-link
          to="/dns"
          class="px-3 py-1.5 rounded-xl hover:text-brand-600 hover:bg-slate-100/80 dark:hover:bg-slate-800/70 dark:hover:text-brand-400 transition-all text-xs lg:text-sm font-medium"
          active-class="!text-brand-600 dark:!text-brand-300 !bg-brand-500/10 dark:!bg-brand-500/15 !font-bold"
        >
          DNS 查询
        </router-link>

        <router-link
          to="/ping"
          class="px-3 py-1.5 rounded-xl hover:text-brand-600 hover:bg-slate-100/80 dark:hover:bg-slate-800/70 dark:hover:text-brand-400 transition-all text-xs lg:text-sm font-medium"
          active-class="!text-brand-600 dark:!text-brand-300 !bg-brand-500/10 dark:!bg-brand-500/15 !font-bold"
        >
          Ping / TCPing
        </router-link>

        <router-link
          to="/ipv6"
          class="px-3 py-1.5 rounded-xl hover:text-brand-600 hover:bg-slate-100/80 dark:hover:bg-slate-800/70 dark:hover:text-brand-400 transition-all text-xs lg:text-sm font-medium"
          active-class="!text-brand-600 dark:!text-brand-300 !bg-brand-500/10 dark:!bg-brand-500/15 !font-bold"
        >
          IPv6 检测
        </router-link>

        <router-link
          to="/ssl"
          class="px-3 py-1.5 rounded-xl hover:text-brand-600 hover:bg-slate-100/80 dark:hover:bg-slate-800/70 dark:hover:text-brand-400 transition-all text-xs lg:text-sm font-medium"
          active-class="!text-brand-600 dark:!text-brand-300 !bg-brand-500/10 dark:!bg-brand-500/15 !font-bold"
        >
          SSL 证书
        </router-link>

        <router-link
          to="/speed"
          class="px-3 py-1.5 rounded-xl hover:text-brand-600 hover:bg-slate-100/80 dark:hover:bg-slate-800/70 dark:hover:text-brand-400 transition-all text-xs lg:text-sm font-medium"
          active-class="!text-brand-600 dark:!text-brand-300 !bg-brand-500/10 dark:!bg-brand-500/15 !font-bold"
        >
          HTTP 测速
        </router-link>

        <router-link
          to="/docs/api"
          class="px-3 py-1.5 rounded-xl hover:text-brand-600 hover:bg-slate-100/80 dark:hover:bg-slate-800/70 dark:hover:text-brand-400 transition-all text-xs lg:text-sm font-medium"
          active-class="!text-brand-600 dark:!text-brand-300 !bg-brand-500/10 dark:!bg-brand-500/15 !font-bold"
        >
          API
        </router-link>
      </nav>

      <!-- Right actions -->
      <div class="flex items-center gap-2">
        <a
          href="https://github.com/sunzhonghui/netip"
          target="_blank"
          rel="noopener noreferrer"
          title="GitHub 源码仓库"
          aria-label="查看 GitHub 开源仓库"
          class="p-2 rounded-lg text-slate-500 hover:text-slate-900 hover:bg-slate-100 dark:text-slate-400 dark:hover:text-white dark:hover:bg-slate-800 transition-colors"
        >
          <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
            <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
          </svg>
        </a>

        <button
          type="button"
          @click="themeStore.toggleTheme"
          aria-label="切换明暗主题"
          class="p-2 rounded-lg text-slate-500 hover:text-slate-900 hover:bg-slate-100 dark:text-slate-400 dark:hover:text-white dark:hover:bg-slate-800 transition-colors"
        >
          <Sun v-if="themeStore.isDark" class="w-5 h-5 text-amber-400" />
          <Moon v-else class="w-5 h-5 text-slate-600" />
        </button>

        <!-- Mobile hamburger -->
        <button
          type="button"
          @click="mobileMenuOpen = !mobileMenuOpen"
          aria-label="切换导航菜单"
          class="md:hidden p-2 rounded-lg text-slate-600 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800 transition-colors"
        >
          <X v-if="mobileMenuOpen" class="w-6 h-6" />
          <Menu v-else class="w-6 h-6" />
        </button>
      </div>
    </div>

    <!-- Mobile Drawer -->
    <div
      v-if="mobileMenuOpen"
      class="md:hidden border-b border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 px-4 py-3 space-y-1 shadow-lg"
    >
      <button
        v-for="item in navLinks"
        :key="item.path"
        @click="navigate(item.path)"
        class="w-full text-left px-3 py-2 rounded-lg text-sm font-medium text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 flex items-center justify-between"
      >
        <span>{{ item.name }}</span>
      </button>
    </div>
  </header>
</template>
