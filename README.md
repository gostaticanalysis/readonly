# readonly

[![godoc.org][godoc-badge]][godoc]

`readonly` checks assignment to package variables.

```
package a

var (
	V1 int = 10 // OK - first assignment
	V2 int
)

func init() {
	V1 = 20  // OK - in init
	V2 = 100 // OK - in init
}

func main() {
	V1 = 30  // NG
	V2 = 200 // NG
	// assign
	V2 = 300 // OK - explicit assign
}
```

<!-- links -->
[godoc]: https://godoc.org/github.com/tenntenn/underlying
[godoc-badge]: https://img.shields.io/badge/godoc-reference-4F73B3.svg?style=flat-square&label=%20godoc.org

