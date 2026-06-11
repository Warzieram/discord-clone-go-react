// Backend base URL (including protocol), configured per environment via Vite env vars.
//
//   - Local dev:   VITE_BACKEND_URL=http://localhost:8080   (.env.development)
//   - Production:  VITE_BACKEND_URL=https://api.example.com  (.env.production / deploy env)
//
// If unset, we fall back to the current page origin, which is handy when the
// frontend is served behind a reverse proxy that also exposes the API.
const rawBackendUrl = import.meta.env.VITE_BACKEND_URL ?? window.location.origin;

// Normalize: drop any trailing slash so callers can safely do `${BACKEND_URL}/api/...`.
export const BACKEND_URL = rawBackendUrl.replace(/\/+$/, "");

// Derive the WebSocket base from the HTTP base: http -> ws, https -> wss.
export const WS_BACKEND_URL = BACKEND_URL.replace(/^http/, "ws");
