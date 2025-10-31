import { formatJustDate } from "../helpers/formats";
import { buildQueryString, isValidQuery, normalizePageNumber, normalizePageSize } from "../helpers/query";
import { sanitizeSQL } from "../helpers/sanatizer";
import { ListParams } from "../interfaces/query";


export const API_URL = window.__ENV__?.VITE_API_URL ||  import.meta.env.VITE_API_URL ||  'http://localhost:8080'

// API Configuration
export const API_CONFIG = {
  BASE_URL: API_URL,
  ENDPOINTS: {
    /**
     * Get paginated list of tickers
     * @param params - List parameters including query, pagination and sorting
     * @returns API endpoint URL
     */
    List: ({ q, page, size, sort = 'asc' }: ListParams = {}): string => {
      const normalizedPage = normalizePageNumber(page);
      const normalizedSize = normalizePageSize(size);
      
      const params: Record<string, string | number> = {
        page: normalizedPage,
        size: normalizedSize,
        sort,
      };

      if (isValidQuery(q)) {
        params.q = sanitizeSQL(q!.trim());
      }

      return `${API_URL}/api/v1/tickers${buildQueryString(params)}`;
    },

    /**
     * Get ticker overview from a specific date
     * @param id - Ticker ID
     * @param from - Start date for overview
     * @returns API endpoint URL
     * @throws Error if Ticker ID is not provided
     */
    Overview: (id: string, from: Date): string => {
      if (!id?.trim()) {
        throw new Error('Ticker ID is required');
      }

      const date = formatJustDate(from);
      return `${API_URL}/api/v1/tickers/${encodeURIComponent(id)}/overview?from=${date}`;
    },

    /**
     * Get ticker logo URL
     * @param id - Ticker ID
     * @returns API endpoint URL
     * @throws Error if Ticker ID is not provided
     */
    Logo: (id: string): string => {
      if (!id?.trim()) {
        throw new Error('Ticker ID is required');
      }
      return `${API_URL}/api/v1/tickers/${encodeURIComponent(id)}/logo`;
    },

    /**
     * Get ticker predictions
     * @param id - Ticker ID
     * @returns API endpoint URL
     * @throws Error if Ticker ID is not provided
     */
    Predictions: (id: string): string => {
      if (!id?.trim()) {
        throw new Error('Ticker ID is required');
      }
      return `${API_URL}/api/v1/tickers/${encodeURIComponent(id)}/predictions`;
    },

    /**
     * Fetch historical prices from the API
     * @param id - Ticker ID
     * @param from - Optional start date. If not provided, fetches all historical prices
     * @returns API endpoint URL
     */
    HistoricalPrices: (id: string, from?: Date): string => {
      if (!id?.trim()) {
        throw new Error('Ticker ID is required');
      }
      
      const basePath = `${API_URL}/api/v1/tickers/${encodeURIComponent(id)}/historical`;
      
      if (!from) {
        return basePath;
      }

      const date = formatJustDate(from);
      return `${basePath}?from=${date}`;
    },

    /**
     * Onboarding endpoint
     * @returns Onboarding base url
     */
    Onboarding: (): string => {
      return `${API_URL}/api/v1/onboarding`;
    },
  },
} as const;