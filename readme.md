# ePay JSON to struct validator
With actual json validation.

# Usage
```
type MyRequest struct {
    Id string `json:"id" validate:"required|int"`
    User User `json:"user" validate:"required|object"`
}

type User struct {
    Name *string `json:"name" validate:"present|nullable"`
}

func main() {
    jsonBytes := (...)
    var request MyRequest
    
    data, err := JsonValidator.Validate[MyRequest](jsonBytes)
}
```