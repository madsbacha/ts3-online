package main

import (
	"fmt"
	"github.com/ziutek/telnet"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type status struct {
	mu        sync.Mutex
	Online    int      `json:"online"`
	Usernames []string `json:"usernames"`
}

func (st *status) set(online int, usernames []string) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.Online = online
	st.Usernames = usernames
}

func (st *status) get() (int, []string) {
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.Online, st.Usernames
}

func emptyStatus() status {
	return status{
		Online:    0,
		Usernames: make([]string, 0),
	}
}

func fetchTsStatus(host, username, password string) status {
	timeout, err := time.ParseDuration("5s")
	if err != nil {
		panic(err)
	}
	conn, err := telnet.DialTimeout("tcp", host, timeout)
	if err != nil {
		log.Println("dial:", err)
		return emptyStatus()
	}
	defer conn.Close()

	err = conn.SkipUntil("command.\n\r")
	if err == io.EOF {
		return emptyStatus()
	} else if err != nil {
		log.Println("skip until:", err)
		return emptyStatus()
	}

	cmd := fmt.Sprintf("login %v %v\nuse 1\nclientlist\nquit\n", username, password)
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		log.Println("write:", err)
		return emptyStatus()
	}

	usernames := make([]string, 0)
	for {
		response, err := conn.ReadString('\r')

		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("read:", err)
			break
		}

		if response == "error id=0 msg=ok\n\r" {
			continue
		}

		matches := userRegex.FindAllStringSubmatch(response, -1)
		for _, match := range matches {
			if strings.HasPrefix(match[1], username) || excludeUsername(match[1]) {
				continue
			}
			usernames = append(usernames, match[1])
		}
	}

	return status{
		Online:    len(usernames),
		Usernames: usernames,
	}
}

func fetchTsStatusCron() {
	host := os.Getenv("TS_HOST")
	username := os.Getenv("TS_USERNAME")
	password := os.Getenv("TS_PASSWORD")

	newStatus := fetchTsStatus(host, username, password)
	currentStatus.set(newStatus.Online, newStatus.Usernames)
	connectedSockets.PushStatus(&currentStatus)
}
