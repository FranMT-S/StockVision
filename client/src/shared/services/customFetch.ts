import { ApiResponse, ErrorType } from "@/shared/models/response";

/**
 * Custom fetch function to handle API responses
 * transform the response in a json format
 * if exist a error is saved inthe error property
 */
export const customFetch = async <T>(input: RequestInfo | URL, init?: RequestInit): Promise<ApiResponse<T>> => {
  let data: ApiResponse<T>;
  let response: Response;
  
  try {
    response = await fetch(input, {
      ...init});
    data = await response.json();
    data.status = response.status;
    if(response.status >= 400){
      data.ok = false;
    }else{
      data.ok = true;
    }
  } catch (error) {
    if (error instanceof Error && error.name === "AbortError") {
      data = {
        ok: false,
        error: "The request was aborted",
        status: 500,
        errorType: ErrorType.ABORT_ERROR,
      }
    }else{
      data = {
        ok: false,
        error: "there was an error, try again later",
        status: 500,
        errorType: ErrorType.INVALID_JSON,
      }
    }
  }

  return data;
}