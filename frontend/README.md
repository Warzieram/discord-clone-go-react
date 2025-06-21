# AUTHENTICATION BOILERPLATE GO / REACT

A full-stack authentication boilerplate featuring email verification with Mailjet integration. Built with Go (Gorilla/mux) backend and React (Vite) frontend.
This was made for personal use and is * (obviously) * subject to improvements

## ✨ Features

- 🔐 Complete user authentication system
- 📧 Email verification using Mailjet

## 🏗️ Architecture

### Frontend
- **React 18** with TypeScript support
- **Vite** for fast development and building

### Backend
- **Go** with Gorilla/mux router
- **go-blueprint** scaffolding
- Automatic user table creation
- RESTful API endpoints
- Mailjet email service integration

## 📋 Prerequisites

- **Node.js** (v16 or higher)
- **pnpm** package manager
- **Go** (v1.19 or higher)
- **Make** utility
- **Mailjet account** for email services

## 🚀 Quick Start

### 1. Clone the repository
```bash
git clone git@github.com:Warzieram/auth-boilerplate-go-react.git
cd auth-boilerplate-go-react
```

### 2. Setup Frontend
```bash
cd frontend/go_auth_test_frontend/
pnpm install
```

### 3. Setup Backend
```bash
cd ../../backend/back/
# Copy environment variables template
cp .env.example .env
# Edit .env with your Mailjet credentials and database settings
```

### 4. Configure Environment Variables
Create a `.env` file in the backend directory using the `.env.example` with:

## Running the Application

### Start Frontend (Terminal 1)
```bash
cd frontend/go_auth_test_frontend/
pnpm dev
```
Frontend will be available at `http://localhost:5173`

### Start Backend (Terminal 2)
```bash
cd backend/back/
make run
```
Backend API will be available at `http://localhost:8080`

## 📚 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/register` | User registration |
| POST | `/api/login` | User login |
| POST | `/api/verify` | Email verification |
| POST | `/api/resend-verification` | Resend verification email (WIP) |
| GET | `/api/profile` | Get user profile (protected) |

## 🤝 Contributing

Contributions are welcome! This project was initially created for personal use, but I'm excited to learn from the community.

### How to Contribute:
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/great-beautiful-awesome-feature`)
3. Commit your changes (`git commit -m 'Add the great beautiful awesom feature'`)
4. Push to the branch (`git push origin feature/great-beautiful-awesome-feature`)
5. Open a Pull Request

### Areas for Improvement:
- [ ] Add comprehensive tests
- [ ] Implement password reset functionality
- [ ] Add social authentication
- [ ] Improve error handling
- [ ] Add rate limiting
- [ ] Enhance security features

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🙏 Acknowledgments

- Built with [go-blueprint](https://github.com/Melkeydev/go-blueprint)
- Email services powered by [Mailjet](https://www.mailjet.com/)
- Frontend bootstrapped with [Vite](https://vitejs.dev/)

---

**Note:** This is a boilerplate project intended for learning and development purposes. Please review and enhance security measures before using in production.
