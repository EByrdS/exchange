package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"

	enginepb "exchange/engine/api/v1"
)

func (p *Printer) processMessage(record *kgo.Record) error {
	parts := strings.Split(record.Topic, ".")

	if len(parts) != 4 {
		return fmt.Errorf("invalid topic %q", record.Topic)
	}

	switch parts[3] {
	case "volumes":
		volumeEvent := &enginepb.VolumeEvent{}
		err := proto.Unmarshal(record.Value, volumeEvent)
		if err != nil {
			return err
		}

		log.Printf("volume event: %v\n", volumeEvent)
	case "orders":
		orderEvent := &enginepb.OrderEvent{}
		err := proto.Unmarshal(record.Value, orderEvent)
		if err != nil {
			return err
		}

		log.Printf("order event: %v\n", orderEvent)
	case "matches":
		matchEvent := &enginepb.MatchEvent{}
		err := proto.Unmarshal(record.Value, matchEvent)
		if err != nil {
			return err
		}

		log.Printf("match event: %v\n", matchEvent)
	default:
		return fmt.Errorf("Unsupported topic suffix %q", parts[3])
	}

	return nil
}

func (p *Printer) Listen(ctx context.Context) {
	go func() {
		for {
			fetches := p.kafka.PollFetches(ctx)
			fetches.EachPartition(func(partition kgo.FetchTopicPartition) {
				for _, record := range partition.Records {
					if err := p.processMessage(record); err != nil {
						fmt.Printf("Error: topic %q: %v\n", record.Topic, err)
					}

					p.kafka.CommitRecords(ctx, record)
				}
			})
		}
	}()
}
