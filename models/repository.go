package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type DynamicRepository struct {
	db     *sql.DB
	schema *Schema
}

func NewDynamicRepository(db *sql.DB, schema *Schema) *DynamicRepository {
	return &DynamicRepository{db: db, schema: schema}
}

func (r *DynamicRepository) Create(data map[string]interface{}) (int64, error) {
	cols := []string{}
	placeholders := []string{}
	values := []interface{}{}

	for _, field := range r.schema.Fields {
		if field.PrimaryKey {
			continue
		}
		if val, ok := data[field.Name]; ok {
			cols = append(cols, field.Name)
			placeholders = append(placeholders, "?")
			values = append(values, val)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		r.schema.TableName,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(values...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *DynamicRepository) Update(id interface{}, data map[string]interface{}) error {
	cols := []string{}
	values := []interface{}{}
	var pkName string

	for _, field := range r.schema.Fields {
		if field.PrimaryKey {
			pkName = field.Name
			continue
		}
		if val, ok := data[field.Name]; ok {
			cols = append(cols, fmt.Sprintf("%s = ?", field.Name))
			values = append(values, val)
		}
	}

	if pkName == "" {
		return fmt.Errorf("nenhuma chave primária definida no schema")
	}

	values = append(values, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		r.schema.TableName,
		strings.Join(cols, ", "),
		pkName,
	)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	return err
}

func (r *DynamicRepository) Delete(id interface{}) error {
	pkName := ""
	for _, field := range r.schema.Fields {
		if field.PrimaryKey {
			pkName = field.Name
			break
		}
	}

	if pkName == "" {
		return fmt.Errorf("nenhuma chave primária definida no schema")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", r.schema.TableName, pkName)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (r *DynamicRepository) FindByID(id interface{}) (map[string]interface{}, error) {
	pkName := ""
	for _, field := range r.schema.Fields {
		if field.PrimaryKey {
			pkName = field.Name
			break
		}
	}
	if pkName == "" {
		return nil, fmt.Errorf("nenhuma chave primária definida no schema")
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", r.schema.TableName, pkName)
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanRowToMap(rows)
	}
	return nil, sql.ErrNoRows
}

func (r *DynamicRepository) FindAll(page, limit int, search string) ([]map[string]interface{}, int, error) {
	var query strings.Builder
	var countQuery strings.Builder
	args := []interface{}{}

	query.WriteString("SELECT * FROM ")
	query.WriteString(r.schema.TableName)
	countQuery.WriteString("SELECT COUNT(*) FROM ")
	countQuery.WriteString(r.schema.TableName)

	if search != "" {
		whereClause := []string{}
		searchLike := fmt.Sprintf("%%%s%%", search)
		for _, field := range r.schema.Fields {
			if field.Type == "string" || field.Type == "text" {
				whereClause = append(whereClause, fmt.Sprintf("%s LIKE ?", field.Name))
				args = append(args, searchLike)
			}
		}
		if len(whereClause) > 0 {
			clause := fmt.Sprintf(" WHERE %s", strings.Join(whereClause, " OR "))
			query.WriteString(clause)
			countQuery.WriteString(clause)
		}
	}

	var totalRecords int
	err := r.db.QueryRow(countQuery.String(), args...).Scan(&totalRecords)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset))

	rows, err := r.db.Query(query.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		rowMap, err := scanRowToMap(rows)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, rowMap)
	}

	return results, totalRecords, nil
}

func scanRowToMap(rows *sql.Rows) (map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	vals := make([]sql.RawBytes, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range vals {
		scans[i] = &vals[i]
	}

	if err := rows.Scan(scans...); err != nil {
		return nil, err
	}

	rowMap := make(map[string]interface{})
	for i, col := range cols {
		if vals[i] == nil {
			rowMap[col] = nil
		} else {
			rowMap[col] = string(vals[i])
		}
	}

	return rowMap, nil
}
