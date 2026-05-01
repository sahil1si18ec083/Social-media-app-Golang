# Redis Notes

This folder contains the Redis cache layer used by the app.

## Popular Redis Commands

List all keys:

```bash
KEYS *
```

Get a value by key:

```bash
GET user:92
```

Check how long a key will live:

```bash
TTL user:92
```

Delete one key:

```bash
DEL user:92
```

Delete multiple keys:

```bash
DEL user:92 user:80
```

Check whether a key exists:

```bash
EXISTS user:92
```

Clear the current Redis database:

```bash
FLUSHDB
```

Clear all Redis databases:

```bash
FLUSHALL
```

Be careful with `FLUSHDB` and `FLUSHALL`.

## Run Redis With Docker

Start only Redis from Docker Compose:

```bash
docker compose up -d redis
```

Start Redis and Redis Commander:

```bash
docker compose up -d redis redis-commander
```

Check running containers:

```bash
docker ps
```

Stop Redis:

```bash
docker compose stop redis
```

## Open Redis CLI

Open Redis CLI inside the running container:

```bash
docker exec -it redis-cache redis-cli
```

Meaning:

- `docker exec` runs a command inside an already running container
- `-it` opens an interactive terminal
- `redis-cache` is the container name
- `redis-cli` is the Redis command-line client

## Check Cached Values

After opening `redis-cli`, run:

```bash
KEYS *
GET user:92
TTL user:92
```

You can also run commands directly without entering the Redis shell:

```bash
docker exec -it redis-cache redis-cli KEYS '*'
```

```bash
docker exec -it redis-cache redis-cli GET user:92
```

```bash
docker exec -it redis-cache redis-cli TTL user:92
```

## Redis Commander

If Redis Commander is running, open:

```text
http://localhost:8081
```

This gives a browser UI to inspect keys and values.
