package helper

import "math/rand"

func GenerateFilename(originalFilename string) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	kode := make([]byte, 35)
	for i := range kode {
		kode[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	fileName := string(kode) + "-" + originalFilename

	return fileName
}