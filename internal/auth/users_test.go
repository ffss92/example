package auth_test

import (
	"testing"

	"github.com/ffss92/example/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	assert := assert.New(t)
	service := newTestService()

	testCases := []struct {
		name  string
		id    int64
		fails bool
	}{
		{
			name:  "should return a user for a existing id",
			id:    seedUser.ID,
			fails: false,
		},
		{
			name:  "should fail for a not existing id",
			id:    999,
			fails: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := service.Get(tc.id)
			if tc.fails {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)
			assert.Equal(seedUser.ID, user.ID)
			assert.Equal(seedUser.Username, user.Username)
		})
	}

}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	service := newTestService()
	testCases := []struct {
		name     string
		username string
		password string
		fails    bool
	}{
		{
			name:     "should create for valid data",
			username: "john.doe",
			password: "sekreT123!",
			fails:    false,
		},
		{
			name:     "should fail for invalid data (missing username)",
			password: "sekreT123!",
			fails:    true,
		},
		{
			name:     "should fail for invalid data (missing password)",
			username: "john.doe",
			fails:    true,
		},
		{
			name:     "should fail for already existing user",
			username: seedUser.Username,
			password: "sekreT123!",
			fails:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := service.CreateUser(auth.CreateUserParams{
				Username: tc.username,
				Password: tc.password,
			})

			if tc.fails {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)
			assert.Equal(tc.username, user.Username)
			assert.NotEqual(tc.password, user.PasswordHash)
			assert.NotZero(user.ID)
			assert.NotZero(user.CreatedAt)
			assert.NotZero(user.UpdatedAt)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	assert := assert.New(t)
	service := newTestService()
	testCases := []struct {
		name     string
		username string
		password string
		fails    bool
	}{
		{
			name:     "should authenticate a valid user",
			username: seedUser.Username,
			password: "password",
			fails:    false,
		},
		{
			name:     "should fail to authenticate for invalid password",
			username: seedUser.Username,
			password: "wrongpassword",
			fails:    true,
		},
		{
			name:     "should fail to authenticate for invalid username",
			username: "not-a-user",
			password: "password",
			fails:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := service.Authenticate(auth.CredentialsParam{
				Username: tc.username,
				Password: tc.password,
			})

			if tc.fails {
				assert.ErrorIs(err, auth.ErrInvalidCredentials)
				return
			}

			assert.NoError(err)
			assert.Equal(auth.ScopeAuthentication, token.Scope)
			assert.NotZero(token.PlainText)
			assert.NotZero(token.Expiry)
		})
	}
}

func TestGetUserForToken(t *testing.T) {
	assert := assert.New(t)
	service := newTestService()
	testCases := []struct {
		name  string
		token string
		scope auth.Scope
		fails bool
	}{
		{
			name:  "should return a user for valid token",
			token: seedToken.PlainText,
			scope: seedToken.Scope,
			fails: false,
		},
		{
			name:  "should not return a user for valid token",
			token: "some-invalid-token",
			scope: auth.ScopeAuthentication,
			fails: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := service.GetUserForToken(tc.token, tc.scope)
			if tc.fails {
				assert.ErrorIs(err, auth.ErrInvalidToken)
				return
			}

			assert.NoError(err)
			assert.NotNil(user)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)
	service := newTestService()

	testCases := []struct {
		name  string
		id    int64
		fails bool
	}{
		{
			name:  "should delete for existing user",
			id:    seedUser.ID,
			fails: false,
		},
		{
			name:  "should failed for non existing user",
			id:    9999,
			fails: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.DeleteUser(tc.id)
			if tc.fails {
				assert.ErrorIs(auth.ErrNotFound, err)
				return
			}
			assert.NoError(err)
		})
	}
}
