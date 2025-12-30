FROM golang:1.24-alpine

# CGO_ENABLED=0 ile statik binary üretiyoruz, alpine için gerekli
ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

# Sadece go.mod ve go.sum kopyalayarak layer caching avantajını kullanıyoruz
COPY go.mod go.sum ./
RUN go mod download

# Air (hot-reload) yüklemesi
RUN go install github.com/cosmtrek/air@v1.52.0

COPY . .

# Air ile uygulamayı başlatıyoruz
CMD ["air", "-c", ".air.toml"]