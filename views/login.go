package views

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/go-ldap/ldap"
)

func Login(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	var data map[string]string
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		panic(err)
	}
	session := sessions.Default(c)
	username := data["username"]
	password := data["password"]

	l, err := ldap.Dial("tcp", "10.96.140.61:389")
	if err != nil {
		fmt.Println("连接AD域失败",err)
	}
	err = l.Bind("cn=zabbix,cn=Users,dc=open,dc=com,dc=cn", "ZbxOpen09")
	if err != nil {
		fmt.Println("管理员认证失败",err)
	}
	searchRequest := ldap.NewSearchRequest(
		"dc=open,dc=com,dc=cn",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(sAMAccountName=%s))", username),
		[]string{"dn"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"code": 500, "message": "user error"})
		return
	}

	if len(sr.Entries) != 1 {
		c.JSON(http.StatusInternalServerError,
			gin.H{"code": 500, "message": "user not found"})
		return
	}
	userdn := sr.Entries[0].DN
	err = l.Bind(userdn, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"code": 500, "message": "username and password aren't match"})
		return
	}
	session.Set("username", username)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"code": 500, "message": "session存储错误"})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{"code": 200, "message": "Successfully authenticated user"})
}