package json

import (
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
)

type JsonMetadataSerializer struct {
}

func NewJsonMetadataSerializer() *JsonMetadataSerializer {
	return &JsonMetadataSerializer{}
}

func (s *JsonMetadataSerializer) Serialize(meta core.Metadata) ([]byte, error) {
	if meta == nil {
		return nil, nil
	}

	marshal, err := jsonSerializer.Marshal(meta)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to marshal metadata")
	}

	return marshal, nil
}

func (s *JsonMetadataSerializer) Deserialize(bytes []byte) (core.Metadata, error) {
	if bytes == nil {
		return nil, nil
	}

	var meta core.Metadata

	err := jsonSerializer.Unmarshal(bytes, &meta)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to unmarshal metadata")
	}

	return meta, nil
}
