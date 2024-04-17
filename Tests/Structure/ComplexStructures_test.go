package Structure

// func Test_it_can_validate_valid_complex_structures(t *testing.T) {
// 	// Setup
// 	cases := []struct {
// 		jsonString     []byte
// 		shouldFail     bool
// 		expectedErrors map[string][]string
// 	}{
// 		{[]byte(`{"Data": [1,2,3]}`), false, map[string][]string{}},
// 	}
//
// 	type Subscription struct {
// 		Id                 *string `json:"id" validation:"nullable|uuid"`
// 		Type               string  `json:"type" validation:"in:recurring,unscheduled"`
// 		IsFirstTransaction bool    `json:"isFirstTransaction" validation:"present|bool"`
// 	}
//
// 	type Transaction struct {
// 		Id                    string  `json:"id" validation:"required|uuid"`
// 		Amount                int     `json:"amount" validation:"required|int|min:0"`
// 		Currency              string  `json:"currency" validation:"required|minLen:2|len:3"`
// 		MerchantTransactionId *string `json:"merchantTransactionId" validation:"present|nullable"`
// 		Reference             *string `json:"reference" validation:"present|nullable"`
// 		TextOnStatement       *string `json:"textOnStatement" validation:"present|nullable"`
// 	}
//
// 	type Merchant struct {
// 		Id         string  `json:"id" validation:"required|string"`
// 		WebsiteUrl *string `json:"websiteUrl" validation:"nullable|url"`
// 	}
//
// 	type Customer struct {
// 		Ip *string `json:"ip" validation:"nullable|ip"`
// 	}
//
// 	type Card struct {
// 		Pan          string `json:"pan" validation:"required|minLen:8|maxLen:19"`
// 		ExpireMonth  string `json:"expireMonth" validation:"required|len:2"`
// 		ExpireYear   string `json:"expireYear" validation:"required|len:2"`
// 		SecurityCode string `json:"csc" validation:"required|minLen:3|maxLen:4"`
// 	}
//
// 	type Token struct {
// 		Framework   string `json:"framework" validation:"required|in:VTS,M4M"`
// 		Tan         string `json:"tan" validation:"required|minLen:8|maxLen:19"`
// 		ExpireMonth string `json:"expire_month" validation:"required|len:2"`
// 		ExpireYear  string `json:"expire_year" validation:"required|len:2"`
// 		Eci         string `json:"eci" validation:"required|len:2"`
// 		Cryptogram  string `json:"cryptogram" validation:"required|len28"`
// 	}
//
// 	type PaymentMethod struct {
// 		Card  *Card  `url:"card" validation:"requiredWithout:token|nullable|object"`
// 		Token *Token `url:"token" validation:"requiredWithout:card|nullable|object"`
// 	}
//
// 	type Exemption struct {
// 		Moment string `json:"moment" validation:"required|string"`
// 		Type   string `json:"type" validation:"required|string"`
// 	}
//
// 	type SCA struct {
// 		Exemption *Exemption `json:"exemption" validation:"present|nullable|object"`
// 	}
//
// 	type CitAuthorizationRequest struct {
// 		Transaction   Transaction   `json:"transaction" validation:"required"`
// 		Customer      Customer      `json:"customer" validation:"required"`
// 		Merchant      Merchant      `json:"merchant" validation:"required"`
// 		PaymentMethod PaymentMethod `json:"paymentMethod" validation:"required"`
// 		Subscription  *Subscription `json:"subscription"`
// 		SCA           *SCA          `json:"SCA"`
// 	}
//
// 	for i, testCase := range cases {
// 		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
// 			// Arrange
// 			var errorBag *JsonValidator.ErrorBag
//
// 			// Act
// 			_, err := JsonValidator.Validate[testData](testCase.jsonString)
// 			_ = errors.As(err, &errorBag)
//
// 			// Assert
// 			if testCase.shouldFail {
// 				require.True(t, errorBag != nil)
// 				require.True(t, errorBag.HasFailedKeyAndRule("Data", "array"))
// 				require.Equal(t, 1, errorBag.CountErrors())
// 			} else {
// 				require.True(t, errorBag == nil)
// 				require.Equal(t, 0, errorBag.CountErrors())
// 			}
// 		})
// 	}
// }
