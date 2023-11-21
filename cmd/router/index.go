package router

import (
	"admin/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
)

// Route 路由结构体
type Route struct {
	path       string        //url路径
	httpMethod string        //http方法 get post
	Method     reflect.Value //方法路由
}

// Routes 路由集合
var Routes []Route

func InitRouter(e *gin.Engine) *gin.Engine {
	Bind(e)
	return e
}

// Register 注册控制器
// @message 注册形成的路由方法是采用小驼峰的方式，组采用的是下划线
func Register(target any) bool {
	ctrlName := reflect.TypeOf(target).String()
	module := ctrlName
	if strings.Contains(ctrlName, ".") {
		module = ctrlName[strings.Index(ctrlName, ".")+1:]
	}

	// 将module改为下划线
	module = utils.ToUnderscore(module)

	v := reflect.ValueOf(target)
	//遍历方法
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		action := v.Type().Method(i).Name

		action = utils.ToCamelCase(action)

		path := "/" + module + "/" + action
		httpMethod := http.MethodGet
		if strings.Contains(action, "_") {
			httpMethod = action[strings.LastIndex(action, "_")+1:]
			httpMethod = strings.ToUpper(http.MethodGet)
		}

		route := Route{path: path, Method: method, httpMethod: httpMethod}
		Routes = append(Routes, route)
	}
	return true
}

func Bind(e *gin.Engine) {
	for _, route := range Routes {
		switch route.httpMethod {
		case http.MethodGet:
			e.GET(route.path, match(route.path, route))
		case http.MethodPost:
			e.POST(route.path, match(route.path, route))
		case http.MethodPut:
			e.PUT(route.path, match(route.path, route))
		case http.MethodDelete:
			e.DELETE(route.path, match(route.path, route))
		default:
			log.Fatalf("httpMothod不支持，目前只支持%v", []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			})
		}
	}
}

// 根据path匹配对应的方法
func match(path string, route Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := strings.Split(path, "/")
		if len(fields) < 3 {
			slog.With("fields", fields).Info("路由规则不合法")
			return
		}

		if len(Routes) > 0 {
			arguments := make([]reflect.Value, 1)
			arguments[0] = reflect.ValueOf(c) // *gin.Context
			route.Method.Call(arguments)
		}
	}
}
