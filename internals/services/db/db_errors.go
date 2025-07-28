package db

import "fmt"

type DBError struct {
	Table   string
	Item    string
	Message string
}

func (e *DBError) Error() string {
	return fmt.Sprintf("Error %s %s: %s", e.Table, e.Item, e.Message)
}

func NewDBError(table, item, message string) *DBError {
	return &DBError{
		item,
		table,
		message,
	}
}

type AlreadyExistsError struct {
	Table string
	Item  string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("Item %s already exists in table %s", e.Item, e.Table)
}

func NewAlreadyExistsError(table, item, message string) *AlreadyExistsError {
	return &AlreadyExistsError{
		table,
		item,
	}
}

type NotFoundError struct {
	Resource string
	Item     string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Item %s not found in table %s", e.Item, e.Resource)
}

func NewNotFoundError(Resource, item string) *NotFoundError {
	return &NotFoundError{
		Resource,
		item,
	}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error for field %s: %s", e.Field, e.Message)
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		field,
		message,
	}
}

type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Internal error: %s", e.Message)
}

func NewInternalError(message string) *InternalError {
	return &InternalError{
		message,
	}
}
