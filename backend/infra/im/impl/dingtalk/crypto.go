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

package dingtalk

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func verifySignature(token, timestamp, nonce, encrypted, signature string) bool {
	items := []string{token, timestamp, nonce, encrypted}
	sort.Strings(items)
	sum := sha1.Sum([]byte(strings.Join(items, "")))
	return strings.EqualFold(hex.EncodeToString(sum[:]), signature)
}

func decryptBody(aesKey, encrypted string) ([]byte, error) {
	if strings.TrimSpace(aesKey) == "" {
		return nil, fmt.Errorf("missing dingtalk aes key")
	}
	if strings.TrimSpace(encrypted) == "" {
		return nil, fmt.Errorf("empty dingtalk encrypt payload")
	}

	key, err := base64.StdEncoding.DecodeString(aesKey + "=")
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid dingtalk aes key length: %d", len(key))
	}

	cipherText, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid dingtalk cipher text length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plain := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, key[:aes.BlockSize]).CryptBlocks(plain, cipherText)
	plain, err = pkcs7Unpad(plain, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	if len(plain) < 20 {
		return nil, fmt.Errorf("invalid dingtalk plaintext length")
	}

	msgLen := binary.BigEndian.Uint32(plain[16:20])
	if int(20+msgLen) > len(plain) {
		return nil, fmt.Errorf("invalid dingtalk message length")
	}

	return plain[20 : 20+msgLen], nil
}

func pkcs7Unpad(src []byte, blockSize int) ([]byte, error) {
	if len(src) == 0 || len(src)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded bytes")
	}

	padding := int(src[len(src)-1])
	if padding == 0 || padding > blockSize || padding > len(src) {
		return nil, fmt.Errorf("invalid padding size")
	}
	if !bytes.Equal(src[len(src)-padding:], bytes.Repeat([]byte{byte(padding)}, padding)) {
		return nil, fmt.Errorf("invalid padding content")
	}

	return src[:len(src)-padding], nil
}
