package postgresql

import (
	"app/models"
	"context"
	"testing"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		Name    string
		Input   *models.CreateUser
		Output  string
		WantErr bool
	}{
		{
			Name: "Case 1",
			Input: &models.CreateUser{
				FirstName:   "Shokhrukh",
				LastName:    "Safarov",
				Login:       "webslinger",
				Password:    "12345",
				PhoneNumber: "+998900976035",
			},
			WantErr: false,
		},
	}

	for _, test := range tests {
		id, err := userTestRepo.Create(context.Background(), test.Input)
		if id == "" {
			t.Errorf("%s: got: %v", test.Name, err)
			return
		}
	}
}

func TestGetByIdUser(t *testing.T) {
	tests := []struct {
		Name    string
		Input   *models.UserPrimaryKey
		Output  *models.User
		WantErr bool
	}{
		{
			Name: "Case 1",
			Input: &models.UserPrimaryKey{
				Id: "e83d40fa-520a-41d4-83cb-47a78afe6398",
			},
			Output: &models.User{
				Id:          "e83d40fa-520a-41d4-83cb-47a78afe6398",
				FirstName:   "Shokhrukh",
				LastName:    "Safarov",
				Login:       "webslinger",
				Password:    "12345",
				PhoneNumber: "+998900976035",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			user, err := userTestRepo.GetById(context.Background(), test.Input)

			if user.Id != test.Output.Id ||
				user.FirstName != test.Output.FirstName ||
				user.LastName != test.Output.LastName ||
				user.Login != test.Output.Login ||
				user.Password != test.Output.Password ||
				user.PhoneNumber != test.Output.PhoneNumber {
				t.Errorf("%s: got: %+v, expected: %+v", test.Name, *user, *test.Output)
				return
			}

			if err != nil {
				t.Errorf("%s: got: %v", test.Name, err)
				return
			}

			if test.WantErr {
				t.Errorf("%s: got: %v", test.Name, err)
				return
			}
		})
	}
}
