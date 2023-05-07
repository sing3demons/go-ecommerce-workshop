package main

import (
	"os"

	"github.com/sing3demons/shop/config"
	"github.com/sing3demons/shop/modules/servers"
	"github.com/sing3demons/shop/pkg/database"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env.dev"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := database.DbConnect(cfg.DB())
	// defer db.Close()

	servers.NewServer(cfg, db).Start()

	// type User struct {
	// 	Id       string `json:"id"`
	// 	Email    string `json:"email"`
	// 	Username string `json:"username"`
	// 	RoleId   string `json:"role_id"`
	// 	Password string `json:"password"`
	// }
	// var u User

	// body := User{
	// 	Email:    "sing@dev.com",
	// 	Username: "sing3demons",
	// 	RoleId:   "1",
	// 	Password: "123456",
	// }
	// row := db.QueryRow("INSERT INTO users (email, username, password, role_id) VALUES ($1, $2, $3, $4) RETURNING id", body.Email, body.Username, body.Password, body.RoleId)

	// err = row.Scan(&body.Id)
	// if err != nil {
	// 	log.Fatalf("insert user failed: %v", err)
	// }

	// fmt.Println(body.Id)

	// var users []User

	// query := `SELECT "id", "email", "username", "role_id" FROM "users"`

	// result, err := db.Query(query)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// defer result.Close()

	// for result.Next() {
	// 	err := result.Scan(&u.Id, &u.Email, &u.Username, &u.RoleId)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	users = append(users, u)
	// }

	// fmt.Println(users)
}
