package codec

import (
	"github.com/google/go-cmp/cmp"
	"github.com/kamkali/go-timeline/internal/server/schema"
	"github.com/kamkali/go-timeline/internal/timeline"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHTTPToDomainUser(t *testing.T) {
	schemaUser := &schema.User{
		Username: "test@email.com",
		Password: "pass",
	}

	want := &timeline.User{
		Email:    "test@email.com",
		Password: "pass",
	}

	got, err := HTTPToDomainUser(schemaUser)
	require.NoError(t, err)

	if !cmp.Equal(got, want) {
		cmp.Diff(got, want)
		t.Fail()
	}
}
