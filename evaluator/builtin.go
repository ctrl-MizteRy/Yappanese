package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
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
	//There are still problem appending another array into a 2d array
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
				if _, ok := arg.Elements[0].(*object.Array); ok {
					if arr, ok := args[1].(*object.Array); ok {
						newArr.Elements = append(newArr.Elements, arr)
						return newArr
					}
				}
				if arr, ok := args[1].(*object.Array); ok {
					for _, elemnt := range arr.Elements {
						newArr.Elements = append(newArr.Elements, elemnt)
					}

					return newArr
				} else if reflect.TypeOf(arg.Elements[0]) == reflect.TypeOf(args[1]) {
					var elements []object.Object
					elements = append(elements, arg.Elements...)
					elements = append(elements, args[1])
					return &object.Array{Elements: elements}
				} else {
					return newError("Appending error: cannot append an array of %T with %T", arg.Elements[0], args[1])
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
				switch arg.Inspect() {
				case "\\t":
					msg = append(msg, "\t")
				case "\\n":
					msg = append(msg, "\n")
				default:
					msg = append(msg, arg.Inspect())
				}
			}
			fmt.Print(strings.Join(msg, " "))
			return nil
		},
	},

	"pop": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 2 {
				return newError("Unexpect amount of arguement, expect=2 (Array, index), or 1 (Array)")
			}
			if len(args) == 2 {
				if arr, ok := args[0].(*object.Array); ok {
					if idx, ok := args[1].(*object.Integer); ok {
						if len(arr.Elements) <= int(idx.Value) {
							return newError("Error: index out of range, array contain=%d elements",
								len(arr.Elements))
						}
						obj := arr.Elements[idx.Value]
						left := arr.Elements[0:idx.Value]
						right := arr.Elements[idx.Value+1:]
						arr.Elements = left
						for _, element := range right {
							arr.Elements = append(arr.Elements, element)
						}
						return obj
					} else {
						return newError("Unexpected type error: expect= Integer, got= %T (%+v)",
							args[1], args[1])
					}
				} else {
					return newError("Unexpected type error: expect= Array, got= %T", args[0])
				}
			}
			if arr, ok := args[0].(*object.Array); ok {
				index := len(arr.Elements) - 1
				obj := arr.Elements[index]
				arr.Elements = arr.Elements[:index]
				return obj
			} else {
				return newError("Unexpected type error: expect= Array, got= %T", args[0])
			}
		},
	},

	"keys": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Unexpected argument length, expect=1, got=%d", len(args))
			}

			hash, ok := args[0].(*object.Hash)
			if !ok {
				return newError("Unexpected argument type: expect= HASH, got= %s", args[0].Type())
			}

			keys := []object.Object{}
			for _, key := range hash.Keys {
				keys = append(keys, key)
			}
			return &object.Array{Elements: keys}
		},
	},

	"values": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Argument length error: expect= 1, got= %d", len(args))
			}

			hash, ok := args[0].(*object.Hash)
			if !ok {
				return newError("Argument type error: expect= HASH, got= %s", args[0].Type())
			}

			value := []object.Object{}

			for _, val := range hash.Pairs {
				value = append(value, val.Value)
			}

			return &object.Array{Elements: value}
		},
	},
}
