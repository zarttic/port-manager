<script lang="ts">
  import { fly, fade } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import ProcessCard from '../ProcessCard/ProcessCard.svelte';

  export let ports: any[];

  $: filteredPorts = ports.sort((a, b) => a.port - b.port);
</script>

<div class="max-h-[600px] overflow-y-auto custom-scrollbar">
  <div class="divide-y divide-dark-700">
    {#each filteredPorts as port (port.port)}
      <div
        class="p-4 hover:bg-dark-700/50 transition-all duration-200 cursor-pointer animate-slide-up"
        animate:fly={{ y: 20, duration: 300, easing: cubicOut }}
        transition:fade={{ duration: 200 }}
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <div class="text-2xl font-mono font-bold text-primary-400">
              {port.port}
            </div>
            <div>
              <div class="text-sm text-dark-300">{port.processName || 'Unknown'}</div>
              <div class="flex items-center space-x-2 mt-1">
                <span class="px-2 py-0.5 text-xs font-medium bg-primary-900/30 text-primary-400 rounded">
                  {port.protocol.toUpperCase()}
                </span>
                <span class="px-2 py-0.5 text-xs font-medium bg-dark-700 text-dark-300 rounded">
                  PID: {port.pid}
                </span>
                <span class="px-2 py-0.5 text-xs font-medium bg-dark-700 text-dark-300 rounded">
                  {port.state}
                </span>
              </div>
            </div>
          </div>

          <div class="flex items-center space-x-2">
            <button
              class="px-3 py-1 text-sm bg-red-900/20 hover:bg-red-900/40 text-red-400 rounded transition-colors duration-200"
            >
              杀死进程
            </button>
          </div>
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 6px;
  }

  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }

  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #4b5563;
    border-radius: 3px;
  }

  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #6b7280;
  }
</style>
