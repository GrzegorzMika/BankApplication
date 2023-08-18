package db

import (
	"context"
	"testing"

	"BankApplication/internal/util"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) (CreateUserParams, User) {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	return arg, user
}

func TestCreateUser(t *testing.T) {
	arg, user := createRandomUser(t)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Zero(t, user.PasswordChangedAt.Time)
	require.NotZero(t, user.CreatedAt)
}

func TestGetUser(t *testing.T) {
	_, user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user, user2)
}
