<template>
  <div class="chart-view">
    <!-- 图表类型选择 -->
    <div v-if="showTypeSelector" class="chart-toolbar">
      <el-radio-group v-model="chartType" size="small">
        <el-radio-button value="bar">柱状图</el-radio-button>
        <el-radio-button value="line">折线图</el-radio-button>
        <el-radio-button value="pie">饼图</el-radio-button>
        <el-radio-button value="radar">雷达图</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 加载状态 -->
    <el-skeleton v-if="loading || !echartsReady" :rows="8" animated />

    <!-- 图表容器 -->
    <div v-show="!loading && echartsReady" ref="chartRef" class="chart-container" :style="{ height: height }" />

    <!-- 图例/统计信息 -->
    <div v-if="showLegend && legendData.length > 0" class="chart-legend">
      <div
        v-for="(item, index) in legendData"
        :key="index"
        class="legend-item"
      >
        <span
          class="legend-color"
          :style="{ background: colors[index % colors.length] }"
        />
        <span class="legend-label">{{ item.name }}</span>
        <span class="legend-value">{{ item.value }}</span>
        <span v-if="'percentage' in item" class="legend-percent">{{ item.percentage }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, shallowRef, nextTick } from 'vue'

// 图表数据项
export interface ChartDataItem {
  name: string
  value: number
  percentage?: string
}

// 图表系列配置
export interface ChartSeries {
  name: string
  data: number[]
  type?: string
  color?: string
}

interface Props {
  // 基础配置
  height?: string
  loading?: boolean
  showTypeSelector?: boolean
  showLegend?: boolean
  
  // 数据
  data?: ChartDataItem[]
  series?: ChartSeries[]
  categories?: string[]
  
  // 图表配置
  title?: string
  defaultType?: 'bar' | 'line' | 'pie' | 'radar'
  colors?: string[]
  
  // 柱状图/折线图配置
  xAxisLabel?: string
  yAxisLabel?: string
  showDataLabels?: boolean
  
  // 饼图配置
  showPieLabel?: boolean
  pieRadius?: string | [string, string]
  
  // 雷达图配置
  radarIndicators?: { name: string; max: number }[]
}

const props = withDefaults(defineProps<Props>(), {
  height: '400px',
  loading: false,
  showTypeSelector: true,
  showLegend: true,
  defaultType: 'bar',
  colors: () => [
    '#409eff', '#67c23a', '#e6a23c', '#f56c6c',
    '#909399', '#b37feb', '#ff85c0', '#36cfc9',
  ],
  xAxisLabel: '',
  yAxisLabel: '',
  showDataLabels: true,
  showPieLabel: true,
  pieRadius: '60%',
})

const emit = defineEmits<{
  chartClick: [params: any]
}>()

const chartRef = shallowRef<HTMLElement>()
const chartInstance = shallowRef<any>()
const chartType = ref(props.defaultType)
const echartsReady = ref(false)

// echarts 实例（延迟加载）
let echartsModule: any = null

// 图例数据
const legendData = computed(() => {
  if (props.data) return props.data
  if (props.series && props.series.length > 0) {
    return props.series.map(s => ({
      name: s.name,
      value: s.data.reduce((a, b) => a + b, 0),
    }))
  }
  return []
})

// 生成图表配置
function getChartOption(): any {
  const baseOption: any = {
    color: props.colors,
    tooltip: {
      trigger: chartType.value === 'pie' ? 'item' : 'axis',
    },
    title: props.title ? { text: props.title, left: 'center' } : undefined,
  }

  switch (chartType.value) {
    case 'bar':
    case 'line':
      return {
        ...baseOption,
        xAxis: {
          type: 'category',
          data: props.categories || props.data?.map(d => d.name) || [],
          name: props.xAxisLabel,
        },
        yAxis: {
          type: 'value',
          name: props.yAxisLabel,
        },
        series: props.series || [{
          name: '数据',
          type: chartType.value,
          data: props.data?.map(d => d.value) || [],
          label: props.showDataLabels ? { show: true, position: 'top' } : undefined,
          smooth: chartType.value === 'line',
        }],
      }

    case 'pie':
      return {
        ...baseOption,
        series: [{
          type: 'pie',
          radius: props.pieRadius,
          data: props.data?.map(d => ({ name: d.name, value: d.value })) || [],
          label: props.showPieLabel ? { show: true, formatter: '{b}: {d}%' } : { show: false },
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)',
            },
          },
        }],
      }

    case 'radar':
      return {
        ...baseOption,
        radar: {
          indicator: props.radarIndicators || props.data?.map(d => ({
            name: d.name,
            max: Math.max(d.value * 1.2, 100),
          })) || [],
        },
        series: [{
          type: 'radar',
          data: props.series?.map((s: any) => ({
            name: s.name,
            value: s.data,
          })) || [{
            name: '数据',
            value: props.data?.map(d => d.value) || [],
          }],
        }],
      }

    default:
      return baseOption
  }
}

// 初始化图表
async function initChart() {
  if (!chartRef.value) return

  // 动态加载 echarts
  if (!echartsModule) {
    try {
      echartsModule = await import('echarts')
      echartsReady.value = true
    } catch (e) {
      console.warn('echarts not installed, chart view disabled')
      return
    }
  }

  await nextTick()

  if (chartInstance.value) {
    chartInstance.value.dispose()
  }

  chartInstance.value = echartsModule.init(chartRef.value)
  chartInstance.value.setOption(getChartOption())

  chartInstance.value.on('click', (params: any) => {
    emit('chartClick', params)
  })
}

// 更新图表
function updateChart() {
  if (chartInstance.value && echartsModule) {
    chartInstance.value.setOption(getChartOption(), true)
  }
}

// 监听变化
watch([() => props.data, () => props.series, chartType], updateChart, { deep: true })

// 窗口大小变化
function handleResize() {
  chartInstance.value?.resize()
}

onMounted(() => {
  initChart()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance.value?.dispose()
})

// 暴露方法
defineExpose({
  chartInstance,
  chartType,
  updateChart,
  resize: () => chartInstance.value?.resize(),
})
</script>

<style scoped>
.chart-view {
  width: 100%;
}

.chart-toolbar {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

.chart-container {
  width: 100%;
}

.chart-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.legend-label {
  color: #606266;
  font-size: 14px;
}

.legend-value {
  color: #303133;
  font-weight: 600;
  font-size: 14px;
}

.legend-percent {
  color: #909399;
  font-size: 12px;
}
</style>
