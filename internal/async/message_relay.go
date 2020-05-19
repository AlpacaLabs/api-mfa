package async

import (
	"context"
	"errors"
	"fmt"
	"time"

	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"

	hermesTopics "github.com/AlpacaLabs/api-hermes/pkg/topic"

	"github.com/AlpacaLabs/api-mfa/internal/configuration"
	"github.com/AlpacaLabs/api-mfa/internal/db"
	goKafka "github.com/AlpacaLabs/go-kafka"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

func RelayMessagesForSend(config configuration.Config, dbClient db.Client) {
	brokers := []string{
		fmt.Sprintf("%s:%d", config.KafkaConfig.Host, config.KafkaConfig.Port),
	}

	writerForEmail := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   hermesTopics.TopicForSendEmailRequest,
	})
	defer writerForEmail.Close()

	writerForSms := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   hermesTopics.TopicForSendSmsRequest,
	})
	defer writerForSms.Close()

	writers := make(map[string]*kafka.Writer)
	writers[hermesTopics.TopicForSendEmailRequest] = writerForEmail
	writers[hermesTopics.TopicForSendSmsRequest] = writerForSms

	for {
		ctx := context.TODO()
		fn := relayMessageForSend(writers)
		err := dbClient.RunInTransaction(ctx, fn)
		if err != nil {
			log.Errorf("message relay encountered error... sleeping for a bit... %v", err)
			time.Sleep(time.Second * 2)
		}
	}
}

func relayMessageForSend(writers map[string]*kafka.Writer) db.TransactionFunc {
	return func(ctx context.Context, tx db.Transaction) error {
		e, err := tx.ReadEvent(ctx)
		if err != nil {
			return fmt.Errorf("failed to read event from transactional outbox for sending emails: %w", err)
		}

		topic, err := getTopicFromDeliverCodeRequest(e.Catalyst)
		if err != nil {
			return fmt.Errorf("failed to infer topic name from catalyst protobuf: %w", err)
		}

		msg, err := goKafka.NewMessage(e.TraceInfo, e.EventInfo, e.Payload)
		if err != nil {
			return fmt.Errorf("failed to create event for topic: %s: %w", topic, err)
		}

		writer := writers[topic]

		if err := writer.WriteMessages(ctx, msg.ToKafkaMessage()); err != nil {
			return fmt.Errorf("failed to send error to topic: %s: %w", topic, err)
		}

		return tx.MarkEventAsSent(ctx, e.EventId)
	}
}

func getTopicFromDeliverCodeRequest(in mfaV1.DeliverCodeRequest) (string, error) {
	if in.GetEmailAddressId() != "" {
		return hermesTopics.TopicForSendEmailRequest, nil
	} else if in.GetPhoneNumberId() != "" {
		return hermesTopics.TopicForSendSmsRequest, nil
	}
	return "", errors.New("found DeliverCodeRequest with empty oneof")
}
