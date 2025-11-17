# AAR

Simple API authentication result

## Usage
```go
aar, err := aar.New("UNIQUE-NAME-IN-YOUR-APP%s%d", "abc", 123)
if err != nil {
    return
}

token, err := aar.SetDurcation(10).Read()
if err != nil {
    // Get token
    token = app.GetToken()
    aar.Write(token)
}
// Use token in you app
```
