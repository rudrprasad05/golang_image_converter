import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: false,
};

module.exports = {
  output: "standalone",
  images: {
    domains: ["mctechfiji.s3.us-east-1.amazonaws.com"],
  },
};

export default nextConfig;
