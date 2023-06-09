package pubsub

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
)

type SubscribeHandler struct {
	config *fxconfig.Config
	client *pubsub.Client
}

func NewSubscribeHandler(config *fxconfig.Config, client *pubsub.Client) *SubscribeHandler {
	return &SubscribeHandler{
		config: config,
		client: client,
	}
}

func (h *SubscribeHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		subscription, err := h.getSubscription(c)
		if err != nil {
			return err
		}

		var messages []*pubsub.Message
		var messagesCount int32

		err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			messages = append(messages, msg)
			atomic.AddInt32(&messagesCount, 1)
			msg.Ack()
		})
		if err != nil {
			c.Logger().Errorf("error during subscription: %v", err)
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"count":    messagesCount,
			"messages": messages,
		})
	}
}

func (h *SubscribeHandler) getSubscription(c echo.Context) (*pubsub.Subscription, error) {

	topicName := h.config.GetString("pubsub.topics.test")
	topic := h.client.Topic(topicName)

	topicExists, err := topic.Exists(c.Request().Context())
	if err != nil {
		c.Logger().Errorf("cannot check if topic exist: %v", err)
		return nil, err
	}

	if !topicExists {
		c.Logger().Infof("topic %s does not exist, creating it", topicName)
		topic, err = h.client.CreateTopic(c.Request().Context(), topicName)
		if err != nil {
			c.Logger().Errorf("cannot create topic: %v", err)
			return nil, err
		}
	}

	subscriptionName := h.config.GetString("pubsub.subscriptions.test")
	subscription := h.client.Subscription(subscriptionName)

	subscriptionExists, err := subscription.Exists(c.Request().Context())
	if err != nil {
		c.Logger().Errorf("cannot check if subscription exist: %v", err)
		return nil, err
	}

	if !subscriptionExists {
		c.Logger().Infof("subscription %s does not exist, creating it", subscriptionName)
		subscription, err = h.client.CreateSubscription(
			c.Request().Context(),
			h.config.GetString("pubsub.subscriptions.test"),
			pubsub.SubscriptionConfig{
				Topic:       topic,
				AckDeadline: 10 * time.Second,
			},
		)
		if err != nil {
			c.Logger().Errorf("cannot create subscription: %v", err)
			return nil, err
		}
	}

	return subscription, nil
}
