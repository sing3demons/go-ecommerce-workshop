package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load .env failed: %v", err)
	}

	return &config{
		app: &app{
			port:    envMap["APP_PORT"],
			host:    envMap["APP_HOST"],
			name:    envMap["APP_NAME"],
			version: envMap["APP_VERSION"],
			readTimeout: func() time.Duration {
				r, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
				if err != nil {
					log.Fatal("load read timeout failed")
				}
				return time.Duration(int64(r) * int64(math.Pow10(9)))
			}(),
			writeTimeout: func() time.Duration {
				r, err := strconv.Atoi(envMap["APP_WRITE_TIMEOUT"])
				if err != nil {
					log.Fatal("load write timeout failed")
				}
				return time.Duration(int64(r) * int64(math.Pow10(9)))
			}(),
			bodyLimit: func() int {
				r, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
				if err != nil {
					log.Fatal("load body limit failed")
				}
				return r
			}(),
			fileLimit: func() int {
				r, err := strconv.Atoi(envMap["APP_FILE_LIMIT"])
				if err != nil {
					log.Fatal("load file limit failed")
				}
				return r
			}(),
		},
		db: &db{
			host:     envMap["DB_HOST"],
			port:     envMap["DB_PORT"],
			username: envMap["DB_USERNAME"],
			password: envMap["DB_PASSWORD"],
			database: envMap["DB_DATABASE"],
			sslmode:  envMap["DB_SSLMODE"],
			maxConnection: func() int {
				r, err := strconv.Atoi(envMap["DB_MAX_CONNECTION"])
				if err != nil {
					log.Fatal("load max connection failed")
				}
				return r
			}(),
		},
		jwt: &jwt{
			adminKey:  envMap["JWT_ADMIN_KEY"],
			secretKey: envMap["JWT_SECRET_KEY"],
			apiKey:    envMap["JWT_API_KEY"],
			accessExpire: func() int {
				r, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRE"])
				if err != nil {
					log.Fatal("load access expire failed")
				}
				return r
			}(),
			refreshExpire: func() int {
				r, err := strconv.Atoi(envMap["JWT_REFRESH_EXPIRE"])
				if err != nil {
					log.Fatal("load refresh expire failed")
				}
				return r
			}(),
		},
	}
}

func (c *config) App() IAppConfig {
	return c.app
}

func (c *config) DB() IDbConfig {
	return c.db
}

func (c *config) JWT() IJwtConfig {
	return c.jwt
}

type IConfig interface {
	App() IAppConfig
	DB() IDbConfig
	JWT() IJwtConfig
}

type IAppConfig interface {
	Url() string
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	FileLimit() int
}

func (a *app) Url() string {
	return fmt.Sprintf("%s:%s", a.host, a.port)
}
func (a *app) Name() string {
	return a.name
}
func (a *app) Version() string {
	return a.version
}
func (a *app) ReadTimeout() time.Duration {
	return a.readTimeout
}
func (a *app) WriteTimeout() time.Duration {
	return a.writeTimeout
}
func (a *app) BodyLimit() int {
	return a.bodyLimit
}
func (a *app) FileLimit() int {
	return a.fileLimit
}

type IDbConfig interface {
	Url() string
	MaxOpenConns() int
}

func (d *db) Url() string {
	// return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", d.username, d.password,d.host,d.port,d.database,d.sslmode)
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.host, d.port, d.username, d.password, d.database, d.sslmode)
}
func (d *db) MaxOpenConns() int {
	return d.maxConnection
}

type IJwtConfig interface {
	SecretKey() []byte
	AdminKey() []byte
	ApiKey() []byte
	AccessExpire() int
	RefreshExpire() int
	SetJwtAccessExpire(int)
	SetJwtRefreshExpire(int)
}

func (j *jwt) SecretKey() []byte {
	return []byte(j.secretKey)
}

func (j *jwt) AdminKey() []byte {
	return []byte(j.adminKey)
}

func (j *jwt) ApiKey() []byte {
	return []byte(j.apiKey)
}

func (j *jwt) AccessExpire() int {
	return j.accessExpire
}

func (j *jwt) RefreshExpire() int {
	return j.refreshExpire
}

func (j *jwt) SetJwtAccessExpire(expire int) {
	j.accessExpire = expire
}

func (j *jwt) SetJwtRefreshExpire(expire int) {
	j.refreshExpire = expire
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

type app struct {
	port         string
	host         string
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int // in bytes
	fileLimit    int // in bytes
}
type db struct {
	host          string
	port          string
	username      string
	password      string
	database      string
	sslmode       string
	maxConnection int
}
type jwt struct {
	adminKey      string
	secretKey     string
	apiKey        string
	accessExpire  int // in seconds
	refreshExpire int // in seconds
}
