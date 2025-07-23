/**
 * Next.js Configuration
 *
 * This configuration enables API proxying from the frontend to the backend server.
 * All requests to /api/* are automatically forwarded to the backend API server,
 * allowing the frontend to make API calls without CORS issues during development.
 *
 * The API URL is configurable via the NEXT_PUBLIC_API_URL environment variable,
 * defaulting to http://localhost:8080 for local development.
 *
 * For production, set NEXT_PUBLIC_API_URL to your production API server URL.
 */

/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    // Use environment variable for API URL, fallback to localhost for development
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

    return [
      {
        source: '/api/:path*',
        destination: `${apiUrl}/api/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
