package panicretry

import (
	"fmt"
	"log"
)

// Do calls fn and retry that if panics.
func Do(fn func() error) error {
	return defaultRetrier.Do(fn)
}

var defaultRetrier Retrier

// Retrier executes specified function and retry that if panics.
type Retrier struct {
	// Number of retry attempts.
	// zero value means infinite retry.
	MaxRetry int
	// For logging panic message
	LoggerFunc LoggerFunc
}

// Do calls fn and retry that if panics.
func (r *Retrier) Do(fn func() error) error {
	loggerFunc := r.LoggerFunc
	if loggerFunc == nil {
		loggerFunc = DefaultLoggerFunc
	}

	attempts := 1
	for {
		err := wrap(fn)
		if err == nil {
			return nil
		}
		perr, ok := err.(*panicError)
		if !ok {
			return err
		}
		loggerFunc(perr)
		if r.MaxRetry == 0 {
			continue
		}
		if r.MaxRetry < attempts {
			return perr
		}
		attempts++
	}
}

func wrap(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &panicError{message: fmt.Sprintf("panicretry: %+v", r)}
		}
	}()
	return fn()
}

type panicError struct {
	message string
}

func (e *panicError) Error() string {
	return e.message
}

// LoggerFunc is a type of function for logging panic message
type LoggerFunc func(panicErr error)

// DefaultLoggerFunc is a default LoggerFunc
var DefaultLoggerFunc LoggerFunc = func(panicErr error) {
	log.Println(panicErr.Error())
}
