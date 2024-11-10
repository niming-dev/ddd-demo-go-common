package expression

type Expression interface {
	Evaluate(ctx ExecuteContext) (*Data, error)
	Dump(deep int) string
}
