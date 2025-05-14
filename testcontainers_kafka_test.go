package testcontainers_go_kafka_2670

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"golang.org/x/sync/errgroup"
)

var (
	imageFlag = flag.String("image", "confluentinc/confluent-local:7.5.0", "image with Kafka")
	runsFlag  = flag.Uint("runs", 10, "number of parallel runs")
)

func TestKafkaContainerStart(t *testing.T) {
	image := *imageFlag
	ctx := context.Background()
	group := errgroup.Group{}
	for left := *runsFlag; left > 0; left-- {
		group.Go(func() error {
			t.Logf("starting container using image: %s", image)
			container, err := kafka.Run(ctx, image)
			if err != nil {
				return fmt.Errorf("container start failed: %w", err)
			}
			defer func() {
				t.Logf("terminating container: %s", container.GetContainerID())
				err := container.Terminate(ctx)
				if err != nil {
					t.Errorf("container termination failed: %v", err)
				}
			}()
			t.Logf("successfully started container: %s", container.GetContainerID())
			return nil
		})
	}
	err := group.Wait()
	if err != nil {
		t.Fatalf("failed to run Kafka container: %v", err)
	}
}
