package testcontainermongo

import (
	"context"
	"path/filepath"
	"regexp"
	"testing"
)

// TestNew runs an example mongodb container
func TestNew(t *testing.T) {
	schemaPath, err := filepath.Abs("init.sh")
	if err != nil {
		t.Fatalf("%v", err)
	}
	ctx := context.Background()
	c, conn, err := New(ctx, "latest", schemaPath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer Terminate(ctx, c)

	expected := regexp.QuoteMeta("mongodb://root:example@localhost:") + "[0-9]+" + regexp.QuoteMeta("/test_db")
	rx, err := regexp.Compile(expected)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !rx.MatchString(conn) {
		t.Errorf("Expected a connection string that looks like: %v, got: %v", expected, conn)
	}
}
