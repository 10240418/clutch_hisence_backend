package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type KeyManagementService struct{}

func NewKeyManagementService() (IKeyManagementService, error) {
	return &KeyManagementService{}, nil
}

func (kms *KeyManagementService) GeneratePublicKeyFromDeviceID(deviceID string) (string, error) {
	if deviceID == "" {
		return "", fmt.Errorf("deviceID cannot be empty")
	}

	// 使用SHA256哈希算法对DeviceID进行哈希处理，生成公钥
	hash := sha256.Sum256([]byte(deviceID))
	publicKey := hex.EncodeToString(hash[:])

	return publicKey, nil
}

func (kms *KeyManagementService) ValidateDeviceIDAndPublicKey(deviceID, publicKey string) bool {
	if deviceID == "" || publicKey == "" {
		return false
	}

	// 重新生成公钥并与提供的公钥进行比较
	expectedPublicKey, err := kms.GeneratePublicKeyFromDeviceID(deviceID)
	if err != nil {
		return false
	}

	return expectedPublicKey == publicKey
}
