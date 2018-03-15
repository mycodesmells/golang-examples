package checks

import (
	"errors"
	"testing"
)

func TestProcessRequest_NewWay(t *testing.T) {
	type check func(*Response, error, *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasError := func(exp error) check {
		return func(_ *Response, err error, t *testing.T) {
			if exp != err {
				t.Errorf("Expected error %v, got: %v", exp, err)
			}
		}
	}
	hasNoError := func() check {
		return func(_ *Response, err error, t *testing.T) {
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}
		}
	}
	hasTotalWeight := func(exp int) check {
		return func(r *Response, _ error, t *testing.T) {
			if r.Shipment.Weight != exp {
				t.Errorf("Expected shipment weight %d, got: %d", exp, r.Shipment.Weight)
			}
		}
	}
	hasNPackages := func(n int) check {
		return func(r *Response, _ error, t *testing.T) {
			if len(r.Shipment.Packages) != n {
				t.Errorf("Expected packages count to be %d, got: %d", n, len(r.Shipment.Packages))
			}
		}
	}
	hasPackageWeight := func(index, exp int) check {
		return func(r *Response, _ error, t *testing.T) {
			weight := r.Shipment.Packages[index].Weight
			if weight != exp {
				t.Errorf("Expected package %d weight %d, got: %d", index, exp, weight)
			}
		}
	}
	hasPackageDescription := func(index int, exp string) check {
		return func(r *Response, _ error, t *testing.T) {
			desc := r.Shipment.Packages[index].Description
			if desc != exp {
				t.Errorf("Expected package %d description %s, got: %s", index, exp, desc)
			}
		}
	}

	testCases := []struct {
		desc   string
		in     *Request
		checks []check
	}{
		{
			desc: "",
			in:   &Request{},
			checks: checks(
				hasNoError(),
				hasError(errors.New("aaa")),
				hasTotalWeight(1000),
				hasNPackages(2),
				hasPackageWeight(0, 500),
				hasPackageDescription(1, "some desc"),
			),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			out, err := ProcessRequest(tC.in)
			for _, ch := range tC.checks {
				ch(out, err, t)
			}
		})
	}
}
