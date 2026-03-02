package wikistream

import (
	"github.com/r3labs/sse/v2"
)

type Stream struct {
	client *sse.Client
	events chan []byte
	errs   chan error
}

func CreateStream(url string) *Stream {
	client := sse.NewClient(url)
	wikiStream := &Stream{
		client: client,
		events: make(chan []byte, 100), // buffered
	}

	go func() {
		// SubscribeRaw blocks and continuously pushes events
		client.SubscribeRaw(func(e *sse.Event) {
			wikiStream.events <- e.Data
		})
	}()

	return wikiStream
}

func (stream *Stream) ReadStream() ([]byte, error) {
	data := <-stream.events
	return data, nil
}
