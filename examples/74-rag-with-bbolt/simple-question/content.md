### Calling A Plug-in's Exports

This plug-in was written in Rust and it does one thing, it counts vowels in a string. As such, it exposes one "export" function: `count_vowels`. We can call exports using [extism.Plugin.Call](https://pkg.go.dev/github.com/extism/go-sdk#Plugin.Call).
Let's add that code to our main func:

```go
func main() {
    // ...

	data := []byte("Hello, World!")
	exit, out, err := plugin.Call("count_vowels", data)
	if err != nil {
		fmt.Println(err)
		os.Exit(int(exit))
	}

	response := string(out)
	fmt.Println(response)
    // => {"count": 3, "total": 3, "vowels": "aeiouAEIOU"}
}
```

Running this should print out the JSON vowel count report:

```bash
$ go run main.go
# => {"count":3,"total":3,"vowels":"aeiouAEIOU"}
```