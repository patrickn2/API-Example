package service

import (
	"context"
	"testing"

	"github.com/patrickn2/api-challenge/mocks"
	"github.com/patrickn2/api-challenge/schema"
	"github.com/stretchr/testify/assert"
)

func Test_Clerks(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	tcases := map[string]struct {
		getClerksOutput []*schema.User
		params          *schema.GetClerksParams
		expectedError   error
		expectedOutput  *ClerksResponse
	}{
		"Empty Params no Next Page": {
			getClerksOutput: []*schema.User{{Id: 123}, {Id: 234}, {Id: 345}, {Id: 456}, {Id: 567}, {Id: 678}},
			params:          &schema.GetClerksParams{},
			expectedError:   nil,
			expectedOutput: &ClerksResponse{
				Users:      []*schema.User{{Id: 123}, {Id: 234}, {Id: 345}, {Id: 456}, {Id: 567}, {Id: 678}},
				TotalUsers: 6,
			},
		},
		"Limit 6. Next Page": {
			getClerksOutput: []*schema.User{{Id: 123}, {Id: 234}, {Id: 345}, {Id: 456}, {Id: 567}, {Id: 678}, {Id: 789}},
			params:          &schema.GetClerksParams{Limit: func(u uint) *uint { return &u }(6)},
			expectedError:   nil,
			expectedOutput: &ClerksResponse{
				Users:      []*schema.User{{Id: 123}, {Id: 234}, {Id: 345}, {Id: 456}, {Id: 567}, {Id: 678}},
				TotalUsers: 6,
				NextPage:   func(u uint) *uint { return &u }(678),
			},
		},
		"Limit 6. Ending Before no Previous Page": {
			getClerksOutput: []*schema.User{{Id: 123}, {Id: 234}, {Id: 345}, {Id: 456}, {Id: 567}, {Id: 678}},
			params:          &schema.GetClerksParams{Limit: func(u uint) *uint { return &u }(6), EndingBefore: func(u uint) *uint { return &u }(789)},
			expectedError:   nil,
			expectedOutput: &ClerksResponse{
				Users:      []*schema.User{{Id: 123}, {Id: 234}, {Id: 345}, {Id: 456}, {Id: 567}, {Id: 678}},
				TotalUsers: 6,
			},
		},
		"Limit 6. Ending Before, Previous Page": {
			getClerksOutput: []*schema.User{{Id: 123}, {Id: 234}, {Id: 345}, {Id: 456}, {Id: 567}, {Id: 678}, {Id: 789}},
			params:          &schema.GetClerksParams{Limit: func(u uint) *uint { return &u }(6), EndingBefore: func(u uint) *uint { return &u }(890)},
			expectedError:   nil,
			expectedOutput: &ClerksResponse{
				Users:        []*schema.User{{Id: 234}, {Id: 345}, {Id: 456}, {Id: 567}, {Id: 678}, {Id: 789}},
				TotalUsers:   6,
				PreviousPage: func(u uint) *uint { return &u }(234),
			},
		},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			repo.EXPECT().GetClerks(ctx, tcase.params).Return(tcase.getClerksOutput, tcase.expectedError)
			userService := NewUserService(repo)
			users, err := userService.Clerks(ctx, tcase.params)
			assert.Equal(t, tcase.expectedError, err, "wrong error at %s", name)
			assert.EqualValues(t, tcase.expectedOutput, users, "wrong users output at %s", name)
		})
	}
}
