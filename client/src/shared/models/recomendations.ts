
export interface Recommendation {
    id: number; 
    ticker_id: string; 
    target_from: string; 
    target_to: string; 
    action: string; 
    rating_from: string; 
    rating_to: string; 
    time: string;
    brokerage: Brokerage
}

export interface Brokerage {
    id: number;
    name: string;
}

export interface Ticker {
  company:string,
  id:string,
  sentiment: string;
  recommendations?: Recommendation[];
}


export interface CompanyData {
  symbol: string;
  price: number;
  marketCap: number;
  beta: number;
  lastDividend: number;
  change: number;
  changePercentage: number;
  volume: number;
  averageVolume: number;
  companyName: string;
  exchangeFullName: string;
  exchange: string;
  industry: string;
  website: string;
  sector: string;
  country: string;
  image: string;
  ceo: string;
}

export interface HistoricalPrice{
  symbol: string;
  date: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
  change: number;
  changePercent: number;
  vwap: number;
}

export interface StockHLOC {
    time: string;
    open: number;
    high: number;
    low: number;
    close: number;
    volume: number;
    change: number;
    changePercent: number;
    vwap: number;
    date: Date;
}

export interface CompanyNew {
  id: number;
  category: string;
  datetime: number;
  datetimeUtc: string;
  headline: string;
  image: string;
  related: string;
  source: string;
  summary: string;
  url: string;
}

export interface CompanyOverview {
  companyData:CompanyData,
  recommendations:Recommendation[],
  historicalPrices:HistoricalPrice[],
  companyNews:CompanyNew[]
  advice: string;
}

export interface TickerListResponse{
  ticker:Ticker,
  companyData:CompanyData
  advice: string;
}