Kita bisa buat dokumentasi API yang bisa dibuka seperti:
https://budgeting-production.vercel.app/api/docs
atau kalau mau lebih simpel:
https://budgeting-production.vercel.app/docs
Rekomendasi
Pakai pendekatan OpenAPI YAML manual dulu, bukan generate otomatis dari komentar Go.
Alasannya:
Gratis.
Lebih stabil untuk project kecil-menengah.
Tidak perlu banyak dependency.
Cocok dengan Gin + serverless Vercel.
Dokumentasi bisa langsung dites di Swagger UI dengan tombol “Try it out”.
Lebih mudah dikontrol response, schema, auth bearer token, dan contoh body.
Plan Implementasi
Inventaris endpoint dari router
Ambil semua route dari backend/shared/router/router.go, termasuk:
GET  /api/health
GET  /api/ready
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/me
PUT  /api/settings/profile
PUT  /api/settings/password
GET  /api/dashboard
GET  /api/categories
GET  /api/categories/:id
POST /api/categories
PUT  /api/categories/:id
DELETE /api/categories/:id
GET  /api/transactions
GET  /api/transactions/:id
POST /api/transactions
PUT  /api/transactions/:id
DELETE /api/transactions/:id
GET  /api/budgets
GET  /api/budgets/:id
POST /api/budgets
PUT  /api/budgets/:id
DELETE /api/budgets/:id
GET  /api/reports/summary

Buat file OpenAPI
Tambahkan file misalnya:
docs/openapi.yaml
Isinya mendefinisikan:
servers
securitySchemes untuk JWT Bearer
request body login/register/category/transaction/budget/settings
response sukses dan error
schema reusable seperti User, Category, Transaction, Budget, Dashboard, ApiError

Tambahkan Swagger UI gratis
Ada dua opsi:
Opsi A, paling simpel:
Serve Swagger UI dari CDN lewat file HTML statis:
docs/swagger.html
Opsi B, lebih Go-native:
Pakai package seperti:
github.com/swaggo/gin-swagger
github.com/swaggo/files
Saya lebih pilih Opsi A untuk Vercel serverless karena lebih ringan dan minim risiko deploy.

Expose dokumentasi di Vercel
Kita bisa buat route rewrite agar:
/api/docs
membuka Swagger UI, dan Swagger UI membaca:
/openapi.yaml
Atau lebih sederhana:
/docs
/openapi.yaml

Tambahkan contoh auth di Swagger
Di Swagger UI nanti ada tombol Authorize. User tinggal paste token dari login:
Bearer eyJ...
Setelah itu endpoint protected seperti /dashboard, /categories, /transactions bisa dites langsung dari dokumentasi.

Tambahkan dokumentasi Postman opsional
Dari openapi.yaml, Postman bisa import langsung. Jadi sekali buat OpenAPI, kamu dapat:
Swagger UI
Postman collection
dokumentasi API yang rapi

Verifikasi
Setelah implementasi:
go test ./...
npm run build
Lalu cek:
https://budgeting-production.vercel.app/docs
https://budgeting-production.vercel.app/openapi.yaml

Struktur yang saya sarankan
docs/
  openapi.yaml
  swagger.html
Lalu update vercel.json supaya static docs bisa diakses.
Urutan terbaik
Buat docs/openapi.yaml.
Buat docs/swagger.html.
Update vercel.json.
Test lokal/build.
Deploy.
Import openapi.yaml ke Postman kalau mau collection otomatis.
Saya sarankan kita pakai OpenAPI YAML + Swagger UI CDN. Itu ringan, gratis, cocok untuk Vercel, dan dokumentasinya tetap profesional.