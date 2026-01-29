# 🚀 Dockerizing a Golang App: From 496MB to 28MB using Multi‑Stage Builds

This README explains **how we dockerized a Golang application**, starting with a basic (single‑stage) Dockerfile and then optimizing it using a **multi‑stage Docker build**, resulting in a **~94% reduction in image size**.

This is written as a learning + reference guide, not just instructions.

---

## 🧩 Project Overview

- **Language:** Golang  
- **Go version:** 1.25.x (as defined in `go.mod`)  
- **Entry point:** `cmd/api/main.go`  
- **Port:** 8080  

---

## 🤔 Why Docker?

Docker helps us:
- Run the app consistently across environments
- Avoid “works on my machine” problems
- Package code + dependencies together
- Prepare the app for production deployment

---

## 🟢 Approach 1: Single‑Stage Docker Build (Baseline)

### Dockerfile (Single‑Stage)

```dockerfile
FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/api

EXPOSE 8080
CMD ["./app"]
```

---

### 🔄 Flow (Single‑Stage)

```
┌──────────────────────────┐
│ golang:1.25-alpine       │
│  - Go compiler           │
│  - Go stdlib             │
│  - Source code           │
│  - Dependencies          │
│                          │
│  go build → app binary   │
│  CMD ./app               │
└──────────────────────────┘
```

---

### ❌ Problems with Single‑Stage Build

- Final image contains:
  - Go compiler
  - Go standard library
  - Source code
  - Build tools
- Larger attack surface
- Slower image pulls & deployments

📦 **Measured image size:** ~496 MB  

This works, but it’s **not production‑friendly**.

---

## 🔵 Approach 2: Multi‑Stage Docker Build (Optimized)

Multi‑stage builds separate:
- **Build environment**
- **Runtime environment**

---

### Dockerfile (Multi‑Stage)

```dockerfile
# ---------- BUILD STAGE ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api


# ---------- RUNTIME STAGE ----------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
```

---

### 🔄 Flow (Multi‑Stage Build)

```
STAGE 1: BUILDER
┌─────────────────────────────┐
│ golang:1.25-alpine          │
│  - Go compiler              │
│  - Dependencies             │
│  - Source code              │
│                             │
│  go build → app binary      │
└──────────────┬──────────────┘
               │ (copy binary only)
               ▼
STAGE 2: RUNTIME
┌─────────────────────────────┐
│ alpine:latest               │
│  - No Go                    │
│  - No source code           │
│                             │
│  ./app (executable only)    │
└─────────────────────────────┘
```

---

## 📉 Image Size Comparison (Measured)

| Build Type      | Image Size |
|-----------------|------------|
| Single‑stage    | ~496 MB    |
| Multi‑stage     | **~28 MB** |

✅ **~94% reduction in size**

Same application.  
Same behavior.  
Massive optimization.

---

## 🧠 Why Multi‑Stage Builds Are Better

- Smaller images → faster pulls & deployments
- No source code in production image
- No compiler or build tools in runtime
- Reduced attack surface
- Industry‑standard best practice

---

## ▶️ Build & Run Commands

### Build single‑stage image
```bash
docker build -t social-app .
```

### Build multi‑stage image
```bash
docker build -f Dockerfile.multi -t social-app-multi .
```

### Run the container
```bash
docker run -p 8080:8080 social-app-multi
```

Open in browser:
```
http://localhost:8080
```

---

## 🎯 Key Takeaway

> Multi‑stage Docker builds allow us to use a heavy image for compilation and a lightweight image for runtime, keeping the final image small, secure, and production‑ready.

---

## ✅ Recommendation

- Use **single‑stage Dockerfiles** for learning or quick prototypes
- Use **multi‑stage Dockerfiles** for production systems

---

## 🏁 Final Thought

This optimization is not theoretical — it was **measured**.
Understanding *why* this works is far more valuable than copy‑pasting a Dockerfile.

Happy shipping 🚀
