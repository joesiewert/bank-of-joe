package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

type Account struct {
	Id        int       `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"firstName"`
	LastName  string    `db:"last_name" json:"lastName"`
	Balance   int       `db:"balance" json:"balance"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func index(c *gin.Context) {
	content := gin.H{"Hello": "World"}
	c.JSON(200, content)
}

var dbmap = initDb()

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

func GetAccounts(c *gin.Context) {
	var accounts []Account
	_, err := dbmap.Select(&accounts, "select * from accounts")

	if err == nil {
		c.JSON(200, accounts)
	} else {
		c.JSON(404, gin.H{"error": "no account(s) found"})
	}
}

func GetAccount(c *gin.Context) {
	id := c.Params.ByName("id")
	var account Account
	// err := dbmap.Get(Account{}, id)
	err := dbmap.SelectOne(&account, "select * from accounts where id=$1", id)

	if err == nil {
		account_id, _ := strconv.Atoi(id)

		content := &Account{
			Id:        account_id,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Balance:   account.Balance,
		}
		c.JSON(200, content)
	} else {
		c.JSON(404, gin.H{"error": "account not found"})
	}
}

func CreateAccount(c *gin.Context) {
	var json Account
	c.Bind(&json)

	account := Account{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Balance:   json.Balance,
	}

	if account.FirstName != "" && account.LastName != "" {
		err := dbmap.Insert(&account)
		checkErr(err, "Insert failed")

		if account.FirstName == json.FirstName {
			content := gin.H{
				"result": "Success",
			}
			c.JSON(201, content)
		} else {
			c.JSON(500, gin.H{"result": "An error occured"})
		}
	} else {
		c.JSON(400, gin.H{"error": "fields are empty"})
	}
}

func UpdateAccount(c *gin.Context) {
	id := c.Params.ByName("id")
	var account Account
	err := dbmap.SelectOne(&account, "select * from accounts where id=$1", id)

	if err == nil {
		var json Account
		c.Bind(&json)

		account_id, _ := strconv.Atoi(id)

		account := Account{
			Id:        account_id,
			UpdatedAt: time.Now(),
			FirstName: json.FirstName,
			LastName:  json.LastName,
			Balance:   json.Balance,
		}

		if account.FirstName != "" && account.LastName != "" {
			_, err = dbmap.Update(&account)

			if err == nil {
				c.JSON(200, account)
			} else {
				checkErr(err, "Update failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
}

func DeleteAccount(c *gin.Context) {
	id := c.Params.ByName("id")
	var account Account
	err := dbmap.SelectOne(&account, "select * from accounts where id=$1", id)

	if err == nil {
		_, err = dbmap.Delete(&account)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
}
