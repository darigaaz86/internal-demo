package main

import (
	"fmt"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
)

const(
	API = "/api/v1"
	Port = 8080
	Rm1Port = 8081
	Rm2Port = 8082
	dtmServer = "http://localhost:36789/api/dtmsvr"
)

var server1 = fmt.Sprintf("http://localhost:%d%s", Rm1Port, API)
var server2 = fmt.Sprintf("http://localhost:%d%s", Rm2Port, API)

func main() {
	StartSvr()
}

type Trans struct{
	Amount int
}

func StartSvr() {
	app := gin.New()

	app.POST(API+"/Trans", func(c *gin.Context) {
		err := FireDtmRequest()
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		c.JSON(200, "success submit trans req with amount: 30")
	})
	app.POST(API+"/BarrierTransV1", func(c *gin.Context) {
		err := FireDtmBarrierRequestV1()
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		c.JSON(200, "success submit trans req with amount: 30")
	})
	app.POST(API+"/BarrierTransV2F", func(c *gin.Context) {
		err := FireDtmWithoutBarrierRequestV2()
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		c.JSON(200, "success submit trans req with amount: 30")
	})
	app.POST(API+"/BarrierTransV2S", func(c *gin.Context) {
		err := FireDtmWithBarrierRequestV2()
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		c.JSON(200, "success submit trans req with amount: 30")
	})

	app.Run(fmt.Sprintf(":%d", Port))
}

func FireDtmRequest() error {
	req := Trans{Amount: 30}
	gid := shortuuid.New()
	saga := dtmcli.NewSaga(dtmServer, gid).
		Add(server1+"/TransOut", server1+"/TransOutCompensate", req).
		Add(server2+"/TransIn", server2+"/TransInCompensate", req)
	if err := saga.Submit(); err != nil {
		return err
	}

	return nil
}

func FireDtmBarrierRequestV1() error {
	req := Trans{Amount: 30}
	gid := shortuuid.New()
	saga := dtmcli.NewSaga(dtmServer, gid).
		Add(server1+"/BarrierTransOutV1", server1+"/BarrierTransOutV1Compensate", req).
		Add(server2+"/BarrierTransInV1", server2+"/BarrierTransInV1Compensate", req)
	if err := saga.Submit(); err != nil {
		return err
	}

	return nil
}

func FireDtmWithoutBarrierRequestV2() error {
	req := Trans{Amount: 30}
	gid := shortuuid.New()
	saga := dtmcli.NewSaga(dtmServer, gid).
		Add(server1+"/WithoutBarrierTransOutV2", server1+"/BarrierTransOutV2Compensate", req).
		Add(server2+"/BarrierTransInV2", server2+"/BarrierTransInV2Compensate", req)
	saga.Concurrent = true
	if err := saga.Submit(); err != nil {
		return err
	}

	return nil
}

func FireDtmWithBarrierRequestV2() error {
	req := Trans{Amount: 30}
	gid := shortuuid.New()
	saga := dtmcli.NewSaga(dtmServer, gid).
		Add(server1+"/WithBarrierTransOutV2", server1+"/BarrierTransOutV2Compensate", req).
		Add(server2+"/BarrierTransInV2", server2+"/BarrierTransInV2Compensate", req)
	saga.Concurrent = true
	if err := saga.Submit(); err != nil {
		return err
	}

	return nil
}