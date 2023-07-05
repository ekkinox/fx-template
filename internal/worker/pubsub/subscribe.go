package pubsub

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
)

type SubscribeWorker struct {
	config *fxconfig.Config
	logger *fxlogger.Logger
	client *pubsub.Client
}

func NewSubscribeWorker(config *fxconfig.Config, logger *fxlogger.Logger, client *pubsub.Client) *SubscribeWorker {
	return &SubscribeWorker{
		config: config,
		logger: logger,
		client: client,
	}
}

func (w *SubscribeWorker) Run(ctx context.Context) error {

	subscription, err := w.getSubscription(ctx)
	if err != nil {
		w.logger.Error().Err(err).Msg("cannot get subscription")

		return err
	}

	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		w.logger.Info().Msgf("received new message, data: %v", string(msg.Data))

		msg.Ack()
	})
	if err != nil {
		w.logger.Error().Err(err).Msg("error during subscription")
		return err
	}

	return nil
}

func (w *SubscribeWorker) getSubscription(ctx context.Context) (*pubsub.Subscription, error) {

	topicName := w.config.GetString("pubsub.topics.test")
	topic := w.client.Topic(topicName)

	topicExists, err := topic.Exists(ctx)
	if err != nil {
		w.logger.Error().Err(err).Msg("cannot check if topic exist")
		return nil, err
	}

	if !topicExists {
		w.logger.Info().Msgf("topic %s does not exist, creating it", topicName)
		topic, err = w.client.CreateTopic(ctx, topicName)
		if err != nil {
			w.logger.Error().Err(err).Msg("cannot create topic")
			return nil, err
		}
	}

	subscriptionName := w.config.GetString("pubsub.subscriptions.test")
	subscription := w.client.Subscription(subscriptionName)

	subscriptionExists, err := subscription.Exists(ctx)
	if err != nil {
		w.logger.Error().Err(err).Msg("cannot check if subscription exist")
		return nil, err
	}

	if !subscriptionExists {
		w.logger.Info().Msgf("subscription %s does not exist, creating it", subscriptionName)
		subscription, err = w.client.CreateSubscription(
			ctx,
			w.config.GetString("pubsub.subscriptions.test"),
			pubsub.SubscriptionConfig{
				Topic:       topic,
				AckDeadline: 10 * time.Second,
			},
		)
		if err != nil {
			w.logger.Error().Err(err).Msg("cannot create subscription")
			return nil, err
		}
	}

	return subscription, nil
}
