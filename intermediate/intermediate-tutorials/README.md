---
id: intermediate-tutorials
title: "Stage 2: Real-World Services & Architecture"
sidebar_label: "Stage 2: Architecture"
---

# 🟡 Stage 2: Real-World Services & Architecture

Stage 1 gave you the building blocks — Go, SQL, and GraphQL. This stage puts them together into **production-shaped applications**. You'll build a CLI tool, architect a layered service, and deploy with containers.

---

## 🎯 What You Will Learn

By completing this stage, you will:

- Understand the **HTTP protocol** deeply by building a client from scratch.
- Design services using **layered architecture** (persistent layer + API layer).
- Work with **PostgreSQL** as a production database.
- **Containerize** applications with Docker and orchestrate them with Docker Compose.
- Think about software design **before** writing code.

---

## 🗺️ Study Roadmap

1. **[HTTP CLI Tool](./1-http-cmdline.md)**: Build a command-line HTTP client to understand the protocol from the client side.
2. **[Pokedex Service](./2-pokedex-graphql-sqlite.md)**: Design and build a complete web service with Go-Chi, GraphQL, and SQLite — using proper software layers.
3. **[Containers & Deployment](./3-container.md)**: Migrate to PostgreSQL, containerize your app with Docker, and orchestrate with Docker Compose.

---

## 🧠 Mental Model Check
- **Design before code**: Think about layers, schemas, and data flow before writing a single line.
- **Layers separate concerns**: Persistent layer handles data, API layer handles users.
- **Containers are the deployment unit**: What runs locally in Docker Compose runs the same way in production.
- **You should be able to**: Design, build, and deploy a multi-layer backend service — and explain the architectural decisions behind each layer.

Next: [HTTP CLI Tool](./1-http-cmdline.md)
