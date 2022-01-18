package models

type StatusType string

const (
	StatusTypeOk      StatusType = "ok"
	StatusTypeError   StatusType = "error"
	StatusTypeWarning StatusType = "warning"
)
