/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** Backend base URL including protocol, e.g. "http://localhost:8080". */
  readonly VITE_BACKEND_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
