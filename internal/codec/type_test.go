package codec

import (
	"github.com/google/go-cmp/cmp"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"github.com/stretchr/testify/require"
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
