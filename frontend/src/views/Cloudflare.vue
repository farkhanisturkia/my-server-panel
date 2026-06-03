<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router' // <-- Import router
import axios from 'axios'

const router = useRouter()
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

interface TunnelInfo {
  id: string;
  name: string;
  status: string;
}

const tunnels = ref<TunnelInfo[]>([])
const loading = ref<boolean>(false)
const errorMessage = ref<string>('')

const fetchTunnels = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await axios.get<TunnelInfo[]>(`${API_BASE_URL}/cloudflare/tunnels`)
    tunnels.value = response.data
  } catch (error: any) {
    errorMessage.value = 'Gagal memuat data Cloudflare Tunnel dari server.'
    console.error(error)
  } finally {
    loading.value = false
  }
}

// Fungsi navigasi ke halaman konfigurasi rute
const configureTunnel = (id: string) => {
  router.push(`/cloudflare/${id}`)
}

onMounted(() => {
  fetchTunnels()
})
</script>

<template>
  <div class="absolute top-0 right-1/4 w-96 h-96 bg-orange-500/5 rounded-full blur-[120px] pointer-events-none"></div>
  
  <header class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8 border-b border-slate-800 pb-5 relative z-10">
    <div>
      <h1 class="text-xl md:text-2xl font-black bg-clip-text text-transparent bg-linear-to-r from-white via-slate-200 to-slate-400 tracking-tight uppercase">Cloudflare Edge Networks</h1>
      <p class="text-xs text-orange-400/70 font-mono mt-1 tracking-wide">// Argo Tunnel subsystem configuration</p>
    </div>
    
    <div class="w-full sm:w-auto">
      <button @click="fetchTunnels" :disabled="loading" class="w-full sm:w-auto bg-slate-900 hover:bg-slate-800 text-slate-300 border border-slate-800 hover:border-slate-700 font-mono text-xs px-5 py-2 rounded-lg transition duration-300">
        {{ loading ? 'TUNNEL_SYNCING...' : '⚡ REFRESH TUNNELS' }}
      </button>
    </div>
  </header>

  <div v-if="errorMessage" class="mb-6 p-3.5 bg-slate-950/80 border border-rose-500/20 rounded-xl text-xs text-rose-400 font-mono">
    SYSTEM ERROR: {{ errorMessage }}
  </div>

  <div class="w-full bg-slate-950/40 backdrop-blur-sm border border-slate-800 rounded-2xl shadow-[0_20px_50px_rgba(0,0,0,0.4)] overflow-hidden">
    
    <table class="hidden md:table w-full text-left border-collapse table-auto">
      <thead>
        <tr class="bg-slate-950/80 border-b border-slate-800 text-slate-400 font-mono text-[10px] uppercase tracking-[0.15em]">
          <th class="py-4 px-5 w-1/3">Tunnel ID</th>
          <th class="py-4 px-5">Tunnel Name</th>
          <th class="py-4 px-5">Network Status</th>
          <th class="py-4 px-5 text-right w-32">Configuration</th> </tr>
      </thead>
      <tbody class="divide-y divide-slate-900 text-xs">
        <tr v-for="tunnel in tunnels" :key="tunnel.id" class="hover:bg-orange-500/2 border-l-2 border-l-transparent hover:border-l-orange-500 transition-all duration-200 group">
          <td class="py-4 px-5 font-mono text-slate-500 group-hover:text-orange-400 transition">{{ tunnel.id }}</td>
          <td class="py-4 px-5 font-bold text-slate-200 tracking-wide">{{ tunnel.name }}</td>
          <td class="py-4 px-5">
            <div class="flex items-center">
              <span :class="tunnel.status === 'healthy' || tunnel.status === 'active' ? 'bg-emerald-500 shadow-[0_0_10px_#10b981]' : 'bg-rose-500 shadow-[0_0_10px_#f43f5e]'" class="w-2 h-2 rounded-full inline-block mr-2"></span>
              <span :class="tunnel.status === 'healthy' || tunnel.status === 'active' ? 'text-emerald-400' : 'text-rose-400'" class="text-[10px] font-mono font-bold uppercase tracking-wider">{{ tunnel.status }}</span>
            </div>
          </td>
          <td class="py-4 px-5 text-right">
            <button @click="configureTunnel(tunnel.id)" class="bg-slate-950 hover:bg-orange-500/10 border border-slate-800 hover:border-orange-500/40 text-orange-400 font-mono px-3 py-1.5 rounded-md transition duration-200 text-[10px] uppercase tracking-wider">
              Configure
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <div class="block md:hidden p-4 space-y-4">
      <div v-for="tunnel in tunnels" :key="tunnel.id" class="p-4 bg-slate-900/60 border border-slate-800/80 rounded-xl flex flex-col justify-between">
        <div class="flex justify-between items-center mb-2">
          <span class="text-[10px] font-mono text-slate-500">{{ tunnel.id.substring(0, 8) }}...</span>
          <div class="flex items-center bg-slate-950 px-2 py-0.5 rounded border border-slate-800">
            <span :class="tunnel.status === 'healthy' || tunnel.status === 'active' ? 'bg-emerald-500' : 'bg-rose-500'" class="w-1.5 h-1.5 rounded-full inline-block mr-1.5"></span>
            <span :class="tunnel.status === 'healthy' || tunnel.status === 'active' ? 'text-emerald-400' : 'text-rose-400'" class="text-[9px] font-mono font-bold uppercase">{{ tunnel.status }}</span>
          </div>
        </div>
        <div class="mb-4">
          <h3 class="text-sm font-bold text-slate-200">{{ tunnel.name }}</h3>
          <p class="text-[10px] font-mono text-slate-600 mt-1">ID: {{ tunnel.id }}</p>
        </div>
        <div class="border-t border-slate-800/60 pt-3">
          <button @click="configureTunnel(tunnel.id)" class="w-full bg-slate-950 border border-slate-700 hover:border-orange-500/40 text-orange-400 font-mono py-2 rounded-md text-[10px] uppercase tracking-wider transition">
            Configure Route Matrix
          </button>
        </div>
      </div>
    </div>

    <div v-if="tunnels.length === 0 && !loading" class="text-center py-12 text-slate-500 font-mono text-xs tracking-wide">
      // NO ACTIVE ARGO TUNNELS CONFIGURATION FOUND
    </div>
  </div>
</template>