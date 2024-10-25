package util

// supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func isSupportedCurrency(cur string) bool {
	switch cur {
	case USD, EUR, CAD:
		return true
	}
	return false
}
