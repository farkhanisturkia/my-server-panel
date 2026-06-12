<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const tunnelId = route.params.id as string

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

interface RouteInfo {
  hostname: string;
  service?: string;
}

const routesList = ref<RouteInfo[]>([])
const loading = ref<boolean>(false)
const errorMessage = ref<string>('')

const fetchTunnelRoutes = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await axios.get<RouteInfo[]>(`${API_BASE_URL}/cloudflare/tunnels/${tunnelId}/routes`)
    routesList.value = response.data
  } catch (error: any) {
    errorMessage.value = `Gagal memuat konfigurasi rute untuk objek MATRIX_ID [${tunnelId.substring(0,8)}].`
    console.error(error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchTunnelRoutes()
})
</script>

<template>
  <div class="absolute top-0 right-1/4 w-96 h-96 bg-amber-500/5 rounded-full blur-[120px] pointer-events-none"></div>
  
  <header class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8 border-b border-slate-800 pb-5 relative z-10">
    <div>
      <div class="flex items-center space-x-2 text-[10px] font-mono text-amber-500 uppercase tracking-widest mb-1">
        <span class="cursor-pointer hover:underline" @click="router.push('/cloudflare')">Cloudflare</span> 
        <span>/</span> 
        <span class="text-slate-500">TUNNEL_ID: {{ tunnelId }}</span>
      </div>
      <h1 class="text-xl md:text-2xl font-black bg-clip-text text-transparent bg-linear-to-r from-white via-slate-200 to-slate-400 tracking-tight uppercase">Application Traffic Routes</h1>
      <p class="text-xs text-amber-400/70 font-mono mt-1 tracking-wide">// Live ingress rule mappings and ingress targets</p>
    </div>
    
    <div class="flex items-center gap-2 w-full sm:w-auto">
      <button @click="router.push('/cloudflare')" class="bg-slate-950 hover:bg-slate-900 text-slate-400 border border-slate-800 font-mono text-xs px-4 py-2 rounded-lg transition">
        BACK
      </button>
      <button @click="fetchTunnelRoutes" :disabled="loading" class="bg-slate-900 hover:bg-slate-800 text-slate-300 border border-slate-800 font-mono text-xs px-5 py-2 rounded-lg transition">
        {{ loading ? 'ROUTES SYNCING...' : '⚡ REFRESH ROUTES' }}
      </button>
    </div>
  </header>

  <div v-if="errorMessage" class="mb-6 p-3.5 bg-slate-950/80 border border-rose-500/20 rounded-xl text-xs text-rose-400 font-mono shadow-md">
    SYSTEM ERROR: {{ errorMessage }}
  </div>

  <div class="w-full bg-slate-950/40 backdrop-blur-sm border border-slate-800 rounded-2xl shadow-[0_20px_50px_rgba(0,0,0,0.4)] overflow-hidden">
    
    <table class="hidden md:table w-full text-left border-collapse table-auto">
      <thead>
        <tr class="bg-slate-950/80 border-b border-slate-800 text-slate-400 font-mono text-[10px] uppercase tracking-[0.15em]">
          <th class="py-4 px-5 w-1/3">Hostname</th>
          <th class="py-4 px-5">Target Service</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-900 text-xs">
        <tr v-for="(route, index) in routesList" :key="index" class="hover:bg-amber-500/2 border-l-2 border-l-transparent hover:border-l-amber-500 transition-all duration-200 group">
          <td class="py-4 px-5 font-bold text-slate-200 tracking-wide font-mono">{{ route.hostname }}</td>
          <td class="py-4 px-5 font-mono text-emerald-400">{{ route.service }}</td>
        </tr>
      </tbody>
    </table>

    <div class="block md:hidden p-4 space-y-4">
      <div v-for="(route, index) in routesList" :key="index" class="p-4 bg-slate-900/60 border border-slate-800/80 rounded-xl flex flex-col gap-2">
        <div class="flex justify-between items-center">
          <span class="text-[10px] font-mono text-white">{{ route.hostname }}</span>
          <span class="text-[10px] font-mono text-emerald-400">{{ route.service }}</span>
        </div>
      </div>
    </div>

    <div v-if="routesList.length === 0 && !loading" class="text-center py-12 text-slate-500 font-mono text-xs tracking-wide">
      // NO APPS ROUTE INTERFACES BOUND TO THIS ENGINE TERMINAL
    </div>
  </div>
</template>