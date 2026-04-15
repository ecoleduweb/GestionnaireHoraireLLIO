package controllers

import (
	"llio-api/services"

	"github.com/gin-gonic/gin"
)

func ExportExcel(c *gin.Context) {
	// Headers download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=activities.xlsx")
	c.Header("Transfer-Encoding", "chunked")

    err := services.GenerateExcel(c.Writer)
    if err != nil {
        handleError(c, err, "ExportExcel")
        return
    }
}