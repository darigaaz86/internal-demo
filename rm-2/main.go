package main

import (
	"fmt"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/dtm-labs/client/dtmcli"
	_ "github.com/lib/pq"
)

const (
	API = "/api/v1"
	Port = 8082

	DbHost = "localhost"
	DbPort = 5432
	DbUser = "postgres"
	DbPwd = "example"
	DbName = "postgres"
	Query = "update trans_in SET amount = amount + %d WHERE id = 1;"
)

var Db *sql.DB = DB()

func DB() (db *sql.DB) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DbHost, DbPort, DbUser, DbPwd, DbName)
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        panic(err)
    }
    return
}

type Trans struct{
	Amount int
}

func StartSvr() {
	app := gin.New()
	dtmcli.SetBarrierTableName("public.barrier")
	dtmcli.SetCurrentDBType("postgres")
	AddRoute(app)
	app.Run(fmt.Sprintf(":%d", Port))
}

func AddRoute(app *gin.Engine) {
	app.POST(API+"/TransIn", func(c *gin.Context) {
		var trans Trans
		if err := c.ShouldBind(&trans); err != nil {
			c.JSON(400, err.Error())
			return
		}
		query := fmt.Sprintf(Query, trans.Amount)
		Db.Exec(query)
		log.Printf("TransIn amount: %d", trans.Amount)
		c.JSON(200, "")
		// c.JSON(409, "") // Status 409 for Failure. Won't be retried
	})
	app.POST(API+"/TransInCompensate", func(c *gin.Context) {
		var trans Trans
		if err := c.ShouldBind(&trans); err != nil {
			c.JSON(400, err.Error())
			return
		}
		query := fmt.Sprintf(Query, -trans.Amount)
		Db.Exec(query)
		log.Printf("TransInCompensate amount: %d", -trans.Amount)
		c.JSON(200, "")
	})

	// barrier
	app.POST(API+"/BarrierTransInV1", func(c *gin.Context) {
		barrier, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
		if err != nil {
			panic(err)
		}
		barrier.CallWithDB(Db, func(tx *sql.Tx) error {
			var trans Trans
			if err := c.ShouldBind(&trans); err != nil {
				c.JSON(400, err.Error())
				return err
			}
			query := fmt.Sprintf(Query, trans.Amount)
			log.Printf("Barrier TransInV1 amount: %d", trans.Amount)
			_, err := tx.Exec(query)
			return err
		})
		c.JSON(200, "")
	})
	app.POST(API+"/BarrierTransInV1Compensate", func(c *gin.Context) {
		barrier, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
		if err != nil {
			panic(err)
		}
		barrier.CallWithDB(Db, func(tx *sql.Tx) error {
			var trans Trans
			if err := c.ShouldBind(&trans); err != nil {
				c.JSON(400, err.Error())
				return err
			}
			query := fmt.Sprintf(Query, -trans.Amount)
			log.Printf("Barrier TransInV1Compensate amount: %d", -trans.Amount)
			_, err := tx.Exec(query)
			return err
		})
		c.JSON(200, "")
	})

	app.POST(API+"/BarrierTransInV2", func(c *gin.Context) {
		var trans Trans
		if err := c.ShouldBind(&trans); err != nil {
			c.JSON(400, err.Error())
			return
		}
		query := fmt.Sprintf(Query, trans.Amount)
		Db.Exec(query)
		log.Printf("Barrier TransInV2 amount: %d", trans.Amount)
		c.JSON(200, "")
	})
	app.POST(API+"/BarrierTransInV2Compensate", func(c *gin.Context) {
		barrier, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
		if err != nil {
			panic(err)
		}
		barrier.CallWithDB(Db, func(tx *sql.Tx) error {
			var trans Trans
			if err := c.ShouldBind(&trans); err != nil {
				c.JSON(400, err.Error())
				return err
			}
			query := fmt.Sprintf(Query, -trans.Amount)
			log.Printf("Barrier TransInV2Compensate amount: %d", -trans.Amount)
			_, err := tx.Exec(query)
			return err
		})
		c.JSON(200, "")
	})
}

func main() {
	StartSvr()
}
