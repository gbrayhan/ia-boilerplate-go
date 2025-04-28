# 🚀 IA Golang Boilerplate 2025

> **The ultimate Go microservice template**—optimized for speed, reliability, and seamless AI-driven collaboration.

![Go](https://img.shields.io/badge/Go-1.24-blue) ![Gin](https://img.shields.io/badge/Gin-Framework-brightgreen) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.4-blue) ![Docker](https://img.shields.io/badge/Docker-Compose-orange) ![LLM‑Friendly](https://img.shields.io/badge/LLM--Friendly-lightgrey)

---

## Security Checks using Trivy

https://github.com/aquasecurity/trivy?tab=readme-ov-file

command:
```bash
trivy fs . 
```

---

## 📖 Contents

- [✨ Key Features](#-key-features)
- [⚙️ Prerequisites](#️-prerequisites)
- [🔧 Installation & Setup](#-installation--setup)
- [🐳 Docker & Compose](#-docker--compose)
- [💻 Local Development](#-local-development)
- [📡 Main Endpoints](#-main-endpoints)
- [🛡️ Security & Auth](#️-security--auth)
- [📝 Environment Variables](#-environment-variables)
- [🤝 Contributing](#-contributing)
- [📜 License](#-license)

---

## ✨ Key Features

- 🏗️ **RESTful API** with Gin in Release mode.
- 🔒 **JWT Authentication** for access & refresh tokens.
- 🗄️ **GORM ORM** with automatic migrations & seeding (admin role + initial user).
- ⏰ **cron jobs** via robfig/cron (daily tasks).
- 🤖 **LLM‑Friendly** code structure—designed for easy snippet sharing, AI-assisted edits, and smooth integration with
  large language models.
- 🐳 **Full containerization**: Docker + Distroless + Compose with healthchecks.
- 📜 **Clean architecture**: clear separation of layers (db, handlers, middleware, infra, utils).

---

## ⚙️ Prerequisites

- Go ≥ 1.24
- Docker ≥ 20.10 & Docker Compose v2
- PostgreSQL ≥ 17.4
- Bash or Make (optional)

---

## 🔧 Installation & Setup

1. **Clone the repo**
   ```bash
   git clone https://github.com/your-org/ia-boilerplate.git
   cd ia-boilerplate
   ```

2. **Copy & edit your `.env`**
   ```bash
   cp .env.example .env
   # Fill in your credentials and secrets
   ```

3. **Build the Go binary**
   ```bash
   go mod download
   go build -o ia-boilerplate .
   ```

---

## 🐳 Docker & Compose

Bring up your full stack (API + DB) in one command:

```bash
docker compose up --build
```

- **db-ia-boilerplate** → PostgreSQL on port **5432**
- **ia-boilerplate** → API on port **8080**
- Healthchecks ensure DB readiness before app startup.

_Teardown:_

```bash
docker compose down --volumes
```

---

## 💻 Local Development

Skip Docker and run directly:

```bash
export $(grep -v '^#' .env | xargs)
go run main.go
```

Open your browser at `http://localhost:8080` and you’re ready!

---

## 📡 Main Endpoints

| Method | Route                             | Description                                |
|:------:|-----------------------------------|--------------------------------------------|
|  POST  | `/login`                          | Authenticate: returns access & refresh JWT |
|  POST  | `/access-token/refresh`           | Refresh access token with refresh token    |
|  GET   | `/api/device`                     | Device info (requires JWT)                 |
|  GET   | `/api/users`                      | List users                                 |
|  GET   | `/api/users/:id`                  | Get user by ID                             |
|  POST  | `/api/users`                      | Create new user                            |
|  PUT   | `/api/users/:id`                  | Update user                                |
| DELETE | `/api/users/:id`                  | Delete user                                |
|  GET   | `/api/medicines/search-paginated` | Paginated medicine search                  |
|  GET   | `/api/icd-cie/search-paginated`   | Paginated ICD‑CIE code search              |

> 🔎 Explore additional endpoints for roles, devices, ICD‑CIE, etc., under `/api`.

---

## 🛡️ Security & Auth

- Use header `Authorization: Bearer <token>` for protected routes.
- Token lifetimes controlled by `ACCESS_TOKEN_TTL` & `REFRESH_TOKEN_TTL` (minutes).
- Secrets managed entirely via environment variables.

---

## 📝 Environment Variables

| Variable             | Description                  | Example                |
|----------------------|------------------------------|------------------------|
| `DB_HOST`            | PostgreSQL host              | `db-ia-boilerplate`    |
| `DB_PORT`            | PostgreSQL port              | `5432`                 |
| `DB_USER`            | DB username                  | `app_user`             |
| `DB_PASSWORD`        | DB password                  | `yourpassword`         |
| `DB_NAME`            | DB name                      | `ia-boilerplate`       |
| `DB_SSLMODE`         | SSL mode (disable/require)   | `disable`              |
| `APP_PORT`           | API port                     | `8080`                 |
| `ACCESS_SECRET_KEY`  | JWT access token secret      | `yourAccessSecretKey`  |
| `REFRESH_SECRET_KEY` | JWT refresh token secret     | `yourRefreshSecretKey` |
| `ACCESS_TOKEN_TTL`   | Access token TTL (minutes)   | `15`                   |
| `REFRESH_TOKEN_TTL`  | Refresh token TTL (minutes)  | `10080`                |
| `JWT_ISSUER`         | JWT issuer                   | `my-app`               |
| `IMGUR_CLIENT_ID`    | (Optional) Imgur integration | `yourImgurClientId`    |
| `START_USER_EMAIL`   | Seed admin user email        | `gbrayhan@gmail.com`   |
| `START_USER_PW`      | Seed admin user password     | `qweqwe`               |

---

## 🤝 Contributing

1. Fork the repo.
2. Create your feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```
3. Commit your changes:
   ```bash
   git commit -am "Add awesome feature"
   ```
4. Push to branch:
   ```bash
   git push origin feature/my-feature
   ```
5. Open a Pull Request—let’s make this the best boilerplate of 2025 together! 🚀

---

## 📜 License

This project is licensed under the **MIT License**. Open source, future‑proof, and AI‑ready! 😊

---

> “The best way to predict the future is to create it.” – Peter Drucker 🤖

