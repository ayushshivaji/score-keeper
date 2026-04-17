import { APIResponse } from "@/types/api";

// Always call the API same-origin. next.config.ts rewrites /api/v1/* to the
// backend service (driven by BACKEND_URL at Next server startup), so the
// browser never talks to the backend directly. This keeps cookies on the
// frontend origin and avoids CORS entirely.
const API_BASE = "/api/v1";

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<APIResponse<T>> {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  const json = await res.json();

  if (!res.ok && !json.error) {
    return {
      data: null,
      error: { code: "NETWORK_ERROR", message: `HTTP ${res.status}` },
    };
  }

  return json;
}

export const api = {
  get: <T>(path: string) => request<T>(path),

  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "POST", body: JSON.stringify(body) }),

  put: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PUT", body: JSON.stringify(body) }),

  delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),
};
