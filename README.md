# E-Wallet App - Backend

[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)
<br>
Backend REST API project for E-Wallet application by Muh. Ilham Mursidi (Koda Batch 7 Fullstack Web Developer).

## Technologies Used

- [![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev/)
- [![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?logo=go&logoColor=white)](https://gin-gonic.com/)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.10-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
- [![Redis](https://img.shields.io/badge/Redis-8.8.0-FF4438?logo=redis&logoColor=white)](https://redis.io/)
- [![JWT](https://img.shields.io/badge/JWT-Auth-000000?logo=jsonwebtokens&logoColor=white)](https://jwt.io/)
- [![Swagger](https://img.shields.io/badge/Swagger-Docs-85EA2D?logo=swagger&logoColor=white)](https://swagger.io/)
- [![Docker](https://img.shields.io/badge/Docker-29.4.2-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

## Features

- User Authentication (Register, Login, Logout)
- PIN Management (Create & Verify)
- JWT-based Authorization
- Top Up via Payment Method
- Wallet Dashboard (Balance & Summary)
- Forgot Password via Email Verification
- Fund Transfer between Users
- Profile Management with Avatar Upload
- Transaction History & Reports
- API Documentation via Swagger

## API Endpoints

| Method | Endpoint                             | Description            |
| ------ | ------------------------------------ | ---------------------- |
| POST   | `/auth/register`                     | Register new user      |
| POST   | `/auth/login`                              | Login                  |
| POST | `/auth/logout`                       | Logout                 |
| POST   | `/auth/forgot-password` | Send link verify password       |
| POST   | `/auth/verify-reset-token`        | Verify link token reset password         |
| POST   | `/auth/reset-password`        | Set new password         |
| PATCH   | `/auth/enter-pin`        | Set new pin for new user         |
| GET    | `/users/dashboard`                      | Get dashboard info            |
| GET    | `/users/profile`                      | Get user profile            |
| PATCH  | `/users/profile`                      | Edit user profile         |
| PATCH  | `/users/profile/change-password`                     | Change user password        |
| PATCH  | `/users/profile/change-pin`                          | Change user PIN             |
| GET   | `/users/transaction-report`            | Get transaction report             |
| GET    | `/users/transactions`                       | Get transactions history     |
| GET    | `/transaction/receivers`                      | Get and Find receivers |
| POST    | `/transaction/topup`             | Top up wallet balance
| GET    | `/transaction/transfer`               | Transaction history    |

Full interactive docs available at `/swagger/index.html` after running the server.  

## Usage Instruction

### Running the Application (Local Development)

1. Clone this repository:

```bash
$ git clone https://github.com/Ilhammursidi/Ewallet-Backend.git
```

2. Install dependencies:

```bash
$ go mod tidy
```

3. Run database migrations:

```bash
$ migrate -path db/migrations -database "postgres://myuser:yourpassword@localhost:5432/mydb?sslmode=disable" up
```

4. Run the development server:

```bash
$ go run cmd/main.go
```

### Running with Docker Compose

Make sure you are in the root `deployment/` directory:

```bash
$ docker compose up --build
```

Then run migrations inside the backend container:

```bash
$ docker compose exec backend sh -c "migrate -path db/migrations -database 'postgres://myuser:yourpassword@db:5432/mydb?sslmode=disable' up"
```

## Changelog

| Version | Description                                                                                                                            |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| 1.0     | Setup Docker multi-stage build and docker-compose orchestration with PostgreSQL & Redis by [ilhammursidi](https://github.com/ilhammursidi) |

## How to Contribute

- Fork this repository
- Create your changes
- Commit your changes (Please strictly follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard: `feat:`, `fix:`, `chore:`, `docs:`)
- Push to the branch
- Open a Pull Request

## License

This project is licensed under the MIT License

## Related Project

[Frontend E-Wallet Repository](https://github.com/Ilhammursidi/React-Ewallet-Project.git)