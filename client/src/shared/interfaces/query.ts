// Type definitions
export interface PaginationParams {
  page?: number;
  size?: number;
  sort?: 'asc' | 'desc';
}

export interface ListParams extends PaginationParams {
  q?: string;
}
