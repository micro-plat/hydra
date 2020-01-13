package db

import (
	"database/sql"

	"github.com/micro-plat/lib4go/db/tpl"
)

type MDBTrans struct {
	db       *sql.Tx
	tpl      tpl.ITPLContext
	provider string
}

func NewMDBTrans(provider string, db *sql.Tx) *MDBTrans {
	tpl, _ := tpl.GetDBContext(provider)
	return &MDBTrans{
		provider: provider,
		db:       db,
		tpl:      tpl,
	}
}
func (m *MDBTrans) Query(sql string, input map[string]interface{}) (data QueryRows, query string, args []interface{}, err error) {
	query, args = m.tpl.GetSQLContext(sql, input)
	rows, err := m.db.Query(sql, args...)
	if err != nil {
		return nil, query, args, err
	}
	data, _, err = resolveRows(rows, 0)
	return
}
func (m *MDBTrans) Scalar(sql string, input map[string]interface{}) (data interface{}, query string, args []interface{}, err error) {
	query, args = m.tpl.GetSQLContext(sql, input)
	rows, err := m.db.Query(sql, args...)
	if err != nil {
		return nil, query, args, err
	}
	result, colus, err := resolveRows(rows, 0)
	if err != nil || len(result) == 0 || len(result[0]) == 0 || len(colus) == 0 {
		return
	}
	data = result[0][colus[0]]
	return
}
func (m *MDBTrans) Execute(sql string, input map[string]interface{}) (row int64, query string, args []interface{}, err error) {
	query, args = m.tpl.GetSQLContext(sql, input)
	result, err := m.db.Exec(query, args...)
	if err != nil {
		return
	}
	row, err = result.RowsAffected()
	return
}
func (m *MDBTrans) Executes(sql string, input map[string]interface{}) (insertID int64, row int64, query string, args []interface{}, err error) {
	query, args = m.tpl.GetSQLContext(sql, input)
	result, err := m.db.Exec(query, args...)
	if err != nil {
		return
	}
	insertID, err = result.LastInsertId()
	row, err = result.RowsAffected()
	return
}
func (m *MDBTrans) Rollback() error {
	return m.db.Rollback()

}
func (m *MDBTrans) Commit() error {
	return m.db.Commit()
}
