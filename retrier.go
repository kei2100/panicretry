package panicretry

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
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
		perr, ok := err.(*panicRetry)
		if !ok {
			return err
		}
		loggerFunc(perr)
		if r.MaxRetry == 0 {
			continue
		}
		if r.MaxRetry < attempts {
			panic(perr.message)
		}
		attempts++
	}
}

func wrap(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			message := fmt.Sprintf("%+v", r)
			frame := make([]string, 0)
			for depth := 2; depth < 10; depth++ {
				_, file, line, ok := runtime.Caller(depth)
				if !ok {
					break
				}
				frame = append(frame, fmt.Sprintf("    %v:%d", file, line))
			}
			err = &panicRetry{
				message: message,
				frame:   frame,
			}
		}
	}()
	return fn()
}

type panicRetry struct {
	message string
	frame   []string
}

func (e *panicRetry) Error() string {
	return e.message
}

func (e *panicRetry) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "panicretry: %s\n%s", e.message, strings.Join(e.frame, "\n"))
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// LoggerFunc is a type of function for logging panic message
type LoggerFunc func(panicErr error)

// DefaultLoggerFunc is a default LoggerFunc
var DefaultLoggerFunc LoggerFunc = func(panicErr error) {
	log.Printf("%+v", panicErr)
}
