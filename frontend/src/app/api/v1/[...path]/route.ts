import { NextRequest } from "next/server";

// Same-origin proxy to the backend. Runs per-request so BACKEND_URL is read
// from the Node process environment on every call — no build-time baking.
// Replaces the old next.config.ts rewrites(), which was evaluated at
// `next build` time and froze the wrong URL into the manifest.

export const dynamic = "force-dynamic";
export const runtime = "nodejs";

async function proxy(
  req: NextRequest,
  ctx: { params: Promise<{ path: string[] }> }
): Promise<Response> {
  const backend = process.env.BACKEND_URL || "http://localhost:8080";
  const { path } = await ctx.params;
  const upstreamURL = `${backend}/api/v1/${path.join("/")}${req.nextUrl.search}`;

  // Forward request headers, except ones that refer to the proxy hop itself.
  const fwdHeaders = new Headers();
  req.headers.forEach((value, key) => {
    const k = key.toLowerCase();
    if (k === "host" || k === "connection" || k === "content-length") return;
    fwdHeaders.set(key, value);
  });

  const hasBody = req.method !== "GET" && req.method !== "HEAD";
  const upstream = await fetch(upstreamURL, {
    method: req.method,
    headers: fwdHeaders,
    body: hasBody ? await req.arrayBuffer() : undefined,
    redirect: "manual",
  });

  // Copy response headers, preserving multiple Set-Cookie entries correctly.
  // `new Headers(upstream.headers)` collapses duplicates; using getSetCookie
  // + append is the supported way to forward them all.
  const resHeaders = new Headers();
  upstream.headers.forEach((value, key) => {
    if (key.toLowerCase() !== "set-cookie") {
      resHeaders.set(key, value);
    }
  });
  for (const cookie of upstream.headers.getSetCookie()) {
    resHeaders.append("set-cookie", cookie);
  }

  return new Response(upstream.body, {
    status: upstream.status,
    statusText: upstream.statusText,
    headers: resHeaders,
  });
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const PATCH = proxy;
export const DELETE = proxy;
export const OPTIONS = proxy;
