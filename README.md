# Discord Clone (Go / React)

A full-stack real-time chat application. The backend is written in Go (Gorilla/mux)
and exposes a REST API plus a WebSocket endpoint for live messaging. The frontend is
a React (Vite + TypeScript) single-page app. User authentication includes email
verification through Mailjet.

This project started as a personal learning exercise and is subject to ongoing
improvements.

## Features

- User authentication (register, login) with JWT
- Email verification using Mailjet
- Real-time messaging over WebSocket
- Message history with pagination
- Message deletion broadcast to all connected clients

## Architecture

### Frontend
- React 19 with TypeScript
- Vite for development and production builds
- Redux Toolkit for state management
- React Router for navigation
- Backend URL configured per environment via Vite env variables

### Backend
- Go with the Gorilla/mux router
- Gorilla WebSocket for real-time messaging
- JWT-based authentication (HTTP header for REST, subprotocol for WebSocket)
- Mailjet email service integration
- Scaffolded with go-blueprint

## Prerequisites

- Node.js (v18 or higher)
- npm or pnpm
- Go (v1.21 or higher)
- Make utility
- A Mailjet account for email services
- Docker and Docker Compose (optional, for containerized frontend)

## Quick Start

### 1. Clone the repository
```bash
git clone git@github.com:Warzieram/discord-clone-go-react.git
cd discord-clone-go-react
```

### 2. Set up the backend
```bash
cd backend/back/
# Copy the environment template and fill in your values
cp .env.example .env
```

Configure the `.env` file with your settings:

| Variable | Description |
|----------|-------------|
| `PORT` | Port the API listens on (e.g. `8080`) |
| `APP_ENV` | Application environment (e.g. `local`) |
| `DOMAIN_NAME` | Domain used in verification emails |
| `BLUEPRINT_DB_HOST` | Database host |
| `BLUEPRINT_DB_PORT` | Database port |
| `BLUEPRINT_DB_DATABASE` | Database name |
| `BLUEPRINT_DB_USERNAME` | Database user |
| `BLUEPRINT_DB_PASSWORD` | Database password |
| `BLUEPRINT_DB_SCHEMA` | Database schema |
| `JWT_SECRET` | Secret used to sign JWTs |
| `MJ_APIKEY_PUBLIC` | Mailjet public API key |
| `MJ_APIKEY_PRIVATE` | Mailjet private API key |

### 3. Set up the frontend
```bash
cd frontend/go_auth_test_frontend/
npm install
```

The frontend reads the backend URL from a Vite environment variable,
`VITE_BACKEND_URL`. Defaults are provided per mode:

- `.env.development` -> `http://localhost:8080` (used by `npm run dev`)
- `.env.production` -> override with your real API origin (used by `npm run build`)

To override locally without editing committed files, copy `.env.example` to
`.env.local` and set `VITE_BACKEND_URL`. The WebSocket URL is derived automatically
from this value (`http` becomes `ws`, `https` becomes `wss`).

## Running the Application

### Start the backend (Terminal 1)
```bash
cd backend/back/
make run
```
The API will be available at `http://localhost:8080`.

### Start the frontend (Terminal 2)
```bash
cd frontend/go_auth_test_frontend/
npm run dev
```
The frontend will be available at `http://localhost:5173`.

## Running the Frontend with Docker

The frontend includes a production Dockerfile (multi-stage build served by nginx)
and a development Dockerfile with hot-reload.

```bash
cd frontend/go_auth_test_frontend/

# Production build served by nginx on http://localhost:8080
VITE_BACKEND_URL=https://api.example.com docker compose build frontend
docker compose up frontend

# Development server with hot-reload on http://localhost:5173
docker compose --profile dev watch frontend-dev
```

Note: `VITE_BACKEND_URL` is baked into the static bundle at build time, so a
production image targets a single backend origin. Rebuild with a different value to
point at a different backend.

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/register` | Public | User registration |
| POST | `/api/login` | Public | User login, returns a JWT |
| GET | `/api/verify` | Public | Email verification |
| GET | `/api/profile` | Bearer token | Get the current user profile |
| GET | `/api/messages` | Bearer token | Retrieve message history (`limit`, `offset` query params) |
| GET | `/api/message` | WebSocket subprotocol | Real-time messaging channel |

### WebSocket authentication

The messaging endpoint authenticates via a WebSocket subprotocol rather than an
`Authorization` header (browsers cannot set custom headers on WebSocket
connections). The client opens the connection with a subprotocol of the form
`auth.<jwt>`:

```ts
new WebSocket(`${WS_BACKEND_URL}/api/message`, [`auth.${token}`]);
```

The server reads the token from the subprotocol, validates it, and echoes the
subprotocol back in the handshake response.

## Contributing

Contributions are welcome. This project was initially created for personal use, but
improvements from the community are appreciated.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -m 'Add your feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a Pull Request

### Areas for Improvement
- [ ] Add comprehensive tests
- [ ] Implement password reset functionality
- [ ] Add social authentication
- [ ] Improve error handling
- [ ] Add rate limiting
- [ ] Enhance security features

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- Built with [go-blueprint](https://github.com/Melkeydev/go-blueprint)
- Email services powered by [Mailjet](https://www.mailjet.com/)
- Frontend bootstrapped with [Vite](https://vitejs.dev/)

---

Note: This project is intended for learning and development purposes. Review and
enhance the security measures before using it in production.
