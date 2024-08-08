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
	schema "github.com/patrickn2/Clerk-Challenge/schemas"
)

func (h *Handler) Populate(ctx *gin.Context) {
	cx, cancel := context.WithCancel(context.Background())
	go checkConnectionStatus(cancel, ctx)
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
			ctx.Status(500)
			return
		}
		users = append(users, *u...)
	}
	rows, err := model.InsertNewUsers(h.db, &users)
	if err != nil {
		fmt.Printf("error inserting users %v\n", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, map[string]int{"rows_affected": rows})
}

func checkConnectionStatus(cancel context.CancelFunc, ctx *gin.Context) {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	for i := 0; i < 10; i++ {
		select {
		case <-t.C:
			continue
		case <-ctx.Request.Context().Done():
			fmt.Println("Connection Aborted")
			cancel()
			ctx.AbortWithError(http.StatusGone, ctx.Request.Context().Err())
			// client gave up
			return
		}
	}
}

func getRandomUser(ctx context.Context, u chan *[]schema.User, wg *sync.WaitGroup, try int) {
	if try == 0 {
		fmt.Printf("server is not responding as expected after some requests")
		u <- nil
		wg.Done()
		return
	}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://randomuser.me/api/?results=1000", http.NoBody)
	if err != nil {
		fmt.Println("error creating new request", err)
		u <- nil
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error getting user", err)
		u <- nil
		return
	}
	defer resp.Body.Close()

	var randomUsers schema.RandomUser

	if resp.StatusCode == http.StatusTooManyRequests {
		if ctx.Err() != nil {
			fmt.Println("Request Canceled")
			return
		}
		delay := int(rand.Float64()*100) + 400
		log.Printf("Too Many Requests. Trying again in %d ms", delay)
		time.Sleep(time.Millisecond * time.Duration(delay))
		getRandomUser(ctx, u, wg, try)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&randomUsers)

	if err != nil {
		log.Printf("error decoding json %v - Fetch Failed - Trying Again... ", err)
		time.Sleep(time.Millisecond * (time.Duration(int(rand.Float64()*100)) + 100))
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
