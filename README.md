# 🚀 Go E-Commerce API (Fiber)

Bu proje, **Go Fiber** framework'ü ve **Clean Architecture** prensipleriyle geliştirilmiş, Docker desteğine sahip profesyonel bir e-ticaret backend çözümüdür. Ölçeklenebilir, test edilebilir ve modern standartlara uygun bir yapı sunar.

## 🛠 Teknolojik Stack
- **Dil:** Go 1.24
- **Framework:** Fiber v2
- **ORM:** GORM
- **Logging:** Uber-Zap Logger
- **Geliştirme:** Air (Hot-Reload)
- **Konteynerleştirme:** Docker & Docker Compose
- **Dokümantasyon:** Swagger (OpenAPI 3.0)

## 🚀 Çalıştırma Talimatları

### A. Docker ile Çalıştırma (Tavsiye Edilen)
Projeyi tüm bağımlılıkları ve veritabanı ile birlikte ayağa kaldırmak için:
**make up**

### B. Manuel Çalıştırma
Eğer yerel makinenizde Go yüklüyse:
1. **go mod tidy**
2. **go run cmd/api/main.go**

## 📜 Makefile Komutları
Geliştirme sürecini yönetmek için aşağıdaki kısa yolları kullanabilirsiniz:
- **make up**: Konteynerleri build eder ve arka planda başlatır.
- **make down**: Konteynerleri durdurur ve tüm kaynakları siler.
- **make restart**: Uygulama konteynerini yeniden başlatır.
- **make logs**: Uygulama loglarını (Zap) canlı olarak izler.
- **make test**: Projedeki tüm birim (unit) testlerini çalıştırır.

## 📁 Proje Yapısı
Proje, bağımlılıkların yönetimi için katmanlı mimari (Clean Architecture) kullanmaktadır:

```text
.
├── cmd/api             # Uygulama giriş noktası (main.go)
├── internal/           # Uygulama çekirdek kodları
│   ├── delivery/http   # Handler'lar, Rotalar ve Middleware'ler
│   ├── domain/         # Modeller (Product, Order, User vb.) ve Interface'ler
│   ├── repository/     # Veritabanı erişim katmanı (GORM)
│   └── service/        # İş mantığı (Business Logic)
├── pkg/                # Yardımcı araçlar
│   └── utils/          # Logger (Zap), Pagination vb.
├── docs/               # Otomatik üretilen Swagger (OpenAPI) dokümanları
├── Dockerfile          # Docker build yapılandırması
├── docker-compose.yml  # Servislerin (App, DB) orkestrasyonu
├── Makefile            # Otomasyon komutları
└── .air.toml           # Air (Hot-reload) yapılandırması
```

## 🔒 Güvenlik ve Yetkilendirme
- JWT: Kullanıcı login işlemleri sonrası verilen token ile yetkilendirme sağlanır.
- RBAC (Role Based Access Control): Admin, Seller ve Buyer rolleri için özel middleware kontrolleri (RoleCheck) uygulanmaktadır.

## 📡 API Dokümantasyonu
Uygulama çalıştıktan sonra aşağıdaki adresten interaktif Swagger dokümantasyonuna erişebilirsiniz:
👉 http://localhost:8080/swagger/index.html
