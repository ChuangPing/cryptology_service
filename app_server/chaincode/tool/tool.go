package tool

import (
	"encoding/binary"
	"fmt"
)

// Int64ToBytes int64 to byte
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// BytesToInt64 byte to int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// ConstructUserKey user info id
func ConstructUserKey(userId string) string {
	return fmt.Sprintf("user_%s", userId)
}

// ConstructSymmetricKeyEncKey info id   -- use this id can getState
func ConstructSymmetricKeyEncKey(userId string) string {
	return fmt.Sprintf("symmetricKeyEnc_%s", userId)
}

// ConstructKeywordIndexKey get keywordIndex id
func ConstructKeywordIndexKey() string {
	return "keywordIndexInfo"
}

// ConstructMetadataKey metadata key
func ConstructMetadataKey(metaId string) string {
	return fmt.Sprintf("metadata_%s", metaId)
}

// ConstructMetadataLogKey  metadata log key
func ConstructMetadataLogKey() string {
	return "metadataLogKey"
}

