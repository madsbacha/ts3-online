package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/ziutek/telnet"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

type status struct {
	mu        sync.Mutex
	online    int
	usernames []string
}

var currentStatus status

var userRegex = regexp.MustCompile(`\sclient_nickname=(.*?)\s`)

func fetchTsStatus(host, username, password string) status {
	timeout, err := time.ParseDuration("5s")
	if err != nil {
		panic(err)
	}
	conn, err := telnet.DialTimeout("tcp", host, timeout)

	if err != nil {
		panic(err)
	}

	err = conn.SkipUntil("command.\n\r")
	if err != nil {
		panic(err)
	}

	cmd := fmt.Sprintf("login %v %v\nuse 1\nclientlist\nquit\n", username, password)
	_, err = conn.Write([]byte(cmd))

	usernames := make([]string, 0)
	for {
		response, err := conn.ReadString('\r')

		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}

		if response == "error id=0 msg=ok\n\r" {
			continue
		}

		matches := userRegex.FindAllStringSubmatch(response, -1)
		for _, match := range matches {
			if match[1] == username {
				continue
			}
			usernames = append(usernames, match[1])
		}
	}

	return status{
		online:    len(usernames),
		usernames: usernames,
	}
}

func fetchTsStatusCron() {
	host := os.Getenv("TS_HOST")
	username := os.Getenv("TS_USERNAME")
	password := os.Getenv("TS_PASSWORD")

	newStatus := fetchTsStatus(host, username, password)
	currentStatus.set(newStatus.online, newStatus.usernames)
}

func (st *status) set(online int, usernames []string) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.online = online
	st.usernames = usernames
}

func (st *status) get() (int, []string) {
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.online, st.usernames
}

func main() {
	currentStatus.set(0, make([]string, 0))

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, err := s.Every(5).Seconds().Do(fetchTsStatusCron)
	if err != nil {
		panic(err)
	}
	s.StartAsync()

	r := gin.Default()
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
	log.Panic(r.Run())
}
