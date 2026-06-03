<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import axios from 'axios'
import TerminalModal from './components/Terminal.vue'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

interface ContainerInfo {
  id: string;
  names: string[];
  image: string;
  state: string;
  status: string;
  working_dir: string;
}

const containers = ref<ContainerInfo[]>([])
const loading = ref<boolean>(false)
const message = ref<string>('')

const openDropdownId = ref<string | null>(null)
const isMobileMenuOpen = ref<boolean>(false)

const isLocalTerminalOpen = ref<boolean>(false)
const activeTerminalContainer = ref<ContainerInfo | null>(null)

const fetchContainers = async () => {
  loading.value = true
  try {
    const response = await axios.get<ContainerInfo[]>(`${API_BASE_URL}/docker/containers`)
    containers.value = response.data
    // message.value = ''
  } catch (error: any) {
    message.value = 'Gagal memuat data dari server backend.'
    console.error(error)
  } finally {
    loading.value = false
  }
}

let messageTimer: ReturnType<typeof setTimeout> | null = null

const handleAction = async (action: 'start' | 'stop' | 'down', id: string) => {
  loading.value = true
  openDropdownId.value = null 

  if (messageTimer) {
    clearTimeout(messageTimer)
  }

  try {
    message.value = `Menjalankan perintah ${action.toUpperCase()} untuk kontainer...`
    const response = await axios.post<{ message: string }>(`${API_BASE_URL}/docker/containers/${id}/${action}`)
    await fetchContainers()
    message.value = response.data.message
    messageTimer = setTimeout(() => {
      message.value = ''
    }, 3000)
  } catch (error: any) {
    message.value = error.response?.data?.error || `Gagal melakukan aksi ${action}`
    messageTimer = setTimeout(() => {
      message.value = ''
    }, 3000)
  } finally {
    loading.value = false
  }
}

const toggleDropdown = (id: string, event: Event) => {
  event.stopPropagation()
  if (openDropdownId.value === id) {
    openDropdownId.value = null
  } else {
    openDropdownId.value = id
  }
}

const closeAllDropdowns = () => {
  openDropdownId.value = null
}

const openProjectTerminal = (item: ContainerInfo) => {
  openDropdownId.value = null
  activeTerminalContainer.value = item
}

onMounted(() => {
  fetchContainers()
  window.addEventListener('click', closeAllDropdowns)
})

onBeforeUnmount(() => {
  window.removeEventListener('click', closeAllDropdowns)
})
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
            <button class="bg-linear-to-r from-blue-600 to-cyan-600 text-white border border-cyan-400/20 px-4 py-1.5 rounded-lg text-xs font-bold uppercase tracking-wider shadow-[0_0_15px_rgba(6,182,212,0.2)] flex items-center space-x-1.5">
              <span>🐳 Docker</span>
            </button>
            <div class="text-gray-500 bg-gray-900/40 opacity-40 border border-dashed border-gray-800 px-4 py-1.5 rounded-lg text-xs font-bold uppercase tracking-wider cursor-not-allowed select-none">🌐 DNS</div>
            <div class="text-gray-500 bg-gray-900/40 opacity-40 border border-dashed border-gray-800 px-4 py-1.5 rounded-lg text-xs font-bold uppercase tracking-wider cursor-not-allowed select-none">🚇 Tunnel</div>
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
        <button class="w-full text-left bg-linear-to-r from-blue-600 to-cyan-600 text-white px-4 py-2.5 rounded-xl text-xs font-bold uppercase tracking-wider flex items-center space-x-2">🐳 Docker</button>
        <div class="w-full text-left text-gray-600 bg-gray-900/20 opacity-30 border border-dashed border-gray-800 px-4 py-2.5 rounded-xl text-xs font-bold uppercase tracking-wider cursor-not-allowed">🌐 DNS</div>
        <div class="w-full text-left text-gray-600 bg-gray-900/20 opacity-30 border border-dashed border-gray-800 px-4 py-2.5 rounded-xl text-xs font-bold uppercase tracking-wider cursor-not-allowed">🚇 Tunnel</div>
      </div>
    </nav>

    <main class="flex-1 p-4 md:p-8 overflow-y-auto relative max-w-7xl w-full mx-auto">
      <div class="absolute top-0 right-1/4 w-96 h-96 bg-cyan-500/5 rounded-full blur-[120px] pointer-events-none"></div>
      
      <header class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8 border-b border-slate-800 pb-5 relative z-10">
        <div>
          <h1 class="text-xl md:text-2xl font-black bg-clip-text text-transparent bg-linear-to-r from-white via-slate-200 to-slate-400 tracking-tight uppercase">Docker Node Matrix</h1>
          <p class="text-xs text-cyan-400/70 font-mono mt-1 tracking-wide">// Real-time container subsystem monitoring</p>
        </div>
        
        <div class="flex flex-col sm:flex-row items-center gap-2.5 w-full sm:w-auto">
          <button 
            @click="isLocalTerminalOpen = true"
            class="w-full sm:w-auto bg-slate-950 hover:bg-cyan-500/10 text-cyan-400 border border-cyan-500/20 hover:border-cyan-500/50 font-mono text-xs px-5 py-2 rounded-lg transition duration-200 flex items-center justify-center space-x-2 shadow-[0_0_15px_rgba(6,182,212,0.02)]"
          >
            <span>💻</span> <span>LOCAL_SHELL</span>
          </button>

          <button @click="fetchContainers" :disabled="loading" class="w-full sm:w-auto bg-slate-900 hover:bg-slate-800 text-slate-300 border border-slate-800 hover:border-slate-700 font-mono text-xs px-5 py-2 rounded-lg transition duration-300">
            {{ loading ? 'SYNCING...' : '⚡ REFRESH NODE' }}
          </button>
        </div>
      </header>

      <Transition name="fade">
        <div v-if="message" class="mb-6 p-3.5 bg-slate-950/80 border border-cyan-500/20 rounded-xl text-xs text-cyan-300 font-mono shadow-[0_0_15px_rgba(6,182,212,0.05)]">
          SYSTEM LOG: {{ message }}
        </div>
      </Transition>

      <div class="w-full bg-slate-950/40 backdrop-blur-sm border border-slate-800 rounded-2xl shadow-[0_20px_50px_rgba(0,0,0,0.4)] overflow-visible">
        
        <table class="hidden md:table w-full text-left border-collapse table-auto">
          <thead>
            <tr class="bg-slate-950/80 border-b border-slate-800 text-slate-400 font-mono text-[10px] uppercase tracking-[0.15em]">
              <th class="py-4 px-5 w-28">HEX_ID</th>
              <th class="py-4 px-5">Subsystem Identity</th>
              <th class="py-4 px-5">Source Image</th>
              <th class="py-4 px-5">Cluster Status</th>
              <th class="py-4 px-5 text-right w-32">Operation</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-900 text-xs">
            <tr v-for="item in containers" :key="item.id" class="hover:bg-cyan-500/2 border-l-2 border-l-transparent hover:border-l-cyan-500 transition-all duration-200 group">
              <td class="py-4 px-5 font-mono text-slate-500 group-hover:text-cyan-400 transition">#{{ item.id.substring(0, 8) }}</td>
              <td class="py-4 px-5 font-bold text-slate-200 tracking-wide">{{ item.names[0] }}</td>
              <td class="py-4 px-5 font-mono text-slate-400 max-w-60 truncate opacity-80" :title="item.image">{{ item.image }}</td>
              <td class="py-4 px-5">
                <div class="flex items-center space-x-3">
                  <div class="flex items-center">
                    <span :class="item.state === 'running' ? 'bg-emerald-500 shadow-[0_0_10px_#10b981]' : 'bg-rose-500 shadow-[0_0_10px_#f43f5e]'" class="w-2 h-2 rounded-full inline-block mr-2"></span>
                    <span :class="item.state === 'running' ? 'text-emerald-400' : 'text-rose-400'" class="text-[10px] font-mono font-bold uppercase tracking-wider">{{ item.state }}</span>
                  </div>
                  <span class="text-slate-500 text-[11px] font-mono truncate max-w-45 border-l border-slate-800 pl-3">{{ item.status }}</span>
                </div>
              </td>
              <td class="py-4 px-5 text-right">
                <div class="inline-block relative text-left">
                  <button @click="(e) => toggleDropdown(item.id, e)" :class="openDropdownId === item.id ? 'border-cyan-400 text-cyan-400 bg-slate-900' : 'border-slate-700 text-slate-300 bg-slate-950'" class="hover:bg-slate-900 border font-mono font-bold px-3 py-1 rounded-md transition duration-200 text-[10px] uppercase tracking-wider inline-flex items-center space-x-1.5 focus:outline-none">
                    <span>Manage</span>
                    <span :class="openDropdownId === item.id ? 'rotate-180' : ''" class="text-[8px] transition-transform duration-200">▼</span>
                  </button>
                  
                  <div v-if="openDropdownId === item.id" class="absolute right-0 top-full mt-2 w-36 bg-slate-950/95 backdrop-blur-md border border-slate-800 rounded-xl shadow-[0_10px_30px_rgba(0,0,0,0.5)] z-50 overflow-hidden ring-1 ring-cyan-500/10 animate-[fadeIn_0.15s_ease-out]">
                    <button @click="handleAction('start', item.id)" :disabled="item.state === 'running' || loading" class="w-full text-left px-4 py-2.5 text-emerald-400 hover:bg-emerald-500/10 text-[11px] font-mono font-bold transition flex items-center space-x-2 disabled:opacity-20">⚡ UP</button>
                    <button @click="handleAction('stop', item.id)" :disabled="item.state !== 'running' || loading" class="w-full text-left px-4 py-2.5 text-amber-400 hover:bg-amber-500/10 text-[11px] font-mono font-bold transition flex items-center space-x-2 disabled:opacity-20">🛑 STOP</button>
                    <button @click="handleAction('down', item.id)" :disabled="loading" class="w-full text-left px-4 py-2.5 text-rose-400 hover:bg-rose-500/10 text-[11px] font-mono font-bold transition flex items-center space-x-2 disabled:opacity-30">💥 DOWN</button>
                    <button @click="openProjectTerminal(item)" class="w-full text-left px-4 py-2.5 text-cyan-400 hover:bg-cyan-500/10 text-[11px] font-mono font-bold transition flex items-center space-x-2 border-t border-slate-900">💻 SHELL</button>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <div class="block md:hidden p-4 space-y-4">
          <div v-for="item in containers" :key="item.id" class="p-4 bg-slate-900/60 border border-slate-800/80 rounded-xl flex flex-col justify-between relative overflow-visible">
            <div class="flex justify-between items-center mb-3">
              <span class="text-[10px] font-mono text-slate-500">#{{ item.id.substring(0, 8) }}</span>
              <div class="flex items-center bg-slate-950 px-2 py-0.5 rounded border border-slate-800">
                <span :class="item.state === 'running' ? 'bg-emerald-500' : 'bg-rose-500'" class="w-1.5 h-1.5 rounded-full inline-block mr-1.5"></span>
                <span :class="item.state === 'running' ? 'text-emerald-400' : 'text-rose-400'" class="text-[9px] font-mono font-bold uppercase">{{ item.state }}</span>
              </div>
            </div>
            <div class="mb-4">
              <h3 class="text-sm font-bold text-slate-200 mb-1">{{ item.names[0] }}</h3>
              <p class="text-[11px] text-slate-400 font-mono truncate"><span class="text-slate-600">IMG:</span> {{ item.image }}</p>
              <p class="text-[10px] text-slate-500 font-mono mt-1 italic">{{ item.status }}</p>
            </div>
            <div class="border-t border-slate-800/60 pt-3 flex justify-end relative">
              <div class="relative w-full text-right">
                <button @click="(e) => toggleDropdown(item.id, e)" class="w-full bg-slate-950 border border-slate-700 text-slate-300 font-mono font-bold px-4 py-1.5 rounded-md transition text-[10px] uppercase tracking-wider inline-flex items-center justify-center space-x-2">
                  <span>Manage</span>
                  <span :class="openDropdownId === item.id ? 'rotate-180' : ''" class="text-[8px] transition-transform duration-200">▼</span>
                </button>
                
                <div v-if="openDropdownId === item.id" class="absolute right-0 bottom-full mb-2 w-full bg-slate-950 border border-slate-700 rounded-xl shadow-2xl z-50 overflow-hidden ring-1 ring-cyan-500/20">
                  <button @click="handleAction('start', item.id)" :disabled="item.state === 'running' || loading" class="w-full text-left px-4 py-3 text-emerald-400 hover:bg-emerald-500/10 text-[11px] font-mono font-bold transition flex items-center space-x-3 disabled:opacity-20">⚡ UP</button>
                  <button @click="handleAction('stop', item.id)" :disabled="item.state !== 'running' || loading" class="w-full text-left px-4 py-3 text-amber-400 hover:bg-amber-500/10 text-[11px] font-mono font-bold transition flex items-center space-x-3 disabled:opacity-20">🛑 STOP</button>
                  <button @click="handleAction('down', item.id)" :disabled="loading" class="w-full text-left px-4 py-3 text-rose-400 hover:bg-rose-500/10 text-[11px] font-mono font-bold transition flex items-center space-x-3 disabled:opacity-30">💥 DOWN</button>
                  <button @click="openProjectTerminal(item)" class="w-full text-left px-4 py-3 text-cyan-400 hover:bg-cyan-500/10 text-[11px] font-mono font-bold transition flex items-center space-x-3 border-t border-slate-900">💻 SHELL</button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="containers.length === 0 && !loading" class="text-center py-12 text-slate-500 font-mono text-xs tracking-wide">
          // NO ACTIVE REPOSITORIES DETECTED ON LOCAL ENGINE
        </div>
      </div>
    </main>

    <TerminalModal 
      v-if="isLocalTerminalOpen"
      container-id="local-machine"
      @close="isLocalTerminalOpen = false"
    />

    <TerminalModal 
      v-if="activeTerminalContainer"
      :container-id="activeTerminalContainer.id"
      @close="activeTerminalContainer = null"
    />

  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-5px); }
  to { opacity: 1; transform: translateY(0); }
}
@keyframes slideDown {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.4s ease-in-out;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>