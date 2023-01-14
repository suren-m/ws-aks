package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func main() {
	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/host", func(c *gin.Context) {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "not available"
		}

		nodename := os.Getenv("NODE_NAME")
		result := hostname + " on " + nodename
		c.String(http.StatusOK, result)

	})

	r.GET("/connections", func(c *gin.Context) {
		activeSessions, _ := m.Sessions()
		c.String(http.StatusOK, strconv.Itoa(len(activeSessions)))
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	r.Run(":5000")
}
