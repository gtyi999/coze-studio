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

package feishu

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func verifySignature(signature, timestamp, nonce, encryptKey string, body []byte) error {
	if signature == "" {
		return fmt.Errorf("missing feishu signature")
	}
	if timestamp == "" || nonce == "" {
		return fmt.Errorf("missing feishu signature headers")
	}
	if encryptKey == "" {
		return fmt.Errorf("missing feishu encrypt key")
	}

	sum := sha256.Sum256([]byte(timestamp + nonce + encryptKey + string(body)))
	expected := hex.EncodeToString(sum[:])
	if expected != signature {
		return fmt.Errorf("invalid feishu signature")
	}

	return nil
}

func decryptBody(encryptKey, encrypt string) ([]byte, error) {
	if encryptKey == "" {
		return nil, fmt.Errorf("missing feishu encrypt key")
	}
	if encrypt == "" {
		return nil, fmt.Errorf("empty feishu encrypt payload")
	}

	key := sha256.Sum256([]byte(encryptKey))
	encryptedData, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return nil, err
	}
	if len(encryptedData) < aes.BlockSize || len(encryptedData)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid encrypted payload length")
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	iv := encryptedData[:aes.BlockSize]
	cipherText := encryptedData[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	plainText, err = pkcs7Unpad(plainText, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded data")
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize || padding > len(data) {
		return nil, fmt.Errorf("invalid padding size")
	}
	if !bytes.Equal(bytes.Repeat([]byte{byte(padding)}, padding), data[len(data)-padding:]) {
		return nil, fmt.Errorf("invalid padding")
	}

	return data[:len(data)-padding], nil
}
