package port

type HelloSerivcePort interface {
	GenerateGreet(name string) string
}
