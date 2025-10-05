package model

import "github.com/google/uuid"

// ======================================
// Model for user presence
// ======================================

// For client send to server
type (
	// For initiating a offer all
	CallOffer struct {
		SenderId   uuid.UUID `json:"sender_id"`
		ReceiverId uuid.UUID `json:"receiver_id"`
		CallType   CallType  `json:"call_type"`
		SdpOffer   string    `json:"sdp_offer"`
	}

	// For call answer
	CallAnswer struct {
		SenderId   uuid.UUID `json:"sender_id"`
		ReceiverId uuid.UUID `json:"receiver_id"`
		CallType   CallType  `json:"call_type"`
		SdpAnswer  string    `json:"sdp_answer"`
	}

	// Call ICE Candidate
	CallIceCandidate struct {
		SenderId   uuid.UUID `json:"sender_id"`
		ReceiverId uuid.UUID `json:"receiver_id"`
		Candidate  string    `json:"candidate"`
	}

	// Call end
	CallEnd struct {
		SenderId   uuid.UUID     `json:"sender_id"`
		ReceiverId uuid.UUID     `json:"receiver_id"`
		Reason     CallEndReason `json:"reason"` // 0 - hang_up | 1 - hang_up | 2 - declined | 3- missed | 4 - connection_failed
	}
)

// For server send to clients
type (
	// Call end
	CallEndSend struct {
		SenderId   uuid.UUID     `json:"sender_id"`
		ReceiverId uuid.UUID     `json:"receiver_id"`
		Reason     CallEndReason `json:"reason"` // 0 - hang_up | 1 - hang_up | 2 - declined | 3- missed | 4 - connection_failed
	}

	// Call ICE Candidate
	CallIceCandidateSend struct {
		SenderId   uuid.UUID `json:"sender_id"`
		ReceiverId uuid.UUID `json:"receiver_id"`
		Candidate  string    `json:"candidate"`
	}

	// For initiating a offer call
	CallOfferSend struct {
		SenderId   uuid.UUID `json:"sender_id"`
		ReceiverId uuid.UUID `json:"receiver_id"`
		CallType   CallType  `json:"call_type"`
		SdpOffer   string    `json:"sdp_offer"`
	}

	// For call answer
	CallAnswerSend struct {
		SenderId   uuid.UUID `json:"sender_id"`
		ReceiverId uuid.UUID `json:"receiver_id"`
		CallType   CallType  `json:"call_type"`
		SdpAnswer  string    `json:"sdp_answer"`
	}
)
