<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'

const props = defineProps<{
  containerId: string;
}>()

const emit = defineEmits(['close'])

const terminalRef = ref<HTMLElement | null>(null)
let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let ws: WebSocket | null = null

onMounted(() => {
  // 1. Setup Xterm dengan style monokrom Matrix/Cyberpunk
  term = new Terminal({
    cursorBlink: true,
    cursorStyle: 'underline',
    fontSize: 13,
    fontFamily: 'Fira Code, JetBrains Mono, Courier New, monospace',
    letterSpacing: 0.5,
    theme: {
      background: '#020617', // slate-950
      foreground: '#22d3ee', // cyan-400
      cursor: '#22d3ee',
      selectionBackground: 'rgba(6, 182, 212, 0.3)',
      black: '#020617',
      red: '#ef4444',
      green: '#10b981',
      yellow: '#f59e0b',
      blue: '#3b82f6',
      magenta: '#d946ef',
      cyan: '#06b6d4',
      white: '#cbd5e1'
    }
  })

  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)

  if (terminalRef.value) {
    term.open(terminalRef.value)
    fitAddon.fit()
  }

  // Ambil base URL WebSocket dari env
  const WS_BASE_URL = import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8080/api'

  // Gabungkan dengan path spesifik kontainer
  const wsUrl = `${WS_BASE_URL}/docker/containers/${props.containerId}/terminal`

  // 2. Hubungkan ke WebSocket Backend (Menyesuaikan Port dev server Go kamu)
  ws = new WebSocket(wsUrl)

  // Alirkan output dari backend ke layar terminal xterm
  ws.onmessage = (event) => {
    term?.write(event.data)
  }

  // Tangkap ketikan user di terminal, kirim langsung ke stdin backend
  term.onData((data) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(data)
    }
  })

  ws.onclose = () => {
    term?.writeln('\r\n\x1b[1;31m[MATRIX_SHELL]: CONNECTION CLOSED BY HOST ENGINE.\x1b[0m')
  }

  ws.onerror = () => {
    term?.writeln('\r\n\x1b[1;31m[ERROR]: FAILED TO ESTABLISH WEBSOCKET STREAM.\x1b[0m')
  }

  // Handle auto-fit layout saat ukuran jendela browser berubah
  window.addEventListener('resize', handleResize)
})

const handleResize = () => {
  fitAddon?.fit()
}

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  ws?.close()
  term?.dispose()
})
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-3 md:p-6 bg-black/80 backdrop-blur-md animate-[fadeIn_0.2s_ease-out]">
    
    <div class="w-full max-w-4xl bg-slate-950 border border-slate-800/80 rounded-2xl shadow-[0_0_50px_rgba(6,182,212,0.15)] flex flex-col h-[85vh] max-h-150 overflow-hidden relative">
      
      <div class="bg-slate-900/90 border-b border-slate-800/60 px-4 py-3 flex items-center justify-between">
        <div class="flex items-center space-x-2.5 font-mono text-xs text-cyan-400">
          <span class="w-2 h-2 rounded-full bg-cyan-500 shadow-[0_0_8px_#06b6d4] animate-pulse"></span>
          <span class="font-bold uppercase tracking-wider">TTY_CONSOLE</span>
        </div>
        
        <button 
          @click="emit('close')"
          class="bg-slate-950 hover:bg-rose-500/10 border border-slate-800 hover:border-rose-500/30 text-slate-400 hover:text-rose-400 px-2.5 py-1 rounded-md text-xs font-mono transition duration-200"
        >
          CLOSE
        </button>
      </div>

      <div ref="terminalRef" class="flex-1 bg-[#020617] p-3 overflow-hidden terminal-container"></div>
      
      <div class="bg-slate-900/40 border-t border-slate-800/40 px-4 py-1.5 flex justify-between items-center text-[10px] font-mono text-slate-500">
        <span>BAUD_RATE: 115200</span>
        <span>STREAM: WEBSOCKET_RAW</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
:deep(.xterm-viewport::-webkit-scrollbar) {
  width: 6px;
}
:deep(.xterm-viewport::-webkit-scrollbar-track) {
  background: #020617;
}
:deep(.xterm-viewport::-webkit-scrollbar-thumb) {
  background: #1e293b;
  border-radius: 3px;
}
:deep(.xterm-viewport::-webkit-scrollbar-thumb:hover) {
  background: #06b6d4;
}
</style>