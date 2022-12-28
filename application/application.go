package application

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Application struct {
	router *gin.Engine
}

var appInstance *Application
var once sync.Once

func GetApplicationInstance() *Application {
	once.Do(func() {
		appInstance = &Application{}
	})
	return appInstance
}

func (a *Application) Run() {
	a.router = gin.Default()
	a.router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome To Rate Server")
	})

	a.router.GET("/rates", GetRates)
	e := a.router.Run(":8080")
	if e != nil {
		panic(e)
	}
}
