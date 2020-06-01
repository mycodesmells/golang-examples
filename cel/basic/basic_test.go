package basic

import "testing"

func TestPhoneNumberValidator(t *testing.T) {
	cases := []struct {
		desc   string
		input  string
		rules  []string
		result bool
	}{
		{
			desc:  "needs to start with '+48'",
			input: "+48 236",
			rules: []string{
				`number.starts_with("+48")`,
			},
			result: true,
		},
		{
			desc:  "needs to have 11 digits",
			input: "12345678901",
			rules: []string{
				`has_digits(number, 11)`,
			},
			result: true,
		},
		{
			desc:  "needs to start with '+48', have 11 digits - bad length",
			input: "+48 236",
			rules: []string{
				`number.starts_with("+48")`,
				`has_digits(number, 11)`,
			},
			result: false,
		},
		{
			desc:  "needs to start with '+48', have 11 digits - bad prefix",
			input: "12345678901",
			rules: []string{
				`number.starts_with("+48")`,
				`has_digits(number, 11)`,
			},
			result: false,
		},
		{
			desc:  "needs to start with '+48', have 11 digits - correct",
			input: "+48 786 234 283",
			rules: []string{
				`number.starts_with("+48")`,
				`has_digits(number, 11)`,
			},
			result: true,
		},
		{
			desc:  "has 9 or 11 digits - correct",
			input: "+48 786 234 283",
			rules: []string{
				`has_digits(number, 9) || has_digits(number, 11)`,
			},
			result: true,
		},
		{
			desc:  "has length of 8",
			input: "qwertyui",
			rules: []string{
				`number.size() == 8`,
			},
			result: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			validator, err := NewPhoneNumberValidator(tc.rules)
			if err != nil {
				t.Fatalf("Failed to create validator: %v", err)
			}

			valid, err := validator.IsValid(tc.input)
			if err != nil {
				t.Fatalf("Failed to check phone number %q validity: %v", tc.input, err)
			}

			if want, got := tc.result, valid; want != got {
				t.Errorf("Expected %s validity to be %v, got: %v", tc.input, want, got)
			}
		})
	}
}
