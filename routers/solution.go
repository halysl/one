package routers

import (
	"github.com/gin-gonic/gin"
	. "github.com/halysl/one/log"
	"net/http"
	"os"
	"path"
)

func GetMinerInfoView(c *gin.Context) {
	_view := GetfilScountData()
	c.JSON(http.StatusOK, _view)
}

func GetMinerReportDownload(c *gin.Context) {
	minerList := GetfilScountData()
	overview, err := GetTotalNetInfo()
	if err != nil {
		ErrorLogger.Println(err)
	}
	file := renderReport(minerList, overview)
	filePath := "daily-report.xlsx"
	if err = file.SaveAs(filePath); err != nil {
		ErrorLogger.Println(err)
	}
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		c.Redirect(http.StatusFound, "/404")
		ErrorLogger.Println(err)
	}

	fileName:=path.Base(filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")


	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")

	c.File(filePath)
	return
}

func GetOrphanBlockView(c *gin.Context) {
	_view := orphanBlockView(1, 20)
	c.JSON(http.StatusOK, _view)
}

func GetOrphanBlockViewLast5Block(c *gin.Context) {
	_view := orphanBlockView(1, 5)
	c.JSON(http.StatusOK, _view)
}

func GetOrphanBlockViewHuman(c *gin.Context) {
	_view := orphanBlockView(1, 20)
	res := humanReadOrphanBlock(_view)
	c.String(http.StatusOK, res)
}

func GetOrphanBlockViewLast5BlockHuman(c *gin.Context) {
	_view := orphanBlockView(1, 5)
	res := humanReadOrphanBlock(_view)
	c.String(http.StatusOK, res)
}
