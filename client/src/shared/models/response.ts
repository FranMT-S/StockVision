export enum ErrorType {
  ABORT_ERROR = "AbortError",
  INVALID_JSON = "InvalidJSON",
}

export type ApiResponse<T> = {
    ok: true;
    data: T;
    status: number;
    total?: number;
  }
| {
    ok: false;
    error: string;  
    status: number;
    errorType?: ErrorType;
  };

