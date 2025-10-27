<template>
  <div class="chart-wrapper">
    <div v-if="isGeneratingScreenshot"  data-html2canvas-ignore="true" class="tw-z-50 tw-fixed tw-inset-0 tw-bg-black/30 tw-text-white tw-flex tw-items-center tw-justify-center">
      <p class="tw-text-[18px] tw-font-medium tw-text-blue">Generating Screenshot...</p>
    </div>
    <!-- Controles superiores -->
    <div class="chart-controls">
      <div class="controls-left">

        <v-btn :color="chartType === ChartType['candlestick'] ? 'primary' : 'secondary'" variant="tonal" @click="toggleChartType(ChartType['candlestick'])" class=" py-2 px-4">
          <v-icon icon="mdi-chart-waterfall" size="18" /> 
        </v-btn>
        
        <v-btn :color="chartType === ChartType['area'] ? 'primary' : 'secondary'" variant="tonal" @click="toggleChartType(ChartType['area'])" class=" py-2 px-4">
          <v-icon icon="mdi-chart-line" size="18" /> 
        </v-btn>


        <v-btn 
          v-if="chartType === ChartType['candlestick']"
          :color="isShowPredict ? 'primary' : 'secondary'" variant="tonal" @click="toggleShowPredict()" class=" py-2 px-4">
          <v-icon icon="mdi-eye" size="18" /> 
          <span v-if="!isShowPredict">{{isPredictLoading ? 'Loading...' : 'Get Vision'}}</span>
          <span v-else>Hide Vision</span>
        </v-btn>
      </div>
      
      <div class="controls-right">
        <div class="tw-flex tw-flex-row tw-items-center tw-gap-1 timeframe">
          <button 
          :class="{ 'tw-text-primary tw-bg-[#f5f5f5] tw-rounded': timeframe === option.value }"
          v-for="option in timeframeOptions" :key="option.value" @click="timeframe = option.value" class="tw-text-[#717171] tw-p-2 tw-py-1 tw-text-[12px] tw-font-medium tw-pa-0 tw-capitalize">{{ option.label }}</button>
        </div>  
        <button @click="takeScreenshot" class="btn py-2 px-4">
          <v-icon icon="mdi-camera" size="18" /> 
        </button>
      </div>
    </div>

    <!-- Chard container -->
    <div ref="chartContainer" class="chart-container"/>

 
    <!-- crosshair tooltip -->
    <FloatingTooltip :visible="selectedCrosshairData !== null && !isTouchDevice">
      <div class="tw-bg-black/70 tw-rounded tw-p-2 tw-text-white">
        <CrossHairDetails :data="selectedCrosshairData!" />
      </div>  
    </FloatingTooltip>
  </div>
</template>


<script setup lang="ts">
import FloatingTooltip from '@/shared/components/FloatingTooltip.vue';
import { useDebounce } from '@/shared/composables/useDebounce';
import { Timeframe } from '@/shared/enums/timeFrame';
import { CompanyNew, HistoricalPrice, StockHLOC } from '@/shared/models/recomendations';
import html2canvas from 'html2canvas';
import {
  AreaData,
  AreaSeries,
  AreaSeriesPartialOptions,
  CandlestickData,
  CandlestickSeries,
  CandlestickSeriesPartialOptions,
  createChart,
  CrosshairMode,
  HistogramData,
  HistogramSeries,
  HistogramSeriesPartialOptions,
  IChartApi,
  ISeriesApi,
  LineStyle,
  MouseEventParams,
  SeriesMarker,
  Time
} from 'lightweight-charts';
import { computed, ComputedRef, onMounted, onUnmounted, ref, shallowReactive, watch } from 'vue';
import CrossHairDetails from './CrossHairDetails.vue';

interface Props {
  ticker: string;
  historicalData?: HistoricalPrice[];
  width?: number;
  height?: number;
  initialInterval?: Timeframe;
  predictNextWeek?: HistoricalPrice[];
  predictError?: string | null;
  isPredictLoading?: boolean;
  isTouchDevice?: boolean;
}

type TSeriesType = 'Candlestick' | 'Line' | 'Area' | 'Histogram' | 'Predict' | 'PredictHistogram';
type CandleChardColor = {
  upColor: string,
  downColor: string,
  upVolumeColor: string,
  downVolumeColor: string
}

type CharData = {
  mainChartData: CandlestickData[], 
  volumeChartData: HistogramData[],
  areaChartData: AreaData[]
  stockHLOC: StockHLOC[]
}

enum Scales {
  'volume' = 'volume',
  'right' = 'right'
}

enum ChartType {
  'candlestick' = 'candlestick',
  'line' = 'line',
  'area' = 'area'
}

const props = withDefaults(defineProps<Props>(), {
  historicalData: () => [],
  companyNews: () => [],
  width: 800,
  height: 500,
  initialInterval: Timeframe['1M'],
  predictNextWeek: () => [],
  predictError: () => "",
  isPredictLoading: () => false
});

const emit = defineEmits(['update:timeframe','update:predict'])
const {debounced} = useDebounce()
const chartContainer = ref<HTMLElement | null>(null);
const chartType = ref<ChartType>(ChartType['candlestick']);
const timeframe = ref(props.initialInterval);
const visibleRange = ref<{ from: Time; to: Time } | null>(null);
const isShowPredict = ref<boolean>(false);
const isGeneratingScreenshot = ref<boolean>(false);
const selectedCrosshairData = ref<StockHLOC | null>(null);

const timeframeOptions =ref([
    { value: Timeframe['1M'], label: '1M' },
    { value: Timeframe['3M'], label: '3M' },
    { value: Timeframe['6M'], label: '6M' },
    { value: Timeframe['1Y'], label: '1Y' },
    { value: Timeframe['All'], label: 'All' }
  ]);


// styles of the elements in charts
const CandleSeriesTheme: CandlestickSeriesPartialOptions = {
  priceScaleId: Scales.right,
  upColor: '#089981',
  downColor: '#f23645',
  wickUpColor: '#089981',
  wickDownColor: '#f23645',
  borderVisible: false,
}

const CandlePredictSeriesTheme: CandlestickSeriesPartialOptions = {
  priceScaleId: Scales.right,
  upColor: '#1976d2',      
  downColor: '#ff9800',     
  borderColor: '#1976d2',
  wickUpColor: '#1976d2',
  wickDownColor: '#ff9800',
  borderVisible: false,   
}

const VolumeSeriesTheme: HistogramSeriesPartialOptions = {
  priceScaleId: Scales.volume, 
  color: '#2962FF',   
};

const PredictVolumeSeriesTheme: HistogramSeriesPartialOptions = {
  priceScaleId: Scales.volume, 
  color: '#ff9800',
}

const LinearAreaSeriesTheme: AreaSeriesPartialOptions = {
  priceScaleId: Scales.right,
  lineColor: '#2962FF', 
  topColor: '#2962FF', 
  bottomColor: 'rgba(41, 98, 255, 0.28)' 
}

let chart: IChartApi | null = null;

const series = shallowReactive<Record<TSeriesType, ISeriesApi<any> | null>>({
  Candlestick: null,
  Line: null,
  Area: null,
  Histogram: null,
  Predict: null,
  PredictHistogram: null
});

const predictChartData: ComputedRef<{
  mainChartData: CandlestickData[], 
  volumeChartData: HistogramData[],
  stockHLOC: StockHLOC[]
}> = computed(() => {
  if (!props.predictNextWeek) return {mainChartData: [], volumeChartData: [], stockHLOC: []};
  const filteredData = props.predictNextWeek.filter((stock) => {
    const stockDate = new Date(stock.date);
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    stockDate.setHours(0, 0, 0, 0);
    tomorrow.setHours(0, 0, 0, 0);

    return stockDate >= tomorrow;
  });

  const orderedData = creteaStockHLOC(filteredData);
  return createCandleCharData(orderedData,{
    upColor: '#1976d2',      
    downColor: '#ff9800',
    upVolumeColor: '#1976d280',
    downVolumeColor: '#ff980080'
  });
}); 

const chartData: ComputedRef<{
  mainChartData: CandlestickData[], 
  volumeChartData: HistogramData[]
  areaChartData: AreaData[]
  stockHLOC: StockHLOC[]
}> = computed(() => {
  if (!props.historicalData) 
    return {mainChartData: [], volumeChartData: [], areaChartData: [], stockHLOC: []};
  
  const orderedData = creteaStockHLOC(props.historicalData);
  return createCandleCharData(orderedData);
}); 

const creteaStockHLOC = (data: HistoricalPrice[]) : StockHLOC[] => {
  if (!data) return [];
  return data.map(item => {
    return {
      ...item,
      time: item.date,
      date: new Date(item.date)
    }
  }).sort((a, b) => a.date.getTime() - b.date.getTime()) 
}

const createCandleCharData = (data: StockHLOC[],color?:CandleChardColor) : CharData => {
  if (!data) return {mainChartData: [], volumeChartData: [], areaChartData: [], stockHLOC: []};
  const mainChartData: CandlestickData[] = [];
  const volumeChartData: HistogramData[] = [];
  const areaChartData: AreaData[] = [];

  if (!color) {
    color = {
      upColor: '#26a69a',
      downColor: '#ef5350',
      upVolumeColor: '#26a69a80',
      downVolumeColor: '#ef535080'
    }
  }

  // default colors if not send
  const {upColor = '#26a69a', downColor = '#ef5350', upVolumeColor = '#26a69a80', downVolumeColor = '#ef535080'} = color;

  for (let index = 0; index < data.length; index++) {
    const element = data[index];
    mainChartData.push({
      ...element,
      color: element.close > element.open ? upColor : downColor,
    })

    volumeChartData.push({
      time: element.time,
      value: element.volume,
      color: element.close > element.open ? upVolumeColor : downVolumeColor
    })

    areaChartData.push({
      time: element.time,
      value: element.high,
      lineColor: '#2962FF', 
      topColor: '#2962FF', 
      bottomColor: 'rgba(41, 98, 255, 0.28)' 
    })
  }

  return { mainChartData, volumeChartData,areaChartData, stockHLOC:data }
}



// build the container to visualize the chart
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
      mode: CrosshairMode.Normal,
      vertLine: {
          width: 4,
          color: '#C3BCDB44',
          style: LineStyle.Solid,
          labelBackgroundColor: '#9B7DFF',
      },

      horzLine: {
          color: '#9B7DFF',
          labelBackgroundColor: '#9B7DFF',
      },
    },
    timeScale: {
      borderColor: '#cccccc',
      timeVisible: true,
      secondsVisible: false,
      fixRightEdge: true,
      fixLeftEdge: true,
      
    },
    rightPriceScale: {
      visible: true,
      borderColor: '#cccccc',
    },
    localization: {
      priceFormatter: (price: number) => `${price.toFixed(2)}$`,
    },
    handleScroll:!props.isTouchDevice,
    handleScale:!props.isTouchDevice,
  });

  createMainSeries();

  // crosshair to show modal when hover over the chart
  chart.subscribeCrosshairMove(handlerCossHairMove);

};


// hadnlers
const handlerCossHairMove = (param: MouseEventParams) => {

    const index = param.logical;
    if(index === undefined || index < 0){
      selectedCrosshairData.value = null;
      return;
    } 
    
    console.log(index)
    console.log( chartData.value.stockHLOC[index])
    if(index <= props.historicalData.length){
      selectedCrosshairData.value = chartData.value.stockHLOC[index];
    }

    if(index > props.historicalData.length){
      selectedCrosshairData.value = predictChartData.value.stockHLOC[index - props.historicalData.length + 1];
    }
    
  }

// create the way to show the chart
const createMainSeries = () => {
  if (!chart) return;

  if (chartType.value === ChartType['candlestick']) {
    buildCandleChart()
  } else {
    buildLineAreaChart()
  }
  
  updateScales()
};

const buildCandleChart = () => {
  updateChartSerie('Area', 'remove')
  updateChartSerie('Candlestick', 'add')
  updateChartSerie('Histogram', 'add')
   
  const predictAction = isShowPredict.value ? 'add' : 'remove'
  updateChartSerie('Predict', predictAction)
  updateChartSerie('PredictHistogram', predictAction)
}

const buildLineAreaChart = () => {
  updateChartSerie('Area', 'add')
  
  updateChartSerie('Histogram', 'remove')
  updateChartSerie('Candlestick', 'remove')
  updateChartSerie('Predict', 'remove')
  updateChartSerie('PredictHistogram', 'remove')

}

// update the chart series, add or remove
// add just work if the series is not already added
// remove just work if the series is already added
const updateChartSerie = (type: TSeriesType, action: 'add' | 'remove') =>{
  if(!chart) return;

  if(action === 'remove'){ 
    if(series[type] === null) 
      return;
    
    chart.removeSeries(series[type])
    series[type] = null;
    return;  
  }

  if(action === 'add' && series[type] !== null){
    return;
  }

  switch (type) {
    case 'Candlestick':
      series['Candlestick'] = chart.addSeries(CandlestickSeries, CandleSeriesTheme);
      series['Candlestick'].setData(chartData.value.mainChartData);
      break;
    case 'Predict':
      series['Predict'] = chart.addSeries(CandlestickSeries, CandlePredictSeriesTheme);
      series['Predict'].setData(predictChartData.value.mainChartData);
      break;
    case 'Area':
      series['Area'] = chart.addSeries(AreaSeries, LinearAreaSeriesTheme);
      series['Area'].setData(chartData.value.areaChartData);
      break;
    case 'Histogram':
      series['Histogram'] = chart.addSeries(HistogramSeries,VolumeSeriesTheme);
      series['Histogram'].setData(chartData.value.volumeChartData);
      break;
    case 'PredictHistogram':
      series['PredictHistogram'] = chart.addSeries(HistogramSeries,PredictVolumeSeriesTheme);
      series['PredictHistogram'].setData(predictChartData.value.volumeChartData);
      break;
    default:
      break;
  }
}

const toggleShowPredict = async () =>{
  if(!chart || props.isPredictLoading) return;
  
  if(!props.predictNextWeek || props.predictNextWeek?.length === 0){
    emit('update:predict', !isShowPredict.value)
    return;
  }

  isShowPredict.value = !isShowPredict.value;
  if(chartType.value === ChartType['candlestick'])
    buildCandleChart()
}
// togle the view of the chart
const toggleChartType = (_chartType: ChartType) => {
  if (!chart) 
    return;

  if(chartType.value === _chartType)
    return;

  chartType.value = _chartType;
  createMainSeries();
  zoomBasedInIntervalTime()
  GoToLastStock();
};

const updateScales = () => {
  if (!chart) return;
  const existVolumeScale = series['Histogram'] !== null; 
  
  const volumeScale = {
    visible: false,
    scaleMargins: {
      top: 0.8,
      bottom: 0,
    },
  }

  const mainSerieScale = {
    visible: true,
    scaleMargins: {
      top: 0,  
      bottom: 0.2, 
    },
  }

  if (!existVolumeScale) {
    mainSerieScale.scaleMargins = {
      top: 0,
      bottom: 0,
    }
  }

  if(existVolumeScale){
    chart.priceScale(Scales.volume).applyOptions(volumeScale); 
  }
  
  chart.priceScale(Scales.right).applyOptions(mainSerieScale);
  handleResize()
}

// set the scale of the charts and do zoom  
const zoomBasedInIntervalTime = () => {
 if (!chart || chartData.value.mainChartData.length === 0) return;
  const data = chartData.value.mainChartData;
  const lastIndex = data.length - 1;
  
  let fromIndex = 0;

  switch (timeframe.value) {
    case Timeframe['All']:
      fromIndex = 0;
      break;
    default:
      fromIndex = Math.max(0, lastIndex - timeframe.value);
      break;
  }

  visibleRange.value = { from: data[fromIndex].time, to: data[lastIndex].time };
  GoToLastStock(); 
}


const handleResize = () => {
  debounced(() =>{
    if (chart && chartContainer.value) {
      chart.applyOptions({
        width: chartContainer.value.clientWidth,
      });
    }
  }, 100)
}

const GoToLastStock = () => {
  if (!chart) return;
    chart.timeScale().scrollToRealTime();
}

const takeScreenshot = async () => {
  try {
    if ( !chartContainer.value) return;
  
    isGeneratingScreenshot.value = true;
    const lightWeightLogo = document.querySelector('#tv-attr-logo')
    if(!lightWeightLogo?.getAttribute('data-html2canvas-ignore'))
      lightWeightLogo?.setAttribute('data-html2canvas-ignore', 'true')

    // await to allow show the background
    await new Promise((r) => setTimeout(r, 10));

    const canvas = await html2canvas(chartContainer.value, {
      backgroundColor: "#ffffff",
      useCORS: true,
      scale: 2
    });

    const link = document.createElement("a");
    link.href = canvas.toDataURL("image/png");
    link.download = `chart-${Date.now()}.png`;
    link.click();
  } finally {
    isGeneratingScreenshot.value = false;
  }
};

watch(timeframe, async () => {
  emit('update:timeframe', timeframe.value)
  zoomBasedInIntervalTime()
});

watch(() => props.predictNextWeek, async () => {

  if(props.predictNextWeek?.length > 0){
    isShowPredict.value = true;
  }
  else{
    isShowPredict.value = false;
  }
  createMainSeries()
});

watch(() => props.isTouchDevice, (newTouchable) => {
  if(chart){
    chart.applyOptions({
      handleScroll: !newTouchable,
      handleScale: !newTouchable,
    });
  }

  createMainSeries();
});

watch(visibleRange, () => {
  if (chart) {
    if (visibleRange.value) {
      // show a segment of line of the chart
      chart.timeScale().setVisibleRange(visibleRange.value);
    } else {
      // show all the chart
      chart.timeScale().scrollToRealTime();
    }
  }
});


// recreate the chart when the data changes
watch(() => props.historicalData, (newData) => {
  if (!chart) return;

  if(chartType.value === ChartType['candlestick']){
    series['Candlestick']?.setData(chartData.value.mainChartData)
    series['Histogram']?.setData(chartData.value.volumeChartData)
    series['Predict']?.setData(predictChartData.value.mainChartData)
    series['PredictHistogram']?.setData(predictChartData.value.volumeChartData)
  }

  if(chartType.value === ChartType['area']){
    series['Area']?.setData(chartData.value.areaChartData)
  }
  
  updateScales()
  zoomBasedInIntervalTime()
  chart.timeScale().fitContent();
  GoToLastStock();
});


// Lifecycle
onMounted(() => {
  initChart();
  GoToLastStock(); 
  zoomBasedInIntervalTime()
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  if (chart) {
    chart.unsubscribeCrosshairMove(handlerCossHairMove);
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

.info-row {
  margin-bottom: 4px;
}

.info-row:last-child {
  margin-bottom: 0;
}


</style>