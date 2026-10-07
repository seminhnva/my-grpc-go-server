package port

type HelloSerivcePort interface {
	GenerateGreet(name string) string
}

type BankServicePort interface {
	FindCurrentBalance(account string) float64
	GetLatestExchaneRate() float64
}
