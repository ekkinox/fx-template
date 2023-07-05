package pubsub

import (
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
)

type PublishHandler struct {
	config *fxconfig.Config
	client *pubsub.Client
}

func NewPublishHandler(config *fxconfig.Config, client *pubsub.Client) *PublishHandler {
	return &PublishHandler{
		config: config,
		client: client,
	}
}

func (h *PublishHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		data := time.Now().String()
		dataParam := c.QueryParam("data")
		if dataParam != "" {
			data = dataParam
		}

		topic, err := h.getTopic(c)
		if err != nil {
			return err
		}

		message := &pubsub.Message{
			Data: []byte(data),
		}

		if _, err = topic.Publish(c.Request().Context(), message).Get(c.Request().Context()); err != nil {
			c.Logger().Errorf("cannot publish message: %v", err)
			return err
		}

		return c.JSON(http.StatusCreated, map[string]bool{
			"success": true,
		})
	}
}

func (h *PublishHandler) getTopic(c echo.Context) (*pubsub.Topic, error) {

	topicName := h.config.GetString("modules.pubsub.topics.test")
	topic := h.client.Topic(topicName)

	exists, err := topic.Exists(c.Request().Context())
	if err != nil {
		c.Logger().Errorf("cannot check if topic exist: %v", err)
		return nil, err
	}

	if !exists {
		c.Logger().Infof("topic %s does not exist, creating it", topicName)
		topic, err = h.client.CreateTopic(c.Request().Context(), topicName)
		if err != nil {
			c.Logger().Errorf("cannot create topic: %v", err)
			return nil, err
		}
	}

	return topic, nil
}
