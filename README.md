

# 🟢 GoAuth

### Warning: Still working .

GoAuth is a **Go authentication library** that provides ready-to-use HTTP handlers for **Fiber, Gin, Echo, and Fasthttp**. It supports **JWT authentication** as well as **OAuth login** via GitHub and Google.

---

## ⚙ How It Works

1. **Database Setup**:
   When initialized, GoAuth automatically ensures that all required tables and indexes exist (users, accounts, sessions, password resets, email verification).

2. **JWT Authentication**:
   If enabled, GoAuth checks that the `GOAUTH_JWT_SECRET` environment variable is set.

3. **OAuth Providers**:
   Optional GitHub and Google OAuth can be configured by passing credentials. GoAuth validates that the necessary environment variables exist.

4. **Framework Handlers**:
   Provides prebuilt HTTP handlers for:

    * Fiber
    * Gin
    * Echo
    * Fasthttp

---

## 📦 Installation

```bash
go get github.com/SwanHtetAungPhyo/go-auth
```

---

## 🛠 Usage Example

```go
package main

import (
	"github.com/SwanHtetAungPhyo/go-auth"
	"github.com/gofiber/fiber/v3"
)

func main() {
	goauthConfig := goauth.NewGoAuth("postgres://user:password@localhost:5432/mydb",
		goauth.WithJwtAuth(true),
		goauth.WithGithubOauth("GITHUB_CLIENT_ID", "GITHUB_CLIENT_SECRET", "http://localhost:3000/callback", ""),
	)

	app := fiber.New()

	// Initialize GoAuth Fiber handler
	goAuthFiberHandler := goauth.NewGOAuthFiberHandler(goauthConfig)

	// Register routes
	app.Post("/api/register", goAuthFiberHandler.Register)
	app.Post("/api/login", goAuthFiberHandler.Login)
	app.Post("/api/logout", goAuthFiberHandler.Logout)
	app.Get("/api/me", goAuthFiberHandler.Profile)

	// Start server
	app.Listen(":3000")
}
```

---

### 🔹 Configuration Options

| Option        | Description                                                                   |
| ------------- | ----------------------------------------------------------------------------- |
| `JWTAuth`     | Enable JWT authentication. Requires `GOAUTH_JWT_SECRET` environment variable. |
| `GithubOauth` | Enable GitHub OAuth login. Requires client ID, secret, and redirect URL.      |
| `GoogleOauth` | Enable Google OAuth login. Requires client ID, secret, and redirect URL.      |
| `DSN`         | Database connection string (Postgres supported).                              |

---

### 🔹 OAuth Environment Variables

**GitHub**

```
GOAUTH_GITHUB_CLIENT_ID
GOAUTH_GITHUB_CLIENT_SECRET
GOAUTH_GITHUB_REDIRECT_URL
```

**Google**

```
GOAUTH_GOOGLE_CLIENT_ID
GOAUTH_GOOGLE_CLIENT_SECRET
GOAUTH_GOOGLE_REDIRECT_URL
```

---

### 🔹 Fiber Handler Methods

| Method     | Description                                |
| ---------- | ------------------------------------------ |
| `Register` | Registers a new user with email/password.  |
| `Login`    | Logs in a user and returns JWT or session. |
| `Logout`   | Logs out a user, invalidating session/JWT. |
| `Profile`  | Returns current logged-in user's profile.  |

> Handlers for **Gin**, **Echo**, and **Fasthttp** follow a similar pattern, using `NewGOAuthGinHandler`, `NewGOAuthEchoHandler`, and `NewGOAuthFastHTTPHandler`.

---

### 🔹 Example: Gin

```go
r := gin.Default()
goAuthGinHandler := goauth.NewGOAuthGinHandler(goauthConfig)

r.POST("/api/register", goAuthGinHandler.Register)
r.POST("/api/login", goAuthGinHandler.Login)
r.GET("/api/me", goAuthGinHandler.Profile)
r.Run(":3000")
```

---

### 🔹 Example: Echo

```go
e := echo.New()
goAuthEchoHandler := goauth.NewGOAuthEchoHandler(goauthConfig)

e.POST("/api/register", goAuthEchoHandler.Register)
e.POST("/api/login", goAuthEchoHandler.Login)
e.GET("/api/me", goAuthEchoHandler.Profile)
e.Start(":3000")
```

---

### 🔹 Example: Fasthttp

```go
app := goauth.NewGOAuthFastHTTPHandler(goauthConfig)
fasthttp.ListenAndServe(":3000", app.Handler)
```

---

### ✅ Key Features

* JWT-based authentication
* OAuth login (GitHub & Google)
* Prebuilt handlers for Fiber, Gin, Echo, Fasthttp
* Automatic database table & index creation
* Modular and extensible for custom middleware or auth providers

---
