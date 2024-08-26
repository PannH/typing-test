# Typing Test
This is a simple typing test program written in Go that measures your typing speed in WPM (Words Per Minute) and accuracy.

## Usage
You can run the program using `go run`
```bash
go run .\main\main.go
```

You can also compile it using `go build`
```bash
go build -o typingtest.exe .\main\main.go
```

You can control the amount of words in the test using the `-len` flag
```bash
go run .\main\main.go -len 25
```
*This will generate a test with 25 words.*

## Example
![Usage Example](./assets/usage_example.gif)