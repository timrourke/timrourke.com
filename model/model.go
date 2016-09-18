package model

import (
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

// Model interface defines base model interface
type Model struct {
}

func (m Model) GetID() string {
	return ""
}
