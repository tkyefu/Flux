import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    // In development, proxy /api/* to the Go backend to avoid CORS issues
    if (process.env.NODE_ENV === "development") {
      return [
        {
          source: "/api/:path*",
          destination: "http://localhost:8080/api/v1/:path*",
        },
      ];
    }
    return [];
  },
};

export default nextConfig;
