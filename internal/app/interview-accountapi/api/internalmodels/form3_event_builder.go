package internalmodels

import (
	"context"
	"encoding/json"
	"time"

	"github.com/form3tech/go-security/security"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/settings"
	"github.com/go-openapi/strfmt"
	uuid "github.com/google/uuid"
)

type Form3EventBuilder struct {
	event Form3Event
}

func NewForm3EventBuilder() *Form3EventBuilder {
	builder := &Form3EventBuilder{
		event: Form3Event{
			Record: AuditRecord{},
		},
	}
	return builder.
		RandomId().
		CurrentActionTime()
}
func (b *Form3EventBuilder) RandomId() *Form3EventBuilder {
	b.event.Id = uuid.New()
	return b
}
func (b *Form3EventBuilder) CurrentActionTime() *Form3EventBuilder {
	b.event.Record.ActionTime = strfmt.DateTime(time.Now().UTC())
	return b
}
func (b *Form3EventBuilder) Created() *Form3EventBuilder {
	return b.EventType("created").Description("Record inserted")
}
func (b *Form3EventBuilder) Updated() *Form3EventBuilder {
	return b.EventType("updated").Description("Record updated")
}
func (b *Form3EventBuilder) RecordId(id uuid.UUID) *Form3EventBuilder {
	b.event.Record.RecordId = id
	return b
}
func (b *Form3EventBuilder) OrganisationId(organisationId uuid.UUID) *Form3EventBuilder {
	b.event.OrganisationId = organisationId
	return b
}
func (b *Form3EventBuilder) Description(description string) *Form3EventBuilder {
	b.event.Record.Description = description
	return b
}
func (b *Form3EventBuilder) EventType(eventType string) *Form3EventBuilder {
	b.event.EventType = eventType
	return b
}
func (b *Form3EventBuilder) Version(version int64) *Form3EventBuilder {
	b.event.Version = version
	return b
}
func (b *Form3EventBuilder) RecordType(recordType string) *Form3EventBuilder {
	b.event.RecordType = recordType
	b.event.DataRecordType = recordType
	b.event.Record.RecordType = recordType
	return b
}
func (b *Form3EventBuilder) AfterData(data interface{}) *Form3EventBuilder {
	if data != nil {
		payload, _ := json.Marshal(data)
		b.event.Record.AfterData = payload
		b.event.Data = payload
	}
	return b
}
func (b *Form3EventBuilder) BeforeData(data interface{}) *Form3EventBuilder {
	if data != nil {
		payload, _ := json.Marshal(data)
		b.event.Record.BeforeData = payload
	}
	return b
}
func (b *Form3EventBuilder) Build(ctx *context.Context) *Form3Event {
	if security.IsApplicationContext(*ctx) {
		b.event.Record.ActionedBy = uuid.MustParse(settings.UserID)
	} else {
		b.event.Record.ActionedBy = uuid.MustParse((*ctx).Value("user_id").(string))
	}
	return &b.event
}
