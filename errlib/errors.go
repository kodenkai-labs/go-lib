package errlib

// AppError is an interface that abstracts application layer specific errors.
type AppError interface {
	Error() string
	Code() Code
	Slug() Slug
}

type appError struct {
	err  error
	code Code
	slug Slug
}

// NewAppError creates a new application error.
func NewAppError(err error, code Code, slug Slug) AppError {
	return appError{
		err:  err,
		code: code,
		slug: slug,
	}
}

func (e appError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return ""
}

func (e appError) Code() Code {
	return e.code
}

func (e appError) Slug() Slug {
	return e.slug
}
