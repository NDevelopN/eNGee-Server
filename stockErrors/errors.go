package stockErrors

import "fmt"

type EmptyValueError struct {
	Field string
}

func (e *EmptyValueError) Error() string {
	return fmt.Sprintf("no value provided for '%s'", e.Field)
}

type InvalidValueError[T any] struct {
	Field string
	Value T
}

func (e *InvalidValueError[T]) Error() string {
	return fmt.Sprintf("invalid value provided for '%s': %v", e.Field, e.Value)
}

type MatchNotFoundError[T any] struct {
	Space string
	Field string
	Value T
}

func (e *MatchNotFoundError[T]) Error() string {
	if e.Space == "" {
		return fmt.Sprintf("no match found for provided %s: %v", e.Field, e.Value)
	}

	return fmt.Sprintf("no match found in %s for provided %s: %v", e.Space, e.Field, e.Value)
}

type MatchFoundError[T any] struct {
	Space string
	Field string
	Value T
}

func (e *MatchFoundError[T]) Error() string {
	if e.Space == "" {
		return fmt.Sprintf("existing match found for provided %s: %v", e.Field, e.Value)

	}

	return fmt.Sprintf("%s already exists in %s for value: %v", e.Field, e.Space, e.Value)
}

type EmptySetError struct {
	Space string
	Field string
}

func (e *EmptySetError) Error() string {
	return fmt.Sprintf("no %s found in %s", e.Field, e.Space)
}

type HttpRequestError struct {
	Call string
	Code int
}

func (e *HttpRequestError) Error() string {
	return fmt.Sprintf("http request %q failed. Returned error code: %d", e.Call, e.Code)
}

var (
	EV_ERR  *EmptyValueError
	IV_ERR  *InvalidValueError[string]
	MNF_ERR *MatchNotFoundError[string]
	MF_ERR  *MatchFoundError[string]
	ES_ERR  *EmptySetError
	HR_ERR  *HttpRequestError
)
