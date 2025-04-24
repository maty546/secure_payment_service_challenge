package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestAccountRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)

	t.Run("happy path", func(t *testing.T) {
		expectedAccountID := "abc123"
		expectedBalance := uint(500)
		now := time.Now()

		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"\."id" LIMIT \$2`).
			WithArgs(expectedAccountID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
				AddRow(expectedAccountID, expectedBalance, now, now),
			)

		repo := NewAccountRepository(db)

		ginCtx, _ := gin.CreateTestContext(nil)
		account, err := repo.GetByID(ginCtx, expectedAccountID)

		assert.NoError(t, err)
		assert.Equal(t, expectedAccountID, account.ID)
		assert.Equal(t, expectedBalance, account.Balance)
		assert.WithinDuration(t, now, account.CreatedAt, time.Second)
		assert.WithinDuration(t, now, account.UpdatedAt, time.Second)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		expectedAccountID := "abc123"

		mock.ExpectQuery(`SELECT \* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"\."id" LIMIT \$2`).
			WithArgs(expectedAccountID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		repo := NewAccountRepository(db)

		ginCtx, _ := gin.CreateTestContext(nil)
		account, err := repo.GetByID(ginCtx, expectedAccountID)

		assert.Error(t, err)
		assert.Equal(t, models.Account{}, account)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
