package views

import (
	orm "BillSystem/database"
	mode "BillSystem/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CenterDepart(c *gin.Context) {
	var data []mode.CenterDepart
	db := orm.BillDb
	if err := db.Model(&mode.CenterDepart{}).Order("id desc").Find(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "抱歉未找到相关信息",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"info": map[string]interface{}{
			"data" : data,
		},
	})
}