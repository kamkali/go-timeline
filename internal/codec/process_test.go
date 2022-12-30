package codec

import (
	"github.com/google/go-cmp/cmp"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"reflect"
	"testing"
	"time"
)

func TestHTTPToDomainProcess(t *testing.T) {
	tests := []struct {
		name    string
		e       *schema.Process
		want    *domain.Process
		wantErr bool
	}{
		{
			name: "successful conversion",
			e: &schema.Process{
				Name:                "Test Process",
				StartTime:           "2022-01-01T20:00:00.000+00:00",
				EndTime:             "2022-02-03T10:00:59.000+00:00",
				ShortDescription:    "Just testing",
				DetailedDescription: "Just testing",
				Graphic:             "some-base64-url-encoded-string",
				TypeID:              1,
			},
			want: &domain.Process{
				Name:                "Test Process",
				StartTime:           time.Date(2022, 1, 1, 20, 0, 0, 0, time.UTC),
				EndTime:             time.Date(2022, 2, 3, 10, 0, 59, 0, time.UTC),
				ShortDescription:    "Just testing",
				DetailedDescription: "Just testing",
				Graphic:             "some-base64-url-encoded-string",
				TypeID:              1,
			},
			wantErr: false,
		},
		{
			name: "invalid start time",
			e: &schema.Process{
				Name:                "Invalid Process",
				StartTime:           "2022-01-01",
				EndTime:             "2022-02-03T10:00:59.000+00:00",
				ShortDescription:    "Just testing",
				DetailedDescription: "Just testing",
				Graphic:             "some-base64-url-encoded-string",
				TypeID:              1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid end time",
			e: &schema.Process{
				Name:                "Invalid Process",
				StartTime:           "2022-01-01T20:00:00.000+00:00",
				EndTime:             "2022-02-03",
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
			got, err := HTTPToDomainProcess(tt.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPToDomainProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				cmp.Diff(got, tt.want)
				t.Fail()
			}
		})
	}
}

func TestHTTPFromDomainProcess(t *testing.T) {
	p1 := &domain.Process{
		ID:                  1,
		Name:                "Test Process",
		StartTime:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:             time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
		ShortDescription:    "Test process",
		DetailedDescription: "Test process description",
		Graphic:             "test.jpg",
		TypeID:              1,
	}

	httpProcess, err := HTTPFromDomainProcess(p1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := &schema.Process{
		ID:                  1,
		Name:                "Test Process",
		StartTime:           "2020-01-01T00:00:00Z",
		EndTime:             "2020-01-01T12:00:00Z",
		ShortDescription:    "Test process",
		DetailedDescription: "Test process description",
		Graphic:             "test.jpg",
		TypeID:              1,
	}
	if !reflect.DeepEqual(httpProcess, expected) {
		t.Errorf("Expected %+v, got %+v", expected, httpProcess)
	}
}
