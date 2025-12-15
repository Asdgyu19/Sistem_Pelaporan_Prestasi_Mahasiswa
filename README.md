# Sistem Pelaporan Prestasi Mahasiswa - Backend API

**Nama:** Aryo Prabowo  
**NIM:** 434231027  
**Mata Kuliah:** Pemrograman Backend Lanjut (Praktikum)

---

## Tentang Aplikasi

Sistem Pelaporan Prestasi Mahasiswa adalah sebuah REST API backend yang dirancang untuk mengelola pencatatan dan verifikasi prestasi akademik mahasiswa. Aplikasi ini memfasilitasi mahasiswa dalam melaporkan prestasi mereka, dosen wali dalam memverifikasi laporan, dan administrator dalam mengelola keseluruhan sistem.

Fitur utama mencakup:
- Sistem autentikasi dan otorisasi berbasis peran (mahasiswa, dosen wali, admin)
- Manajemen pencatatan prestasi dengan workflow verifikasi
- Penyimpanan file pendukung (sertifikat, bukti prestasi)
- Laporan dan statistik performa mahasiswa
- Manajemen data pengguna dan pembimbing akademik

---

## Stack Teknologi

**Backend:**
- **Go 1.24** - Bahasa pemrograman utama
- **Gin** - Web framework dengan performa tinggi
- **PostgreSQL** - Database relasional untuk data terstruktur
- **MongoDB** - NoSQL database untuk penyimpanan file dan data fleksibel
- **JWT** - Autentikasi token-based

**Tools & Libraries:**
- **Swagger/OpenAPI** - Dokumentasi API interaktif
- **swaggo/swag** - Generator dokumentasi otomatis
- **golang-jwt** - Library JWT authentication
- **lib/pq** - PostgreSQL driver
- **mongo-go-driver** - MongoDB client

---

## Struktur Folder

```
.
├── app/                          # Konfigurasi aplikasi utama
├── config/                       # File konfigurasi (database, server, JWT)
├── database/                     # Setup koneksi database
│   └── migrations/              # SQL migration files
├── docs/                         # Swagger documentation (auto-generated)
├── helper/                       # HTTP handler functions
│   ├── auth_helper.go           # Login, register, profile
│   ├── achievement_helper.go    # CRUD prestasi dan verifikasi
│   ├── student_helper.go        # Data mahasiswa
│   ├── lecturer_helper.go       # Data dosen wali
│   ├── report_helper.go         # Laporan dan statistik
│   ├── admin_user_helper.go     # Manajemen user (admin)
│   └── user_helper.go           # Profil user
├── middleware/                   # Middleware (auth, logging, error handling)
├── route/                        # Routing dan endpoint configuration
├── service/                      # Business logic layer
├── utils/                        # Utility functions (JWT, password hashing)
├── main.go                       # Entry point aplikasi
├── go.mod & go.sum              # Dependency management
└── .env.example                 # Template environment variables
```

---

## Instalasi & Setup

### Prerequisites
- Go 1.24 atau lebih tinggi
- PostgreSQL 12+
- MongoDB 4.4+
- Git

### Langkah Instalasi

1. **Clone Repository**
```bash
git clone https://github.com/Asdgyu19/Sistem_Pelaporan_Prestasi_Mahasiswa.git
cd Sistem_Pelaporan_Prestasi_Mahasiswa
```

2. **Setup Environment Variables**
```bash
cp .env.example .env
# Edit .env dengan konfigurasi database dan JWT secret Anda
```

3. **Install Dependencies**
```bash
go mod download
go mod tidy
```

4. **Setup Database**
```bash
# PostgreSQL
createdb prestasi_mahasiswa

# Jalankan migration (opsional - bisa via script atau manual)
psql prestasi_mahasiswa < database/migrations/001_init.sql
```

5. **Build & Run**
```bash
go build -o prestasi-mahasiswa.exe
./prestasi-mahasiswa.exe
```

Server akan berjalan di `http://localhost:8080`

---

## API Documentation

Dokumentasi lengkap API tersedia melalui Swagger UI:

```
http://localhost:8080/swagger/index.html
```

### Endpoint Utama

#### Authentication (5.1)
- `POST /auth/login` - Login dengan email dan password
- `POST /auth/register` - Daftar akun baru
- `POST /auth/logout` - Logout dan revoke token
- `GET /auth/profile` - Ambil profil user saat ini
- `POST /auth/refresh` - Refresh access token

#### Prestasi (5.4)
- `GET /achievements` - Daftar prestasi (filtered by role)
- `POST /achievements` - Buat prestasi baru (mahasiswa)
- `GET /achievements/:id` - Detail prestasi
- `PUT /achievements/:id` - Edit prestasi (mahasiswa)
- `DELETE /achievements/:id` - Hapus prestasi (mahasiswa)
- `POST /achievements/:id/submit` - Submit untuk verifikasi
- `POST /achievements/:id/verify` - Verifikasi (dosen/admin)
- `POST /achievements/:id/reject` - Tolak (dosen/admin)
- `POST /achievements/:id/files` - Upload file pendukung
- `GET /achievements/:id/files` - Lihat file
- `DELETE /achievements/:id/files/:fileId` - Hapus file

#### Mahasiswa (5.5)
- `GET /students` - Daftar semua mahasiswa
- `GET /students/:id` - Detail mahasiswa
- `GET /students/:id/achievements` - Prestasi mahasiswa
- `PUT /students/:id/advisor` - Tentukan pembimbing

#### Dosen (5.5)
- `GET /lecturers` - Daftar dosen wali
- `GET /lecturers/:id` - Detail dosen
- `GET /lecturers/:id/advisees` - Mahasiswa bimbing

#### Laporan (5.8)
- `GET /reports/statistics` - Statistik sistem
- `GET /reports/student/:id` - Laporan mahasiswa

#### User Management (5.2)
- `GET /admin/users` - Daftar user (admin)
- `GET /admin/users/:id` - Detail user
- `POST /admin/users` - Buat user
- `PUT /admin/users/:id` - Edit user
- `DELETE /admin/users/:id` - Hapus user
- `PUT /admin/users/:id/role` - Ubah role

---

## Autentikasi

Sistem menggunakan **JWT Bearer Token**. Untuk akses endpoint yang dilindungi:

```bash
Authorization: Bearer <your_jwt_token>
```

### Role & Permission

- **mahasiswa**: Buat dan kelola prestasi mereka sendiri
- **dosen_wali**: Verifikasi atau tolak prestasi mahasiswa bimbing
- **admin**: Kelola user, role, dan akses semua fitur

---

## Environment Variables

Buat file `.env` berdasarkan template:

```env
# Server
SERVER_PORT=8080
SERVER_MODE=debug

# Database - PostgreSQL
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=prestasi_mahasiswa

# Database - MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=prestasi_mahasiswa

# JWT
JWT_SECRET=your_super_secret_key_min_32_chars
JWT_EXPIRE_HOURS=24

# File Upload
MAX_FILE_SIZE=10485760
UPLOAD_DIR=./uploads
```

---

## Development & Testing

### Compile
```bash
go build -o prestasi-mahasiswa.exe
```

### Generate Swagger Docs
```bash
swag init -g main.go
```

### Run Tests (jika ada)
```bash
go test ./...
```

---

## Database Schema

### Tabel Utama

**users**
- id, nim, name, email, password (hashed), role, advisor_id, is_active

**achievements**
- id, mahasiswa_id, title, description, category, status, achievement_date, verified_by, created_at

**achievement_files**
- id, achievement_id, file_name, file_path, uploaded_by, created_at

**refresh_tokens**
- id, user_id, token_hash, expires_at, is_revoked, ip_address, user_agent

---

## Deployment

### Production Checklist

- [ ] Set `SERVER_MODE=release` di .env
- [ ] Gunakan password database yang kuat
- [ ] Konfigurasi JWT_SECRET dengan nilai random yang panjang
- [ ] Setup HTTPS/SSL certificate
- [ ] Backup database secara berkala
- [ ] Monitor error logs dan performance

### Deploy ke Server

```bash
# Build untuk production
go build -o prestasi-mahasiswa ./

# Run dengan nohup atau supervisor
./prestasi-mahasiswa &
```

---

## Troubleshooting

**Database Connection Error**
- Pastikan PostgreSQL dan MongoDB sudah running
- Cek konfigurasi host, port, dan credentials di .env
- Verify firewall rules

**Token Invalid/Expired**
- Pastikan JWT_SECRET sama di setup dan runtime
- Gunakan endpoint `/auth/refresh` untuk mendapat token baru
- Check expiration time di JWT_EXPIRE_HOURS

**File Upload Error**
- Pastikan folder `uploads/` tersedia dan writable
- Check MAX_FILE_SIZE sesuai kebutuhan
- Verify disk space

---

## License

Proyek ini merupakan tugas akhir mata kuliah Pemrograman Backend Lanjut.

---

## Contact

**Aryo Prabowo** | 434231027  
Advanced Backend Programming - Praktikum

