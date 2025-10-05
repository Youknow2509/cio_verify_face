package mq

import (
	"context"
	"encoding/json"

	libConstants "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/mq"
	libsDomainMq "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/mq"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
)

/**
 *
 */
type SendEventToKafka struct{}

// UpgradeStatusTypingUser implements mq.ISendEventToKafka.
func (s *SendEventToKafka) UpgradeStatusTypingUser(ctx context.Context, input model.UpgradeStatusTypingUser) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessageRequireAck(
		ctx,
		libConstants.KAFKA_TOPIC_USER_TYPING_EVENTS,
		"",
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message to kafka")
		return err
	}

	return nil
}

// UserCallAnswer implements mq.ISendEventToKafka.
func (s *SendEventToKafka) UserCallAnswer(ctx context.Context, input model.UserCallAnswer) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataObj := make(map[string]interface{})
	// Copy fields from input to dataObj
	inputBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal input")
	}
	err = json.Unmarshal(inputBytes, &dataObj)
	if err != nil {
		global.Logger.Error("Error when unmarshal input to map")
	}
	dataObj["type"] = "call_answer" // TODO: define type with int
	dataBytes, err := json.Marshal(dataObj)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessageRequireAllAck(
		ctx,
		libConstants.KAFKA_TOPIC_USER_CALL_EVENTS,
		"",
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message to kafka")
		return err
	}

	return nil
}

// UserCallEnd implements mq.ISendEventToKafka.
func (s *SendEventToKafka) UserCallEnd(ctx context.Context, input model.UserCallEnd) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataObj := make(map[string]interface{})
	// Copy fields from input to dataObj
	inputBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal input")
	}
	err = json.Unmarshal(inputBytes, &dataObj)
	if err != nil {
		global.Logger.Error("Error when unmarshal input to map")
	}
	dataObj["type"] = "call_end" // TODO: define type with int
	dataBytes, err := json.Marshal(dataObj)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessageRequireAllAck(
		ctx,
		libConstants.KAFKA_TOPIC_USER_CALL_EVENTS,
		"",
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message to kafka")
		return err
	}

	return nil
}

// UserCallIceCandidate implements mq.ISendEventToKafka.
func (s *SendEventToKafka) UserCallIceCandidate(ctx context.Context, input model.UserCallIceCandidate) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataObj := make(map[string]interface{})
	// Copy fields from input to dataObj
	inputBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal input")
	}
	err = json.Unmarshal(inputBytes, &dataObj)
	if err != nil {
		global.Logger.Error("Error when unmarshal input to map")
	}
	dataObj["type"] = "call_ice_candidate" // TODO: define type with int
	dataBytes, err := json.Marshal(dataObj)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessageRequireAllAck(
		ctx,
		libConstants.KAFKA_TOPIC_USER_CALL_EVENTS,
		"",
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message to kafka")
		return err
	}

	return nil
}

// UserCallOfferInitilize implements mq.ISendEventToKafka.
func (s *SendEventToKafka) UserCallOfferInitilize(ctx context.Context, input model.UserCallOfferInitilize) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataObj := make(map[string]interface{})
	// Copy fields from input to dataObj
	inputBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal input")
	}
	err = json.Unmarshal(inputBytes, &dataObj)
	if err != nil {
		global.Logger.Error("Error when unmarshal input to map")
	}
	dataObj["type"] = "call_offer_initialize" // TODO: define type with int
	dataBytes, err := json.Marshal(dataObj)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessageRequireAllAck(
		ctx,
		libConstants.KAFKA_TOPIC_USER_CALL_EVENTS,
		"",
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message to kafka")
		return err
	}

	return nil
}

// UserDeleteMessage implements mq.ISendEventToKafka.
func (s *SendEventToKafka) UserDeleteMessage(ctx context.Context, input model.UserDeleteMessage) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessageRequireAllAck(
		ctx,
		libConstants.KAFKA_TOPIC_MESSAGE_DELETE,
		"",
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message to kafka")
		return err
	}

	return nil
}

// UserEditMessage implements mq.ISendEventToKafka.
func (s *SendEventToKafka) UserEditMessage(ctx context.Context, input model.UserEditMessage) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessageRequireAllAck(
		ctx,
		libConstants.KAFKA_TOPIC_MESSAGE_EDIT,
		"",
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message to kafka")
		return err
	}

	return nil
}

// UserReactMessage implements mq.ISendEventToKafka.
func (s *SendEventToKafka) UserReactMessage(ctx context.Context, input model.UserReactMessage) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessageRequireAck(
		ctx,
		libConstants.KAFKA_TOPIC_MESSAGE_REACT,
		"",
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message react to kafka")
		return err
	}

	return nil
}

// UserReadMessageStatus implements mq.ISendEventToKafka.
func (s *SendEventToKafka) UserReadMessageStatus(ctx context.Context, input model.UserReadMessageStatus) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessage(
		ctx,
		libConstants.KAFKA_TOPIC_MESSAGE_READ_STATUS,
		"",
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message read status to kafka")
		return err
	}

	return nil
}

// WriteNewMessage implements mq.ISendEventToKafka.
func (s *SendEventToKafka) WriteNewMessage(ctx context.Context, input model.WriteNewMessage) error {
	kaffaWriter, err := libsDomainMq.GetKafkaWriteService()
	if err != nil {
		global.Logger.Error("Error when get instance kafka writer")
		return err
	}
	dataBytes, err := json.Marshal(input)
	if err != nil {
		global.Logger.Error("Error when marshal data")
	}
	if err := kaffaWriter.WriteMessageRequireAllAck(
		ctx,
		libConstants.KAFKA_TOPIC_NEW_MESSAGE,
		input.ConversationId.String(),
		dataBytes,
	); err != nil {
		global.Logger.Error("Error when write message to kafka")
		return err
	}

	return nil
}

/**
 * New struct and implementation for sending events to Kafka
 */
func NewSendEventToKafka() domainMq.ISendEventToKafka {
	return &SendEventToKafka{}
}
