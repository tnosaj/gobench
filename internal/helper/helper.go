package helper

import (
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal/db"
)

// GenerateRow returns a row
func GenerateRow(rand internal.Random) db.Row {
	return db.Row{
		K:   rand.Intn(2147483647),
		C:   randomString(120, rand),
		Pad: randomString(60, rand),
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int, rand internal.Random) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
