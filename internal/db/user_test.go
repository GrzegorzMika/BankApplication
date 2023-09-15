package db

import (
	"context"
	"testing"

	"BankApplication/internal/util"

	"github.com/jackc/pgx/v5/pgtype"
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

func TestUpdateUserOnyFullName(t *testing.T) {
	_, oldUsers := createRandomUser(t)
	newFullName := util.RandomOwner()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: pgtype.Text{},
		FullName:       pgtype.Text{String: newFullName, Valid: true},
		Email:          pgtype.Text{},
		Username:       oldUsers.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUsers.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUsers.Email, updatedUser.Email)
	require.Equal(t, oldUsers.Username, updatedUser.Username)
	require.Equal(t, oldUsers.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnyEmail(t *testing.T) {
	_, oldUsers := createRandomUser(t)
	newEmail := util.RandomEmail()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: pgtype.Text{},
		FullName:       pgtype.Text{},
		Email:          pgtype.Text{String: newEmail, Valid: true},
		Username:       oldUsers.Username,
	})
	require.NoError(t, err)
	require.Equal(t, oldUsers.FullName, updatedUser.FullName)
	require.NotEqual(t, oldUsers.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUsers.Username, updatedUser.Username)
	require.Equal(t, oldUsers.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnyPassword(t *testing.T) {
	_, oldUsers := createRandomUser(t)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: pgtype.Text{String: newHashedPassword, Valid: true},
		FullName:       pgtype.Text{},
		Email:          pgtype.Text{},
		Username:       oldUsers.Username,
	})
	require.NoError(t, err)
	require.Equal(t, oldUsers.FullName, updatedUser.FullName)
	require.Equal(t, oldUsers.Email, updatedUser.Email)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUsers.Username, updatedUser.Username)
	require.NotEqual(t, oldUsers.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserAllFieds(t *testing.T) {
	_, oldUsers := createRandomUser(t)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: pgtype.Text{String: newHashedPassword, Valid: true},
		FullName:       pgtype.Text{String: newFullName, Valid: true},
		Email:          pgtype.Text{String: newEmail, Valid: true},
		Username:       oldUsers.Username,
	})
	require.NoError(t, err)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUsers.Username, updatedUser.Username)
	require.NotEqual(t, oldUsers.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUsers.Email, updatedUser.HashedPassword)
	require.NotEqual(t, oldUsers.FullName, updatedUser.HashedPassword)
}
