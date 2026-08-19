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
      <router-link to="/" class="flex items-center gap-2.5 font-bold text-lg text-slate-900 dark:text-white tracking-tight">
        <div class="w-8 h-8 rounded-lg bg-brand-600 flex items-center justify-center text-white shadow-sm shadow-brand-500/20">
          <Activity class="w-5 h-5" />
        </div>
        <span class="text-xl">Net<span class="text-brand-600 dark:text-brand-400">IP</span></span>
      </router-link>

      <!-- Desktop Navigation -->
      <nav class="hidden md:flex items-center gap-1 text-sm font-medium text-slate-600 dark:text-slate-300">
        <router-link
          to="/"
          class="px-3 py-1.5 rounded-lg hover:text-brand-600 hover:bg-slate-100 dark:hover:bg-slate-800 dark:hover:text-brand-400 transition-colors"
          active-class="text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-950/50 font-semibold"
        >
          首页
        </router-link>

        <router-link
          to="/ip"
          class="px-3 py-1.5 rounded-lg hover:text-brand-600 hover:bg-slate-100 dark:hover:bg-slate-800 dark:hover:text-brand-400 transition-colors"
          active-class="text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-950/50 font-semibold"
        >
          IP 查询
        </router-link>

        <router-link
          to="/dns"
          class="px-3 py-1.5 rounded-lg hover:text-brand-600 hover:bg-slate-100 dark:hover:bg-slate-800 dark:hover:text-brand-400 transition-colors"
          active-class="text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-950/50 font-semibold"
        >
          DNS 查询
        </router-link>

        <router-link
          to="/ping"
          class="px-3 py-1.5 rounded-lg hover:text-brand-600 hover:bg-slate-100 dark:hover:bg-slate-800 dark:hover:text-brand-400 transition-colors"
          active-class="text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-950/50 font-semibold"
        >
          Ping / TCPing
        </router-link>

        <router-link
          to="/ipv6"
          class="px-3 py-1.5 rounded-lg hover:text-brand-600 hover:bg-slate-100 dark:hover:bg-slate-800 dark:hover:text-brand-400 transition-colors"
          active-class="text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-950/50 font-semibold"
        >
          IPv6 检测
        </router-link>

        <router-link
          to="/ssl"
          class="px-3 py-1.5 rounded-lg hover:text-brand-600 hover:bg-slate-100 dark:hover:bg-slate-800 dark:hover:text-brand-400 transition-colors"
          active-class="text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-950/50 font-semibold"
        >
          SSL 证书
        </router-link>

        <router-link
          to="/speed"
          class="px-3 py-1.5 rounded-lg hover:text-brand-600 hover:bg-slate-100 dark:hover:bg-slate-800 dark:hover:text-brand-400 transition-colors"
          active-class="text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-950/50 font-semibold"
        >
          HTTP 测速
        </router-link>

        <router-link
          to="/docs/api"
          class="px-3 py-1.5 rounded-lg hover:text-brand-600 hover:bg-slate-100 dark:hover:bg-slate-800 dark:hover:text-brand-400 transition-colors"
          active-class="text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-950/50 font-semibold"
        >
          API
        </router-link>
      </nav>

      <!-- Right actions -->
      <div class="flex items-center gap-2">
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
