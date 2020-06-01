package basic

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func NewPhoneNumberValidator(rules []string) (*PhoneNumberValidator, error) {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("number", decls.String),
			decls.NewFunction("starts_with",
				decls.NewInstanceOverload("starts_with",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool,
				), 
			),
			decls.NewFunction("has_digits",
				decls.NewOverload("has_digits",
					[]*exprpb.Type{decls.String, decls.Int},
					decls.Bool,
				),
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}

	funcs := cel.Functions(
		&functions.Overload{
			Operator: "starts_with",
			Binary: func(lhs, rhs ref.Val) ref.Val {
				val := lhs.(types.String)
				strVal := val.Value().(string)

				prefix := rhs.(types.String)
				strPrefix := prefix.Value().(string)

				return types.Bool(strings.HasPrefix(strVal, strPrefix))
			},
		},
		&functions.Overload{
			Operator: "has_digits",
			Binary: func(lhs, rhs ref.Val) ref.Val {
				val := lhs.(types.String)
				strVal := val.Value().(string)

				length := rhs.(types.Int)
				intLength := length.Value().(int64)

				runes := []rune(strVal)

				var digits int64 = 0
				for _, r := range runes {
					if unicode.IsDigit(r) {
						digits += 1
					}
				}

				return types.Bool(digits == intLength)
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

	return &PhoneNumberValidator{
		programs: programs,
	}, nil
}

type PhoneNumberValidator struct {
	programs []cel.Program
}

func (v *PhoneNumberValidator) IsValid(number string) (bool, error) {
	for idx, program := range v.programs {
		val, _, err := program.Eval(map[string]interface{}{"number": number})
		if err != nil {
			return false, fmt.Errorf("failed to evaluate rule %d: %v", idx, err)
		}
		bVal, ok := val.Value().(bool)
		if !bVal || !ok {
			return false, nil
		}
	}
	return true, nil
}
