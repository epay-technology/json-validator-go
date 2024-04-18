package Structure

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_complex_structures(t *testing.T) {
	// Setup
	type ErrorList map[string][]string

	cases := []struct {
		jsonString     []byte
		expectedErrors ErrorList
	}{
		{
			jsonString:     []byte(`{"transaction": {"id": "fd72503c-84d9-46f1-896f-4e9d774229dc", "amount": 100, "currency": "DKK"}, "customer": {"ip": null}, "merchant": {"id": "123456"}, "paymentMethod": {"card": {"pan": "4111111111111111", "expireMonth": "01", "expireYear": "25", "csc": "987"}}}`),
			expectedErrors: ErrorList{},
		},
		{
			jsonString: []byte(`{"transaction": {"id": "fd72503c-84d9-46f1-896f-4e9d774229dc", "amount": 100, "currency": "DKK"}, "customer": {"ip": null}, "merchant": {"id": "123456"}, "paymentMethod": {"card": {"pan": "4111111111111111", "expireMonth": "01", "expireYear": "25", "csc": null}}}`),
			expectedErrors: ErrorList{
				"paymentMethod.card.csc": []string{"required"},
			},
		},
		{
			jsonString: []byte(`{"transaction": {"id": "fd72503c-84d9-46f1-896f-4e9d774229dc", "amount": 100, "currency": "DKK"}, "customer": {"ip": null}, "merchant": {"id": "123456"}, "paymentMethod": {"card": {"pan": "4111111111111111", "expireMonth": "01", "expireYear": "25"}}}`),
			expectedErrors: ErrorList{
				"paymentMethod.card.csc": []string{"required"},
			},
		},
		{
			jsonString: []byte(`{"transaction": {"id": "fd72503c-84d9-46f1-896f-4e9d774229dc", "amount": 100, "currency": "DKK"}, "customer": {"ip": null}, "merchant": {"id": "123456"}, "paymentMethod": {}}`),
			expectedErrors: ErrorList{
				"paymentMethod.card":  []string{"requiredWithout"},
				"paymentMethod.token": []string{"requiredWithout"},
			},
		},
		{
			jsonString: []byte(`{}`),
			expectedErrors: ErrorList{
				"transaction":   []string{"required"},
				"customer":      []string{"required"},
				"merchant":      []string{"required"},
				"paymentMethod": []string{"required"},
			},
		},
		{
			jsonString: []byte(`{"transaction": {"id": "fd72503c-84d9-46f1-896f-4e9d774229dc", "amount": 100, "currency": "DKK"}, "customer": {"ip": "1.0.1"}, "merchant": {"id": "123456"}, "paymentMethod": {"card": {"pan": "4111111111111111", "expireMonth": "01", "expireYear": "25", "csc": "987"}}}`),
			expectedErrors: ErrorList{
				"customer.ip": []string{"ip"},
			},
		},
	}

	type Subscription struct {
		Id                 *string `json:"id" validation:"nullable|uuid"`
		Type               string  `json:"type" validation:"in:recurring,unscheduled"`
		IsFirstTransaction bool    `json:"isFirstTransaction" validation:"present|bool"`
	}

	type Transaction struct {
		Id                    string `json:"id" validation:"required|uuid"`
		Amount                int    `json:"amount" validation:"required|int|min:0"`
		Currency              string `json:"currency" validation:"required|minLen:2|len:3"`
		MerchantTransactionId string `json:"merchantTransactionId" validation:"nullable|minLen:1"`
		Reference             string `json:"reference" validation:"nullable|minLen:1"`
		TextOnStatement       string `json:"textOnStatement" validation:"nullable|minLen:1"`
	}

	type Merchant struct {
		Id         string  `json:"id" validation:"required|string"`
		WebsiteUrl *string `json:"websiteUrl" validation:"nullable|url"`
	}

	type Customer struct {
		Ip *string `json:"ip" validation:"nullable|ip"`
	}

	type Card struct {
		Pan          string `json:"pan" validation:"required|minLen:8|maxLen:19"`
		ExpireMonth  string `json:"expireMonth" validation:"required|len:2"`
		ExpireYear   string `json:"expireYear" validation:"required|len:2"`
		SecurityCode string `json:"csc" validation:"required|lenBetween:3,4"`
	}

	type Token struct {
		Framework   string `json:"framework" validation:"required|in:VTS,M4M"`
		Tan         string `json:"tan" validation:"required|minLen:8|maxLen:19"`
		ExpireMonth string `json:"expire_month" validation:"required|len:2"`
		ExpireYear  string `json:"expire_year" validation:"required|len:2"`
		Eci         string `json:"eci" validation:"required|len:2"`
		Cryptogram  string `json:"cryptogram" validation:"required|len28"`
	}

	type PaymentMethod struct {
		Card  *Card  `json:"card" validation:"requiredWithout:Token|nullable|object"`
		Token *Token `json:"token" validation:"requiredWithout:Card|nullable|object"`
	}

	type Exemption struct {
		Moment string `json:"moment" validation:"required|string"`
		Type   string `json:"type" validation:"required|string"`
	}

	type SCA struct {
		Exemption *Exemption `json:"exemption" validation:"present|nullable|object"`
	}

	type CitAuthorizationRequest struct {
		Transaction   Transaction   `json:"transaction" validation:"required"`
		Customer      Customer      `json:"customer" validation:"required"`
		Merchant      Merchant      `json:"merchant" validation:"required"`
		PaymentMethod PaymentMethod `json:"paymentMethod" validation:"required"`
		Subscription  *Subscription `json:"subscription"`
		SCA           *SCA          `json:"SCA"`
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Arrange
			var errorBag *JsonValidator.ErrorBag

			// Act
			var data CitAuthorizationRequest
			err := JsonValidator.NewValidator().Validate(testCase.jsonString, &data)
			_ = errors.As(err, &errorBag)

			// Assert
			if len(testCase.expectedErrors) == 0 {
				require.NoError(t, err)
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}

			totalExpectedErrors := 0

			for path, rules := range testCase.expectedErrors {
				totalExpectedErrors += len(rules)

				for _, rule := range rules {
					require.True(t, errorBag.HasFailedKeyAndRule(path, rule))
				}
			}

			require.Equal(t, totalExpectedErrors, errorBag.CountErrors(), errorBag.Error())
		})
	}
}
