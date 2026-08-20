<script setup lang="ts">
import { onMounted } from 'vue'
import Navbar from '@/components/Navbar.vue'
import Footer from '@/components/Footer.vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

onMounted(() => {
  themeStore.initTheme()
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-ambient-light dark:bg-ambient-dark">
    <Navbar />
    <main class="flex-1 max-w-6xl w-full mx-auto px-4 sm:px-6 py-8">
      <router-view />
    </main>
    <Footer />
  </div>
</template>
