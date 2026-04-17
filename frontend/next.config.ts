import type { NextConfig } from "next";

// No rewrites() here — they're evaluated at `next build` time and the
// resolved destinations are frozen into the routes manifest, so env vars set
// at container runtime are ignored. The /api/v1/* proxy lives in
// src/app/api/v1/[...path]/route.ts, which runs per-request and reads
// BACKEND_URL fresh from process.env.

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "lh3.googleusercontent.com",
      },
    ],
  },
};

export default nextConfig;
