package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/halysl/one/db"
	. "github.com/halysl/one/log"
	"github.com/halysl/one/routers"
)

func main() {
	initDB := flag.Bool("initdb", false, "init db")
	address := flag.String("address", "0.0.0.0:8080", "assign service ip and port,default 0.0.0.0:8080")
	flag.Parse()
	if *initDB {
		db.InitDB()
	}

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	minerRouter := router.Group("/miner")
	{
		minerRouter.GET("/info/view", routers.GetMinerInfoView)
		minerRouter.GET("/report/download", routers.GetMinerReportDownload)
	}

	orphanBlockRouter := router.Group("/orphanblock")
	{
		orphanBlockRouter.GET("/view/all", routers.GetOrphanBlockView)
		orphanBlockRouter.GET("/view/all/human", routers.GetOrphanBlockViewHuman)
		orphanBlockRouter.GET("/view/last5block", routers.GetOrphanBlockViewLast5Block)
		orphanBlockRouter.GET("/view/last5block/human", routers.GetOrphanBlockViewLast5BlockHuman)
	}

	if err := router.Run(*address); err != nil {
		ErrorLogger.Println(err)
		return
	}
}
