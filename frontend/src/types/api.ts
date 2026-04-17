export interface APIResponse<T> {
  data: T | null;
  error: APIError | null;
  meta?: Meta;
}

export interface APIError {
  code: string;
  message: string;
}

export interface Meta {
  page: number;
  per_page: number;
  total: number;
}
