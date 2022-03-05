package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/websocket"
)

var currentStatus status
var connectedSockets SocketStore

var userRegex = regexp.MustCompile(`\sclient_nickname=(.*?)\s`)

var upgrader = websocket.Upgrader{}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ctx.Close()

	sConn := NewSocketConn()
	connectedSockets.AddSocketConn(sConn)
	var res []byte
	for {
		res = <-sConn.C

		//mt, message, err := ctx.ReadMessage()
		//if err != nil {
		//	log.Println("read:", err)
		//	break
		//}
		//log.Printf("recv: %s", message)
		err = ctx.WriteMessage(websocket.BinaryMessage, res)
		if err != nil {
			log.Println("write:", err)
			sConn.SafeClose()
			break
		}
	}
}

func main() {
	currentStatus.set(0, make([]string, 0))

	pullInterval := 5
	if val, ok := os.LookupEnv("PULL_INTERVAL"); ok {
		parsedVal, err := strconv.Atoi(val)
		if err != nil {
			panic(fmt.Errorf("Invalid PULL_INTERVAL value"))
		}
		pullInterval = parsedVal
	}

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, err := s.Every(pullInterval).Seconds().Do(fetchTsStatusCron)
	if err != nil {
		panic(err)
	}
	s.StartAsync()

	r := gin.Default()
	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/index.tmpl")
	r.GET("/api", func(c *gin.Context) {
		online, usernames := currentStatus.get()

		c.JSON(200, gin.H{
			"count": online,
			"users": usernames,
		})
	})
	r.GET("/", func(c *gin.Context) {
		online, usernames := currentStatus.get()

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"count": online,
			"users": usernames,
		})
	})
	r.GET("/websocket", func(c *gin.Context) {
		websocketEndpoint(c.Writer, c.Request)
	})
	log.Panic(r.Run())
}
