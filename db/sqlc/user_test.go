package db

import (
	"context"
	"testing"
	"time"

	"github.com/devder/grpc-b/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword, // to hash
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
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
	user2, err := testStore.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			Valid:  true,
			String: newFullName,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text{
			Valid:  true,
			String: newEmail,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			Valid:  true,
			String: hashedPassword,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, hashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			Valid:  true,
			String: hashedPassword,
		},
		FullName: pgtype.Text{
			Valid:  true,
			String: newFullName,
		},
		Email: pgtype.Text{
			Valid:  true,
			String: newEmail,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, hashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)
}
