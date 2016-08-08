package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

type Account struct {
	Id        int64  `db:"id" json:"id"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
	Balance   int64  `db:"balance" json: "balance"`
}

func index(c *gin.Context) {
	content := gin.H{"Hello": "World"}
	c.JSON(200, content)
}

func main() {
	router := gin.Default()
	router.GET("/", index)
	v1 := router.Group("api/v1")
	{
		v1.GET("/accounts", GetAccounts)
		v1.GET("/accounts/:id", GetAccount)
		v1.POST("/accounts", CreateAccount)
		v1.PUT("/accounts/:id", UpdateAccount)
		v1.DELETE("/accounts/:id", DeleteAccount)
	}
	router.Run()
}

func initDb() *gorp.DbMap {
	dburl := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dburl)
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Account{}, "accounts").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
