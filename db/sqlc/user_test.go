package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Annongkhanh/Simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	randomPassword := util.RandomString(int(util.RandomInt(32, 6)))
	hashedPassword, err := util.HashPassword(randomPassword)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Fullname:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.NoError(t, util.CheckPassword(randomPassword, hashedPassword))
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreateAt)

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
	require.Equal(t, user1, user2)
}

// func TestDeleteAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)

// 	err := testQueries.DeleteAccount(context.Background(), account1.ID)

// 	require.NoError(t, err)

// 	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, account2)
// }

// func TestUpdateAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)

// 	arg := UpdateAccountParams{
// 		ID:      account1.ID,
// 		Balance: util.RandomMoney(),
// 	}

// 	account2, err := testQueries.UpdateAccount(context.Background(), arg)

// 	require.NoError(t, err)
// 	require.Equal(t, account2.ID, arg.ID)
// 	require.Equal(t, account2.Balance, arg.Balance)
// 	require.Equal(t, account1.Owner, account2.Owner)
// 	require.Equal(t, account1.Currency, account2.Currency)
// 	require.Equal(t, account1.CreatedAt, account2.CreatedAt)

// }

// func TestListAccounts(t *testing.T) {

// 	for i := 0; i < 10; i++ {
// 		createRandomAccount(t)
// 	}
// 	arg := ListAccountsParams{
// 		Limit:  5,
// 		Offset: 5,
// 	}

// 	accounts, err := testQueries.ListAccounts(context.Background(), arg)

// 	require.NoError(t, err)
// 	require.Len(t, accounts, 5)

// 	for _, account := range accounts {
// 		require.NotEmpty(t, account)
// 	}
// }

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user1.Username,
		Fullname: sql.NullString{
			String: util.RandomOwner(),
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.NotEqual(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.Equal(t, user1.CreateAt, user2.CreateAt)

}
