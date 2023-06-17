package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type siteRegionMappingRepo struct {
	db *sql.DB
}

func NewSiteRegionMappingRepo(db *sql.DB) port.SiteRegionMappingRepoPort {
	return &siteRegionMappingRepo{db: db}
}

func (r siteRegionMappingRepo) GetAllSiteRegionMapping() ([]port.SiteRegionMapping, error) {
	queryString := `SELECT id, code, name, area, created_at, updated_at
					FROM site_region_mapping
					WHERE code NOT LIKE 'EMPTY-%'`
	row, err := r.db.Query(queryString)
	if err != nil {
		return []port.SiteRegionMapping{}, err
	}
	results := make([]port.SiteRegionMapping, 0)
	for row.Next() {
		var result port.SiteRegionMapping
		err := row.Scan(
			&result.ID,
			&result.Code,
			&result.Name,
			&result.Area,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return []port.SiteRegionMapping{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r siteRegionMappingRepo) GetSiteRegionMappingWithPagination(limit int, offset int) ([]port.SiteRegionMapping, error) {
	queryString := `SELECT id, code, name, area, created_at, updated_at
					FROM site_region_mapping
					WHERE code NOT LIKE 'EMPTY-%'
					LIMIT ? OFFSET ?`
	row, err := r.db.Query(queryString, limit, offset)
	if err != nil {
		return []port.SiteRegionMapping{}, err
	}
	results := make([]port.SiteRegionMapping, 0)
	for row.Next() {
		var result port.SiteRegionMapping
		err := row.Scan(
			&result.ID,
			&result.Code,
			&result.Name,
			&result.Area,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return []port.SiteRegionMapping{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r siteRegionMappingRepo) GetAreaNullCity() ([]port.SiteRegionMapping, error) {
	queryString := `SELECT id, code, name, area, created_at, updated_at
					FROM site_region_mapping
					WHERE code LIKE 'EMPTY-%' AND area NOT NULL`
	row, err := r.db.Query(queryString)
	if err != nil {
		return []port.SiteRegionMapping{}, err
	}
	results := make([]port.SiteRegionMapping, 0)
	for row.Next() {
		var result port.SiteRegionMapping
		err := row.Scan(
			&result.ID,
			&result.Code,
			&result.Name,
			&result.Area,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return []port.SiteRegionMapping{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r siteRegionMappingRepo) CreateCity(siteRegionMapping port.SiteRegionMapping) error {
	queryString := `INSERT INTO site_region_mapping (code, name)
					VALUES (?, ?)`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		siteRegionMapping.Code,
		siteRegionMapping.Name,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r siteRegionMappingRepo) UpdateCity(id int, siteRegionMapping port.SiteRegionMapping) error {
	queryString := `UPDATE site_region_mapping
					SET code = ?,
						name = ?
					WHERE id = ?`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		siteRegionMapping.Code,
		siteRegionMapping.Name,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r siteRegionMappingRepo) UpdateCityToNullArea(area string) error {
	queryString := `UPDATE site_region_mapping
					SET area = null
					WHERE area = ?`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		area,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r siteRegionMappingRepo) DeleteCity(id int) error {
	queryString := `DELETE FROM site_region_mapping
					WHERE id = ?`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (r siteRegionMappingRepo) CreateBatchSiteRegionMappings(siteRegionMappings []port.SiteRegionMapping) error {
	if len(siteRegionMappings) == 0 {
		err := errors.New("site region mapping list must not be empty")
		return err
	}

	valueQueries := make([]string, 0)
	values := make([]interface{}, 0)

	for _, siteRegionMapping := range siteRegionMappings {
		valueQueries = append(valueQueries, "(?, ?, ?)")
		values = append(values, siteRegionMapping.Code, siteRegionMapping.Name, siteRegionMapping.Area)
	}

	queryString := fmt.Sprintf("INSERT OR REPLACE INTO site_region_mapping (code, name, area) VALUES %s", strings.Join(valueQueries, ","))

	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	return nil
}

func (r siteRegionMappingRepo) UpdateSiteRegionMapping(area string, cityCodeListString string) error {
	queryString := fmt.Sprintf(`UPDATE site_region_mapping
					SET area = '%s'
					WHERE code IN (%s)`, strings.ToUpper(area), cityCodeListString)
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r siteRegionMappingRepo) TotalSiteRegionMapping() (int, error) {
	queryString := `SELECT COUNT(*) 
					FROM site_region_mapping
					WHERE code NOT LIKE 'EMPTY-%'`

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
