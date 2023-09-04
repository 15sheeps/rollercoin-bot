package utils

import (
    "strconv"
    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "encoding/base64"
    "sort"
    "fmt"
    "bytes"
)

func hashifyUserid(userId string) string {
    digits := []rune{}
    chars := []rune(userId)
    sumDigits := 0

    for _, char := range chars {
        if digit, err := strconv.Atoi(string(char)); err == nil {
            sumDigits += digit
            digits = append(digits, char)
        }
    }

    if (sumDigits == 0) {
        return "1509fa0bcf93e7bfe8038fc0886e84d4"
    }

    sort.Slice(digits, func(i int, j int) bool {
        return digits[i] < digits[j] 
    })

    concatenated := fmt.Sprintf("%s%s%s", string(digits), string(chars[:len(digits)]), strconv.Itoa(sumDigits))

    hash := md5.Sum([]byte(concatenated))

    return fmt.Sprintf("%x", hash)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func EncryptCmd(cmd, userId string) string {
    plaintext := []byte(cmd)
    key := []byte(hashifyUserid(userId))
    iv := []byte("dYQ9R99bkKLsLHad")

    plaintext = PKCS7Padding(plaintext, aes.BlockSize)

    block, _ := aes.NewCipher(key)

    ciphertext := make([]byte, len(plaintext))

    mode := cipher.NewCBCEncrypter(block, iv)

    mode.CryptBlocks(ciphertext, plaintext)

    return base64.StdEncoding.EncodeToString(ciphertext)
}