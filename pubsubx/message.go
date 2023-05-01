// Package pubsubx provides common utilities and formats for working with pubsub systems
package pubsubx

import (
	"time"

	"go.infratographer.com/x/gidx"
)

// Message contains the data structure expected to be received when picking
// an event from a message queue
//
// Deprecated: Message exists for backwards compatibility and should be migrated to use
// ChangeMessage or EventMessage depending on the message type
type Message struct {
	// SubjectURN is a string representing the identity of the topic of this message
	SubjectURN string `json:"subject_urn"`
	// EventType describes the type of event that has triggered this message
	EventType string `json:"event_type"`
	// AdditionalSubjectURNs is a group of strings representing additional identities associated with this message
	AdditionalSubjectURNs []string `json:"additional_subjects"`
	// ActorURN is a string representing the identity of the actor that created this message
	ActorURN string `json:"actor_urn"`
	// Source is a string representing the identity of the source system that created the message
	Source string `json:"source"`
	// Timestamp is the time representing when the message was created
	Timestamp time.Time `json:"timestamp"`
	// SubjectFields is a map of additional descriptors for this message
	SubjectFields map[string]string `json:"fields"`
	// AdditionalData is a field to store any addition information that may be important to include with your message
	AdditionalData map[string]interface{} `json:"additional_data"`
}

// ChangeMessage contains the data structure expected to be received when picking
// an event from a changes message queue
type ChangeMessage struct {
	// SubjectID is the PrefixedID representing the node of the topic of this message
	SubjectID gidx.PrefixedID `json:"subjectID"`
	// EventType describes the type of event that has triggered this message
	EventType string `json:"eventType"`
	// AdditionalSubjectURNs is a group of PrefixedIDs representing additional nodes associated with this message
	AdditionalSubjectIDs []gidx.PrefixedID `json:"additionalSubjects"`
	// ActorID is the PrefixedID representing the identity of the actor that caused this message to be triggered
	ActorID gidx.PrefixedID `json:"actorID"`
	// Source is a string representing the identity of the source system that created the message
	Source string `json:"source"`
	// Timestamp is the time representing when the message was created
	Timestamp time.Time `json:"timestamp"`
	// TraceID is the ID of the trace for this event
	TraceID string `json:"traceID"`
	// SpanID is the ID of the span that additional traces should based off of
	SpanID string `json:"spanID"`
	// SubjectFields is a map of the fields on the subject
	SubjectFields map[string]string `json:"subjectFields"`
	// Changedset is an optional map of the fields that changed triggering this message, should be provided if the source can provide
	Changedset map[string]string `json:"changeset"`
	// AdditionalData is a field to store any addition information that may be important to include with your message
	AdditionalData map[string]interface{} `json:"additionalData"`
}

// EventMessage contains the data structure expected to be received when picking
// an event from an events message queue
type EventMessage struct {
	// SubjectID is the PrefixedID representing the node of the topic of this message
	SubjectID gidx.PrefixedID `json:"subject_id"`
	// EventType describes the type of event that has triggered this message
	EventType string `json:"event_type"`
	// AdditionalSubjectURNs is a group of PrefixedIDs representing additional nodes associated with this message
	AdditionalSubjectIDs []gidx.PrefixedID `json:"additionalSubjects"`
	// Source is a string representing the identity of the source system that created the message
	Source string `json:"source"`
	// Timestamp is the time representing when the message was created
	Timestamp time.Time `json:"timestamp"`
	// TraceID is the ID of the trace for this event
	TraceID string `json:"traceID"`
	// SpanID is the ID of the span that additional traces should based off of
	SpanID string `json:"spanID"`
	// data is a field to store any information that may be important to include about the event
	Data map[string]interface{} `json:"data"`
}
