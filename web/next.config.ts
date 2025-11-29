import type { NextConfig } from "next";
import { config } from "dotenv";

// Load env.local instead of .env.local
config({ path: "env.local" });

const nextConfig: NextConfig = {
  /* config options here */
};

export default nextConfig;
