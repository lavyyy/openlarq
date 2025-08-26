# OpenLARQ

**An open-source, self-hosted API + web dashboard for your LARQ products.**

LARQ currently only provides a mobile app and no official API. OpenLARQ reverse-engineers the app to give you both a **REST API** and a **web dashboard** to track your hydration data — all running locally and under your control.

> This is an **unofficial project** and is not affiliated with LARQ.

![openlarq-preview](https://cdn.barking.dev/openlarq-preview.png)

---

## Features

- Web dashboard for viewing your stats over time
- Self-hosted setup — your data stays private
- Open-source and extendable

---

## Requirements

- Go `1.23`
- Node.js `22`
- Docker (optional)

## Quick Start (Docker)

The easiest way to get started is using Docker Compose, which will set up both the API backend and web frontend automatically.

### Prerequisites

- Docker installed on your system
- LARQ account credentials

### Setup with Docker Compose

1. Clone the repository:

   ```bash
   git clone https://github.com/lavyyy/openlarq.git
   cd openlarq
   ```

2. Configure the environment credentials in the `.docker-compose.yml`:

   ```env
   LARQ_USERNAME=your_username
   LARQ_PASSWORD=your_password
   ```

3. Start all services:

   ```bash
   docker-compose up -d
   ```

4. Access your services:

   - **Web Dashboard**: http://localhost:3000
   - **API Backend**: http://localhost:8080

5. To stop all services:
   ```bash
   docker-compose down
   ```

### Docker Compose Services

- **`api`**: Go backend API server (port 8080)
- **`web`**: Svelte frontend dashboard (port 3000)
- **`db`**: Optional PostgreSQL database (commented out by default)

### Development with Docker Compose

For development, you can use the volume mounts to enable hot reloading:

```bash
# Start in development mode with live reloading
docker-compose up

# View logs
docker-compose logs -f

# Rebuild and restart a specific service
docker-compose up --build api
docker-compose up --build web
```

---

## Setup (Manual)

1. Clone the repository:

   ```bash
   git clone https://github.com/lavyyy/openlarq.git
   cd openlarq
   ```

2. Configure your LARQ credentials:

   ```env
   LARQ_USERNAME=your_username
   LARQ_PASSWORD=your_password
   ```

3. Run the API server:

   ```bash
   make run dev
   ```

4. Start the web dashboard:

   ```bash
   cd site
   pnpm install
   pnpm run dev
   ```

5. Open your browser and go to:
   ```
   http://localhost:3000
   ```

---

## Configuration

| Variable        | Description                      | Required |
| --------------- | -------------------------------- | -------- |
| `LARQ_USERNAME` | Your LARQ account email/username | ✅       |
| `LARQ_PASSWORD` | Your LARQ account password       | ✅       |
| `PORT`          | API server port (default: 8080)  | ❌       |

---

## Contributing

Contributions are welcome! Feel free to open issues, suggest features, or submit pull requests.

---

## License

OpenLARQ is under the MIT License. See [here](LICENSE) for more info.

---

## Disclaimer

This project is **not affiliated, associated, authorized, endorsed by, or in any way officially connected** with LARQ.
