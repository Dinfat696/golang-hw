package sqlstorage

import (
	"context"
    "strconv"
	"fmt"
	"strings"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/models"
	"github.com/jmoiron/sqlx"
)

type SQLStorage struct {
	db           *sqlx.DB
	dbDriverName string
	dsn          string
}

func New(dbDriverName string, dsn string) *SQLStorage {
	dbDriverName = strings.ToLower(dbDriverName)
	dsn = strings.TrimSpace(dsn)
	return &SQLStorage{dbDriverName: dbDriverName, dsn: dsn}
}

func (r *SQLStorage) Connect() error {
	r.db = sqlx.MustConnect(r.dbDriverName, r.dsn)

	r.db.SetMaxIdleConns(5)
	r.db.SetMaxOpenConns(20)
	r.db.SetConnMaxLifetime(1 * time.Minute)
	r.db.SetConnMaxIdleTime(10 * time.Minute)
	return nil
}

func (r *SQLStorage) Close() error {
	err := r.db.Close()
	if err != nil {
		fmt.Print(err)
	}
	return nil
}

func (r *SQLStorage) Create(event models.Event) (id int64, err error) {
	result, err := r.db.Query(
		"INSERT INTO event (name) VALUES($1) RETURNING id", event.Title)
	if err != nil {
		return 0, err
	}
	if result.Next() {
		err = result.Scan(&id)
		return id, err
	}
	return id, err
}

func (r *SQLStorage) Update(event models.Event) {
	_, _ = r.db.Query(
		"UPDATE event SET title =$1, date_time = $2 WHERE id = $3", event.Title, event.DateTime, event.ID)
}

func (r *SQLStorage) DeleteByID(eventID int64) (err error) {
	_, err = r.db.Query(
		"DELETE from event where id = $1", eventID)
	return err
}

func (r *SQLStorage) FindEventsByDay(date time.Time) (res []models.Event, err error) {
	err = r.db.Select(&res,
		"SELECT * from event WHERE date_time = $1", date.Format("2006-01-02"))
	return res, err
}

func (r *SQLStorage) FindEventsByWeek(date time.Time) (res []models.Event, err error) {
	err = r.db.Select(&res,
		"SELECT * from event WHERE date_time = $1 AND $2", date.Format("2006-01-02"), date.AddDate(0, 0, 7))
	return res, err
}

func (r *SQLStorage) FindEventsByMonth(date time.Time) (res []models.Event, err error) {
	err = r.db.Select(&res,
		"SELECT * from event WHERE date_time = $1 AND $2", date.Format("2006-01-02"), date.AddDate(0, 0, 30))
	return res, err
}

func (r *SQLStorage) FindByID(id int64) (res, err error) {
	err = r.db.Get(&res,
		"SELECT * from event where id = $1", id)
	return res, err
}

func (r *SQLStorage) FindAll() (res []models.Event, err error) {
	err = r.db.Select(&res,
		"SELECT * from event")
	return res, err
}

func (r *SQLStorage) FindByIDs(ids []int64) (res []models.Event, err error) {
	query, args, err := sqlx.In(
		"SELECT * from event where id IN(?)", ids)
	if err == nil {
		query = r.db.Rebind(query)
		var rows *sqlx.Rows
		rows, err = r.db.Queryx(query, args...)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var event models.Event
			err = rows.StructScan(&event)
			if err != nil {
				return nil, err
			}
			res = append(res, event)
		}
		return res, err
	}
	return make([]models.Event, 0), err
}

func (r *SQLStorage) DeleteByIDs(ids []int64) (err error) {
	query, args, err := sqlx.In(
		"DELETE from event where id IN(?)", ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Query(query, args...)
	return err
}

func (r *SQLStorage) ExecuteQuery(query string) {
	_, err := r.db.Exec(query)
	if err != nil {
		return
	}
}



func (s *SQLStorage) ListEvents(ctx context.Context, from, to time.Time) ([]*models.Event, error) {
    var events []*models.Event
    query := `SELECT id, title, date_time FROM events WHERE date_time BETWEEN $1 AND $2`
    err := s.db.SelectContext(ctx, &events, query, from, to)
    return events, err
}

// DeleteEvent удаляет событие по ID (строка)
func (s *SQLStorage) DeleteEvent(ctx context.Context, id string) error {
    intID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        return err
    }
    query := `DELETE FROM events WHERE id = $1`
    _, err = s.db.ExecContext(ctx, query, intID)
    return err
}