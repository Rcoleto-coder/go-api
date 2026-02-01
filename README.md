# Go API

A Go HTTP API with user authentication. It supports **sign up** and **login**, using JWT access and refresh tokens and MongoDB for user storage.

## Features

- **Sign up** (`POST /register`) — Create an account with email and password (passwords are hashed and stored in MongoDB).
- **Login** (`POST /login`) — Authenticate and receive an access token; a refresh token is set via HTTP-only cookie.
- **Refresh** (`POST /refresh`) — Issue a new access token using the refresh token cookie.
- **Protected route** (`GET /`) — Example endpoint that requires a valid JWT in the `Authorization` header.

## Tech

- **Go** — Standard library `net/http`, no framework.
- **MongoDB** — User storage.
- **JWT** — Access and refresh tokens.
- **CORS** — Enabled for cross-origin requests.

## Setup

1. Copy `.env.example` to `.env` and set:
   - `PORT` (default `8080`)
   - `MONGO_URI` — MongoDB connection string
   - `DB_NAME` — Database name
   - `JWT_SECRET` — Secret for signing tokens

2. Run:
   ```bash
   go run ./cmd/server
   ```

The API will listen on `http://localhost:<PORT>`.
