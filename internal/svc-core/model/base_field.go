package model

import (
	"encoding/json"

	"github.com/go-konsultin/sqlk"
	"github.com/go-konsultin/timek"
	"github.com/konsultin/project-goes-here/dto"
)

type BaseField struct {
	CreatedAt  timek.Time      `db:"createdAt"`
	UpdatedAt  timek.Time      `db:"updatedAt"`
	ModifiedBy *Subject        `db:"modifiedBy"`
	Version    int64           `db:"version"`
	Metadata   json.RawMessage `db:"metadata"`
}

func NewBaseField(subject *dto.Subject) BaseField {
	t := timek.Now()
	return BaseField{
		CreatedAt:  t,
		UpdatedAt:  t,
		ModifiedBy: NewSubject(subject),
		Version:    1,
		Metadata:   sqlk.EmptyObjectJSON,
	}
}

// NewBaseFieldFromModel creates BaseField from model.Subject (for service layer use)
func NewBaseFieldFromModel(subject *Subject) BaseField {
	t := timek.Now()
	return BaseField{
		CreatedAt:  t,
		UpdatedAt:  t,
		ModifiedBy: subject,
		Version:    1,
		Metadata:   sqlk.EmptyObjectJSON,
	}
}
