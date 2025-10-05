package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	libsUtilsUuid "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/shared/utils/uuid"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/dto"
	model "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/model"
)

/**
 * Interface worker handle
 */
type WorkerHandler struct {
}

/**
 * Get WorkerHandler instance
 */

func GetWorkerHandler() *WorkerHandler {
	return &WorkerHandler{}
}

/**
 * Register new client
 */
func (wh *WorkerHandler) RegisterClient(ctx context.Context, input model.ClientInfo) error {
	service := applicationService.GetMapConnectionService()
	if err := service.RegisterConnection(
		ctx,
		&applicationModel.RegisterConnection{
			ConnectionId: input.ConnectionId.String(),
			UserId:       input.UserId.String(),
			IPAddress:    input.IpAddress,
			ConnectedAt:  time.Now().Format("2006-01-02 15:04:05"),
			UserAgent:    input.UserAgent,
		},
	); err != nil {
		return fmt.Errorf("failed to register connection: %v", err)
	}
	return nil
}

/**
 * Unregister client
 */
func (wh *WorkerHandler) UnregisterClient(ctx context.Context, input model.ClientInfo) error {
	service := applicationService.GetMapConnectionService()
	if err := service.UnregisterConnection(
		ctx,
		&applicationModel.UnregisterConnection{
			ConnectionId: input.ConnectionId.String(),
			UserId:       input.UserId.String(),
		},
	); err != nil {
		return fmt.Errorf("failed to unregister connection: %v", err)
	}
	return nil
}

/**
 * Handle data receive
 */
func (wh *WorkerHandler) HandleDataReceive(ctx context.Context, clientInfo model.ClientInfo, eventType int, data []byte) error {
	switch eventType {
	case int(domainModel.WSEventMessageText):
		var messageText dto.NewMessageText
		if err := json.Unmarshal(data, &messageText); err != nil {
			return err
		}
		// convert data to uuid
		conversationId, err := uuid.Parse(messageText.ConversationId)
		if err != nil {
			return fmt.Errorf("failed to parse conversation ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		var replyToId uuid.UUID
		if messageText.ReplyToId != "" {
			replyToId, err = uuid.Parse(messageText.ReplyToId)
			if err != nil {
				return fmt.Errorf("failed to parse reply to ID: %v for (u::%s | s:: %s | c::%s)",
					err,
					clientInfo.UserId,
					clientInfo.SessionId,
					clientInfo.ConnectionId,
				)
			}
		}
		if err := applicationService.GetSendEventService().NewMessageText(
			ctx,
			applicationModel.NewMessageText{
				SenderId:       clientInfo.ConnectionId,
				ConversationId: conversationId,
				TempId:         messageText.TempId,
				Message:        messageText.Message,
				ReplyToId:      replyToId,
				Timestamp:      time.Now().Unix(),
				SessionId:      clientInfo.SessionId,
				ConnectionId:   clientInfo.ConnectionId,
				IpAddress:      clientInfo.IpAddress,
				UserAgent:      clientInfo.UserAgent,
			},
		); err != nil {
			return fmt.Errorf("failed to send new message text: %v", err)
		}
	case int(domainModel.WSEventEditMessage):
		var messageEdit dto.EditMessage
		if err := json.Unmarshal(data, &messageEdit); err != nil {
			return fmt.Errorf("failed to unmarshal message edit: %v", err)
		}
		// convert uuid
		conversationId, err := uuid.Parse(messageEdit.ConversationId)
		if err != nil {
			return fmt.Errorf("failed to parse conversation ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		messageId, err := libsUtilsUuid.ParseTimeUUID(messageEdit.MessageId)
		if err != nil {
			return fmt.Errorf("failed to parse message ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		// send edit message
		if err := applicationService.GetSendEventService().EditMessage(
			ctx,
			applicationModel.EditMessage{
				ConversationId: conversationId,
				MessageId:      messageId,
				NewContent:     messageEdit.NewContent,
				Timestamp:      time.Now().Unix(),
				UserId:         clientInfo.UserId,
				SessionId:      clientInfo.SessionId,
				ConnectionId:   clientInfo.ConnectionId,
				IpAddress:      clientInfo.IpAddress,
				UserAgent:      clientInfo.UserAgent,
			},
		); err != nil {
			return fmt.Errorf("failed to send edit message: %v", err)
		}
	case int(domainModel.WSEventDeleteMessage):
		var messageDelete dto.DeleteMessage
		if err := json.Unmarshal(data, &messageDelete); err != nil {
			return fmt.Errorf("failed to unmarshal message delete: %v", err)
		}
		// convert to uuid
		conversationId, err := uuid.Parse(messageDelete.ConversationId)
		if err != nil {
			return fmt.Errorf("failed to parse conversation ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		messageId, err := libsUtilsUuid.ParseTimeUUID(messageDelete.MessageId)
		if err != nil {
			return fmt.Errorf("failed to parse message ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		// send delete message
		if err := applicationService.GetSendEventService().DeleteMessage(
			ctx,
			applicationModel.DeleteMessage{
				Type:           messageDelete.Type,
				ConversationId: conversationId,
				MessageId:      messageId,
				Timestamp:      time.Now().Unix(),
				// Info user
				UserId:       clientInfo.UserId,
				SessionId:    clientInfo.SessionId,
				ConnectionId: clientInfo.ConnectionId,
				IpAddress:    clientInfo.IpAddress,
				UserAgent:    clientInfo.UserAgent,
			},
		); err != nil {
			return fmt.Errorf("failed to send delete message: %v", err)
		}
	case int(domainModel.WSEventReactMessage):
		var messageReact dto.ReactMessage
		if err := json.Unmarshal(data, &messageReact); err != nil {
			return fmt.Errorf("failed to unmarshal message react: %v", err)
		}
		// convert to uuid
		conversationId, err := uuid.Parse(messageReact.ConversationId)
		if err != nil {
			return fmt.Errorf("failed to parse conversation ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		messageId, err := uuid.Parse(messageReact.MessageId)
		if err != nil {
			return fmt.Errorf("failed to parse message ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		if err := applicationService.GetSendEventService().ReactMessage(
			ctx,
			applicationModel.ReactMessage{
				Status:         messageReact.Status,
				ConversationId: conversationId,
				MessageId:      messageId,
				Reaction:       messageReact.Reaction,
				Timestamp:      time.Now().Unix(),
				UserId:         clientInfo.UserId,
				SessionId:      clientInfo.SessionId,
				ConnectionId:   clientInfo.ConnectionId,
				IpAddress:      clientInfo.IpAddress,
				UserAgent:      clientInfo.UserAgent,
			},
		); err != nil {
			return fmt.Errorf("failed to send react message: %v", err)
		}
	case int(domainModel.WSEventReadReceipt):
		var readReceipt dto.MessageReadStatus
		if err := json.Unmarshal(data, &readReceipt); err != nil {
			return fmt.Errorf("failed to unmarshal read receipt: %v", err)
		}
		// validate uuid
		msgUuid, err := libsUtilsUuid.ParseTimeUUID(readReceipt.MessageId)
		if err != nil {
			return fmt.Errorf("failed to parse message ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		conversationId, err := uuid.Parse(readReceipt.ConversationId)
		if err != nil {
			return fmt.Errorf("failed to parse conversation ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		if err := applicationService.GetSendEventService().MessageReadStatus(
			ctx,
			applicationModel.MessageReadStatus{
				ConversationId: conversationId,
				MessageId:      msgUuid,
				Timestamp:      time.Now().Unix(),
				UserId:         clientInfo.UserId,
				SessionId:      clientInfo.SessionId,
				ConnectionId:   clientInfo.ConnectionId,
				IpAddress:      clientInfo.IpAddress,
				UserAgent:      clientInfo.UserAgent,
			},
		); err != nil {
			return fmt.Errorf("failed to send read receipt: %v", err)
		}
	case int(domainModel.WSEventTypingStatus):
		var typingStatus dto.UserTypingStatus
		if err := json.Unmarshal(data, &typingStatus); err != nil {
			return fmt.Errorf("failed to unmarshal typing status: %v", err)
		}
		// convert uuid
		conversationId, err := uuid.Parse(typingStatus.ConversationId)
		if err != nil {
			return fmt.Errorf("failed to parse conversation ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		// send typing status
		if err := applicationService.GetSendEventService().UserTypingStatus(
			ctx,
			applicationModel.UserTypingStatus{
				ConversationId: conversationId,
				Status:         typingStatus.IsTyping,
				UserId:         clientInfo.UserId,
				Timestamp:      time.Now().Unix(),
				SessionId:      clientInfo.SessionId,
				ConnectionId:   clientInfo.ConnectionId,
				IpAddress:      clientInfo.IpAddress,
				UserAgent:      clientInfo.UserAgent,
			},
		); err != nil {
			return fmt.Errorf("failed to send typing status: %v", err)
		}
	case int(domainModel.WSEventCallOffer):
		var callOffer dto.CallOffer
		if err := json.Unmarshal(data, &callOffer); err != nil {
			return fmt.Errorf("failed to unmarshal call offer: %v", err)
		}
		// convert uuid
		receiverId, err := uuid.Parse(callOffer.ReceiverId)
		if err != nil {
			return fmt.Errorf("failed to parse receiver ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		if err := applicationService.GetSendEventService().CallOffer(
			ctx,
			applicationModel.CallOffer{
				ReceiverId: receiverId,
				CallType:   callOffer.CallType,
				SdpOffer:   callOffer.SdpOffer,
				// Info user
				SessionId:    clientInfo.SessionId,
				ConnectionId: clientInfo.ConnectionId,
				IpAddress:    clientInfo.IpAddress,
				UserAgent:    clientInfo.UserAgent,
			},
		); err != nil {
			return fmt.Errorf("failed to send call offer: %v", err)
		}
	case int(domainModel.WSEventCallAnswer):
		var callAnswer dto.CallAnswer
		if err := json.Unmarshal(data, &callAnswer); err != nil {
			return fmt.Errorf("failed to unmarshal call answer: %v", err)
		}
		// convert uuid
		receiverId, err := uuid.Parse(callAnswer.ReceiverId)
		if err != nil {
			return fmt.Errorf("failed to parse receiver ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		if err := applicationService.GetSendEventService().CallAnswer(
			ctx,
			applicationModel.CallAnswer{
				CallerId:   callAnswer.CallerId,
				ReceiverId: receiverId,
				CallType:   callAnswer.CallType,
				SdpAnswer:  callAnswer.SdpAnswer,
				// Info user
				SessionId:    clientInfo.SessionId,
				ConnectionId: clientInfo.ConnectionId,
				IpAddress:    clientInfo.IpAddress,
				UserAgent:    clientInfo.UserAgent,
			},
		); err != nil {
			return fmt.Errorf("failed to send call answer: %v", err)
		}
	case int(domainModel.WSEventCallIceCandidate):
		var callIceCandidate dto.CallIceCandidate
		if err := json.Unmarshal(data, &callIceCandidate); err != nil {
			return err
		}
		// convert uuid
		receiverId, err := uuid.Parse(callIceCandidate.ReceiverId)
		if err != nil {
			return fmt.Errorf("failed to parse receiver ID: %v for (u::%s | s:: %s | c::%s)",
				err,
				clientInfo.UserId,
				clientInfo.SessionId,
				clientInfo.ConnectionId,
			)
		}
		if err := applicationService.GetSendEventService().CallIceCandidate(
			ctx,
			applicationModel.CallIceCandidate{
				Candidate:     callIceCandidate.Candidate,
				SdpMid:        callIceCandidate.SdpMid,
				SdpMLineIndex: callIceCandidate.SdpMLineIndex,
				ReceiverId:    receiverId,
				CallId:        callIceCandidate.CallId,
				// Info user
				UserId:       clientInfo.UserId,
				SessionId:    clientInfo.SessionId,
				ConnectionId: clientInfo.ConnectionId,
				IpAddress:    clientInfo.IpAddress,
				UserAgent:    clientInfo.UserAgent,
			},
		); err != nil {
			return err
		}
	case int(domainModel.WSEventCallEnd):
		var callEnd dto.CallEnd
		if err := json.Unmarshal(data, &callEnd); err != nil {
			return err
		}
		if err := applicationService.GetSendEventService().CallEnd(
			ctx,
			applicationModel.CallEnd{
				CallId:   callEnd.CallId,
				Reason:   callEnd.Reason,
				Duration: callEnd.Duration,
				// Info user
				UserId:       clientInfo.UserId,
				SessionId:    clientInfo.SessionId,
				ConnectionId: clientInfo.ConnectionId,
				IpAddress:    clientInfo.IpAddress,
				UserAgent:    clientInfo.UserAgent,
			},
		); err != nil {
			return fmt.Errorf("failed to send call end: %v", err)
		}
	default:
		errorBytes, _ := json.Marshal(fmt.Sprintf("Unknown event type: %d", eventType))
		err := wh.SendDataToClient(ctx, clientInfo.ConnectionId, errorBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * Send data to client
 */
func (wh *WorkerHandler) SendDataToClient(ctx context.Context, clientId uuid.UUID, data []byte) error {
	err := applicationService.GetClientService().SendMessageToClient(
		ctx,
		&applicationModel.SendMessageToClientInput{
			ConnectionId: clientId,
			Message:      data,
		},
	)
	return err
}
