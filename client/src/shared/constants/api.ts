
// windows.__ENV__ is used to get the environment variables in runtime in prod
export const API_CONFIG = {
  BASE_URL: window.__ENV__?.VITE_API_URL || import.meta.env.VITE_API_URL || 'http://localhost:8080',
  ENDPOINTS: {
    List: function(id:string = '',page:number = 1,size: number = 10,sort: 'asc' | 'desc' = 'asc'){
      const idParam = id ? `&id=${id}` : ''

      return `/api/v1/tickers?page=${page}&size=${size}&sort=${sort}${idParam}`
    },
    Overview: function(id: string){
      return `/api/v1/tickers/${id}/overview`
    },
    Logo: function(id: string){
      return `/api/v1/tickers/${id}/logo`
    },
  },
};

export const DEFAULT_PAGINATION = {
  PAGE: 1,
  LIMIT: 10,
};
