package db

import (
	"SimpleBank/util"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"database/sql"
)

func createRandomUser(t *testing.T) User {
	hashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T){
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	testQueries.UpdateUser(context.Background(), UpdateUserParams{
		FullName:  sql.NullString{String: newFullName, Valid: true},
		Username: oldUser.Username,
	})
	newUser, err := testQueries.GetUser(context.Background(), oldUser.Username)
	require.NoError(t, err)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.FullName != newUser.FullName, true)
}

func TestUpdateUserOnlyEmail(t *testing.T){
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email:  sql.NullString{String: newEmail, Valid: true},
		Username: oldUser.Username,
	})
	newUser, err := testQueries.GetUser(context.Background(), oldUser.Username)
	require.NoError(t, err)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, newEmail, newUser.Email)
	require.Equal(t, oldUser.FullName, newUser.FullName)
}

func TestUpdateUserOnlyHashedPassword(t *testing.T){
	oldUser := createRandomUser(t)

	newHashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword:  sql.NullString{String: newHashedPassword, Valid: true},
		Username: oldUser.Username,
	})
	newUser, err := testQueries.GetUser(context.Background(), oldUser.Username)
	require.NoError(t, err)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.FullName, newUser.FullName)
}

func TestUpdateUserAllFields(t *testing.T){
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	newFullName := util.RandomOwner()
	newHashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email:  sql.NullString{String: newEmail, Valid: true},
		FullName:  sql.NullString{String: newFullName, Valid: true},
		HashedPassword:  sql.NullString{String: newHashedPassword, Valid: true},
		Username: oldUser.Username,
	})
	newUser, err := testQueries.GetUser(context.Background(), oldUser.Username)
	require.NoError(t, err)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)
	require.NotEqual(t, oldUser.Email, newEmail)
	require.Equal(t, newEmail, newUser.Email)
	require.Equal(t, newFullName, newUser.FullName)
}