package checks

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

func ProcessRequest(req *Request) (*Response, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("no items provided")
	}

	shipment := &Shipment{
		ID:     uuid.NewV4().String(),
		Weight: req.Weight,
	}

	totalWeight := 0
	for _, item := range req.Items {
		totalWeight += item.Weight
	}
	if shipment.Weight == 0 {
		shipment.Weight = totalWeight
	}

	if !req.OnePackage {
		packages := []*Package{}
		for _, it := range req.Items {
			packages = append(packages, &Package{
				ID:          uuid.NewV4().String(),
				Description: it.Name,
				Weight:      it.Weight,
			})
		}
		shipment.Packages = packages
	} else {
		shipment.Packages = []*Package{
			{
				ID:          uuid.NewV4().String(),
				Description: "Grouped goods",
				Weight:      totalWeight,
			},
		}
	}

	return &Response{
		Shipment: shipment,
	}, nil
}
