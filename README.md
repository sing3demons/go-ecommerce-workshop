# go-ecommerce-workshop
เรียนรู้พื้นการสร้าง REST API ด้วยภาษา Golang โดยใช้ GOFiber + PostgreSQL

# golang-migrate migrate CLI
```link
https://github.com/golang-migrate/migrate
```
## install 
```migrate curl
curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$os-$arch.tar.gz | tar xvz
```
```brew
brew install golang-migrate
```

## create migration file
```migrate create -ext sql -dir pkg/database/migrations -seq shop_db```

# docker 
```
docker run --name postgres -e POSTGRES_USER=sing3demons -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:15
docker exec -it postgres bash
/# psql -U sing3demons
# CREATE DATABASE go_shop_db;
```

```
migrate -source file://pkg/database/migrations -database 'postgres://sing3demons:password@localhost:5432/go_shop_db?sslmode=disable' -verbose up 
```

```
APP_HOST=127.0.0.1
APP_PORT=3000
APP_NAME=go-shop
APP_VERSION=v0.1.0
APP_BODY_LIMIT=10490000
APP_READ_TIMEOUT=60
APP_WRITE_TIMEOUT=60
APP_API_KEY=
APP_ADMIN_KEY=
APP_FILE_LIMIT=2097000

JWT_SECRET_KEY=
JWT_ACCESS_EXPIRE=86400
JWT_REFRESH_EXPIRE=604800

DB_HOST=127.0.0.1
DB_PORT=5432
DB_PROTOCOL=tcp
DB_USERNAME=
DB_PASSWORD=
DB_DATABASE=
DB_SSLMODE=disable
DB_MAX_CONNECTION=25
```
