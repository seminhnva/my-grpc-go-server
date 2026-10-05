package port

type HelloSerivcePort interface {
	GenerateHello(name string) string
}
