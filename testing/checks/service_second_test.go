package checks

import (
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
			desc:   "no items",
			in:     &Request{},
			checks: checks(hasError(ErrNoItems)),
		},
		{
			desc: "weight from request",
			in: &Request{
				Items:  []*Item{{Name: "Product 1", Weight: 300}},
				Weight: 500,
			},
			checks: checks(
				hasNoError(),
				hasTotalWeight(500),
			),
		}, {
			desc: "weight from items",
			in: &Request{
				Items: []*Item{{Name: "Product 1", Weight: 300}},
			},
			checks: checks(
				hasNoError(),
				hasTotalWeight(300),
			),
		}, {
			desc: "package per item",
			in: &Request{
				Items: []*Item{
					{Name: "Product 1", Weight: 300},
					{Name: "Product 2", Weight: 250},
				},
			},
			checks: checks(
				hasNoError(),
				hasNPackages(2),
				hasPackageDescription(0, "Product 1"),
				hasPackageWeight(0, 300),
				hasPackageDescription(1, "Product 2"),
				hasPackageWeight(1, 250),
			),
		}, {
			desc: "one package for all",
			in: &Request{
				Items: []*Item{
					{Name: "Product 1", Weight: 300},
					{Name: "Product 2", Weight: 250},
				},
				OnePackage: true,
			},
			checks: checks(
				hasNoError(),
				hasNPackages(1),
				hasPackageDescription(0, "Grouped goods"),
				hasPackageWeight(0, 550),
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
