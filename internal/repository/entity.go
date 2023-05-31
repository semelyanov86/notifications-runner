package repository

import (
	"Notifications/internal/config"
	"Notifications/internal/domain"
	"context"
	"database/sql"
	"errors"
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
	var query = "SELECT MAX(`crmid`) AS 'crmid' FROM vtiger_crmentity WHERE deleted = 0 AND setype = ? AND (label = ? OR label = ?)"

	var order = 0

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, "VDNotification", "Sent", "Error").Scan(
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
	var query = "SELECT * FROM vtiger_crmentity WHERE crmid > ? AND setype = ? AND label = ?"

	var entity domain.Entity

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err = e.DB.QueryRowContext(ctx, query, last, "VDNotification", "Draft").Scan(&entity.Crmid, &entity.Description, &entity.Label)
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
	var query = "SELECT crmid, label, description, chat_id, notify_datetime, notify_status, notify_type FROM vtiger_crmentity INNER JOIN vtiger_vdnotifications ON vtiger_vdnotifications.vdnotificationsid = vtiger_crmentity.crmid WHERE crmid = ?"

	var entity domain.Entity

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err = e.DB.QueryRowContext(ctx, query, id, "VDNotification").Scan(&entity.Crmid, &entity.Label, &entity.Description, &entity.ChatId, &entity.NotifyDateTime, &entity.NotifyStatus, &entity.NotifyType)
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
	var query = "UPDATE vtiger_vdnotifications SET notify_datetime = ?, notify_status = ?, chat_id = ? WHERE vdnotificationsid = ?"
	var args = []any{
		time.Now(),
		"Sent",
		entity.ChatId,
		entity.Crmid,
	}
	entity.NotifyType = "Sent"
	entity.NotifyDateTime = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var _, err2 = e.DB.ExecContext(ctx, query, args...)
	if err2 != nil {
		return err2
	}

	query = "UPDATE vtiger_crmentity SET label = ? WHERE crmid = ?"
	args = []any{
		"Sent",
		entity.Crmid,
	}
	var _, err3 = e.DB.ExecContext(ctx, query, args...)
	return err3
}

func (e Entity) MarkAsError(entity *domain.Entity, err error) error {
	var query = "UPDATE vtiger_vdnotificatons SET notify_datetime = ?, notify_status = ?, chat_id = ?, error_log = ? WHERE id = ?"
	var args = []any{
		time.Now(),
		"Error",
		entity.ChatId,
		err.Error(),
		entity.Crmid,
	}

	entity.NotifyDateTime = time.Now()
	entity.NotifyStatus = "Error"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var _, err2 = e.DB.ExecContext(ctx, query, args...)
	if err2 != nil {
		return err2
	}
	query = "UPDATE vtiger_crmentity SET label = ? WHERE crmid = ?"
	args = []any{
		"Error",
		entity.Crmid,
	}
	var _, err3 = e.DB.ExecContext(ctx, query, args...)

	return err3
}
