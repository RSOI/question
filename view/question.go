package view

import (
	"github.com/RSOI/question/model"
	"github.com/RSOI/question/ui"
)

// ValidateNewQuestion returns nil if all the required form values are passed
func ValidateNewQuestion(data model.Question) error {
	if data.Title == "" ||
		data.Content == nil ||
		*data.Content == "" ||
		data.AuthorID == 0 {
		return ui.ErrFieldsRequired
	}
	return nil
}

// ValidateDeleteQuestion returns true if parameter to delete found
func ValidateDeleteQuestion(data model.Question) (string, error) {
	if data.ID != 0 {
		return "id", nil
	}

	if data.AuthorID != 0 {
		return "author_id", nil
	}

	return "", ui.ErrFieldsRequired
}

// ValidateUpdateQuestion returns nil if parameter to delete found
func ValidateUpdateQuestion(data model.Question) error {
	if data.ID != 0 {
		return nil
	}
	return ui.ErrFieldsRequired
}
