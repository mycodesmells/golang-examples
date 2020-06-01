package beerbar

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/slomek/playground/cel/models"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type Waiter struct {
	beerRules []cel.Program
	countries []*models.Country
}

func (w *Waiter) WillServeBeer(p *models.Person) bool {
	for idx, program := range w.beerRules {
		val, _, err := program.Eval(map[string]interface{}{
			"countries": w.countries,
			"person":    p,
		})
		if err != nil {
			fmt.Printf("Failed to evaluate rule %d: %v\n", idx, err)
			continue
		}
		bVal, ok := val.Value().(bool)
		if ok {
			return bVal
		}
	}
	return false
}

func NewWaiter(countries []*models.Country, rules ...string) (*Waiter, error) {
	env, err := cel.NewEnv(
		cel.Types(&models.Person{}, &models.Country{}),
		cel.Declarations(
			decls.NewVar("person", decls.NewObjectType("mycodesmells.celgo.models.Person")),
			decls.NewVar("countries", decls.NewListType(decls.NewObjectType("mycodesmells.celgo.models.Country"))),
			decls.NewFunction("older_than",
				decls.NewInstanceOverload("older_than",
					[]*exprpb.Type{decls.NewObjectType("mycodesmells.celgo.models.Person"), decls.Int},
					decls.Bool,
				),
			),
			decls.NewFunction("meets_age_limit",
				decls.NewOverload("meets_age_limit",
					[]*exprpb.Type{
						decls.NewObjectType("mycodesmells.celgo.models.Person"),
						decls.NewObjectType("mycodesmells.celgo.models.Country"),
					},
					decls.Bool,
				),
			),
			decls.NewFunction("country",
				decls.NewOverload("country",
					[]*exprpb.Type{
						decls.NewListType(decls.NewObjectType("mycodesmells.celgo.models.Country")),
						decls.NewObjectType("mycodesmells.celgo.models.Person"),
					},
					decls.NewObjectType("mycodesmells.celgo.models.Country"),
				),
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}

	reg := types.NewRegistry(
		&models.Country{},
		&models.Person{},
	)

	funcs := cel.Functions(
		&functions.Overload{
			Operator: "older_than",
			Binary: func(lhs, rhs ref.Val) ref.Val {
				x, _ := lhs.ConvertToNative(reflect.TypeOf(&models.Person{}))
				person := x.(*models.Person)

				ageTyped := rhs.(types.Int)
				age := ageTyped.Value().(int64)

				return types.Bool(int64(person.Age) >= age)
			},
		},
		&functions.Overload{
			Operator: "meets_age_limit",
			Binary: func(lhs, rhs ref.Val) ref.Val {
				x, _ := lhs.ConvertToNative(reflect.TypeOf(&models.Person{}))
				person := x.(*models.Person)

				y, _ := rhs.ConvertToNative(reflect.TypeOf(&models.Country{}))
				country := y.(*models.Country)

				return types.Bool(person.Age > country.BeerAgeLimit)
			},
		},
		&functions.Overload{
			Operator: "country",
			Binary: func(lhs, rhs ref.Val) ref.Val {
				x, _ := lhs.ConvertToNative(reflect.TypeOf([]*models.Country{}))
				countries := x.([]*models.Country)

				y, _ := rhs.ConvertToNative(reflect.TypeOf(&models.Person{}))
				person := y.(*models.Person)
				countryCode := person.Country

				for _, country := range countries {
					if country.Code == countryCode {
						return reg.NativeToValue(country)
					}
				}

				return nil
			},
		},
	)

	programs := make([]cel.Program, 0, len(rules))
	for _, rule := range rules {
		ast, issues := env.Compile(rule)
		if err := issues.Err(); err != nil {
			return nil, fmt.Errorf("failed to compile rule %q: %w", rule, err)
		}

		program, err := env.Program(ast, funcs)
		if err != nil {
			return nil, fmt.Errorf("failed to program AST for rule %q: %w", rule, err)
		}

		programs = append(programs, program)
	}

	return &Waiter{
		countries: countries,
		beerRules: programs,
	}, nil
}
