package worker

import (
	"awesomeProject/prometheus"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Controller struct {
	c  *Config
	wg *sync.WaitGroup
	ch *chan bool
	s  prometheus.Statist
}

type Config struct {
	Time    int `yaml:"time"`
	Threads int `yaml:"threads"`
	URLs    []struct {
		URL    string `yaml:"url"`
		Method string `yaml:"method"`
		ID     bool   `yaml:"id"`
		Body   bool   `yaml:"body"`
	} `yaml:"urls"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewWorker(c *Config, wg *sync.WaitGroup, ch *chan bool, s prometheus.Statist) *Controller {
	return &Controller{c: c, wg: wg, ch: ch, s: s}
}

func (c *Controller) Start(num int) {
	c.wg.Add(1)
	go func(num int) {
		defer c.wg.Done()

		var i = -1
		var id int

		for {
			select {
			case _, ok := <-*c.ch:
				if !ok {
					return
				}
			default:
			}

			i++
			if i >= len(c.c.URLs) {
				i = 0
			}

			client := http.Client{}

			var url string
			if c.c.URLs[i].ID {
				url = fmt.Sprintf("%s/%d", c.c.URLs[i].URL, id)
			} else {
				url = c.c.URLs[i].URL
			}

			var body io.Reader
			if c.c.URLs[i].Body {
				b, err := json.Marshal(User{
					ID:       id,
					Username: fmt.Sprintf("user%s%d", time.Now().String(), num),
					Email:    fmt.Sprintf("email%s%d", time.Now().String(), num),
				})
				if err != nil {
					log.Printf("json marshal err: %s", err.Error())
					continue
				}
				body = io.NopCloser(bytes.NewBuffer(b))
			}

			req, err := http.NewRequest(c.c.URLs[i].Method, url, body)
			if err != nil {
				log.Printf("new request err: %s", err.Error())
				continue
			}

			start := time.Now()

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("client do err: %s", err.Error())
				continue
			}

			c.s.Add(time.Since(start).Seconds())

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("io readall err: %s", err.Error())
				_ = resp.Body.Close()
				continue
			}

			if string(b) != "" {
				var u User
				err = json.Unmarshal(b, &u)
				if err != nil {
					log.Printf(string(b))
					_ = resp.Body.Close()
					continue
				}

				if u.ID != 0 {
					id = u.ID
				}
			}

			_ = resp.Body.Close()
		}
	}(num)
}
