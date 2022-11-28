// Package hashing позволяет хэшировать данные.
package hashing

// Signer интерфейс хэш-подписи. Позволяет хэшировать с помощью ключа.
type Signer interface {
	Hash(data string, key string) []byte     // Хэширует данные
	HashToHex(hash []byte) string            // Конвертирует хэш в hex
	ValidHash(sign []byte, hash string) bool // Проверяет хэш на подлинность
}
