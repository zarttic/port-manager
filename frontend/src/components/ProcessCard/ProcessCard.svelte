<script lang="ts">
  export let process: any;
  export let onKill: () => void;

  let showConfirm = false;

  function toggleConfirm() {
    showConfirm = !showConfirm;
  }
</script>

<div class="bg-dark-800 border border-dark-700 rounded-lg p-4 animate-scale-in">
  <div class="flex items-center justify-between mb-3">
    <h3 class="font-semibold text-lg">{process.name}</h3>
    <button
      on:click={toggleConfirm}
      class="px-3 py-1 text-sm bg-red-900/20 hover:bg-red-900/40 text-red-400 rounded transition-colors"
    >
      杀死进程
    </button>
  </div>

  {#if showConfirm}
    <div class="mb-3 p-3 bg-red-900/20 border border-red-700 rounded animate-fade-in">
      <p class="text-sm text-red-300 mb-2">确定要杀死此进程吗？</p>
      <div class="flex space-x-2">
        <button
          on:click={() => {
            onKill();
            showConfirm = false;
          }}
          class="px-3 py-1 text-sm bg-red-600 hover:bg-red-700 text-white rounded"
        >
          确认
        </button>
        <button
          on:click={() => showConfirm = false}
          class="px-3 py-1 text-sm bg-dark-700 hover:bg-dark-600 text-dark-300 rounded"
        >
          取消
        </button>
      </div>
    </div>
  {/if}

  <div class="space-y-2 text-sm">
    <div class="flex justify-between">
      <span class="text-dark-400">PID:</span>
      <span class="font-mono">{process.pid}</span>
    </div>
    <div class="flex justify-between">
      <span class="text-dark-400">内存:</span>
      <span>{(process.memory / 1024 / 1024).toFixed(2)} MB</span>
    </div>
    <div class="flex justify-between">
      <span class="text-dark-400">CPU:</span>
      <span>{process.cpu.toFixed(1)}%</span>
    </div>
    <div class="flex justify-between">
      <span class="text-dark-400">用户:</span>
      <span>{process.user || 'N/A'}</span>
    </div>
  </div>
</div>
