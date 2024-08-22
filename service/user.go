package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/patrickn2/api-challenge/interfaces"
	"github.com/patrickn2/api-challenge/schema"
)

type UserService struct {
	user interfaces.UserRepository
}

type ClerksResponse struct {
	Users        []*schema.User `json:"users"`
	TotalUsers   int            `json:"total_users"`
	NextPage     *uint          `json:"next_page,omitempty"`
	PreviousPage *uint          `json:"previous_page,omitempty"`
}

type Registered struct {
	Date time.Time `json:"date"`
}
type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}
type Picture struct {
	Thumbnail string `json:"thumbnail"`
}

type Result struct {
	Name       Name       `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Picture    Picture    `json:"picture"`
	Registered Registered `json:"registered"`
}

type RandomUser struct {
	Results []Result `json:"results"`
}

func NewUserService(user interfaces.UserRepository) *UserService {
	return &UserService{
		user,
	}
}

func (us *UserService) NewUsers(ctx context.Context) (int, error) {
	var wg sync.WaitGroup
	var users []*schema.User
	for i := 0; i < 5; i++ {
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
		return 0, errors.New("no users to insert")
	}
	return us.user.InsertUsers(ctx, users)
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

	var randomUsers RandomUser

	delay := time.Millisecond * time.Duration((6-try)*2000)
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

func (us *UserService) Clerks(ctx context.Context, params *schema.GetClerksParams) (*ClerksResponse, error) {
	if params.Limit != nil {
		// The idea is not return an error if the user use a limit lower than 1 or higher than 100
		if *params.Limit == 0 {
			limit := uint(10)
			params.Limit = &limit
		}
		if *params.Limit > 100 {
			limit := uint(100)
			params.Limit = &limit
		}
	} else {
		// Default limit is 10
		limit := uint(10)
		params.Limit = &limit
	}

	if params.Email != nil {
		email := strings.ToLower(*params.Email)
		params.Email = &email
	}

	users, err := us.user.GetClerks(ctx, params)
	if err != nil {
		return nil, err
	}

	var response = ClerksResponse{
		Users:      users,
		TotalUsers: len(users),
	}
	if len(users) > int(*params.Limit) {
		response.TotalUsers = int(*params.Limit)
	} else {
		response.TotalUsers = len(users)
	}

	// Checking if there is a next page or previous page
	if len(users) > int(*params.Limit) {
		if params.EndingBefore != nil {
			previousPage := users[1].Id
			response.Users = users[1:]
			response.PreviousPage = &previousPage
		} else {
			nextPage := users[len(users)-2].Id
			response.Users = users[:len(users)-1]
			response.NextPage = &nextPage
		}
	}
	return &response, nil
}
