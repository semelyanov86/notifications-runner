package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/semelyanov86/notifications-runner/internal/config"
	"github.com/semelyanov86/notifications-runner/internal/domain"
	"time"
)

type Entity struct {
	DB     *sql.DB
	config config.Config
}

func NewEntityRepository(db *sql.DB, config2 config.Config) Entity {
	return Entity{
		DB:     db,
		config: config2,
	}
}

func (e Entity) GetLastProcessedEntity() (int, error) {
	// TODO: Fix query
	var query = "SELECT MAX(`crmid`) AS 'crmid' FROM vtiger_crmentity WHERE deleted = 0 AND entity = ?"

	var order = 0

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, "VDNotification").Scan(
		&order,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 1, nil
		default:
			return 0, err
		}
	}
	return order, nil
}

func (e Entity) GetNextNotProcessedEntity(last int) (*domain.Entity, error) {
	// TODO: Fix query
	var query = "SELECT * FROM vtiger_crmentity WHERE crmid > ? AND entity = ?"

	var entity domain.Entity

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err = e.DB.QueryRowContext(ctx, query, last, "VDNotification").Scan(&entity.Crmid, &entity.Description, &entity.Label)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &entity, nil
}

func (e Entity) GetEntityById(id int) (*domain.Entity, error) {
	// TODO: Fix query
	var query = "SELECT * FROM vtiger_crmentity WHERE crmid = ? AND entity = ?"

	var entity domain.Entity

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err = e.DB.QueryRowContext(ctx, query, id, "VDNotification").Scan(&entity.Crmid, &entity.Description, &entity.Label)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &entity, nil
}

func (e Entity) MarkAsSent(entity *domain.Entity) error {
	// TODO: Fix query
	var query = "UPDATE vtiger_vdnotificatons SET timestamp = ?, status = ? WHERE id = ?"
	var args = []any{
		entity.Label,
		entity.Description,
		entity.Crmid,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var _, err2 = e.DB.ExecContext(ctx, query, args...)
	return err2
}

func (e Entity) MarkAsError(entity *domain.Entity, err error) error {
	// TODO: FIx query
	var query = "UPDATE vtiger_vdnotificatons SET timestamp = ?, status = ? WHERE id = ?"
	var args = []any{
		entity.Label,
		err.Error(),
		entity.Crmid,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var _, err2 = e.DB.ExecContext(ctx, query, args...)
	return err2
}
