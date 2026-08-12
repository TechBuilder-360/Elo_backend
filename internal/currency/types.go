package currency

type Currency struct {
	ID         string
	Name       string
	Code       string
	Symbol     string
	Multiplier int64
	IsFiat     bool
}
