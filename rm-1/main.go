package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	API = "/api/v1"
	Port = 8081

	DbHost = "localhost"
	DbPort = 5432
	DbUser = "postgres"
	DbPwd = "example"
	DbName = "postgres"
	Query = "update trans_out SET amount = amount + %d WHERE id = 1;"
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
	app.POST(API+"/TransOut", func(c *gin.Context) {
		var trans Trans
		if err := c.ShouldBind(&trans); err != nil {
			c.JSON(400, err.Error())
			return
		}
		query := fmt.Sprintf(Query, -trans.Amount)
		Db.Exec(query)
		log.Printf("TransOut amount: %d", -trans.Amount)
		c.JSON(200, "")
	})
	app.POST(API+"/TransOutCompensate", func(c *gin.Context) {
		var trans Trans
		if err := c.ShouldBind(&trans); err != nil {
			c.JSON(400, err.Error())
			return
		}
		query := fmt.Sprintf(Query, trans.Amount)
		Db.Exec(query)
		log.Printf("TransOutCompensate amount: %d", trans.Amount)
		c.JSON(200, "")
	})

	// barrier
	app.POST(API+"/BarrierTransOutV1", func(c *gin.Context) {
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
			_, err := tx.Exec(query)
			log.Printf("Barrier TransOutV1 amount: %d", -trans.Amount)
			return err
		})
		c.JSON(200, "")
	})
	app.POST(API+"/BarrierTransOutV1Compensate", func(c *gin.Context) {
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
			log.Printf("Barrier TransOutV1Compensate amount: %d", trans.Amount)
			_, err := tx.Exec(query)
			return err
		})
		c.JSON(200, "")
	})

	app.POST(API+"/WithoutBarrierTransOutV2", func(c *gin.Context) {
		var trans Trans
		if err := c.ShouldBind(&trans); err != nil {
			c.JSON(400, err.Error())
			return
		}
		query := fmt.Sprintf(Query, -trans.Amount)
		Db.Exec(query)
		log.Printf("Barrier TransOutV2 amount: %d", -trans.Amount)
		c.JSON(400, "")
	})
	app.POST(API+"/WithBarrierTransOutV2", func(c *gin.Context) {
		barrier, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
		if err != nil {
			panic(err)
		}
		// ret := barrier.CallWithDB(Db, func(tx *sql.Tx) error {
		barrier.CallWithDB(Db, func(tx *sql.Tx) error {
			var trans Trans
			if err := c.ShouldBind(&trans); err != nil {
				c.JSON(400, err.Error())
				return err
			}
			query := fmt.Sprintf(Query, -trans.Amount)
			_, err := tx.Exec(query)
			log.Printf("Barrier TransOutV2 amount: %d", -trans.Amount)
			return err
		})
		// status, res := dtmcli.Result2HttpJSON(ret)
		// c.JSON(status, res)
		c.JSON(400, "")
	})
	app.POST(API+"/BarrierTransOutV2Compensate", func(c *gin.Context) {
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
			log.Printf("Barrier TransOutV2Compensate amount: %d", trans.Amount)
			_, err := tx.Exec(query)
			return err
		})
		c.JSON(200, "")
	})
}

func main() {
	StartSvr()
}