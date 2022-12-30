package codec

import (
	"github.com/google/go-cmp/cmp"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestHTTPToDomainType(t *testing.T) {
	schemaType := &schema.Type{
		Name:  "Test",
		Color: "Green",
	}

	want := &domain.Type{
		Name:  "Test",
		Color: "Green",
	}

	got, err := HTTPToDomainType(schemaType)
	require.NoError(t, err)

	if !cmp.Equal(got, want) {
		cmp.Diff(got, want)
		t.Fail()
	}
}

func TestHTTPFromDomainType(t *testing.T) {
	typ1 := &domain.Type{
		ID:    1,
		Name:  "Test Type",
		Color: "red",
	}
	httpType, err := HTTPFromDomainType(typ1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := &schema.Type{
		ID:    1,
		Name:  "Test Type",
		Color: "red",
	}
	if !reflect.DeepEqual(httpType, expected) {
		t.Errorf("Expected %+v, got %+v", expected, httpType)
	}
}
