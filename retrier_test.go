package panicretry_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kei2100/panicretry"
)

func ExampleDo() {
	counter := 0
	err := panicretry.Do(func() error {
		if counter < 10 {
			counter++
			panic("oops")
		}
		return nil
	})

	fmt.Println(err)
	fmt.Println(counter)

	// Output:
	// <nil>
	// 10
}

func TestRetrier_Do(t *testing.T) {
	panicretry.DefaultLoggerFunc = func(_ error) {
		// noop
	}

	t.Run("panics/infinite retry", func(t *testing.T) {
		pr := panicretry.Retrier{}
		counter := 0
		got := pr.Do(func() error {
			if counter < 10 {
				counter++
				panic("oops")
			}
			return nil
		})

		if g, w := got, error(nil); g != w {
			t.Errorf("err got %v, want %v", g, w)
		}
		if g, w := counter, 10; g != w {
			t.Errorf("counter got %v, want %v", g, w)
		}
	})

	t.Run("panics/10 times retry", func(t *testing.T) {
		pr := panicretry.Retrier{MaxRetry: 10}
		counter := 0
		got := pr.Do(func() error {
			counter++
			panic("oops")
		})

		if got == nil {
			t.Error("got no error, want an error")
			return
		}
		if g, w := got.Error(), fmt.Sprintf("panicretry: oops"); g != w {
			t.Errorf("error message got %v, want %v", g, w)
		}
		if g, w := counter, 11; g != w {
			t.Errorf("counter got %v, want %v", g, w)
		}
	})

	t.Run("no panics/no error", func(t *testing.T) {
		pr := panicretry.Retrier{}
		counter := 0
		got := pr.Do(func() error {
			counter++
			return nil
		})

		if g, w := got, error(nil); g != w {
			t.Errorf("error got %v, want %v", g, w)
		}
		if g, w := counter, 1; g != w {
			t.Errorf("counter got %v, want %v", g, w)
		}
	})

	t.Run("no panics/an error", func(t *testing.T) {
		pr := panicretry.Retrier{}
		var someErr = errors.New("oh")
		counter := 0
		got := pr.Do(func() error {
			counter++
			return someErr
		})

		if g, w := got, someErr; g != w {
			t.Errorf("error got %v, want %v", g, w)
		}
		if g, w := counter, 1; g != w {
			t.Errorf("counter got %v, want %v", g, w)
		}
	})
}
