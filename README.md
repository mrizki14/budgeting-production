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

---

## 📘 Arsitektur & Hosting

Dokumentasi di bawah ini disusun untuk membantu menjelaskan arsitektur deployment proyek ini pada slide presentasi atau dokumentasi sistem.

### 1. Struktur Foldering Vercel & Routing Serverless
Proyek ini mengadopsi model **monorepo hybrid** (React frontend + Go API backend) yang dioptimalkan untuk hosting di Vercel:

* **Folder `/api` (Vercel Serverless Function):**
  * Di Vercel, semua berkas di dalam direktori `/api` akan dikompilasi secara otomatis menjadi *serverless functions*.
  * **[api/index.go](file:///d:/budgeting-monorepo/api/index.go)**: Merupakan satu-satunya titik masuk (*monolithic serverless entrypoint*). Berkas ini mengimplementasikan fungsi `Handler(w http.ResponseWriter, r *http.Request)` yang meneruskan semua lalu lintas HTTP API ke Gin Router (`router.ServeHTTP`).
* **[vercel.json](file:///d:/budgeting-monorepo/vercel.json) (Routing Configuration):**
  * `/api/(.*)` didestinasikan ke `/api/index.go` agar semua panggilan endpoint backend ditangani oleh serverless function Go.
  * `/(.*)` didestinasikan ke `/index.html` agar rute frontend React SPA (Single Page Application) dapat di-render langsung oleh client-side router (mencegah error *404 Not Found* saat halaman di-*refresh*).

---

### 2. Koneksi Database TiDB
Aplikasi terhubung ke database **TiDB** (distributed SQL database) yang fully-compatible dengan MySQL menggunakan GORM:

* **Deteksi & Normalisasi DSN:**
  * Koneksi dibaca dari environment variable `DATABASE_URL` (atau manual `DB_HOST`, `DB_USER`, dll).
  * URL mysql (`mysql://...`) dinormalisasi secara dinamis ke dalam format DSN driver MySQL (`user:pass@tcp(host)/db`).
* **Secure TLS Connection (TiDB Cloud):**
  * Karena TiDB Cloud mewajibkan koneksi terenkripsi, sistem mendeteksi domain `tidbcloud.com` dan otomatis menyetel parameter `tls=tidb`.
  * Fungsi `registerDatabaseTLS` akan mendaftarkan konfigurasi `tls.Config{MinVersion: tls.VersionTLS12}` ke driver MySQL agar koneksi terenkripsi dengan aman.
* **Serverless Connection Pooling:**
  * Di lingkungan serverless, instance dapat tumbuh horizontal dengan cepat. Oleh karena itu, koneksi pool dibatasi sangat ketat (`MaxOpenConns = 2`, `MaxIdleConns = 1`) agar tidak melebihi kapasitas koneksi TiDB.

---


