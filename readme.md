# ePay JSON to struct validator
With actual JSON validation included!

This JSON validator operates differently from popular struct validators like [`gookit/validate`](https://github.com/gookit/validate) and [`go-playground/validator`](https://github.com/go-playground/validator). It validates JSON directly against a target data structure, rather than after unmarshalling.

The validator distinguishes between a missing key and a key with a Go zero value. This redefines the semantics of the `required` validation rule and introduces the new rule `present`. In traditional struct validators, `required` means "not the data type's default zero value" but in this JSON validator, it means "the key exists in JSON and the value is not null."

This allows requiring a field from a client while permitting the data type's default zero value.

# Usage
```go
type MyRequest struct {
    Id string `json:"id" validate:"required|string"`
    User User `json:"user" validate:"required|object"`
}

type User struct {
    Name *string `json:"name" validate:"present|nullable"`
}

func main() {
    jsonBytes := (...)
    
    myRequest, err := JsonValidator.Validate[MyRequest](jsonBytes)
}
```

# Rules
| Name                             | Description                                                                                                                         |
|----------------------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| `nullable/nilable`               | Explicitly allows the field to be nullable. <br/>(Disables other validation rules when value is null)                               |
| `required`                       | The field key must be both present in the JSON and have a non-null value.                                                           |
| `requiredWith:{x}`               | The field key must be both present in the JSON and have a non-null value if sibling field `{x}` is present.                         |
| `requiredWithout:{x}`            | The field key must be both present in the JSON and have a non-null value if sibling field `{x}` is not present.                     |
| `requiredWithAny:{x},{z},...`    | The field key must be both present in the JSON and have a non-null value if any of the sibling fields `{x},{z},...` is present.     |
| `requiredWithoutAny:{x},{z},...` | The field key must be both present in the JSON and have a non-null value if any of the sibling fields `{x},{z},...` is not present. |
| `requiredWithAll:{x},{z},...`    | The field key must be both present in the JSON and have a non-null value if all of the sibling fields `{x},{z},...` is present.     |
| `requiredWithoutAll:{x},{z},...` | The field key must be both present in the JSON and have a non-null value if all of the sibling fields `{x},{z},...` is not present. |
| `present`                        | The field key must be present in the JSON.                                                                                          |
| `len:{n}`                        | Checks value is countable and length is exactly `{n}` (`string`, `array`, `object`)                                                 |
| `lenMax:{n}`                     | Checks value is countable and length is at most `{n}` (`string`, `array`, `object`)                                                 |
| `lenMin:{n}`                     | Checks value is countable and length is at least `{n}` (`string`, `array`, `object`)                                                |
| `lenBetween:{n},{m}`             | Checks value is countable and length is between `{n}` and `{m}` inclusively (`string`, `array`, `object`)                           |
| `array`                          | Checks value is an `array`/`slice`                                                                                                  |
| `object`                         | Checks value is an `object`/`map`/`struct`                                                                                          |
| `string`                         | Checks value is a `string`                                                                                                          |
| `int/integer`                    | Checks value is an `integer`                                                                                                        |
| `float`                          | Checks value is a `float`                                                                                                           |
| `bool/boolean`                   | Checks value is a `boolean`                                                                                                         |
| `in:{a},{b},...`                 | Checks that the value is within the the list `{a},{b},...`                                                                          |
| `uuid`                           | Checks that the value is an UUID string.                                                                                            |
| `regex:{x}`                      | Checks that the value is a string that is matched by the `{x}` regex definition                                                     |
| `between:{x},{z}`                | Checks that the value is a number between `{x}` and `{z}` inclusively                                                               |
| `url`                            | Checks that the value is a non-empty URL string                                                                                     |
| `ip`                             | Checks that the value is a non-empty IP string                                                                                      |
| `email`                          | Checks that the value is a non-empty email string                                                                                   |
