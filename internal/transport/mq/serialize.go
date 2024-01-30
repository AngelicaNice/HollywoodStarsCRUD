package mq

import (
	"bytes"
	"encoding/json"

	audit "github.com/AngelicaNice/auditlog_mq/pkg/domain"
)

func Serialize(log audit.LogItem) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(log)
	return b.Bytes(), err
}
