package random

import (
	"fmt"
	"math/rand/v2"
	"time"
)


func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.IntN(len(letters))]
	}
	return string(result)
}

func RandomPhone() string {
	return fmt.Sprintf("+%d%d", 970, rand.IntN(1000000000))
}

func RandomZip() string {
	return fmt.Sprintf("%06d", rand.IntN(999999))
}

func RandomEmail() string {
	domains := []string{"gmail.com", "yahoo.com", "example.com", "test.com"}
	return RandomString(8) + "@" + domains[rand.IntN(len(domains))]
}

func RandomCity() string {
	cities := []string{"New York", "Paris", "London", "Tokyo", "Berlin", "Kiryat Mozkin", "Moscow"}
	return cities[rand.IntN(len(cities))]
}

func RandomRegion() string {
	regions := []string{"California", "Texas", "Bavaria", "Ontario", "Kraiot", "Quebec"}
	return regions[rand.IntN(len(regions))]
}

func RandomBrand() string {
	brands := []string{"Nike", "Adidas", "Vivienne Sabo", "Apple", "Samsung"}
	return brands[rand.IntN(len(brands))]
}

func RandomProvider() string {
	providers := []string{"paypal", "stripe", "wbpay", "alpha"}
	return providers[rand.IntN(len(providers))]
}

func RandomDate() string {
	t := time.Now().AddDate(0, 0, -rand.IntN(365))
	return t.Format(time.RFC3339)
}
