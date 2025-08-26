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

---

## Setup

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
