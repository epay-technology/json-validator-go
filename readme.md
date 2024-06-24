# ePay JSON to struct validator
With actual JSON validation included!

```bash
go get github.com/epay-technology/json-validator-go
```

This JSON validator operates differently from popular struct validators like [`gookit/validate`](https://github.com/gookit/validate) and [`go-playground/validator`](https://github.com/go-playground/validator). 
**It validates JSON directly against a target data structure, rather than after unmarshalling**.

The validator distinguishes between a missing key and a given value with a data type's default zero value. 
This redefines the semantics of the `required` validation rule and introduces the new rule `present`. 
In traditional struct validators, `required` means "not the data type's default zero value" 
but in this JSON validator, it means "the key exists in JSON and the value is not null."

This allows requiring a field from a client while permitting the data type's default zero value.

## Advantages over Traditional Struct Validators
### Better Handling of Invalid Value Types
Traditional struct validators struggle to handle cases where the JSON value type doesn't match the struct field type. 
For instance, if a struct requires an integer but receives a string in the JSON, the default Go JSON unmarshaler would return a generic error, making it hard for struct validators to inform the client about the mismatch. 
Consequently, these validators often return a generic error code 400 with a vague message like "invalid payload," leaving clients unsure how to correct their code.

In contrast, the JsonValidator enables APIs to return more informative error responses. 
By validating the JSON directly, it allows for precise error messages, such as returning a 422 status code with a validation error stating "The field must be an integer." 
This clarity empowers clients to align their code with API specifications more effectively.

### Introduction of the "Present" Rule
Another significant benefit is the introduction of the "present" rule, which is not feasible with traditional struct validators. 
This rule mandates that a field must be present within the JSON, ensuring that clients send values for all required fields rather than relying on the default zero values of Go struct data types.

This capability enables APIs to differentiate between a client not sending any value and a client intentionally sending the default zero value. 
Consequently, it allows for scenarios where numeric fields can accept values like 0 while enforcing the requirement for clients to provide explicit values.

## Usage
```go
type MyRequest struct {
    Id string `json:"id" validation:"required|string|uuid"` // Id must be present with non-null uuid string
    User *User `json:"user" validation:"required|object"`   // User must be present with non-null object value
}

type User struct {
    Name *string `json:"name" validation:"present|nullable|lenMax:255"` // Name must be present, but can be null or a string with a maximum length of 255 chars
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
| `requireOneInGroup:{groupName}`  | Exactly one field with the `{groupName}` must be present and have a non-null value                                                  |
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
| `uuid`                           | Checks that the value is a valid non-zero UUID string.                                                                              |
| `zeroableUuid`                   | Checks that the value is any valid UUID string.                                                                                     |
| `regex:{x}`                      | Checks that the value is a string that is matched by the `{x}` regex definition                                                     |
| `between:{x},{z}`                | Checks that the value is a number between `{x}` and `{z}` inclusively                                                               |
| `min:{x}`                        | Checks that the value is a number greater than or equal to `{x}`                                                                    |
| `max:{x}`                        | Checks that the value is a number less than or equal to `{x}`                                                                       |
| `url`                            | Checks that the value is a non-empty URL string                                                                                     |
| `ip`                             | Checks that the value is a non-empty IP string                                                                                      |
| `email`                          | Checks that the value is a non-empty email string                                                                                   |
