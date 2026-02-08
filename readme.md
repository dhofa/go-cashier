# Cashier API - Golang

API Kasir sederhana dengan fitur manajemen produk, kategori, keranjang belanja, transaksi, dan laporan penjualan.

## Inisialisasi Proyek
```bash
go mod init go-cashier
go build
```

## Menjalankan Aplikasi (Lokal)
Pastikan file `.env` sudah terisi dengan `PORT` dan `DB_CONN`.
```bash
# Debug mode dengan IPv4 force (opsional)
GODEBUG=netdns=go+v4 go run main.go
```

---

## Dokumentasi API

Semua request body menggunakan format **JSON**.

### 1. Produk (Products)
| Method | Endpoint | Deskripsi | Query Params |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/products` | List produk | `search`, `page`, `limit` |
| `POST` | `/api/products` | Tambah produk | - |
| `GET` | `/api/products/{id}` | Detail produk | - |
| `PUT` | `/api/products/{id}` | Update produk | - |
| `DELETE` | `/api/products/{id}` | Hapus produk | - |

**Contoh Payload (POST/PUT):**
```json
{
  "name": "Kopi Susu",
  "price": 15000,
  "stock": 50,
  "category_id": 1
}
```

### 2. Kategori (Categories)
| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/api/categories` | List kategori |
| `POST` | `/api/categories` | Tambah kategori |
| `GET` | `/api/categories/{id}` | Detail kategori |
| `PUT` | `/api/categories/{id}` | Update kategori |
| `DELETE` | `/api/categories/{id}` | Hapus kategori |

**Contoh Payload (POST/PUT):**
```json
{
  "name": "Minuman",
  "description": "Segala jenis minuman dingin dan panas"
}
```

### 3. Keranjang (Carts)
| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/cart/items` | Tambah ke keranjang (Bisa buat cart baru) |
| `GET` | `/api/carts/{id}` | Lihat isi keranjang |
| `PUT` | `/api/cart/items` | Update quantity item |
| `DELETE` | `/api/carts/{id}/items/{product_id}` | Hapus item dari keranjang |

**Tambah Item (POST):**
Jika `cart_id` adalah `0` atau tidak dikirim, sistem akan membuat cart baru.
```json
{
  "cart_id": 0,
  "product_id": 1,
  "quantity": 2
}
```

**Update Item (PUT):**
```json
{
  "cart_id": 1,
  "product_id": 1,
  "quantity": 5
}
```

### 4. Transaksi & Checkout
| Method | Endpoint | Deskripsi | Query Params |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/checkout` | Proses Checkout | - |
| `GET` | `/api/transactions` | List transaksi | `search`, `page`, `limit` |
| `GET` | `/api/transactions/{id}` | Detail transaksi | - |

**Checkout Payload (POST):**
`cart_items` berisi array of ID item di `cart_items` table.
```json
{
  "cart_id": 1,
  "cart_items": [1, 2, 3]
}
```

### 5. Laporan & Summary
| Method | Endpoint | Deskripsi | Query Params |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/sales/summary` | Ringkasan Penjualan | `start_date`, `end_date` |
| `GET` | `/api/report` | Laporan Lengkap + Produk Terlaris | `start_date`, `end_date` |

**Format Tanggal:** `YYYY-MM-DD` (misal: `2024-02-08`)

---

## Panduan Sales Summary Hari Ini
Untuk mendapatkan data penjualan khusus hari ini, panggil endpoint summary dengan tanggal hari ini pada kedua parameter:
`GET /api/sales/summary?start_date=2024-02-08&end_date=2024-02-08`