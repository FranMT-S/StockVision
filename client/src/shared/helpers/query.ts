import { DEFAULT_PAGINATION } from "../constants/default";

// Helper functions
export const normalizePageNumber = (page?: number): number => {
  return !page || page <= 0 ? DEFAULT_PAGINATION.PAGE : page;
};

export const normalizePageSize = (size?: number): number => {
  return !size || size <= 0 ? DEFAULT_PAGINATION.SIZE : size;
};

export const isValidQuery = (q?: string): boolean => {
  return !!q && q !== 'undefined' && q.trim() !== '';
};

// Build query string from params
// Use URLSearchParams to clean the params
export const buildQueryString = (params: Record<string, string | number>): string => {
  const KeyValuesAsString = Object.entries(params).map(([key, value]) => [key, String(value)])
  const query = new URLSearchParams(KeyValuesAsString).toString();
  return query ? `?${query}` : '';
};