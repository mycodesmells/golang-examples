package checks

type Request struct {
	Items      []*Item
	Weight     int
	OnePackage bool
}

type Response struct {
	Shipment *Shipment
}

type Shipment struct {
	ID       string
	Packages []*Package
	Weight   int
}

type Package struct {
	ID          string
	Description string
	Weight      int
}

type Item struct {
	Name   string
	Weight int
}
