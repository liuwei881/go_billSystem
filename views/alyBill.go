package views

import (
	orm "BillSystem/database"
	mode "BillSystem/models"
	"BillSystem/tools"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func AlyBill(c *gin.Context) {
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
	offset := (page - 1) * pageSize
	if err := db.Model(&mode.AlyBill{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "查询数据异常",
		})
		return
	}
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
			if err := db.Model(&mode.AlyBill{}).Order("id desc").Limit(pageSize).Offset(offset).Find(&data).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "抱歉未找到相关信息",
				})
				return
			}
			if err := db.Table("(?) as o", db.Model(&mode.AlyBill{}).Select("payment_amount").Order("id desc").Limit(pageSize).Offset(offset)).Select("sum(payment_amount) as payCount").Find(&payCount).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": "查询数据异常",
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"info": map[string]interface{}{
				"data":     data,
				"count":    total,
				"payCount": payCount,
				"page":     page,
				"pageSize": pageSize,
			},
		})
	}
}

func TestRoot(c *gin.Context){
	c.String(200, "Hello, World")
}

func TestInfo(c *gin.Context){
	var user []mode.Fest
	var total int64
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	db := orm.BillDb
	if err := db.Model(&mode.Fest{}).Count(&total).Error; err != nil{
		c.JSON(http.StatusOK, gin.H{
			"code" : 500,
			"message" : "查询数据异常",
		})
		return
	}
	offset := (page - 1) * pageSize
	if err := db.Model(&mode.Fest{}).Order("id desc").Limit(pageSize).Offset(offset).Find(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"message": "抱歉未找到相关信息",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"info": map[string]interface{}{
			"data" : user,
			"total": total,
			"page" : page,
			"pageSize": pageSize,
		},
	})
}

func TestPost(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"info": map[string]interface{}{
			"data" : data,
			"total": len(data),
		},
	})
}

func TestUpload(c *gin.Context) {
	var date, dirname string
	Year, Month, _ := time.Now().Date()
	if int(Month) == 1 {
		Month = 12
		Year = Year - 1
		date = fmt.Sprintf("%d%d", Year, int(Month))
	} else {
		date = fmt.Sprintf("%d%d", Year, int(Month))
	}
	dirname = fmt.Sprintf("/data/BillSystem/upload/%s", date)
	if boo := tools.PathExists(dirname); !boo {
		err := os.Mkdir(dirname, 0755)
		if err != nil {
			panic(err)
		}
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"info": map[string]interface{}{
				"message": err.Error(),
			},
		})
	}
	files := form.File["filename"]
	var fileList []string
	for _, file := range files {
		fileList = append(fileList, file.Filename)
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, dirname + "/" + filename); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"info": map[string]interface{}{
					"message": err.Error(),
				},
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"info": map[string]interface{}{
			"data" : fileList,
			"total": len(files),
		},
	})
}