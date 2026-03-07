<script lang="ts">
  import { onMount } from 'svelte';
  import * as echarts from 'echarts';

  export let ports: any[] = [];
  export let stats: any[] = [];

  let chartContainer: HTMLDivElement;
  let chart: echarts.ECharts;

  $: protocolData = {
    tcp: ports.filter(p => p.protocol === 'tcp').length,
    udp: ports.filter(p => p.protocol === 'udp').length
  };

  $: stateData = ports.reduce((acc, port) => {
    acc[port.state] = (acc[port.state] || 0) + 1;
    return acc;
  }, {});

  $: processCount = new Set(ports.map(p => p.pid)).size;

  onMount(() => {
    chart = echarts.init(chartContainer, 'dark');
    updateChart();

    window.addEventListener('resize', () => {
      chart.resize();
    });

    return () => {
      chart.dispose();
    };
  });

  function updateChart() {
    const option = {
      backgroundColor: 'transparent',
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        bottom: '0%',
        textStyle: {
          color: '#9ca3af'
        }
      },
      series: [
        {
          name: '协议分布',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#111827',
            borderWidth: 2
          },
          label: {
            show: true,
            position: 'inside',
            formatter: '{b}\n{c}',
            color: '#fff'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 14,
              fontWeight: 'bold'
            },
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          },
          data: [
            { value: protocolData.tcp, name: 'TCP', itemStyle: { color: '#0ea5e9' } },
            { value: protocolData.udp, name: 'UDP', itemStyle: { color: '#10b981' } }
          ]
        }
      ]
    };

    chart.setOption(option);
  }

  $: {
    if (chart) {
      updateChart();
    }
  }
</script>

<div class="bg-dark-800 rounded-xl border border-dark-700 p-4 animate-fade-in">
  <h2 class="text-lg font-semibold mb-4">统计信息</h2>

  <div bind:this={chartContainer} class="h-64 mb-4"></div>

  <div class="space-y-3">
    <div class="flex justify-between items-center p-3 bg-dark-700/50 rounded-lg">
      <div>
        <div class="text-sm text-dark-400">总端口数</div>
        <div class="text-2xl font-bold text-primary-400">{ports.length}</div>
      </div>
      <div class="text-3xl">🔌</div>
    </div>

    <div class="flex justify-between items-center p-3 bg-dark-700/50 rounded-lg">
      <div>
        <div class="text-sm text-dark-400">TCP 端口</div>
        <div class="text-xl font-semibold text-blue-400">{protocolData.tcp}</div>
      </div>
      <div class="text-2xl">🔷</div>
    </div>

    <div class="flex justify-between items-center p-3 bg-dark-700/50 rounded-lg">
      <div>
        <div class="text-sm text-dark-400">UDP 端口</div>
        <div class="text-xl font-semibold text-green-400">{protocolData.udp}</div>
      </div>
      <div class="text-2xl">🟢</div>
    </div>

    <div class="flex justify-between items-center p-3 bg-dark-700/50 rounded-lg">
      <div>
        <div class="text-sm text-dark-400">活跃进程</div>
        <div class="text-xl font-semibold text-purple-400">
          {processCount}
        </div>
      </div>
      <div class="text-2xl">⚡</div>
    </div>

    {#if stats && stats.length > 0}
      <div class="mt-4 p-4 bg-dark-700/30 rounded-lg">
        <h3 class="text-sm font-semibold text-dark-300 mb-2">Top 使用端口</h3>
        <div class="space-y-2">
          {#each stats.slice(0, 5) as stat}
            <div class="flex justify-between items-center text-sm">
              <span class="text-primary-400 font-mono">:{stat.port}</span>
              <span class="text-dark-400">{stat.usageCount || 0} 次</span>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </div>
</div>
