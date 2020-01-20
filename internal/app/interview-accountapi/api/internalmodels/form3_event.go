package internalmodels

import (
	"encoding/json"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type Form3Event struct {
	Id             uuid.UUID       `json:"id,omitempty"`
	OrganisationId uuid.UUID       `json:"organisation_id,omitempty"`
	Version        int64           `json:"version"`
	EventType      string          `json:"event_type,omitempty"`
	DataRecordType string          `json:"data_record_type,omitempty"`
	RecordType     string          `json:"record_type,omitempty"`
	Data           json.RawMessage `json:"data,omitempty"`
	Record         AuditRecord     `json:"record"`
}

type AuditRecord struct {
	RecordType  string          `json:"record_type,omitempty"`
	RecordId    uuid.UUID       `json:"record_id,omitempty"`
	ActionedBy  uuid.UUID       `json:"actioned_by,omitempty"`
	ActionTime  strfmt.DateTime `json:"action_time,omitempty"`
	Description string          `json:"description,omitempty"`
	BeforeData  json.RawMessage `json:"before_data,omitempty"`
	AfterData   json.RawMessage `json:"after_data,omitempty"`
}
