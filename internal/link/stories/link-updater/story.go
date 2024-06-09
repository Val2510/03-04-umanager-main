package link_updater

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/internal/database"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/scrape"
)

const queueName = "link_queue"

func New(repository repository, consumer amqpConsumer, lg logger) *Story {
	return &Story{repository: repository, consumer: consumer, lg: lg}
}

type Story struct {
	repository repository
	consumer   amqpConsumer
	lg         logger
}

func (s *Story) Run(ctx context.Context) error {
	msgCh, err := s.consumer.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("channel Consume: %w", err)
	}

	go func() {
		for in := range msgCh {
			if err := s.process(ctx, in.Body); err != nil {
				s.lg.Error("link updater process", slog.Any("err", err))
				continue
			}
		}
	}()

	<-ctx.Done()

	return nil
}

func (s *Story) process(ctx context.Context, payload []byte) error {
	type message struct {
		ID string `json:"id"`
	}

	var m message
	if err := json.Unmarshal(payload, &m); err != nil {
		return fmt.Errorf("json Unmarshal: %w", err)
	}

	parsed, err := primitive.ObjectIDFromHex(m.ID)
	if err != nil {
		return fmt.Errorf("primitive ObjectIDFromHex: %w", err)
	}

	l, err := s.repository.FindByID(ctx, parsed)
	if err != nil {
		return fmt.Errorf("links repository FindByID: %w", err)
	}

	p, err := scrape.Parse(ctx, l.URL)
	if err != nil {
		return fmt.Errorf("scrape Parse: %w", err)
	}

	req := database.UpdateLinkReq{
		ID:     parsed,
		URL:    l.URL,
		Title:  l.Title,
		Images: l.Images,
		UserID: l.UserID,
	}

	if len(p.Title) > 0 {
		req.Title = p.Title
	}

	if len(p.Tags) > 0 {
		uniqTags := make(map[string]struct{})

		for _, t := range l.Tags {
			uniqTags[t] = struct{}{}
		}

		for _, t := range p.Tags {
			uniqTags[t] = struct{}{}
		}

		req.Tags = make([]string, 0, len(uniqTags))
		for k := range uniqTags {
			req.Tags = append(req.Tags, k)
		}
	}

	if _, err := s.repository.Update(ctx, req); err != nil {
		return fmt.Errorf("links repository Update: %w", err)
	}

	return nil
}
