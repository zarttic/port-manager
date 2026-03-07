<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import PortList from './components/PortList/PortList.svelte';
  import StatsChart from './components/StatsChart/StatsChart.svelte';
  import { ScanPorts, KillProcess, GetPortStats, GetTopUsedPorts } from '../wailsjs/go/main/App';

  let ports: any[] = [];
  let loading = false;
  let error = '';
  let stats: any = null;
  let selectedPort: any = null;
  let autoRefresh = false;
  let refreshInterval: number;

  onMount(async () => {
    await scanPorts();
    await loadStats();
  });

  onDestroy(() => {
    if (refreshInterval) {
      clearInterval(refreshInterval);
    }
  });

  async function scanPorts() {
    loading = true;
    error = '';
    try {
      ports = await ScanPorts();
      console.log('Scanned ports:', ports.length);
    } catch (e: any) {
      error = e.message || 'Failed to scan ports';
      console.error('Scan error:', e);
    } finally {
      loading = false;
    }
  }

  async function handleKillProcess(pid: number) {
    try {
      await KillProcess(pid);
      // Refresh the port list
      await scanPorts();
    } catch (e: any) {
      error = e.message || 'Failed to kill process';
      console.error('Kill process error:', e);
    }
  }

  async function loadStats() {
    try {
      stats = await GetTopUsedPorts(10);
    } catch (e: any) {
      console.error('Stats error:', e);
    }
  }

  function toggleAutoRefresh() {
    autoRefresh = !autoRefresh;
    if (autoRefresh) {
      refreshInterval = setInterval(scanPorts, 5000);
    } else {
      clearInterval(refreshInterval);
    }
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-dark-900 via-dark-800 to-dark-900">
  <!-- Header -->
  <header class="border-b border-dark-700 bg-dark-800/50 backdrop-blur-sm">
    <div class="max-w-7xl mx-auto px-4 py-4">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold bg-gradient-to-r from-primary-400 to-primary-600 bg-clip-text text-transparent">
            Port Manager
          </h1>
          <p class="text-sm text-dark-400 mt-1">端口管理工具</p>
        </div>

        <div class="flex items-center space-x-3">
          <button
            on:click={toggleAutoRefresh}
            class="px-4 py-2 {autoRefresh ? 'bg-green-600 hover:bg-green-700' : 'bg-dark-700 hover:bg-dark-600'} rounded-lg font-medium transition-all duration-200"
          >
            {autoRefresh ? '🔄 自动刷新中' : '⏸️ 自动刷新'}
          </button>

          <button
            on:click={scanPorts}
            disabled={loading}
            class="px-6 py-2 bg-primary-600 hover:bg-primary-700 disabled:bg-dark-600 rounded-lg font-medium transition-all duration-200 transform hover:scale-105 active:scale-95"
          >
            {loading ? '扫描中...' : '🔍 扫描端口'}
          </button>
        </div>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="max-w-7xl mx-auto px-4 py-8">
    {#if error}
      <div class="mb-4 p-4 bg-red-900/20 border border-red-700 rounded-lg animate-fade-in">
        <p class="text-red-400">{error}</p>
      </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Port List -->
      <div class="lg:col-span-2">
        <div class="bg-dark-800 rounded-xl border border-dark-700 overflow-hidden animate-fade-in">
          <div class="p-4 border-b border-dark-700">
            <h2 class="text-lg font-semibold">端口列表</h2>
            <p class="text-sm text-dark-400">找到 {ports.length} 个端口</p>
          </div>

          {#if loading}
            <div class="p-8 text-center">
              <div class="inline-block animate-spin rounded-full h-8 w-8 border-2 border-primary-500 border-t-transparent"></div>
              <p class="mt-2 text-dark-400">扫描中...</p>
            </div>
          {:else if ports.length > 0}
            <PortList {ports} onKillProcess={handleKillProcess} />
          {:else}
            <div class="p-8 text-center text-dark-400">
              点击"扫描端口"开始扫描
            </div>
          {/if}
        </div>
      </div>

      <!-- Stats Panel -->
      <div>
        <StatsChart {ports} />
      </div>
    </div>
  </main>
</div>
