package reflect_study

import "reflect"

type FuncInfo struct {
	Name        string         // 方法名
	InputTypes  []reflect.Type // 入参类型列表
	OutputTypes []reflect.Type // 输出类型列表
	Result      []any          // 输出结果 (入参值为零值的情况下)
}

// IterateFunc 输出方法信息并且执行调用
// 注意事项：
// • 方法接收器
// • 以结构体作为输入， 那么只能访问到结构体作为接收器的方法
// • 以指针作为输入， 那么能访问到任何接收器的方法
// • 输入的第一个参数， 永远都是接收器本身
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
