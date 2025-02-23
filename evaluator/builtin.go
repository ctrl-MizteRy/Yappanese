package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"yap/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, expect=1, got=%d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},

	"scan": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("the function is not taking in any argument")
			}

			sc := bufio.NewScanner(os.Stdin)
			sc.Scan()
			return &object.String{Value: sc.Text()}
		},
	},

	"append": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of argument, expected=2, got=%d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				var elements []object.Object
				elements = append(elements, arg.Elements...)
				newArr := &object.Array{Elements: elements}
				if arr, ok := args[1].(*object.Array); ok {
					for _, elemnt := range arr.Elements {
						newArr.Elements = append(newArr.Elements, elemnt)
					}

					return &object.Array{Elements: newArr.Elements}
				} else {
					return newError("Appending error: cannot append an array with %T", arr)
				}
			default:
				return newError("Function does not appending of %T and %T", args[0], args[1])
			}
		},
	},

	"yap": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			msg := []string{}

			for _, arg := range args {
				msg = append(msg, arg.Inspect())
			}
			fmt.Println(strings.Join(msg, " "))
			return nil
		},
	},
}
