package mock

// Call struct for mock call
type Call struct {
	FunctionName string
	Params       []interface{}
	Returns      []interface{}
}
