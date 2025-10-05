package model

// ======================================
// Enums use in model
// ======================================

// Ws event types
type WSEventType int

const (
	WSEventMessageText WSEventType = iota
	WSEventMessageAck
	WSEventMessageImage
	WSEventMessageVideo
	WSEventMessageFile
	WSEventMessageAudio
	WSEventTypingStatus
	WSEventReadReceipt
	WSEventEditMessage
	WSEventDeleteMessage
	WSEventReactMessage
	WSEventCreateGroup
	WSEventAddMembersToGroup
	WSEventLeaveGroup
	WSEventUserPresence
	WSEventCallOffer
	WSEventCallAnswer
	WSEventCallIceCandidate
	WSEventCallEnd
)

// DeleteMessageType for delete message action
type DeleteMessageType int

const (
	DeleteForMe DeleteMessageType = iota
	DeleteForEveryone
)

// React message type
type ReactMessageType int

const (
	ReactLove ReactMessageType = iota
	ReactLike
	ReactDislike
	ReactLaugh
	ReactAngry
	// ...
)

// Call type
type CallType int
const (
	CallVideo CallType = iota
	CallAudio
	CallScreenShare
)

// CallEndReason
type CallEndReason int
const (
	CallEndReasonHangUp CallEndReason = iota
	CallEndReasonDeclined
	CallEndReasonMissed
	CallEndReasonConnectionFailed
)
