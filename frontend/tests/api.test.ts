/**
 * Tests for the API client wrapper.
 * Uses a mock fetch to verify request construction.
 */

import { api } from "@/lib/api";

const originalFetch = global.fetch;

beforeEach(() => {
  global.fetch = jest.fn();
});

afterEach(() => {
  global.fetch = originalFetch;
});

function mockFetchResponse(data: unknown, status = 200) {
  (global.fetch as jest.Mock).mockResolvedValueOnce({
    ok: status >= 200 && status < 300,
    status,
    json: async () => data,
  });
}

// ---------------------------------------------------------------------------
// GET requests
// ---------------------------------------------------------------------------

describe("api.get", () => {
  test("sends GET request to correct URL", async () => {
    mockFetchResponse({ data: { id: "1" }, error: null });

    await api.get("/users/1");

    const [url] = (global.fetch as jest.Mock).mock.calls[0];
    expect(url).toContain("/users/1");
  });

  test("includes credentials", async () => {
    mockFetchResponse({ data: null, error: null });

    await api.get("/auth/me");

    expect(global.fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({ credentials: "include" })
    );
  });

  test("returns parsed data on success", async () => {
    const mockData = { id: "abc", name: "Alice" };
    mockFetchResponse({ data: mockData, error: null });

    const result = await api.get("/users/abc");
    expect(result.data).toEqual(mockData);
    expect(result.error).toBeNull();
  });

  test("returns error on API error", async () => {
    mockFetchResponse(
      { data: null, error: { code: "NOT_FOUND", message: "not found" } },
      404
    );

    const result = await api.get("/users/missing");
    expect(result.data).toBeNull();
    expect(result.error?.code).toBe("NOT_FOUND");
  });
});

// ---------------------------------------------------------------------------
// POST requests
// ---------------------------------------------------------------------------

describe("api.post", () => {
  test("sends POST with JSON body", async () => {
    mockFetchResponse({ data: { id: "new" }, error: null });
    const body = { name: "Test Match" };

    await api.post("/matches", body);

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining("/matches"),
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify(body),
      })
    );
  });

  test("sets Content-Type header", async () => {
    mockFetchResponse({ data: null, error: null });

    await api.post("/matches", {});

    expect(global.fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: expect.objectContaining({
          "Content-Type": "application/json",
        }),
      })
    );
  });
});

// ---------------------------------------------------------------------------
// DELETE requests
// ---------------------------------------------------------------------------

describe("api.delete", () => {
  test("sends DELETE request", async () => {
    mockFetchResponse({ data: { message: "deleted" }, error: null });

    await api.delete("/matches/abc");

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining("/matches/abc"),
      expect.objectContaining({ method: "DELETE" })
    );
  });
});

// ---------------------------------------------------------------------------
// Error handling
// ---------------------------------------------------------------------------

describe("api error handling", () => {
  test("returns NETWORK_ERROR when response has no error body", async () => {
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: false,
      status: 500,
      json: async () => ({ data: null }),
    });

    const result = await api.get("/fail");
    expect(result.error?.code).toBe("NETWORK_ERROR");
    expect(result.error?.message).toContain("500");
  });
});
