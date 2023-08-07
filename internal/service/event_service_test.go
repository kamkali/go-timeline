package service

import (
	"github.com/kamkali/go-timeline/internal/mocks"
	"github.com/kamkali/go-timeline/internal/timeline"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestGetEvent(t *testing.T) {
	var (
		ctx      = context.Background()
		repoMock = mocks.NewEventRepository(t)
		validID  = uint(10)
		event    = timeline.Event{
			ID:                  10,
			Name:                "test event",
			EventTime:           time.Now(),
			ShortDescription:    "short-desc",
			DetailedDescription: "long-desc",
			Graphic:             "img:abcd",
			TypeID:              1,
		}
	)

	eventService := NewEventService(nil, repoMock)
	t.Run("happy path", func(t *testing.T) {
		repoMock.On("GetEvent", ctx, validID).
			Return(event, nil).
			Once()

		e, err := eventService.GetEvent(ctx, validID)
		require.NoError(t, err)
		require.Equal(t, validID, e.ID)
	})
}
