/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package im

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func verifyWeComSignature(token, timestamp, nonce, encrypted, signature string) bool {
	items := []string{token, timestamp, nonce, encrypted}
	sort.Strings(items)
	sum := sha1.Sum([]byte(strings.Join(items, "")))
	return strings.EqualFold(hex.EncodeToString(sum[:]), signature)
}

func decryptWeComMessage(encodedAESKey, encrypted string) (plainText string, corpID string, err error) {
	key, err := base64.StdEncoding.DecodeString(encodedAESKey + "=")
	if err != nil {
		return "", "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", "", err
	}

	if len(key) != 32 {
		return "", "", fmt.Errorf("invalid aes key length: %d", len(key))
	}
	if len(cipherText) == 0 || len(cipherText)%aes.BlockSize != 0 {
		return "", "", fmt.Errorf("invalid cipher text length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	plain := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, key[:aes.BlockSize]).CryptBlocks(plain, cipherText)
	plain, err = pkcs7Unpad(plain, aes.BlockSize)
	if err != nil {
		return "", "", err
	}
	if len(plain) < 20 {
		return "", "", fmt.Errorf("invalid plaintext length")
	}

	msgLen := binary.BigEndian.Uint32(plain[16:20])
	if int(20+msgLen) > len(plain) {
		return "", "", fmt.Errorf("invalid message length")
	}

	return string(plain[20 : 20+msgLen]), string(plain[20+msgLen:]), nil
}

func encryptWeComMessage(encodedAESKey, corpID, plainText string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(encodedAESKey + "=")
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", fmt.Errorf("invalid aes key length: %d", len(key))
	}

	raw := make([]byte, 16)
	if _, err = rand.Read(raw); err != nil {
		return "", err
	}

	content := bytes.NewBuffer(raw)
	if err = binary.Write(content, binary.BigEndian, uint32(len(plainText))); err != nil {
		return "", err
	}
	content.WriteString(plainText)
	content.WriteString(corpID)

	plain := pkcs7Pad(content.Bytes(), aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	encrypted := make([]byte, len(plain))
	cipher.NewCBCEncrypter(block, key[:aes.BlockSize]).CryptBlocks(encrypted, plain)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func pkcs7Pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	return append(src, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func pkcs7Unpad(src []byte, blockSize int) ([]byte, error) {
	if len(src) == 0 || len(src)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded bytes")
	}

	padding := int(src[len(src)-1])
	if padding == 0 || padding > blockSize || padding > len(src) {
		return nil, fmt.Errorf("invalid padding size")
	}

	for _, b := range src[len(src)-padding:] {
		if int(b) != padding {
			return nil, fmt.Errorf("invalid padding content")
		}
	}

	return src[:len(src)-padding], nil
}
