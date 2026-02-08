package domain

import "fmt"

type EventType string

const (
	// Evidence events
	EvidenceUploaded      EventType = "EVIDENCE_UPLOADED"
	EvidenceViewed        EventType = "EVIDENCE_VIEWED"
	EvidenceDownloaded    EventType = "EVIDENCE_DOWNLOADED"
	EvidenceShared        EventType = "EVIDENCE_SHARED"
	EvidenceDeleted       EventType = "EVIDENCE_DELETED"
	EvidenceEdited        EventType = "EVIDENCE_EDITED"
	EvidenceRedacted      EventType = "EVIDENCE_REDACTED"
	EvidenceTagged        EventType = "EVIDENCE_TAGGED"
	EvidenceCaseLinked    EventType = "EVIDENCE_CASE_LINKED"
	EvidenceCaseUnlinked  EventType = "EVIDENCE_CASE_UNLINKED"
	EvidenceRetentionSet  EventType = "EVIDENCE_RETENTION_SET"
	EvidencePermissionSet EventType = "EVIDENCE_PERMISSION_SET"

	// User events
	UserSignedIn          EventType = "USER_SIGNED_IN"
	UserSignedOut         EventType = "USER_SIGNED_OUT"
	UserPasswordChanged   EventType = "USER_PASSWORD_CHANGED"
	UserRoleChanged       EventType = "USER_ROLE_CHANGED"
	UserPermissionChanged EventType = "USER_PERMISSION_CHANGED"

	// Device events
	DeviceCameraStarted   EventType = "DEVICE_CAMERA_STARTED"
	DeviceCameraStopped   EventType = "DEVICE_CAMERA_STOPPED"
	DeviceDocked          EventType = "DEVICE_DOCKED"
	DeviceAssigned        EventType = "DEVICE_ASSIGNED"
	DeviceFirmwareUpdated EventType = "DEVICE_FIRMWARE_UPDATED"
)

var validEventTypes = map[EventType]struct{}{
	EvidenceUploaded:      {},
	EvidenceViewed:        {},
	EvidenceDownloaded:    {},
	EvidenceShared:        {},
	EvidenceDeleted:       {},
	EvidenceEdited:        {},
	EvidenceRedacted:      {},
	EvidenceTagged:        {},
	EvidenceCaseLinked:    {},
	EvidenceCaseUnlinked:  {},
	EvidenceRetentionSet:  {},
	EvidencePermissionSet: {},
	UserSignedIn:          {},
	UserSignedOut:         {},
	UserPasswordChanged:   {},
	UserRoleChanged:       {},
	UserPermissionChanged: {},
	DeviceCameraStarted:   {},
	DeviceCameraStopped:   {},
	DeviceDocked:          {},
	DeviceAssigned:        {},
	DeviceFirmwareUpdated: {},
}

func (t EventType) IsValid() bool {
	_, ok := validEventTypes[t]
	return ok
}

func (t EventType) Validate() error {
	if !t.IsValid() {
		return fmt.Errorf("invalid event type: %q", t)
	}
	return nil
}

func (t EventType) Category() string {
	switch {
	case len(t) > 8 && t[:8] == "EVIDENCE":
		return "evidence"
	case len(t) > 4 && t[:4] == "USER":
		return "user"
	case len(t) > 6 && t[:6] == "DEVICE":
		return "device"
	default:
		return "unknown"
	}
}
