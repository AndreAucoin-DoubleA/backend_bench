package redpanda

import (
	"backend_bench/internal/model"
	"bufio"
	"context"
	"fmt"
	"strings"
)

func ConnectRedpandaProducer(scanner *bufio.Scanner, producer *model.Producer, context context.Context) {
	var dataBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			eventData := dataBuilder.String()
			if eventData != "" {
				if err := producer.Produce(context, []byte(eventData)); err != nil {
					fmt.Printf("Error publishing: %v\n", err)
				}
			}
			dataBuilder.Reset()
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			dataBuilder.WriteString(strings.TrimPrefix(line, "data: "))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
	}
}
