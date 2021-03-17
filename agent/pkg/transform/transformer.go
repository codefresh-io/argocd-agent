package transform

type Transformer interface {
	Transform(data interface{}) interface{}
}
