package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const (
	idColumnName       = "id"
	createAtColumnName = "created_at"
)

// crudQ handles common CRUD operations
type crudQ[T data.IEntity[ID], ID comparable] struct {
	db        *pgdb.DB
	tableName string
	sel       sq.SelectBuilder
}

// newCRUDQ creates a new instance of crudQ, where the T is a
// POINTER to an entity type and ID is the ID type such as uuid.UUID or uint64, etc.
func newCRUDQ[T data.IEntity[ID], ID comparable](db *pgdb.DB, tableName string) *crudQ[T, ID] {
	return &crudQ[T, ID]{
		db:        db,
		tableName: tableName,
		sel:       emptySelector(tableName),
	}
}

func emptySelector(tableName string) sq.SelectBuilder {
	return sq.Select("*").From(tableName)
}

// Insert adds a new record to the table
func (r *crudQ[T, ID]) Insert(entity T) error {
	entry := structs.Map(entity)

	return r.db.Get(entity.GetID(),
		sq.Insert(r.tableName).
			SetMap(entry).
			Suffix(fmt.Sprintf("RETURNING %s ", idColumnName)),
	)
}

// Update modifies an existing record
func (r *crudQ[T, ID]) Update(entity T) error {
	return r.db.Exec(
		sq.Update(r.tableName).
			SetMap(structs.Map(entity)).
			Where(sq.Eq{idColumnName: entity.GetID()}),
	)
}

// Get retrieves a single record by ID
func (r *crudQ[T, ID]) Get(result T) (bool, error) {
	if err := r.db.Get(result, r.sel); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed to get row: %w", err)
	}

	return true, nil
}

// Select retrieves multiple records
func (r *crudQ[T, ID]) Select() (result []T, err error) {
	if err = r.db.Select(&result, r.sel); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, nil
		}

		return nil, err
	}

	r.sel = emptySelector(r.tableName)

	return result, nil
}

// Delete removes a record by ID
func (r *crudQ[T, ID]) Delete(id ID) error {
	err := r.db.Exec(sq.Delete(r.tableName).
		Where(sq.Eq{idColumnName: id}))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	}

	return nil
}

func (r *crudQ[T, ID]) Count() (uint64, error) {
	var result uint64

	if err := r.db.Get(&result,
		sq.Select("COUNT(*)").
			FromSelect(r.sel, "filtered_select"),
	); err != nil {
		return 0, err
	}

	return result, nil
}
