<template>
  <div class="chart-wrapper">
    <!-- Controles superiores -->
    <div class="chart-controls">
      <div class="controls-left">
        <button @click="toggleChartType" class="btn">
          {{ chartType === 'candlestick' ? '📊 Candlestick' : '📈 Line' }}
        </button>
        <button @click="addVolume" class="btn" v-if="!showVolume">
          📊 Mostrar Volumen
        </button>
        <button @click="removeVolume" class="btn" v-else>
          ❌ Ocultar Volumen
        </button>
        <button @click="toggleMA" class="btn">
          {{ showMA ? '✅ MA' : '➕ MA' }}
        </button>
      </div>
      
      <div class="controls-right">
        <select v-model="timeframe" @change="onTimeframeChange" class="select">
          <option value="1m">1 Minuto</option>
          <option value="5m">5 Minutos</option>
          <option value="15m">15 Minutos</option>
          <option value="1h">1 Hora</option>
          <option value="1d">1 Día</option>
        </select>
        
        <button @click="zoomIn" class="btn">🔍+</button>
        <button @click="zoomOut" class="btn">🔍-</button>
        <button @click="resetChart" class="btn">🔄 Reset</button>
        <button @click="takeScreenshot" class="btn">📸 Captura</button>
      </div>
    </div>

    <!-- Indicadores de precio -->
    <div class="price-info" v-if="currentPrice">
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
    </div>

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
import { ref, onMounted, onUnmounted, watch, computed } from 'vue';
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
  MouseEventParams
} from 'lightweight-charts';

interface Props {
  data?: CandlestickData[];
  width?: number;
  height?: number;
}

const props = withDefaults(defineProps<Props>(), {
  data: () => [],
  width: 800,
  height: 500
});

const chartContainer = ref<HTMLElement | null>(null);
const chartType = ref<'candlestick' | 'line'>('candlestick');
const showVolume = ref(false);
const showMA = ref(false);
const timeframe = ref('1d');
const currentPrice = ref<CandlestickData | null>(null);
const crosshairData = ref<{ time: string; price: number | null } | null>(null);

let chart: IChartApi | null = null;
let mainSeries: ISeriesApi<'Candlestick'> | ISeriesApi<'Line'> | null = null;
let volumeSeries: ISeriesApi<'Histogram'> | null = null;
let maSeries: ISeriesApi<'Line'> | null = null;

// Datos de ejemplo con volumen
const defaultData: CandlestickData[] = [
  { time: '2024-01-01', open: 100, high: 110, low: 95, close: 105 },
  { time: '2024-01-02', open: 105, high: 115, low: 100, close: 112 },
  { time: '2024-01-03', open: 112, high: 120, low: 108, close: 115 },
  { time: '2024-01-04', open: 115, high: 118, low: 110, close: 113 },
  { time: '2024-01-05', open: 113, high: 125, low: 112, close: 122 },
  { time: '2024-01-08', open: 122, high: 130, low: 120, close: 128 },
  { time: '2024-01-09', open: 128, high: 132, low: 125, close: 127 },
  { time: '2024-01-10', open: 127, high: 135, low: 126, close: 133 },
  { time: '2024-01-11', open: 133, high: 140, low: 130, close: 138 },
  { time: '2024-01-12', open: 138, high: 142, low: 135, close: 137 },
];

const volumeData: HistogramData[] = [
  { time: '2024-01-01', value: 1000000, color: '#26a69a80' },
  { time: '2024-01-02', value: 1500000, color: '#26a69a80' },
  { time: '2024-01-03', value: 1200000, color: '#26a69a80' },
  { time: '2024-01-04', value: 900000, color: '#ef535080' },
  { time: '2024-01-05', value: 1800000, color: '#26a69a80' },
  { time: '2024-01-08', value: 2000000, color: '#26a69a80' },
  { time: '2024-01-09', value: 1100000, color: '#ef535080' },
  { time: '2024-01-10', value: 1600000, color: '#26a69a80' },
  { time: '2024-01-11', value: 1700000, color: '#26a69a80' },
  { time: '2024-01-12', value: 1300000, color: '#ef535080' },
];

// Calcular cambio de precio
const priceChange = computed(() => {
  if (!currentPrice.value) return null;
  const change = ((currentPrice.value.close - currentPrice.value.open) / currentPrice.value.open) * 100;
  return change;
});

// Calcular media móvil
const calculateMA = (data: CandlestickData[], period: number = 7): LineData[] => {
  const ma: LineData[] = [];
  for (let i = period - 1; i < data.length; i++) {
    let sum = 0;
    for (let j = 0; j < period; j++) {
      sum += data[i - j].close;
    }
    ma.push({
      time: data[i].time,
      value: sum / period
    });
  }
  return ma;
};

// Inicializar gráfico
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
    },
    rightPriceScale: {
      borderColor: '#cccccc',
    },
  });

  createMainSeries();

  const dataToUse = props.data.length > 0 ? props.data : defaultData;
  mainSeries?.setData(dataToUse);
  
  if (dataToUse.length > 0) {
    currentPrice.value = dataToUse[dataToUse.length - 1];
  }

  // Evento de crosshair
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

// Crear serie principal
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
  } else {
    mainSeries = chart.addSeries(LineSeries, {
      color: '#2962FF',
      lineWidth: 2,
    });
  }
};

// Cambiar tipo de gráfico
const toggleChartType = () => {
  if (!chart || !mainSeries) return;

  const currentData = props.data.length > 0 ? props.data : defaultData;
  
  chart.removeSeries(mainSeries);
  chartType.value = chartType.value === 'candlestick' ? 'line' : 'candlestick';
  createMainSeries();

  if (chartType.value === 'line') {
    const lineData: LineData[] = currentData.map(d => ({
      time: d.time,
      value: d.close
    }));
    (mainSeries as ISeriesApi<'Line'>).setData(lineData);
  } else {
    mainSeries?.setData(currentData);
  }

  chart.timeScale().fitContent();
};

// Agregar volumen
const addVolume = () => {
  if (!chart || showVolume.value) return;

  volumeSeries = chart.addSeries(HistogramSeries, {
    color: '#26a69a',
    priceFormat: {
      type: 'volume',
    },
    priceScaleId: '',
  });

  volumeSeries.priceScale().applyOptions({
    scaleMargins: {
      top: 0.8,
      bottom: 0,
    },
  });

  volumeSeries.setData(volumeData);
  showVolume.value = true;
};

// Remover volumen
const removeVolume = () => {
  if (!chart || !volumeSeries) return;
  chart.removeSeries(volumeSeries);
  volumeSeries = null;
  showVolume.value = false;
};

// Toggle media móvil
const toggleMA = () => {
  if (!chart) return;

  if (showMA.value && maSeries) {
    chart.removeSeries(maSeries);
    maSeries = null;
    showMA.value = false;
  } else {
    const currentData = props.data.length > 0 ? props.data : defaultData;
    const maData = calculateMA(currentData, 7);
    
    maSeries = chart.addSeries(LineSeries, {
      color: 'rgba(255, 152, 0, 1)',
      lineWidth: 2,
    });
    
    maSeries.setData(maData);
    showMA.value = true;
  }
};

// Zoom
const zoomIn = () => {
  chart?.timeScale().scrollToPosition(2, true);
};

const zoomOut = () => {
  chart?.timeScale().scrollToPosition(-2, true);
};

// Reset
const resetChart = () => {
  chart?.timeScale().fitContent();
};

// Captura de pantalla
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

// Cambio de timeframe
const onTimeframeChange = () => {
  console.log('Timeframe cambiado a:', timeframe.value);
  // Aquí implementarías la lógica para cargar datos del nuevo timeframe
};

// Resize
const handleResize = () => {
  if (chart && chartContainer.value) {
    chart.applyOptions({
      width: chartContainer.value.clientWidth,
    });
  }
};

// Watch data changes
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
  }
});

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