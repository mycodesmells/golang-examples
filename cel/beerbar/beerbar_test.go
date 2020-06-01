package beerbar

import (
	"testing"

	"github.com/slomek/playground/cel/models"
)

func TestBeerBar(t *testing.T) {
	countries := []*models.Country{
		{Code: "PL", BeerLegal: true, BeerAgeLimit: 18},
		{Code: "DE", BeerLegal: true, BeerAgeLimit: 16},
		{Code: "US", BeerLegal: true, BeerAgeLimit: 21},
		{Code: "SOB", BeerLegal: false}, // Soberland.
	}

	// Paul the wine guy!
	// - Won't serve you beer if you're below 30.
	// - Knows age limits, doesn't know legality.
	paul, err := NewWaiter(
		countries,
		`person.older_than(30) && meets_age_limit(person, country(countries, person))`,
	)
	if err != nil {
		t.Fatalf("Failed to hire Paul: %v", err)
	}

	// // Larry O'B.
	// // - Will ask you for ID, according to the law.
	larry, err := NewWaiter(
		countries,
		`country(countries, person).beer_legal && meets_age_limit(person, country(countries, person))`,
	)
	if err != nil {
		t.Fatalf("Failed to hire Larry: %v", err)
	}

	// // Cool Kyle.
	// // - A rebel, will serve anything to anyone.
	kyle, err := NewWaiter(
		countries,
		`true`,
	)
	if err != nil {
		t.Fatalf("Failed to hire Kyle: %v", err)
	}

	// // Lost Slawek.
	// // - A non-English speaker, to be safe he won't sell anything.
	slawek, err := NewWaiter(
		countries,
		`false`,
	)
	if err != nil {
		t.Fatalf("Failed to hire Slawek: %v", err)
	}

	waiters := map[string]*Waiter{
		"paul":   paul,
		"larry":  larry,
		"kyle":   kyle,
		"slawek": slawek,
	}

	cases := []struct {
		desc         string
		customer     *models.Person
		expectations map[string]bool
	}{
		{
			desc: "An older lady",
			customer: &models.Person{
				Name:    "Sally O'Hara",
				Country: "US",
				Age:     60,
			},
			expectations: map[string]bool{
				"paul":   true,
				"larry":  true,
				"kyle":   true,
				"slawek": false,
			},
		},
		{
			desc: "Young but legal bro",
			customer: &models.Person{
				Name:    "Tomek Kolega",
				Country: "PL",
				Age:     22,
			},
			expectations: map[string]bool{
				"paul":   false,
				"larry":  true,
				"kyle":   true,
				"slawek": false,
			},
		},
		{
			desc: "A kid",
			customer: &models.Person{
				Name:    "Kevin Homealone",
				Country: "US",
				Age:     8,
			},
			expectations: map[string]bool{
				"paul":   false,
				"larry":  false,
				"kyle":   true,
				"slawek": false,
			},
		},
		{
			desc: "A no-mans-land man",
			customer: &models.Person{
				Name:    "Viktor Navorski",
				Country: "KR",
				Age:     38,
			},
			expectations: map[string]bool{
				"paul":   false,
				"larry":  false,
				"kyle":   true,
				"slawek": false,
			},
		},
		{
			desc: "A non-drinder",
			customer: &models.Person{
				Name:    "Steve Sober",
				Country: "SOB",
				Age:     32,
			},
			expectations: map[string]bool{
				"paul":   true,
				"larry":  false,
				"kyle":   true,
				"slawek": false,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			for wname, exp := range tc.expectations {
				t.Run(wname, func(t *testing.T) {
					waiter, ok := waiters[wname]
					if !ok {
						t.Fatalf("Failed to find %q waiter!", wname)
					}

					if want, got := exp, waiter.WillServeBeer(tc.customer); want != got {
						t.Errorf("Expected Paul serving a beer to be %v, got: %v", want, got)
					}
				})
			}
		})
	}
}
