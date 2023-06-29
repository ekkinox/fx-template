package fxpubsub

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
)

type PubSubProbe struct {
	config *fxconfig.Config
	client *pubsub.Client
	logger *fxlogger.Logger
}

func NewPubSubProbe(config *fxconfig.Config, client *pubsub.Client, logger *fxlogger.Logger) *PubSubProbe {
	return &PubSubProbe{
		config: config,
		client: client,
		logger: logger,
	}
}

func (p *PubSubProbe) Name() string {
	return "pubsub"
}

func (p *PubSubProbe) Check(ctx context.Context) *fxhealthchecker.HealthCheckerProbeResult {

	success := true
	var messages []string

	for _, topicName := range p.config.GetStringMapString("pubsub.topics") {
		topic := p.client.Topic(topicName)

		exist, err := topic.Exists(ctx)
		if err != nil {
			p.logger.Error().Err(err).Msgf("failed to check if topic %s exists", topicName)
		}

		if exist {
			messages = append(messages, fmt.Sprintf("topic %s exists", topicName))
		} else {
			messages = append(messages, fmt.Sprintf("topic does not %s exist", topicName))
		}

		success = success && exist
	}

	return fxhealthchecker.NewHealthCheckerProbeResult(success, strings.Join(messages, ", "))
}
