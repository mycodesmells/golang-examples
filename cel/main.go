package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/slomek/playground/cel/models"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func main() {
	// Define 'environment' that composes of:
	// - type definitions
	// - variable, function declarations (names used in expressions)
	env, err := cel.NewEnv(
		cel.Types(&models.Person{}),
		cel.Declarations(
			decls.NewVar("person", decls.NewObjectType("mycodesmells.celgo.models.Person")),
			decls.NewFunction("can_drink_beer",
				decls.NewOverload("can_drink_beer_i_bool", []*exprpb.Type{decls.NewObjectType("mycodesmells.celgo.models.Person")}, decls.Bool)),
		),
	)

	// Define function overloads - actual implementations of declarations defined above.
	ff := cel.Functions(
		&functions.Overload{
			Operator: "can_drink_beer",
			Unary: func(val ref.Val) ref.Val {
				x, err := val.ConvertToNative(reflect.TypeOf(&models.Person{}))
				if err != nil {
					return types.NewErr("could not convert type %v into shipping_cost.Cart{}: %v", val.Type(), err)
				}
				person, ok := x.(*models.Person)
				if !ok {
					return types.NewErr("invalid type '%v' to getTotalWeight", val.Type())
				}
				switch {
				case person.Age > 21:
					return types.Bool(true)
				case person.Age > 18:
					return types.Bool(person.GetCountry() != "US")
				default:
					return types.Bool(false)
				}
			},
		},
	)

	ast, issues := env.Compile(`can_drink_beer(person)`)
	if issues != nil && issues.Err() != nil {
		fmt.Printf("Failed to compile expression: %v\n", issues.Errors())
		os.Exit(1)
	}

	prg, err := env.Program(ast, ff)
	if err != nil {
		fmt.Printf("program construction error: %s\n", err)
		os.Exit(1)
	}

	out, details, err := prg.Eval(map[string]interface{}{
		"person": &models.Person{Country: "US", Age: 22},
	})
	fmt.Println(out) // 'true'
	fmt.Println(details)
	fmt.Println(err)
}
