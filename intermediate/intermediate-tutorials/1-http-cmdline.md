---
id: http-cli
title: "HTTP CLI Tool"
sidebar_label: "01. HTTP CLI Tool"
---

# 🟡 HTTP CLI Tool

Building an HTTP client from scratch forces you to understand the HTTP protocol at a fundamental level — methods, headers, query parameters, and request bodies. This knowledge is essential for debugging API issues, writing integration tests, and understanding how your backend services communicate.

---

## 🎯 Goal

Create a command-line HTTP client program (`httpcli`) using Go.

---

## 📋 Requirements

1. Use Go as the programming language.
2. Program name: `httpcli`
3. Implement the following sub-commands:
   - `get` — Send a GET request to a given URL.
   - `post` — Send a POST request to a given URL.
   - `put` — Send a PUT request to a given URL.
   - `delete` — Send a DELETE request to a given URL.
4. Calling the program without a sub-command should default to a GET request.
5. Print the response to STDOUT.

---

## 🔧 Command Specification

### Global Flags

These flags must be available on all sub-commands:

| Flag | Description | Example |
| :--- | :--- | :--- |
| `--help` | Print usage information | `httpcli --help` |
| `--query` | Add query parameters (repeatable) | `httpcli get example.com --query key1=val1 --query key2=val2` |
| `--header` | Add request headers (repeatable) | `httpcli get example.com --header Authorization=Bearer_xyz` |

### Root Command

```sh
httpcli [SUB-COMMAND] <URL> [FLAGS...]
```

Calling without a sub-command sends a GET request: `httpcli google.com`

### `get` Sub-command

```sh
httpcli get <URL> [FLAGS...]
```

### `post` Sub-command

```sh
httpcli post <URL> [FLAGS...]
```

| Flag | Description |
| :--- | :--- |
| `--json` | JSON request body (must validate input) |

```sh
httpcli post example.com --json '{ "key": "value" }'
```

### `put` Sub-command

```sh
httpcli put <URL> [FLAGS...]
```

| Flag | Description |
| :--- | :--- |
| `--json` | JSON request body (must validate input) |

```sh
httpcli put example.com --json '{ "key": "value" }'
```

### `delete` Sub-command

```sh
httpcli delete <URL> [FLAGS...]
```

---

## 📦 Recommended Libraries

> [!TIP]
> Pick one CLI framework and stick with it. Both are excellent — Cobra is more popular in the Go ecosystem; urfave/cli is simpler for smaller tools.

- **CLI Frameworks** (pick one):
  - [urfave/cli](https://cli.urfave.org/) — Lightweight and straightforward.
  - [Cobra](https://cobra.dev/) — Feature-rich, used by kubectl and Hugo.
- **HTTP**: [net/http](https://pkg.go.dev/net/http) — Go's standard library is sufficient for this task.

---

## 📚 Useful Resources

- **[What are REST APIs?](https://www.youtube.com/watch?v=SLwpqD8n3d0)** — Visual explanation of REST concepts.
- **[REST APIs in 100 Seconds](https://www.youtube.com/watch?v=-MTSQjw5DrM)** (Fireship) — Quick overview.
- **[HTTP Protocol Basics](https://www.ibm.com/docs/en/cics-ts/5.3?topic=concepts-http-protocol)** — IBM reference guide.

> [!TIP]
> Before building your own, familiarize yourself with existing HTTP clients like [HTTPie](https://github.com/httpie/httpie). Understanding good CLI UX will make your tool better.

---

## 🧠 Mental Model Check
- **HTTP is a request-response protocol**: Every interaction has a method, URL, headers, and an optional body.
- **Methods express intent**: GET reads, POST creates, PUT updates, DELETE removes.
- **Headers carry metadata**: Authentication, content type, caching directives.
- **Query parameters filter**: They're part of the URL and visible to everyone — never put secrets in them.

Next: [Pokedex Service](./2-pokedex-graphql-sqlite.md)
