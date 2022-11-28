package hmac

import "fmt"

func Example() {
	key := "your_secret_key"
	signer := NewHmacSigner()

	hash := signer.Hash("some data to store", key)
	hexHash := signer.HashToHex(hash)
	valid := signer.ValidHash(hash, hexHash)

	fmt.Println(hash)
	fmt.Println(hexHash)
	fmt.Println(valid)

	// Output:
	// [204 109 33 117 253 133 173 194 243 112 71 84 141 151 92 249 35 107 39 203 147 37 220 15 58 130 39 16 19 161 59 89]
	// cc6d2175fd85adc2f37047548d975cf9236b27cb9325dc0f3a82271013a13b59
	// true
}
