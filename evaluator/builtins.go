package evaluator

import (
	"fmt"
	"github.com/adiletelf/abyss/object"
	"io/ioutil"
	"math"
	"math/rand"
)

var builtins = map[string]*object.Builtin{
	"len": { // len(arr) | len(string)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": { // first(arr) -> first element of array
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": { // last(arr) -> last element of array
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"tail": { // tail(arr) -> new array except first element
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"push": { // push(arr, value)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments, got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"array": { // array(size, defaultValue: INT | FLOAT)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments, got=%d, want=2", len(args))
			}

			if args[0].Type() != object.INTEGER_OBJ {
				return newError("first argument (size) to `array` must be INT, got=%s", args[0].Type())
			}
			size := args[0].(*object.Integer).Value

			switch defaultValue := args[1].(type) {
			case *object.Integer:
				array := make([]object.Object, size)

				var i int64
				for i = 0; i < size; i++ {
					array[i] = &object.Integer{Value: defaultValue.Value}
				}
				return &object.Array{Elements: array}
			case *object.Float:
				array := make([]object.Object, size)

				var i int64
				for i = 0; i < size; i++ {
					array[i] = &object.Float{Value: defaultValue.Value}
				}
				return &object.Array{Elements: array}
			default:
				return newError("second argument to `array` not supported, got=%s", args[1].Type())
			}
		},
	},
	"range": { // range(min, max) -> returns array with elements from min to max
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments, got=%d, want 2", len(args))
			}

			if args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
				return newError("wrong types of arguments to `range`, got=(%s, %s), want=(INT, INT)",
					args[0].Type(), args[1].Type())
			}

			min := args[0].(*object.Integer).Value
			max := args[1].(*object.Integer).Value
			size := int(max - min)
			if size < 0 {
				size *= -1
			}

			if size == 0 {
				return &object.Array{}
			}

			array := make([]object.Object, size)
			for i := 0; i < size; i++ {
				array[i] = &object.Integer{Value: min}
				if min < max {
					// range(0,4) -> [0, 1, 2, 3]
					min++
				} else {
					// range(4,0) -> [4, 3, 2, 1]
					min--
				}
			}

			return &object.Array{Elements: array}
		},
	},
	"print": { // prints every argument, returns null
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"abs": { // abs(value: INT | FLOAT)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				v := arg.Value
				if v < 0 {
					v *= -1
				}
				return &object.Integer{Value: v}
			case *object.Float:
				v := arg.Value
				if v < 0 {
					v *= -1
				}
				return &object.Float{Value: v}
			default:
				return newError("argument to 'abs' not supported, got=%s", arg.Type())
			}
		},
	},
	"pow": { // pow(base: INT|FLOAT, exp: INT|FLOAT)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments, got=%d, want=2", len(args))
			}

			first := args[0]
			second := args[1]
			if first.Type() != object.INTEGER_OBJ && first.Type() != object.FLOAT_OBJ {
				return newError("first argument to `pow` is not INT | FLOAT, got=%s", first.Type())
			}
			if second.Type() != object.INTEGER_OBJ && second.Type() != object.FLOAT_OBJ {
				return newError("second argument to `pow` is not INT | FLOAT, got=%s", second.Type())
			}

			var base, exp float64
			switch f := first.(type) {
			case *object.Integer:
				base = float64(f.Value)
			case *object.Float:
				base = f.Value
			}

			switch s := second.(type) {
			case *object.Integer:
				exp = float64(s.Value)
			case *object.Float:
				exp = s.Value
			}

			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"random": { // random() -> random FLOAT in [0.0, 1.0)
		Fn: func(args ...object.Object) object.Object {
			return &object.Float{Value: rand.Float64()}
		},
	},
	"sqrt": { // sqrt(value: INT | FLOAT)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				v := arg.Value
				return &object.Float{Value: math.Sqrt(float64(v))}
			case *object.Float:
				v := arg.Value
				return &object.Float{Value: math.Sqrt(v)}
			default:
				return newError("argument to `sqrt` not supported, got=%s",
					args[0].Type())
			}
		},
	},
	"open": { // read file as string
		Fn: func(args ...object.Object) object.Object {
			path := ""

			// We need at least one arg
			if len(args) < 1 {
				return newError("wrong number of arguments. got=%d, want=1+", len(args))
			}

			// Get the filename
			switch args[0].(type) {
			case *object.String:
				path = args[0].(*object.String).Value
			default:
				return newError("argument to `file` not supported, got=%s", args[0].Type())

			}

			// Create the object
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return newError(err.Error())
			}

			text := string(content)
			return &object.String{Value: text}
		},
	},
	"type": { // return type of object as string
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to `type`, got=%d, want=1")
			}

			arg := args[0]
			return &object.String{Value: string(arg.Type())}
		},
	},
}
