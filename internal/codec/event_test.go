package codec

import (
	"github.com/google/go-cmp/cmp"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"testing"
	"time"
)

func TestHTTPToDomainEvent(t *testing.T) {
	tests := []struct {
		name    string
		e       *schema.Event
		want    *domain.Event
		wantErr bool
	}{
		{
			name: "successful conversion",
			e: &schema.Event{
				Name:                "Test Event",
				EventTime:           "2022-01-01T20:00:00.000+00:00",
				ShortDescription:    "Just testing",
				DetailedDescription: "Just testing",
				Graphic:             "some-base64-url-encoded-string",
				TypeID:              1,
			},
			want: &domain.Event{
				Name:                "Test Event",
				EventTime:           time.Date(2022, 1, 1, 20, 0, 0, 0, time.UTC),
				ShortDescription:    "Just testing",
				DetailedDescription: "Just testing",
				Graphic:             "some-base64-url-encoded-string",
				TypeID:              1,
			},
			wantErr: false,
		},
		{
			name: "invalid time",
			e: &schema.Event{
				Name:                "Invalid Event",
				EventTime:           "2022-01-01",
				ShortDescription:    "Just testing",
				DetailedDescription: "Just testing",
				Graphic:             "some-base64-url-encoded-string",
				TypeID:              1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HTTPToDomainEvent(tt.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPToDomainEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				cmp.Diff(got, tt.want)
				t.Fail()
			}
		})
	}
}
