import { formatJustDate } from "../helpers/formats";

// windows.__ENV__ is used to get the environment variables in runtime in prod
export const API_CONFIG = {
  BASE_URL: window.__ENV__?.VITE_API_URL || import.meta.env.VITE_API_URL || 'http://localhost:8080',
  ENDPOINTS: {
    List: function(q:string = '',page:number = 1,size: number = 10,sort: 'asc' | 'desc' = 'asc'){
      if(!q || q == 'undefined') {
        return `/api/v1/tickers?page=${page}&size=${size}&sort=${sort}`
      }

      const qParam = q ? `&q=${q}` : ''
      return `/api/v1/tickers?page=${page}&size=${size}&sort=${sort}${qParam}`
    },
    Overview: function(id: string,from:Date){
      const date = formatJustDate(from)
      return `/api/v1/tickers/${id}/overview?from=${date}`
    },
    Logo: function(id: string){
      return `/api/v1/tickers/${id}/logo`
    },
    Predictions: function(id: string){
      return `/api/v1/tickers/${id}/predictions`
    },

    /**Fetch historical prices from the API if from is not provided, fetch all historical prices */
    HistoricalPrices: function(id: string,from?:Date){
      let fromParam = ''
      if(from){
        const date = formatJustDate(from)
        fromParam = `from=${date}`
      }
      
      return `/api/v1/tickers/${id}/historical?${fromParam}`
    },
  },
};


export const DEFAULT_PAGINATION = {
  PAGE: 1,
  LIMIT: 10,
};
