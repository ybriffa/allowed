package allowed

import (
	"testing"
)

type first struct {
	A string `allowed:"POST"`
	B string
}

type second struct {
	C string `allowed:"POST"`
}

type third struct {
	first  *first `allowed:"POST"`
	second *second
}

func TestAllowed(t *testing.T) {
	tests := []struct {
		tagValue    string
		value       interface{}
		shouldError bool
	}{
		// 0 - Invalid data input
		{
			tagValue:    "POST",
			value:       "42",
			shouldError: true,
		},

		// 1 - Valid basic case
		{
			tagValue: "POST",
			value: first{
				A: "toto",
			},
			shouldError: false,
		},

		// 2 - Invalid basic case
		{
			tagValue: "POST",
			value: first{
				B: "toto",
			},
			shouldError: true,
		},

		// 3 - Basic pointer case
		{
			tagValue: "POST",
			value: &first{
				A: "toto",
			},
			shouldError: false,
		},

		// 4 - Invalid Array
		{
			tagValue: "POST",
			value: []string{
				"toto",
			},
			shouldError: true,
		},

		// 5 - Valid Array with pointer to struct
		{
			tagValue: "POST",
			value: []*first{
				&first{
					A: "toto",
				},
			},
			shouldError: false,
		},

		// 6 - Valid Array with pointer to struct with invalid type
		{
			tagValue: "POST",
			value: []*first{
				&first{
					B: "toto",
				},
			},
			shouldError: true,
		},

		// 7 - Invalid map Type
		{
			tagValue: "POST",
			value: map[string]string{
				"key": "value",
			},
			shouldError: true,
		},

		// 8 - Valid map
		{
			tagValue: "POST",
			value: map[string]*first{
				"key": &first{
					A: "toto",
				},
			},
			shouldError: false,
		},

		// 9 - Valid map and invalid field set
		{
			tagValue: "POST",
			value: map[string]*first{
				"key": &first{
					B: "toto",
				},
			},
			shouldError: true,
		},

		// 10 - Valid fields in struct
		{
			tagValue: "POST",
			value: &third{
				first: &first{
					A: "toto",
				},
			},
			shouldError: false,
		},

		// 10 - Invalid fields in field struct
		{
			tagValue: "POST",
			value: &third{
				first: &first{
					B: "toto",
				},
			},
			shouldError: true,
		},

		// 11 - Invalid fields set even if child field is correct
		{
			tagValue: "POST",
			value: &third{
				second: &second{
					C: "toto",
				},
			},
			shouldError: true,
		},
	}

	for i, test := range tests {
		result := Check(test.tagValue, test.value)
		if test.shouldError != (result != nil) {
			t.Fatalf("test #%d failed: expected to crash %v and got error: %v", i, test.shouldError, result)
		}
	}
}
