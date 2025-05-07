package faker

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func init() {
	gofakeit.Seed(0)

	RegisterGenerator("first_name", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.FirstName() }), nil
	})

	RegisterGenerator("surname", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.LastName() }), nil
	})

	RegisterGenerator("full_name", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.Name() }), nil
	})

	RegisterGenerator("title", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.NamePrefix() }), nil
	})

	RegisterGenerator("phone", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.Phone() }), nil
	})

	RegisterGenerator("address", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.Street() }), nil
	})

	RegisterGenerator("zip_code", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.Zip() }), nil
	})

	RegisterGenerator("country", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.Country() }), nil
	})

	RegisterGenerator("url", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.URL() }), nil
	})

	RegisterGenerator("mac_address", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.MacAddress() }), nil
	})

	RegisterGenerator("past_date", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.PastDate().Format(time.RFC3339) }), nil
	})

	RegisterGenerator("price", func(_ string, params map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string {
			return fmt.Sprintf("%s%f%s", params["prefix"], gofakeit.Price(0, 100), params["suffix"])
		}), nil
	})

	RegisterGenerator("negative_price", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return fmt.Sprintf("%f", 0.0-gofakeit.Price(0, 100)) }), nil
	})

	RegisterGenerator("company", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.Company() }), nil
	})

	RegisterGenerator("credit_card", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string { return gofakeit.CreditCardNumber(nil) }), nil
	})

	RegisterGenerator("bool", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string {
			c := cases.Title(language.BritishEnglish)
			return c.String(strconv.FormatBool(gofakeit.Bool()))
		}), nil
	})

	RegisterGenerator("yes_no", func(_ string, _ map[string]string) (FieldGenerator, error) {
		return simpleFaker(func() string {
			if gofakeit.Bool() {
				return "Yes"
			} else {
				return "No"
			}
		}), nil
	})
}

type simpleFaker func() string

func (f simpleFaker) Generate() (any, error) {
	return fmt.Sprintf("%q", f()), nil
}
