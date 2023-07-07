package testcontainermongo_test

import (
	"context"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/thlib/testcontainermongo"
)

// TestNew runs an example mongodb container
func TestNew(t *testing.T) {
	fixtures, err := filepath.Abs("initdb")
	if err != nil {
		t.Fatalf("%v", err)
	}
	ctx := context.Background()
	c, conn, err := testcontainermongo.New(ctx, "latest", fixtures)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer testcontainermongo.Terminate(ctx, c)

	expected := regexp.QuoteMeta("mongodb://root:example@localhost:") + "[0-9]+" + regexp.QuoteMeta("/test_db")
	rx, err := regexp.Compile(expected)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !rx.MatchString(conn) {
		t.Errorf("Expected a connection string that looks like: %v, got: %v", expected, conn)
	}
}
