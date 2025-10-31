package model

// ================================================
// Verdict enum
// ================================================
type Verdict int
const (
	Allowed Verdict = iota
	Denied
)