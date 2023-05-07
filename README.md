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
