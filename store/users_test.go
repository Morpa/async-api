package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/Morpa/async-api/fixtures"
	"github.com/Morpa/async-api/store"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T) {
	env := fixtures.NewTestEnv(t)
	cleanup := env.SetupdDb(t)
	t.Cleanup(func() {
		cleanup(t)
	})

	now := time.Now()
	ctx := context.Background()
	userStore := store.NewUserStore(env.Db)
	user, err := userStore.CreateUser(ctx, "teste@teste.com", "testepassword")
	require.NoError(t, err)
	require.Less(t, now.UnixNano(), user.CreatedAt.UnixNano())

	require.Equal(t, "teste@teste.com", user.Email)
	require.NoError(t, user.ComparePassword("testepassword"))

	user2, err := userStore.ById(ctx, user.Id)
	require.NoError(t, err)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Id, user2.Id)
	require.Equal(t, user.HashedPasswordBase64, user2.HashedPasswordBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user2.CreatedAt.UnixNano())

	user2, err = userStore.ByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Id, user2.Id)
	require.Equal(t, user.HashedPasswordBase64, user2.HashedPasswordBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user2.CreatedAt.UnixNano())
}
