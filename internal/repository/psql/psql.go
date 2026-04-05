// Модуль psql реализует функцонал хранения данных в БД PostgreSQL
package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	psqldrv "github.com/nk87rus/stenographer/internal/db/psql"
	"github.com/nk87rus/stenographer/internal/model"
)

// PSQLDriver - описывает методы необходимые для взаимодействия с БД PostgreSQL
//
//go:generate go run github.com/vektra/mockery/v2@v2.53.6 --name=PSQLDriver --inpackage --testonly
type PSQLDriver interface {
	GetConnConfig() *pgx.ConnConfig
	Insert(ctx context.Context, req string, args ...any) error
	Exec(ctx context.Context, req string, args ...any) error
	SelectRow(ctx context.Context, req string, args ...any) pgx.Row
	SelectRows(ctx context.Context, req string, args ...any) (pgx.Rows, error)
}

// PSQLRepo - структура хранилища PostgreSQL
type PSQLRepo struct {
	db PSQLDriver // подключение к БД PostgreSQL
}

// NewPSQLRepo - инициализирует новое PostgreSQL хранилище
func NewPSQLRepo(ctx context.Context, dsn string) (*PSQLRepo, error) {
	dbDrv, errDB := psqldrv.InitPSQL(ctx, dsn)
	if errDB != nil {
		return nil, errDB
	}

	if err := applyMigrations(ctx, dbDrv.GetConnConfig()); err != nil {
		return nil, err
	}
	return &PSQLRepo{db: dbDrv}, nil
}

func (s *PSQLRepo) RegisterUser(ctx context.Context, userID int64, userName string) error {
	req := `INSERT INTO public.users(id, user_name) VALUES ($1, $2);`
	ctxInsert, cancelInsert := context.WithTimeout(ctx, 5*time.Second)
	defer cancelInsert()
	if err := s.db.Insert(ctxInsert, req, userID, userName); err != nil {
		return fmt.Errorf("RegisterUser: %w", err)
	}
	return nil
}

func (s *PSQLRepo) GetMeetingsList(ctx context.Context) ([]model.MeetingsListItem, error) {
	req := "SELECT id, ts, (SELECT user_name FROM public.users u WHERE u.id = m.user_id) as user_name FROM public.meetings m"
	ctxSelect, cancelSelect := context.WithTimeout(ctx, 5*time.Second)
	defer cancelSelect()

	rawData, errDB := s.db.SelectRows(ctxSelect, req)
	if errDB != nil {
		return nil, errDB
	}

	return pgx.CollectRows(rawData, pgx.RowToStructByName[model.MeetingsListItem])

}

func (s *PSQLRepo) GetMeeting(ctx context.Context, mID int64) (*model.Meeting, error) {
	req := "SELECT id, ts, (SELECT user_name FROM public.users u WHERE u.id = m.user_id) as user_name, data FROM public.meetings m WHERE id = $1;"
	ctxSelect, cancelSelect := context.WithTimeout(ctx, 5*time.Second)
	defer cancelSelect()

	rawData, errDB := s.db.SelectRows(ctxSelect, req, mID)
	if errDB != nil {
		return nil, errDB
	}
	result, err := pgx.CollectOneRow(rawData, pgx.RowToStructByName[model.Meeting])
	return &result, err
}
