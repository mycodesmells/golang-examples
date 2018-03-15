package checks

import "testing"

func TestProcessRequest_OldWay(t *testing.T) {
	testCases := []struct {
		desc string
		in   *Request
		out  *Response
		err  error
	}{
		{
			desc: "",
			in:   &Request{},
			out:  &Response{},
			err:  nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			out, err := ProcessRequest(tC.in)
			if err != nil {
				if err != tC.err {
					t.Errorf("Expected error %v, got: %v", tC.err, err)
				}
				return
			}
			if tC.err != nil {
				t.Errorf("Expected error: %v", tC.err)
				return
			}

			if out.Shipment.Weight != tC.out.Shipment.Weight {
				t.Errorf("Expected shipment weight %d, got: %d", tC.out.Shipment.Weight, out.Shipment.Weight)
			}

			if len(out.Shipment.Packages) != len(tC.out.Shipment.Packages) {
				t.Errorf("Expected packages count to be %d, got: %d", len(tC.out.Shipment.Packages), len(out.Shipment.Packages))
			}

			for i, p := range out.Shipment.Packages {
				exp := tC.out.Shipment.Packages[i]
				if p.Weight != exp.Weight {
					t.Errorf("Expected package %d weight %d, got: %d", i, exp.Weight, p.Weight)
				}
				if p.Description != exp.Description {
					t.Errorf("Expected package %d description %s, got: %s", i, exp.Description, p.Description)
				}
			}
		})
	}
}
