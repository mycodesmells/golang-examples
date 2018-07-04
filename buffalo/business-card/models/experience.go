package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/satori/go.uuid"
)

type Experience struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	Company     string       `json:"company" db:"company"`
	From        time.Time    `json:"from" db:"from"`
	To          time.Time    `json:"to" db:"to"`
	Description nulls.String `json:"description" db:"description"`

	Position string `json:"position" db:"position"`
}

// String is not required by pop and may be deleted
func (e Experience) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Experiences is not required by pop and may be deleted
type Experiences []Experience

// String is not required by pop and may be deleted
func (e Experiences) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Validate gets run every time you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (e *Experience) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: e.Company, Name: "Company"},
		&validators.TimeIsPresent{Field: e.From, Name: "From"},
		&validators.TimeIsPresent{Field: e.To, Name: "To"},
		&validators.StringIsPresent{Field: e.Position, Name: "Position"},
	), nil
}

// ValidateSave gets run every time you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (e *Experience) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (e *Experience) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
