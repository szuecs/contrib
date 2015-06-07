package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szuecs/gin-contrib/gomonitor"
	"gopkg.in/mcuadros/go-monitor.v1/aspects"
)

type CounterAspect struct {
	Count int
}

func (a *CounterAspect) Inc() {
	a.Count++
}

func (a *CounterAspect) GetStats() interface{} {
	return a.Count
}

func (a *CounterAspect) Name() string {
	return "Counter"
}

func (a *CounterAspect) InRoot() bool {
	return false
}

// Counter handler:
func monitor_handler(asp *CounterAspect) gin.HandlerFunc {
	var counter *CounterAspect = asp
	return func(c *gin.Context) {
		counter.Inc()
		c.Next()
	}
}

func main() {
	counterAspect := &CounterAspect{0}
	asps := []aspects.Aspect{counterAspect}
	router := gin.New()
	// curl http://localhost:9000/
	// curl http://localhost:9000/Custom
	router.Use(gomonitor.Metrics(9000, asps))
	// curl http://localhost:9000/Counter
	router.Use(monitor_handler(counterAspect))
	// last middleware
	router.Use(gin.Recovery())

	// each request to all handlers like below will increment the Counter
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"title": "Counter - Hello World"})
	})

	//..
	router.Run(":8080")
}
