package marshaller

type Parameter int

const (
	Param_Amount Parameter = iota + 1 // EnumIndex = 1
	Param_Legal                       // EnumIndex = 2
)

func (p Parameter) String() string {
	return [...]string{"Amount", "Legal"}[p-1]
}

func (p Parameter) EnumIndex() int {
	return int(p)
}
