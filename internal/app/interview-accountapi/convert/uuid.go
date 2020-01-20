package convert

import (
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

func ToUUID(strUuid strfmt.UUID) (uuid.UUID, error) {
	return uuid.Parse(string(strUuid))
}

func ToUUIDs(strArr []string) ([]uuid.UUID, error) {
	var uuids []uuid.UUID

	for _, str := range strArr {
		id, err := uuid.Parse(str)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, id)
	}
	return uuids, nil
}

func FromUUID(id uuid.UUID) strfmt.UUID {
	return strfmt.UUID(id.String())
}

func FromUUIDToPtr(id uuid.UUID) *strfmt.UUID {
	u := FromUUID(id)
	return &u
}
