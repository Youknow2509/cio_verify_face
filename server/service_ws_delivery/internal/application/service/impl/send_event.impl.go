package impl

import (
	"context"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/mq"
)

/**
 * Send event service struct
 */
type SendEventService struct{}

// CallAnswer implements service.ISendEventService.
func (s *SendEventService) CallAnswer(ctx context.Context, input model.CallAnswer) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	var callType domainModel.CallType
	switch input.CallType {
	case int(domainModel.CallVideo):
		callType = domainModel.CallVideo
	case int(domainModel.CallAudio):
		callType = domainModel.CallAudio
	default:
		callType = domainModel.CallAudio
	}
	if err := sendEventService.UserCallAnswer(
		ctx,
		domainModel.UserCallAnswer{
			ReceiverId:   input.ReceiverId,
			CallType:     callType,
			SdpAnswer:    input.SdpAnswer,
			CallAnswerId: input.CallerId,
			SenderId:     input.UserId,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

// CallEnd implements service.ISendEventService.
func (s *SendEventService) CallEnd(ctx context.Context, input model.CallEnd) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	var callEndType domainModel.CallEndReason
	switch input.Reason {
	case int(domainModel.CallEndReasonConnectionFailed):
		callEndType = domainModel.CallEndReasonConnectionFailed
	case int(domainModel.CallEndReasonDeclined):
		callEndType = domainModel.CallEndReasonDeclined
	case int(domainModel.CallEndReasonHangUp):
		callEndType = domainModel.CallEndReasonHangUp
	case int(domainModel.CallEndReasonMissed):
		callEndType = domainModel.CallEndReasonMissed
	default:
		callEndType = domainModel.CallEndReasonHangUp
	}
	if err := sendEventService.UserCallEnd(
		ctx,
		domainModel.UserCallEnd{
			CallId:   input.CallId,
			Reason:   callEndType,
			Duration: input.Duration,
			// Info user
			UserId:       input.UserId,
			SessionId:    input.SessionId,
			ConnectionId: input.ConnectionId,
			IpAddress:    input.IpAddress,
			UserAgent:    input.UserAgent,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

// CallIceCandidate implements service.ISendEventService.
func (s *SendEventService) CallIceCandidate(ctx context.Context, input model.CallIceCandidate) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	if err := sendEventService.UserCallIceCandidate(
		ctx,
		domainModel.UserCallIceCandidate{
			ReceiverId:    input.ReceiverId,
			Candidate:     input.Candidate,
			SdpMid:        input.SdpMid,
			SdpMLineIndex: input.SdpMLineIndex,
			CallId:        input.CallId,
			// Info user
			UserId:       input.UserId,
			SessionId:    input.SessionId,
			ConnectionId: input.ConnectionId,
			IpAddress:    input.IpAddress,
			UserAgent:    input.UserAgent,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

// CallOffer implements service.ISendEventService.
func (s *SendEventService) CallOffer(ctx context.Context, input model.CallOffer) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	var callType domainModel.CallType
	switch input.CallType {
	case int(domainModel.CallVideo):
		callType = domainModel.CallVideo
	case int(domainModel.CallAudio):
		callType = domainModel.CallAudio
	default:
		callType = domainModel.CallAudio
	}
	if err := sendEventService.UserCallOfferInitilize(
		ctx,
		domainModel.UserCallOfferInitilize{
			SenderId:   input.SenderId,
			ReceiverId: input.ReceiverId,
			CallType:   callType,
			SdpOffer:   input.SdpOffer,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

// DeleteMessage implements service.ISendEventService.
func (s *SendEventService) DeleteMessage(ctx context.Context, input model.DeleteMessage) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	// Get delete type
	var deleteType domainModel.DeleteMessageType
	switch input.Type {
	case int(domainModel.DeleteForMe):
		deleteType = domainModel.DeleteForMe
	case int(domainModel.DeleteForEveryone):
		deleteType = domainModel.DeleteForEveryone
	default:
		deleteType = domainModel.DeleteForMe
	}
	// Handle delete type
	if err := sendEventService.UserDeleteMessage(
		ctx,
		domainModel.UserDeleteMessage{
			Type:           deleteType,
			UserId:         input.UserId,
			MessageId:      input.MessageId,
			ConversationId: input.ConversationId,
			Timestamp:      input.Timestamp,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

// EditMessage implements service.ISendEventService.
func (s *SendEventService) EditMessage(ctx context.Context, input model.EditMessage) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	if err := sendEventService.UserEditMessage(
		ctx,
		domainModel.UserEditMessage{
			UserId:     input.UserId,
			MessageId:  input.MessageId,
			NewContent: input.NewContent,
			Timestamp:  input.Timestamp,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

// MessageReadStatus implements service.ISendEventService.
func (s *SendEventService) MessageReadStatus(ctx context.Context, input model.MessageReadStatus) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	if err := sendEventService.UserReadMessageStatus(
		ctx,
		domainModel.UserReadMessageStatus{
			UserId:         input.UserId,
			ConversationId: input.ConversationId,
			MessageId:      input.MessageId,
			Timestamp:      input.Timestamp,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

// NewMessageText implements service.ISendEventService.
func (s *SendEventService) NewMessageText(ctx context.Context, input model.NewMessageText) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	if err := sendEventService.WriteNewMessage(
		ctx,
		domainModel.WriteNewMessage{
			SenderId:       input.SenderId,
			ConversationId: input.ConversationId,
			TempId:         input.TempId,
			Message:        input.Message,
			ReplyToId:      input.ReplyToId,
			Timestamp:      input.Timestamp,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

// ReactMessage implements service.ISendEventService.
func (s *SendEventService) ReactMessage(ctx context.Context, input model.ReactMessage) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	var reactType domainModel.ReactMessageType
	switch input.Reaction {
	case int(domainModel.ReactLike):
		reactType = domainModel.ReactLike
	case int(domainModel.ReactAngry):
		reactType = domainModel.ReactAngry
	case int(domainModel.ReactDislike):
		reactType = domainModel.ReactDislike
	case int(domainModel.ReactLaugh):
		reactType = domainModel.ReactLaugh
	case int(domainModel.ReactLove):
		reactType = domainModel.ReactLove
	default:
		reactType = domainModel.ReactLike
	}
	if err := sendEventService.UserReactMessage(
		ctx,
		domainModel.UserReactMessage{
			Status:         input.Status,
			ConversationId: input.ConversationId,
			MessageId:      input.MessageId,
			Reaction:       reactType,
			Timestamp:      input.Timestamp,
			UserId:         input.UserId,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

// UserTypingStatus implements service.ISendEventService.
func (s *SendEventService) UserTypingStatus(ctx context.Context, input model.UserTypingStatus) error {
	// Get instance use
	sendEventService := domainMq.GetSendEventToKafka()
	if err := sendEventService.UpgradeStatusTypingUser(
		ctx,
		domainModel.UpgradeStatusTypingUser{
			UserId:         input.UserId,
			ConversationId: input.ConnectionId,
			Status:         input.Status,
			Timestamp:      input.Timestamp,
		},
	); err != nil {
		// Send message to kafka
		return err
	}
	return nil
}

/**
 * New send event service and implementation
 */
func NewSendEventService() service.ISendEventService {
	return &SendEventService{}
}
