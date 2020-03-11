# panicretry

Retries function call if panics

```go
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
```