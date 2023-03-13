package reflect

import "reflect"

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      []any
}

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	typ := reflect.TypeOf(entity)
	numMethod := typ.NumMethod()
	res := make(map[string]FuncInfo, numMethod)
	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		fn := method.Func
		numIn := fn.Type().NumIn()
		inputTypes := make([]reflect.Type, 0, numIn)
		inputVals := make([]reflect.Value, 0, numIn)

		inputTypes = append(inputTypes, typ)
		inputVals = append(inputVals, reflect.ValueOf(entity))

		for j := 1; j < numIn; j++ {
			fnInType := fn.Type().In(j)
			inputTypes = append(inputTypes, fnInType)
			inputVals = append(inputVals, reflect.Zero(fnInType))
		}

		numOut := fn.Type().NumOut()
		outputTypes := make([]reflect.Type, 0, numOut)
		outputVals := fn.Call(inputVals)
		result := make([]any, 0, numOut)
		for j := 0; j < numOut; j++ {
			outputTypes = append(outputTypes, fn.Type().Out(j))
			result = append(result, outputVals[j].Interface())
		}
		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  inputTypes,
			OutputTypes: outputTypes,
			Result:      result,
		}
	}

	return res, nil

}

