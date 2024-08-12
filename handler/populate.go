package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickn2/Clerk-Challenge/model"
	"github.com/patrickn2/Clerk-Challenge/schema"
)

func (h *Handler) Populate(ctx *gin.Context) {
	cx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)
	go checkConnectionStatus(cancel, ctx, done)
	var UserChan = make(chan *[]schema.User)
	var wg sync.WaitGroup
	totalAPICalls := 5
	for i := 0; i < totalAPICalls; i++ {
		wg.Add(1)
		go getRandomUser(cx, UserChan, &wg, 5)
		time.Sleep(time.Millisecond * 150)
	}
	var users []schema.User
	for i := 0; i < totalAPICalls; i++ {
		u := <-UserChan
		if u == nil {
			continue
		}
		users = append(users, *u...)
	}
	if len(users) == 0 {
		ctx.JSON(http.StatusServiceUnavailable, map[string]int{"rows_affected": 0})
		return
	}
	rows, err := model.InsertNewUsers(h.db, &users)
	if err != nil {
		log.Printf("error inserting users %v\n", err)
		done <- true
		ctx.Status(http.StatusInternalServerError)
		return
	}
	done <- true
	ctx.JSON(http.StatusOK, map[string]int{"rows_affected": rows})
}

func checkConnectionStatus(cancel context.CancelFunc, ctx *gin.Context, done chan bool) {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			continue
		case <-done:
			return
		case <-ctx.Request.Context().Done():
			log.Println("Connection Aborted")
			cancel() // Cancel all http requests
			ctx.Abort()
			return
		}
	}
}

func getRandomUser(ctx context.Context, u chan *[]schema.User, wg *sync.WaitGroup, try int) {
	if try == 0 {
		log.Println("server is not responding as expected after some requests")
		u <- nil
		wg.Done()
		return
	}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://randomuser.me/api/?results=5000", http.NoBody)
	if err != nil {
		log.Println("error creating new request", err)
		u <- nil
		wg.Done()
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error getting user", err)
		u <- nil
		wg.Done()
		return
	}
	defer resp.Body.Close()

	var randomUsers schema.RandomUser

	delay := time.Millisecond * time.Duration((rand.Float64()*200)+400+((6-float64(try))*2000))
	if resp.StatusCode == http.StatusTooManyRequests {
		log.Printf("randomuser.me is returning too many requests. Trying again in %d ms\n", delay)
		time.Sleep(delay)
		getRandomUser(ctx, u, wg, try-1)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&randomUsers)

	if err != nil {
		log.Printf("error decoding json %v - Fetch Failed - Trying Again... \n", err)
		time.Sleep(delay)
		getRandomUser(ctx, u, wg, try-1)
		return
	}

	var users []schema.User

	for _, user := range randomUsers.Results {
		users = append(users, schema.User{
			Name:      fmt.Sprintf("%s %s", user.Name.First, user.Name.Last),
			Email:     user.Email,
			Phone:     user.Phone,
			Picture:   user.Picture.Thumbnail,
			CreatedAt: user.Registered.Date,
		})
	}
	u <- &users
	wg.Done()
}
