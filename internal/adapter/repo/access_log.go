package repo

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type accessLogRepo struct {
	db *sql.DB
}

func NewAccessLogRepo(db *sql.DB) port.AccessLogRepoPort {
	return &accessLogRepo{
		db: db,
	}
}

func (r accessLogRepo) GetAccessLog(limit int, offset int) ([]port.AccessLog, error) {
	queryString := `SELECT a.id, a.message, a.by_user_id, COALESCE(u.username, "DELETED"), a.created_at
					FROM access_log a
					LEFT JOIN user u
					ON (a.by_user_id = u.id)
					ORDER BY a.id DESC
					LIMIT ? OFFSET ?`

	row, err := r.db.Query(queryString, limit, offset)
	if err != nil {
		return []port.AccessLog{}, err
	}
	results := make([]port.AccessLog, 0)
	for row.Next() {
		var result port.AccessLog
		err := row.Scan(
			&result.ID,
			&result.Message,
			&result.ByUserID,
			&result.ByUsername,
			&result.CreatedAt,
		)
		if err != nil {
			return []port.AccessLog{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r accessLogRepo) CreateAccessLog(accessLog port.AccessLog) error {
	queryString := `INSERT INTO access_log (message, by_user_id)
					VALUES (?, ?)`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		accessLog.Message,
		accessLog.ByUserID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r accessLogRepo) TotalAccessLog() (int, error) {
	queryString := `SELECT COUNT(*) FROM access_log`

	row := r.db.QueryRow(queryString)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var numberOfRecords int
	if err := row.Scan(&numberOfRecords); err != nil {
		return 0, err
	}

	return numberOfRecords, nil
}
