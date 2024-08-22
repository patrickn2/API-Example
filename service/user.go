package service

import (
	"strings"

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

func NewUserService(user interfaces.UserRepository) *UserService {
	return &UserService{
		user,
	}
}

func (us *UserService) NewUsers(users []*schema.User) (int, error) {
	return us.user.InsertUsers(users)
}

func (us *UserService) Clerks(params *interfaces.GetClerksParams) (*ClerksResponse, error) {
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

	users, err := us.user.GetClerks(params)
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
