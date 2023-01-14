package main

import (
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type User struct {
	ID   int
	name string
}

func main() {
	r := gin.Default()
	m := melody.New()
	// allow all origin
	// m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	activeUsers := make(map[*melody.Session]*User)
	lock := new(sync.Mutex)
	counter := 0

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/host", func(c *gin.Context) {
		activeSessions, _ := m.Sessions()
		res := getHost() + " (" + strconv.Itoa(len(activeSessions)) + " active sessions)"
		c.String(http.StatusOK, res)
	})

	r.GET("/connections", func(c *gin.Context) {
		activeSessions, _ := m.Sessions()
		c.String(http.StatusOK, strconv.Itoa(len(activeSessions)))
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		lock.Lock()
		counter += 1
		activeUsers[s] = &User{counter, "User " + strconv.Itoa(counter)}
		s.Write([]byte("Connected to: " + getHost() + " (" + strconv.Itoa(len(activeUsers)) + " active members)"))
		lock.Unlock()
	})

	m.HandleDisconnect(func(s *melody.Session) {
		lock.Lock()
		m.BroadcastOthers([]byte("Disconnected: "+activeUsers[s].name), s)
		delete(activeUsers, s)
		lock.Unlock()
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		lock.Lock()
		prefix := []byte("<" + activeUsers[s].name + "> ")
		res := append(prefix, msg...)
		m.Broadcast(res)
		lock.Unlock()
	})

	// m.HandleMessage(func(s *melody.Session, msg []byte) {
	// 	m.Broadcast(msg)
	// })

	r.Run(":5000")
}

func getHost() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "not available"
	}

	nodename := os.Getenv("NODE_NAME")
	result := hostname + " on " + nodename

	return result
}
