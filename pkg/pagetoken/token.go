package pagetoken

import (
	"encoding/base64"

	"google.golang.org/protobuf/proto"
)

func EncodePageToken[T proto.Message](msg T) (string, error) {
	if proto.Message(msg) == nil {
		return "", nil
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func DecodePageToken[T proto.Message](token string, msg T) error {
	if token == "" {
		return nil
	}

	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return err
	}

	return proto.Unmarshal(data, msg)
}
