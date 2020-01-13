package db

import (
	"database/sql"

	"github.com/micro-plat/lib4go/db/tpl"
)

var _ IDB = &MDB{}

type MDB struct {
	db       *sql.DB
	tpl      tpl.ITPLContext
	provider string
}

func NewMDB(provider string, db *sql.DB) *MDB {
	tpl, _ := tpl.GetDBContext(provider)
	return &MDB{
		provider: provider,
		db:       db,
		tpl:      tpl,
	}
}
func (m *MDB) Query(sql string, input map[string]interface{}) (data QueryRows, query string, args []interface{}, err error) {
	query, args = m.tpl.GetSQLContext(sql, input)
	rows, err := m.db.Query(sql, args...)
	if err != nil {
		return nil, query, args, err
	}
	data, _, err = resolveRows(rows, 0)
	return
}
func (m *MDB) Scalar(sql string, input map[string]interface{}) (data interface{}, query string, args []interface{}, err error) {
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
func (m *MDB) Execute(sql string, input map[string]interface{}) (row int64, query string, args []interface{}, err error) {
	query, args = m.tpl.GetSQLContext(sql, input)
	result, err := m.db.Exec(query, args...)
	if err != nil {
		return
	}
	row, err = result.RowsAffected()
	return
}
func (m *MDB) Executes(sql string, input map[string]interface{}) (insertID int64, row int64, query string, args []interface{}, err error) {
	query, args = m.tpl.GetSQLContext(sql, input)
	result, err := m.db.Exec(query, args...)
	if err != nil {
		return
	}
	insertID, err = result.LastInsertId()
	row, err = result.RowsAffected()
	return
}
func (m *MDB) ExecuteSP(procName string, input map[string]interface{}, output ...interface{}) (row int64, query string, err error) {
	query, args := m.tpl.GetSPContext(procName, input)
	ni := append(args, output...)
	result, err := m.db.Exec(query, ni...)
	if err != nil {
		return 0, query, err
	}
	row, err = result.RowsAffected()
	return
}
func (m *MDB) Begin() (IDBTrans, error) {
	t, err := m.db.Begin()
	if err != nil {
		return nil, err
	}
	return NewMDBTrans(m.provider, t), nil
}
func (m *MDB) Close() {
	m.db.Close()
}
