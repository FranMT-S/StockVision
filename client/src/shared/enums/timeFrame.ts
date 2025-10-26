const daysInMonth = 30;
export enum Timeframe {
  'All' = 0,
  '1M' = daysInMonth*1,
  '3M' = daysInMonth*3,
  '6M' = daysInMonth*6,
  '1Y' = daysInMonth*12,
}
