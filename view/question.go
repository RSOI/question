package view

import "github.com/RSOI/question/model"

// ValidateNewQuestion returns true if all the required form values are passed
func ValidateNewQuestion(data model.Question) (bool, string) {
	var required []string
	if data.Title == "" {
		required = append(required, "title")
	}

	if *data.Content == "" {
		required = append(required, "content")
	}

	if data.AuthorID == 0 {
		required = append(required, "author_id")
	}

	return len(required) == 0, fieldsToString(required)
}

// ValidateDeleteQuestion returns true if parameter to delete found
func ValidateDeleteQuestion(data model.Question) (bool, string) {
	var required []string

	if data.ID == 0 {
		required = append(required, "id")
	} else {
		return true, "id"
	}

	if data.AuthorID == 0 {
		required = append(required, "author_id")
	} else {
		return true, "author_id"
	}

	return false, fieldsToString(required)
}

// ValidateUpdateQuestion returns true if parameter to delete found
func ValidateUpdateQuestion(data model.Question) (bool, string) {
	return data.ID != 0, "id"
}
