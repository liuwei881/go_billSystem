package views

import (
	//orm "BillSystem/database"
	//mode "BillSystem/models"
	//"BillSystem/tools"
//	"encoding/json"
//	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	//"io/ioutil"
	"net/http"
	//"os"
	//"path/filepath"
	//"strconv"
	//"time"
)

//func TestRoot(c *gin.Context){
//	c.String(200, "Hello, World")
//}

//func TestInfo(c *gin.Context){
//	var user []mode.Fest
//	var total int64
//	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
//	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
//	db := orm.BillDb
//	if err := db.Model(&mode.Fest{}).Count(&total).Error; err != nil{
//		c.JSON(http.StatusOK, gin.H{
//			"code" : 500,
//			"message" : "查询数据异常",
//		})
//		return
//	}
//	offset := (page - 1) * pageSize
//	if err := db.Model(&mode.Fest{}).Order("id desc").Limit(pageSize).Offset(offset).Find(&user).Error; err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": -1,
//			"message": "抱歉未找到相关信息",
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "success",
//		"info": map[string]interface{}{
//			"data" : user,
//			"total": total,
//			"page" : page,
//			"pageSize": pageSize,
//		},
//	})
//}
//
//func TestPost(c *gin.Context) {
//	body := c.Request.Body
//	value, err := ioutil.ReadAll(body)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	var data map[string]interface{}
//	if err := json.Unmarshal([]byte(value), &data); err != nil {
//		panic(err)
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "success",
//		"info": map[string]interface{}{
//			"data" : data,
//			"total": len(data),
//		},
//	})
//}
//
//func TestUpload(c *gin.Context) {
//	var date, dirname string
//	Year, Month, _ := time.Now().Date()
//	if int(Month) == 1 {
//		Month = 12
//		Year = Year - 1
//		date = fmt.Sprintf("%d%d", Year, int(Month))
//	} else {
//		date = fmt.Sprintf("%d%d", Year, int(Month))
//	}
//	dirname = fmt.Sprintf("/data/BillSystem/upload/%s", date)
//	if boo := tools.PathExists(dirname); !boo {
//		err := os.Mkdir(dirname, 0755)
//		if err != nil {
//			panic(err)
//		}
//	}
//	form, err := c.MultipartForm()
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code": 400,
//			"info": map[string]interface{}{
//				"message": err.Error(),
//			},
//		})
//	}
//	files := form.File["filename"]
//	var fileList []string
//	for _, file := range files {
//		fileList = append(fileList, file.Filename)
//		filename := filepath.Base(file.Filename)
//		if err := c.SaveUploadedFile(file, dirname + "/" + filename); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"code": 400,
//				"info": map[string]interface{}{
//					"message": err.Error(),
//				},
//			})
//		}
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "success",
//		"info": map[string]interface{}{
//			"data" : fileList,
//			"total": len(files),
//		},
//	})
//}

func Me(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("username")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}