# Budgeting App Frontend

Standalone React frontend for the Go Budgeting App API. The interface ports
the Laravel Blade pages without redesigning them.

## Setup

1. Copy `.env.example` to `.env`.
2. Ensure `VITE_API_BASE_URL` points to the running Go API.
3. Install dependencies with `rtk npm install`.
4. Start development with `rtk npm run dev`.

The default URLs are:

- React: `http://127.0.0.1:5173`
- Go API: `http://127.0.0.1:8080/api`

## Verification

```bash
rtk npm test
rtk npm run build
```

Production hosting must serve `index.html` for unknown frontend paths so that
direct navigation and refresh work on React routes.
