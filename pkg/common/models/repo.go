package models

import (
	"gorm.io/gorm"
)

type Repo struct {
	db          *gorm.DB
	tx          *gorm.DB
	transaction bool
}

func (r *Repo) Connection() *gorm.DB {
	if r.tx != nil && r.transaction {
		return r.tx
	}
	return r.db
}

func (r *Repo) SetConnection(db *gorm.DB) {
	r.db = db
}

func (r *Repo) Transaction() {
	r.tx = r.db.Begin()
	r.transaction = true
}

func (r *Repo) Rollback() {
	if r.tx != nil && r.transaction {
		r.tx.Rollback()
	}
	r.tx = nil
	r.transaction = false
}

func (r *Repo) Commit() {
	if r.tx != nil && r.transaction {
		r.tx.Commit()
	}
	r.tx = nil
	r.transaction = false
}
