---
id: containers
title: "Containers & Deployment"
sidebar_label: "03. Containers & Deployment"
---

# 🟡 Containers & Deployment

Containers are how modern services get from your laptop to production. In this assignment, you'll migrate your Pokedex app from SQLite to PostgreSQL, containerize it with Docker, and orchestrate both services with Docker Compose.

---

## 📋 Prerequisites

This assignment assumes basic knowledge of:
- **Go** (from Stage 1 and previous assignments)
- **Docker** fundamentals

> [!TIP]
> If you need a Docker refresher, read the brief guide to [containers](../../background-and-history/containers.md) and watch [Docker Tutorial for Beginners](https://www.youtube.com/watch?v=3c-iBn73dDE) — a comprehensive 3-hour walkthrough.

---

## 🎯 Goal

- Understand container technology: Docker, Kubernetes, and Helm.
- Be able to create a development environment using containerized services.

---

## 📦 Tools

| Tool | Minimum Version |
| :--- | :--- |
| Go | 1.21+ |
| Docker Engine | 24.0+ |
| Docker Compose | 2.21+ |
| Kubernetes (minikube) | 1.28+ |
| Tilt | 0.33+ |

---

## Assignment 1: Use PostgreSQL Instead of SQLite

You've been using SQLite as your database. In production, we use **PostgreSQL** — a full-featured, client-server database that handles concurrency, replication, and advanced queries.

1. Start a PostgreSQL server using the [official Docker image](https://hub.docker.com/_/postgres).
2. Update your [Pokedex application](./2-pokedex-graphql-sqlite.md) to connect to PostgreSQL instead of SQLite.

> [!NOTE]
> The SQL syntax between SQLite and PostgreSQL is nearly identical for basic operations. The main changes will be in your connection code and driver imports, not your queries.

---

## Assignment 2: Create a Pokedex Docker Image

Build a Docker image for your Pokedex application using a **multi-stage build** to keep the final image small and secure.

### 1. Create a Dockerfile

In your Pokedex application's root directory:

```Dockerfile
# Stage 1: Build the Go binary
FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go build -o pokedex-app

# Stage 2: Create a minimal runtime image
FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/pokedex-app .
CMD ["./pokedex-app"]
```

> [!TIP]
> **Why multi-stage builds?** The Go toolchain and source code are only needed during compilation. The final image only contains the compiled binary, reducing the image size from ~1GB to ~80MB.

### 2. Build the Image

```bash
docker build -t pokedex-app .
```

> [!WARNING]
> You must rebuild the image every time you change your code. This is the most common beginner mistake with containers — running an outdated image and wondering why your changes aren't reflected.

---

## Assignment 3: Deploy with Docker Compose

Docker Compose lets you define and run multi-container applications. Here, you'll orchestrate your Pokedex app alongside PostgreSQL.

### 1. Create `docker-compose.yml`

```yaml
version: '3'
services:
  pokedex-app:
    image: pokedex-app
    ports:
      - '8080:8080'
    depends_on:
      - postgres
  postgres:
    image: postgres:latest
    ports:
      - '5432:5432'
    environment:
      POSTGRES_USER: pokedex
      POSTGRES_PASSWORD: password
      POSTGRES_DB: pokedex_db
```

> [!NOTE]
> The port mapping format is `host:container`. The container port (right side) must match what PostgreSQL listens on internally (5432 by default). The host port (left side) can be any available port on your machine.
>
> You can verify PostgreSQL's default port by checking the [`EXPOSE` directive](https://github.com/docker-library/postgres/blob/2f0ed0c7e8f8b05b294740f150397eec0af8dc50/16/bookworm/Dockerfile) in the official Dockerfile.

### 2. Start the Services

```bash
docker compose up
```

### 3. Verify

Test your API with `curl`:

```bash
curl 'http://localhost:8080/query' \
  -X POST \
  -H 'content-type: application/json' \
  --data-raw '{"query":"query { pokemons { id name description category type { name } abilities { name } } }"}'
```

You can also use GraphQL Playground at [http://localhost:8080/](http://localhost:8080/).

---

## What's Next

Once you've completed these assignments, continue to the Kubernetes training:

**[Next step: Kubernetes with Helm and Tilt](../../container/)**

---

## 📚 Useful Resources

- **[Docker + Go Official Tutorial](https://docs.docker.com/language/golang/build-images/)** — Covers Dockerfiles, building images, and running containers.
- **[10 Minute Docker Compose Tutorial](https://youtu.be/MVIcrmeV_6c)** — Quick walkthrough.
- **[Docker vs VM](https://www.freecodecamp.org/news/docker-vs-vm-key-differences-you-should-know/)** — Optional background reading.

---

## 💡 Note for Beginners

If you're not sure where to start, follow this sequence:

1. **Understand what Docker is and why we need it** — conceptual foundation first.
2. **Watch the [video tutorial](https://www.youtube.com/watch?v=3c-iBn73dDE)** — visual learning helps cement concepts.
3. **Experiment** — Make your own Dockerfile, build an image, run a container. Break things.
4. **Start Assignment 1** — Get PostgreSQL running in a container.
5. **Start Assignment 2** — Turn your Go code into a container image.
6. **Start Assignment 3** — Connect both containers using Docker Compose.

---

## 🧠 Mental Model Check
- **Images are blueprints, containers are instances**: You build an image once and run it many times.
- **Multi-stage builds keep images lean**: Separate build tools from runtime artifacts.
- **Docker Compose orchestrates locally**: It defines how multiple services connect and communicate.
- **Containers mirror production**: What runs in Docker Compose locally runs the same way in Kubernetes.
- **Always rebuild after code changes**: The image captures a snapshot — it doesn't auto-update.
