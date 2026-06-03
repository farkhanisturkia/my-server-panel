<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const isMobileMenuOpen = ref<boolean>(false)
const route = useRoute()

const isCloudflareActive = computed(() => route.path.startsWith('/cloudflare'))
</script>

<template>
  <div class="min-h-screen bg-linear-to-br from-gray-950 via-slate-900 to-gray-950 text-gray-100 font-sans selection:bg-cyan-500/30 flex flex-col">
    
    <nav class="bg-gray-950/80 backdrop-blur-md border-b border-cyan-500/10 shadow-[0_4px_30px_rgba(0,0,0,0.4)] sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          
          <div class="flex items-center space-x-3">
            <div class="w-8 h-8 rounded-lg bg-linear-to-tr from-blue-600 to-cyan-400 flex items-center justify-center shadow-[0_0_15px_rgba(6,182,212,0.4)]">
              <span class="text-sm font-black text-white">MS</span>
            </div>
            <div class="flex flex-col">
              <span class="text-sm font-black tracking-wider text-white uppercase">MS Server Tools</span>
              <span class="text-[9px] uppercase font-bold tracking-[0.15em] text-cyan-400/80 animate-pulse">Core Engine</span>
            </div>
          </div>

          <div class="hidden md:flex items-center space-x-4">
            <router-link 
              to="/" 
              custom 
              v-slot="{ navigate, isActive }"
            >
              <button 
                @click="navigate"
                :class="isActive ? 'bg-linear-to-r from-blue-600 to-cyan-600 text-white shadow-[0_0_15px_rgba(6,182,212,0.2)]' : 'bg-slate-900/50 hover:bg-slate-800 text-slate-400 hover:text-slate-200'"
                class="border border-cyan-400/20 px-4 py-1.5 rounded-lg text-xs font-bold uppercase tracking-wider flex items-center space-x-1.5 transition duration-200"
              >
                <span>🐳 Docker</span>
              </button>
            </router-link>

            <router-link 
              to="/cloudflare" 
              custom 
              v-slot="{ navigate }"
            >
              <button 
                @click="navigate"
                :class="isCloudflareActive ? 'bg-linear-to-r from-orange-600 to-amber-500 text-white shadow-[0_0_15px_rgba(245,158,11,0.2)]' : 'bg-slate-900/50 hover:bg-slate-800 text-slate-400 hover:text-slate-200'"
                class="border border-orange-500/20 px-4 py-1.5 rounded-lg text-xs font-bold uppercase tracking-wider flex items-center space-x-1.5 transition duration-200"
              >
                <span>☁️ Cloudflare</span>
              </button>
            </router-link>
          </div>

          <div class="flex md:hidden">
            <button @click="isMobileMenuOpen = !isMobileMenuOpen" class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-cyan-400 hover:bg-slate-900/50 focus:outline-none border border-slate-800/80 transition duration-200">
              <svg class="h-5 w-5" stroke="currentColor" fill="none" viewBox="0 0 24 24">
                <path v-if="!isMobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

        </div>
      </div>

      <div v-if="isMobileMenuOpen" class="md:hidden bg-slate-950 border-b border-slate-800 px-4 pt-2 pb-4 space-y-2 shadow-2xl animate-[slideDown_0.2s_ease-out]">
        <router-link to="/" custom v-slot="{ navigate, isActive }">
          <button @click="() => { navigate(); isMobileMenuOpen = false; }" :class="isActive ? 'bg-linear-to-r from-blue-600 to-cyan-600 text-white' : 'text-slate-400 bg-slate-900/40'" class="w-full text-left px-4 py-2.5 rounded-xl text-xs font-bold uppercase tracking-wider flex items-center space-x-2">🐳 Docker</button>
        </router-link>
        
        <router-link to="/cloudflare" custom v-slot="{ navigate }">
          <button @click="() => { navigate(); isMobileMenuOpen = false; }" :class="isCloudflareActive ? 'bg-linear-to-r from-orange-600 to-amber-500 text-white' : 'text-slate-400 bg-slate-900/40'" class="w-full text-left px-4 py-2.5 rounded-xl text-xs font-bold uppercase tracking-wider flex items-center space-x-2">☁️ Cloudflare</button>
        </router-link>
      </div>
    </nav>

    <main class="flex-1 p-4 md:p-8 overflow-y-auto relative max-w-7xl w-full mx-auto">
      <router-view />
    </main>

  </div>
</template>

<style scoped>
@keyframes slideDown {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>