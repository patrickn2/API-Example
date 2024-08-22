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
	"github.com/patrickn2/api-challenge/schema"
)

func (h *Handler) Populate(ginCtx *gin.Context) {
	ctx, cancel := context.WithCancel(ginCtx)
	done := make(chan bool)
	go checkConnectionStatus(cancel, ginCtx, done)
	var wg sync.WaitGroup
	totalAPICalls := 5
	var users []*schema.User
	for i := 0; i < totalAPICalls; i++ {
		wg.Add(1)
		go func() {
			u := getRandomUser(ctx, 5)
			if u != nil {
				users = append(users, u...)
			}
			wg.Done()
		}()
		time.Sleep(time.Millisecond * 150)
	}

	wg.Wait()

	if len(users) == 0 {
		ginCtx.JSON(http.StatusServiceUnavailable, map[string]int{"rows_affected": 0})
		return
	}
	rows, err := h.Services.UserService.NewUsers(users)
	if err != nil {
		log.Printf("Error inserting users %v\n", err)
		done <- true
		ginCtx.Status(http.StatusInternalServerError)
		return
	}
	done <- true
	ginCtx.JSON(http.StatusOK, map[string]int{"rows_affected": rows})
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

func getRandomUser(ctx context.Context, try int) []*schema.User {
	if try == 0 {
		log.Println("Server is not responding as expected after some requests")
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://randomuser.me/api/?results=1000", http.NoBody)
	if err != nil {
		log.Println("Error creating new request", err)
		return nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting user", err)
		return nil
	}
	defer resp.Body.Close()

	var randomUsers schema.RandomUser

	delay := time.Millisecond * time.Duration((rand.Float64()*200)+400+((6-float64(try))*2000))
	if resp.StatusCode == http.StatusTooManyRequests {
		log.Printf("randomuser.me is returning too many requests. Trying again in %d ms\n", delay)
		time.Sleep(delay)
		return getRandomUser(ctx, try-1)
	}

	err = json.NewDecoder(resp.Body).Decode(&randomUsers)

	if err != nil {
		log.Printf("Error decoding json %v - Fetch Failed - Trying Again... \n", err)
		time.Sleep(delay)
		return getRandomUser(ctx, try-1)
	}

	var users []*schema.User

	for _, user := range randomUsers.Results {
		users = append(users, &schema.User{
			Name:      fmt.Sprintf("%s %s", user.Name.First, user.Name.Last),
			Email:     user.Email,
			Phone:     user.Phone,
			Picture:   user.Picture.Thumbnail,
			CreatedAt: user.Registered.Date,
		})
	}
	return users
}
