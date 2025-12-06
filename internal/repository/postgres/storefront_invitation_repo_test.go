package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vondi-global/listings/internal/domain"
)

func TestStorefrontInvitationRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	t.Run("create email invitation", func(t *testing.T) {
		email := "test@example.com"
		now := time.Now()
		inv := &domain.StorefrontInvitation{
			StorefrontID: 1,
			Role:         "staff",
			Type:         domain.InvitationTypeEmail,
			InvitedEmail: &email,
			InvitedByID:  100,
			Status:       domain.InvitationStatusPending,
		}

		mock.ExpectQuery(`INSERT INTO storefront_invitations`).
			WithArgs(
				inv.StorefrontID,
				inv.Role,
				inv.Type,
				inv.InvitedEmail,
				inv.InvitedUserID,
				inv.InviteCode,
				inv.ExpiresAt,
				inv.MaxUses,
				inv.CurrentUses,
				inv.InvitedByID,
				inv.Status,
				inv.Comment,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(1, now, now))

		err := repo.Create(ctx, inv)
		require.NoError(t, err)
		assert.Equal(t, int64(1), inv.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create link invitation", func(t *testing.T) {
		code := "sf-abc123"
		expires := time.Now().Add(7 * 24 * time.Hour)
		maxUses := int32(10)
		now := time.Now()

		inv := &domain.StorefrontInvitation{
			StorefrontID: 1,
			Role:         "manager",
			Type:         domain.InvitationTypeLink,
			InviteCode:   &code,
			ExpiresAt:    &expires,
			MaxUses:      &maxUses,
			InvitedByID:  100,
			Status:       domain.InvitationStatusPending,
		}

		mock.ExpectQuery(`INSERT INTO storefront_invitations`).
			WithArgs(
				inv.StorefrontID,
				inv.Role,
				inv.Type,
				inv.InvitedEmail,
				inv.InvitedUserID,
				inv.InviteCode,
				inv.ExpiresAt,
				inv.MaxUses,
				inv.CurrentUses,
				inv.InvitedByID,
				inv.Status,
				inv.Comment,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(2, now, now))

		err := repo.Create(ctx, inv)
		require.NoError(t, err)
		assert.Equal(t, int64(2), inv.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStorefrontInvitationRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "storefront_id", "role", "type",
			"invited_email", "invited_user_id",
			"invite_code", "expires_at", "max_uses", "current_uses",
			"invited_by_id", "status", "comment",
			"created_at", "updated_at", "accepted_at", "declined_at",
		}).AddRow(
			1, 10, "staff", "email",
			email, nil,
			nil, nil, nil, 0,
			100, "pending", "",
			now, now, nil, nil,
		)

		mock.ExpectQuery(`SELECT .+ FROM storefront_invitations WHERE id`).
			WithArgs(1).
			WillReturnRows(rows)

		inv, err := repo.GetByID(ctx, 1)
		require.NoError(t, err)
		assert.NotNil(t, inv)
		assert.Equal(t, int64(1), inv.ID)
		assert.Equal(t, int64(10), inv.StorefrontID)
		assert.Equal(t, "staff", inv.Role)
		assert.Equal(t, domain.InvitationTypeEmail, inv.Type)
		assert.NotNil(t, inv.InvitedEmail)
		assert.Equal(t, email, *inv.InvitedEmail)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM storefront_invitations WHERE id`).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		inv, err := repo.GetByID(ctx, 999)
		assert.Error(t, err)
		assert.Nil(t, inv)
		assert.Contains(t, err.Error(), "invitation not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStorefrontInvitationRepository_GetByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	code := "sf-abc123"
	maxUses := int32(10)
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "storefront_id", "role", "type",
		"invited_email", "invited_user_id",
		"invite_code", "expires_at", "max_uses", "current_uses",
		"invited_by_id", "status", "comment",
		"created_at", "updated_at", "accepted_at", "declined_at",
	}).AddRow(
		1, 10, "manager", "link",
		nil, nil,
		code, nil, maxUses, 3,
		100, "pending", "",
		now, now, nil, nil,
	)

	mock.ExpectQuery(`SELECT .+ FROM storefront_invitations WHERE invite_code`).
		WithArgs(code).
		WillReturnRows(rows)

	inv, err := repo.GetByCode(ctx, code)
	require.NoError(t, err)
	assert.NotNil(t, inv)
	assert.Equal(t, int64(1), inv.ID)
	assert.Equal(t, domain.InvitationTypeLink, inv.Type)
	assert.NotNil(t, inv.InviteCode)
	assert.Equal(t, code, *inv.InviteCode)
	assert.NotNil(t, inv.MaxUses)
	assert.Equal(t, maxUses, *inv.MaxUses)
	assert.Equal(t, int32(3), inv.CurrentUses)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorefrontInvitationRepository_IncrementUses(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`UPDATE storefront_invitations SET current_uses`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.IncrementUses(ctx, 1)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("max uses reached", func(t *testing.T) {
		mock.ExpectExec(`UPDATE storefront_invitations SET current_uses`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.IncrementUses(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max uses reached")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStorefrontInvitationRepository_MarkAsExpired(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	mock.ExpectExec(`UPDATE storefront_invitations SET status = 'expired'`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.MarkAsExpired(ctx, 1)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorefrontInvitationRepository_ExpirePendingInvitations(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	mock.ExpectExec(`UPDATE storefront_invitations SET status = 'expired'`).
		WillReturnResult(sqlmock.NewResult(0, 5))

	affected, err := repo.ExpirePendingInvitations(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(5), affected)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorefrontInvitationRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	now := time.Now()
	inv := &domain.StorefrontInvitation{
		ID:          1,
		Status:      domain.InvitationStatusAccepted,
		CurrentUses: 5,
		AcceptedAt:  &now,
		Comment:     "Accepted by user",
	}

	mock.ExpectQuery(`UPDATE storefront_invitations SET`).
		WithArgs(
			inv.ID,
			inv.Status,
			inv.CurrentUses,
			inv.AcceptedAt,
			inv.DeclinedAt,
			inv.Comment,
		).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(now))

	err = repo.Update(ctx, inv)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorefrontInvitationRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM storefront_invitations WHERE id`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(ctx, 1)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM storefront_invitations WHERE id`).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(ctx, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invitation not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStorefrontInvitationRepository_GetStatsByStorefront(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"status", "count"}).
		AddRow("pending", 5).
		AddRow("accepted", 10).
		AddRow("declined", 2).
		AddRow("expired", 3)

	mock.ExpectQuery(`SELECT status, COUNT`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	stats, err := repo.GetStatsByStorefront(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int32(5), stats["pending"])
	assert.Equal(t, int32(10), stats["accepted"])
	assert.Equal(t, int32(2), stats["declined"])
	assert.Equal(t, int32(3), stats["expired"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorefrontInvitationRepository_CheckInviteCodeExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewStorefrontInvitationRepository(db)
	ctx := context.Background()

	t.Run("exists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs("sf-abc123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := repo.CheckInviteCodeExists(ctx, "sf-abc123")
		require.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not exists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs("sf-xyz789").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		exists, err := repo.CheckInviteCodeExists(ctx, "sf-xyz789")
		require.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
