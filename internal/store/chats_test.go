package store

import (
	"testing"
	"time"
)

func TestUpsertChatNameAndLastMessageTS(t *testing.T) {
	db := openTestDB(t)

	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	if err := db.UpsertChat("123@s.whatsapp.net", "dm", "Alice", t1); err != nil {
		t.Fatalf("UpsertChat: %v", err)
	}
	// Empty name should not clobber.
	if err := db.UpsertChat("123@s.whatsapp.net", "dm", "", t2); err != nil {
		t.Fatalf("UpsertChat empty name: %v", err)
	}
	c, err := db.GetChat("123@s.whatsapp.net")
	if err != nil {
		t.Fatalf("GetChat: %v", err)
	}
	if c.Name != "Alice" {
		t.Fatalf("expected name to stay Alice, got %q", c.Name)
	}
	if !c.LastMessageTS.Equal(t2) {
		t.Fatalf("expected LastMessageTS=%s, got %s", t2, c.LastMessageTS)
	}

	// Older timestamp should not override.
	if err := db.UpsertChat("123@s.whatsapp.net", "dm", "Alice2", t1); err != nil {
		t.Fatalf("UpsertChat older ts: %v", err)
	}
	c, err = db.GetChat("123@s.whatsapp.net")
	if err != nil {
		t.Fatalf("GetChat: %v", err)
	}
	if !c.LastMessageTS.Equal(t2) {
		t.Fatalf("expected LastMessageTS to remain %s, got %s", t2, c.LastMessageTS)
	}
}
