package views

import (
	orm "BillSystem/database"
	mode "BillSystem/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	//"time"
	//"fmt"
)

func AlyApportionBill(c *gin.Context) {
	//var date string
	var data []mode.AlyBill
	var total int64
	var payCount float32
	//Year, Month, _ := time.Now().Date()
	//if int(Month) == 1 {
	//	Month = 12
	//	Year = Year - 1
	//	date = fmt.Sprintf("%d-%02d", Year, int(Month))
	//} else {
	//	date = fmt.Sprintf("%d-%02d", Year, int(Month)-1)
	//}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	db := orm.BillDb
	if err := db.Model(&mode.AlyBill{}).Where("apportion_depart IS NOT NULL").Count(&total).Error; err != nil{
		c.JSON(http.StatusOK, gin.H{
			"code" : 500,
			"message" : "查询数据异常",
		})
		return
	}
	offset := (page - 1) * pageSize
	searchKey := c.DefaultQuery("searchKey", "")
	startTime := c.DefaultQuery("startTime", "")
	endTime := c.DefaultQuery("endTime", "")
	if searchKey != "" {
		if startTime != "" {
			if endTime != "" {
				if startTime <= endTime {
					if err := db.Model(&mode.AlyBill{}).Where("apportion_depart = ? AND billing_cycle BETWEEN ? AND ?", searchKey, startTime, endTime).Order("id desc").Limit(pageSize).Offset(offset).Find(&data).Error; err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"code":    400,
							"message": "抱歉未找到相关信息",
						})
						return
					}
					if err := db.Model(&mode.AlyBill{}).Select("sum(payment_amount) as payCount").Where("apportion_depart = ? AND billing_cycle BETWEEN ? AND ?", searchKey, startTime, endTime).Find(&payCount).Error; err != nil {
						c.JSON(http.StatusOK, gin.H{
							"code":    500,
							"message": "查询数据异常",
						})
						return
					}
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    500,
						"message": "开始查询时间应该小于等于结束查询时间",
					})
					return
				}
			}
		} else {
			if err := db.Model(&mode.AlyBill{}).Where("apportion_depart = ?", searchKey).Order("id desc").Limit(pageSize).Offset(offset).Find(&data).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "抱歉未找到相关信息",
				})
				return
			}
			if err := db.Model(&mode.AlyBill{}).Select("sum(payment_amount) as payCount").Where("apportion_depart = ?", searchKey).Find(&payCount).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": "查询数据异常",
				})
				return
			}
		}
	} else if searchKey == "" {
		if startTime != "" {
			if endTime != "" {
				if err := db.Model(&mode.AlyBill{}).Where("billing_cycle BETWEEN ? AND ?", startTime, endTime).Order("id desc").Limit(pageSize).Offset(offset).Find(&data).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    400,
						"message": "抱歉未找到相关信息",
					})
					return
				}
				if err := db.Model(&mode.AlyBill{}).Select("sum(payment_amount) as payCount").Where("billing_cycle BETWEEN ? AND ?", startTime, endTime).Find(&payCount).Error; err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    500,
						"message": "查询数据异常",
					})
					return
				}
			}
		} else {
			if err := db.Model(&mode.AlyBill{}).Where("apportion_depart IS NOT NULL").Order("id desc").Limit(pageSize).Offset(offset).Find(&data).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "抱歉未找到相关信息",
				})
				return
			}
			if err := db.Table("(?) as o", db.Model(&mode.AlyBill{}).Select("payment_amount").Where("apportion_depart IS NOT NULL").Order("id desc").Limit(pageSize).Offset(offset)).Select("sum(payment_amount) as payCount").Find(&payCount).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": "查询数据异常",
				})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"info": map[string]interface{}{
			"data" : data,
			"count": total,
			"payCount": payCount,
			"page" : page,
			"pageSize": pageSize,
		},
	})
}