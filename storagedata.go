package storagenomad

import (
	"time"
)

// StorageData describes the data that is saved in a Secure Variable
type StorageData struct {
	Value    []byte    `json:"value"`
	Modified time.Time `json:"modified"`
}
