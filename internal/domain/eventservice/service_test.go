package eventservice

import (
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/domain/mocks"
	logmock "github.com/kamkali/go-timeline/internal/logger/mocks"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestGetEvent(t *testing.T) {
	var (
		ctx        = context.Background()
		repoMock   = mocks.NewEventRepository(t)
		loggerMock = logmock.NewLogger(t)
		validID    = uint(10)
		event      = domain.Event{
			ID:                  10,
			Name:                "test event",
			EventTime:           time.Now(),
			ShortDescription:    "short-desc",
			DetailedDescription: "long-desc",
			Graphic:             "img:abcd",
			TypeID:              1,
		}
	)

	eventService := New(loggerMock, repoMock)
	t.Run("happy path", func(t *testing.T) {
		repoMock.On("GetEvent", ctx, validID).
			Return(event, nil).
			Once()

		e, err := eventService.GetEvent(ctx, validID)
		require.NoError(t, err)
		require.Equal(t, validID, e.ID)
	})
}
