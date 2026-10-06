package application

type HelloService struct {
}

func (a *HelloService) GenerateGreet(name string) string {
	return "Hilo " + name
}
