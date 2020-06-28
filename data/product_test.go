package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "grant",
		Price: 100.00,
		SKU:   "ppp-ddd-ttt",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
