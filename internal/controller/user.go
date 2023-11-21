package controller

import (
	"admin/cmd/router"
	"admin/internal/response"
	"github.com/gin-gonic/gin"
)

// 自动注册路由
func init() {
	router.Register(&User{})
}

type PagesRequest struct {
	Id int `json:"id" form:"id"`
}

type User struct{}

// Pages 控制器的方法 分页查询
func (api *User) Pages(c *gin.Context) { //,httpMethod string
	var req PagesRequest
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithCode(response.Invalidate, c)
		return
	}
	users := []int{1, 2, 3}
	response.OkWithData(users, c)
}
