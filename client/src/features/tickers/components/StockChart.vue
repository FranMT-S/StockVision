<template>
  <div class="chart-wrapper">
    <!-- Controles superiores -->
    <div class="chart-controls">
      <div class="controls-left">

        <button @click="toggleChartType('candlestick')" class="btn py-2 px-4">
          <v-icon icon="mdi-chart-waterfall" size="18" /> 
        </button>
        
        <button @click="toggleChartType('line')" class="btn py-2 px-4">
          <v-icon icon="mdi-chart-line" size="18" /> 
        </button>
      </div>
      
      <div class="controls-right">
        <div class="tw-flex tw-flex-row tw-items-center tw-gap-1 timeframe">
          <button 
          :class="{ 'tw-text-primary tw-bg-[#f5f5f5] tw-rounded': timeframe === option }"
          v-for="option in timeframeOptions" :key="option" @click="timeframe = option" class="tw-text-[#717171] tw-p-2 tw-py-1 tw-text-[12px] tw-font-medium tw-pa-0 tw-capitalize">{{ option }}</button>
        </div>
        <button @click="takeScreenshot" class="btn py-2 px-4">
          <v-icon icon="mdi-camera" size="18" /> 
        </button>
      </div>
    </div>

    <!-- Indicadores de precio -->
    <!-- <div class="price-info" v-if="currentPrice">
      <div class="price-item">
        <span class="label">Open:</span>
        <span class="value">{{ currentPrice.open.toFixed(2) }}</span>
      </div>
      <div class="price-item">
        <span class="label">High:</span>
        <span class="value high">{{ currentPrice.high.toFixed(2) }}</span>
      </div>
      <div class="price-item">
        <span class="label">Low:</span>
        <span class="value low">{{ currentPrice.low.toFixed(2) }}</span>
      </div>
      <div class="price-item">
        <span class="label">Close:</span>
        <span :class="['value', currentPrice.close >= currentPrice.open ? 'up' : 'down']">
          {{ currentPrice.close.toFixed(2) }}
        </span>
      </div>
      <div class="price-item" v-if="priceChange">
        <span class="label">Change:</span>
        <span :class="['value', priceChange >= 0 ? 'up' : 'down']">
          {{ priceChange >= 0 ? '+' : '' }}{{ priceChange.toFixed(2) }}%
        </span>
      </div>
    </div> -->

    <!-- Contenedor del gráfico -->
    <div ref="chartContainer" class="chart-container"></div>

    <!-- Crosshair info -->
    <div class="crosshair-info" v-if="crosshairData">
      <div class="info-row">
        <span>Tiempo: {{ crosshairData.time }}</span>
      </div>
      <div class="info-row">
        <span>Precio: {{ crosshairData.price?.toFixed(2) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed, ComputedRef } from 'vue';
import { 
  createChart, 
  IChartApi, 
  ISeriesApi, 
  CandlestickData,
  CandlestickSeries,
  LineSeries,
  HistogramSeries,
  LineData,
  HistogramData,
  MouseEventParams,
  Time
} from 'lightweight-charts';
import { HistoricalPrice,CompanyNew} from '@/shared/models/recomendations';

interface Props {
  historicalData?: HistoricalPrice[];
  newData?: CompanyNew[];
  data?: CandlestickData[];
  width?: number;
  height?: number;
}

enum Timeframe {
  '1D' = '1D',
  '1W' = '1W',
  '1M' = '1M',
  '6M' = '6M',
  '1Y' = '1Y',
  'All' = 'All'
}

const timeframeOptions = [ Timeframe['1D'], Timeframe['1W'], Timeframe['1M'],  Timeframe['6M'], Timeframe['1Y'], Timeframe['All'] ];

const props = withDefaults(defineProps<Props>(), {
  historicalData: () => [],
  newData: () => [],
  data: () => [],
  width: 800,
  height: 500
});

const chartContainer = ref<HTMLElement | null>(null);
const chartType = ref<'candlestick' | 'line'>('candlestick');
const timeframe = ref(Timeframe['All']);
const currentPrice = ref<CandlestickData | null>(null);
const crosshairData = ref<{ time: string; price: number | null } | null>(null);
const visibleRange = ref<{ from: Time; to: Time } | null>(null);

let chart: IChartApi | null = null;
let mainSeries: ISeriesApi<'Candlestick'> | ISeriesApi<'Line'> | null = null;
let volumeSeries: ISeriesApi<'Histogram'> | null = null;



const historicalDataProcessed: ComputedRef<CandlestickData[]> = computed(() => {
  if (!props.historicalData) return [];

  return props.historicalData.map(item => {
    return {
      time: item.date,
      open: item.open,
      high: item.high,
      low: item.low,
      close: item.close,
      date: new Date(item.date)
    }
  }).sort((a, b) => a.date.getTime() - b.date.getTime())
}); 

const chartData: ComputedRef<{seriesData: CandlestickData[], volumeData: HistogramData[]}> = computed(() => {
  if (!props.historicalData) return {seriesData: [], volumeData: []};
  const orderedData = props.historicalData.map((item) => {
    return {
      ...item,
      time: item.date,
      date: new Date(item.date)
    }
  }).sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
  
  const seriesData: CandlestickData[] = [];
  const volumeData: HistogramData[] = [];
  
 

  for (let index = 0; index < orderedData.length; index++) {
    const element = orderedData[index];
    seriesData.push({
      time: element.time,
      open: element.open,
      high: element.high,
      low: element.low,
      close: element.close,
      color: element.close > element.open ? '#26a69a' : '#ef5350'
    })

    volumeData.push({
      time: element.time,
      value: element.volume,
      color: element.close > element.open ? '#26a69a80' : '#ef535080'
    })
  }


  return { seriesData, volumeData }
}); 



const priceChange = computed(() => {
  if (!currentPrice.value) return null;
  const change = ((currentPrice.value.close - currentPrice.value.open) / currentPrice.value.open) * 100;
  return change;
});


const initChart = () => {
  if (!chartContainer.value) return;

  chart = createChart(chartContainer.value, {
    width: props.width,
    height: props.height,
    layout: {
      background: { color: '#ffffff' },
      textColor: '#333',
    },
    grid: {
      vertLines: { color: '#e1e1e1' },
      horzLines: { color: '#e1e1e1' },
    },
    crosshair: {
      mode: 1,
    },
    timeScale: {
      borderColor: '#cccccc',
      timeVisible: true,
      secondsVisible: false,
      fixRightEdge: true,
      fixLeftEdge: true,
      
    },
    rightPriceScale: {
      borderColor: '#cccccc',
    },
  });

  createMainSeries();

  const dataToUse = historicalDataProcessed.value;
  mainSeries?.setData(dataToUse);
  
  if (dataToUse.length > 0) {
    currentPrice.value = dataToUse[dataToUse.length - 1];
  }

  chart.subscribeCrosshairMove((param: MouseEventParams) => {
    if (param.time && param.seriesData.get(mainSeries!)) {
      const data = param.seriesData.get(mainSeries!) as any;
      crosshairData.value = {
        time: param.time as string,
        price: data?.close || data?.value || null
      };
    } else {
      crosshairData.value = null;
    }
  });

  chart.timeScale().fitContent();
};


const createMainSeries = () => {
  if (!chart) return;

  if (chartType.value === 'candlestick') {
    mainSeries = chart.addSeries(CandlestickSeries, {
      upColor: '#26a69a',
      downColor: '#ef5350',
      borderVisible: false,
      wickUpColor: '#26a69a',
      wickDownColor: '#ef5350',
    });

    addVolume()
    
    
  } else {
    mainSeries = chart.addSeries(LineSeries, {
      color: '#2962FF',
      lineWidth: 2,
    });
    
  }
};

const toggleChartType = (_chartType: 'candlestick' | 'line') => {
  if (!chart || !mainSeries) return;

  if(chartType.value === _chartType) return;

  const currentData = chartData.value.seriesData;
  
  chart.removeSeries(mainSeries);
  chartType.value = _chartType;
  createMainSeries();

  if (chartType.value === 'line') {
    ;
    const lineData: LineData[] = currentData.map((d) => ({
      time: d.time,
      value: d.close
    }));
    (mainSeries as ISeriesApi<'Line'>).setData(lineData);
    removeVolume()
  } else {
    mainSeries?.setData(currentData);
  }
  setTimeScaleView()
  chart.timeScale().fitContent();
};

const addVolume = () => {
  if (!chart) return;

  volumeSeries = chart.addSeries(HistogramSeries, {
    priceScaleId: '', // Usa una escala independiente
    priceFormat: { type: 'volume' },
    color: '#26a69a'
  });

  volumeSeries.priceScale().applyOptions({
    scaleMargins: {
      top: 0.8,
      bottom: 0,
    },
  });

  chart.priceScale('').applyOptions({
    scaleMargins: {
      top: 0.9,   // 80% del espacio superior para las velas
      bottom: 0,  // 20% inferior para el volumen
    },
  });
 
  volumeSeries.setData(chartData.value.volumeData);
};

const removeVolume = () => {
  if (!chart || !volumeSeries) return;
  chart.removeSeries(volumeSeries);
  volumeSeries = null;

  chart.priceScale('').applyOptions({
    scaleMargins: {
      top: 1,  
      bottom: 0, 
    },
  });
};

const takeScreenshot = () => {
  if (!chart) return;
  const canvas = chartContainer.value?.querySelector('canvas');
  if (canvas) {
    canvas.toBlob((blob) => {
      if (blob) {
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `chart-${Date.now()}.png`;
        a.click();
        URL.revokeObjectURL(url);
      }
    });
  }
};

watch(timeframe, () => {
  setTimeScaleView()
});

const setTimeScaleView = () => {
 if (!chart || historicalDataProcessed.value.length === 0) return;
  const data = historicalDataProcessed.value;
  const lastIndex = data.length - 1;
  

  let fromIndex = 0;

  switch (timeframe.value) {
    case Timeframe['1D']:
      fromIndex = Math.max(0, lastIndex - 1);
      break;
    case Timeframe['1W']:
      fromIndex = Math.max(0, lastIndex - 7);
      break;
    case Timeframe['1M']:
      fromIndex = Math.max(0, lastIndex - 30);
      break;
    case Timeframe['6M']:
      fromIndex = Math.max(0, lastIndex - 180);
      break;
    case Timeframe['1Y']:
      fromIndex = Math.max(0, lastIndex - 365);
      break;
    case Timeframe['All']:
      fromIndex = 0;
      break;
  }




  visibleRange.value = { from: data[fromIndex].time, to: data[lastIndex].time };



  GoToLastStock(); 
}


const handleResize = () => {
  if (chart && chartContainer.value) {
    chart.applyOptions({
      width: chartContainer.value.clientWidth,
    });
  }
};

watch(visibleRange, () => {
  if (chart) {
    if (visibleRange.value) {
      chart.timeScale().setVisibleRange(visibleRange.value);
    } else {
      chart.timeScale().scrollToRealTime();
    }
  }
});

watch(() => props.data, (newData) => {
  if (mainSeries && newData.length > 0) {
    if (chartType.value === 'line') {
      const lineData: LineData[] = newData.map(d => ({
        time: d.time,
        value: d.close
      }));
      (mainSeries as ISeriesApi<'Line'>).setData(lineData);
    } else {
      mainSeries.setData(newData);
    }
    currentPrice.value = newData[newData.length - 1];
    chart?.timeScale().fitContent();
    GoToLastStock();
  }
});

const GoToLastStock = () => {
  if (!chart) return;
    chart.timeScale().scrollToRealTime();
}

// Lifecycle
onMounted(() => {
  initChart();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  if (chart) {
    chart.remove();
  }
});
</script>

<style scoped>
.timeframe button,
.timeframe button:hover,
.timeframe button:active {
  outline: none;
}

.chart-wrapper {
  width: 100%;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  padding: 16px;
}

.chart-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  gap: 8px;
  flex-wrap: wrap;
}

.controls-left,
.controls-right {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.btn {
  padding: 8px 12px;
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.btn:hover {
  background: #e0e0e0;
}

.btn:active {
  transform: scale(0.98);
}

.select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  font-size: 13px;
}

.price-info {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 4px;
  flex-wrap: wrap;
}

.price-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label {
  font-size: 11px;
  color: #666;
  text-transform: uppercase;
}

.value {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.value.up {
  color: #26a69a;
}

.value.down {
  color: #ef5350;
}

.value.high {
  color: #26a69a;
}

.value.low {
  color: #ef5350;
}

.chart-container {
  width: 100%;
  height: 100%;
  min-height: 400px;
  border-radius: 4px;
  overflow: hidden;
}

.crosshair-info {
  position: absolute;
  top: 80px;
  left: 32px;
  background: rgba(255, 255, 255, 0.95);
  padding: 8px 12px;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  font-size: 12px;
  pointer-events: none;
}

.info-row {
  margin-bottom: 4px;
}

.info-row:last-child {
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .chart-controls {
    flex-direction: column;
    align-items: stretch;
  }
  
  .controls-left,
  .controls-right {
    justify-content: center;
  }
  
  .price-info {
    justify-content: center;
  }
}
</style>