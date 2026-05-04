package app

import (
	"context"
	"strings"
	"testing"
	"time"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waHistorySync"
	"go.mau.fi/whatsmeow/proto/waWeb"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func TestSyncStopsAtMaxMessages(t *testing.T) {
	a := newTestApp(t)
	f := newFakeWA()
	a.wa = f

	chat := types.JID{User: "123", Server: types.DefaultUserServer}
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	f.connectEvents = []interface{}{historySyncWithTextMessages(chat, base, "m1", "m2", "m3")}

	res, err := a.Sync(context.Background(), SyncOptions{
		Mode:        SyncModeFollow,
		AllowQR:     false,
		MaxMessages: 2,
	})
	if err == nil || !strings.Contains(err.Error(), "sync storage limit reached: message is 2, limit is 2") {
		t.Fatalf("Sync error = %v", err)
	}
	if res.MessagesStored != 2 {
		t.Fatalf("MessagesStored = %d, want 2", res.MessagesStored)
	}
	if n, err := a.db.CountMessages(); err != nil || n != 2 {
		t.Fatalf("db messages = %d, err=%v; want 2", n, err)
	}
}

func TestSyncRejectsExistingDBOverSizeLimit(t *testing.T) {
	a := newTestApp(t)
	f := newFakeWA()
	a.wa = f

	_, err := a.Sync(context.Background(), SyncOptions{
		Mode:           SyncModeOnce,
		AllowQR:        false,
		MaxDBSizeBytes: 1,
	})
	if err == nil || !strings.Contains(err.Error(), "sync storage limit reached: database size") {
		t.Fatalf("Sync error = %v", err)
	}
	if f.IsConnected() {
		t.Fatal("sync connected before rejecting oversized DB")
	}
}

func TestSyncWarnsWhenStorageUncapped(t *testing.T) {
	a := newTestApp(t)
	f := newFakeWA()
	a.wa = f

	out := captureStderr(t, func() {
		_, err := a.Sync(context.Background(), SyncOptions{
			Mode:         SyncModeOnce,
			AllowQR:      false,
			IdleExit:     time.Millisecond,
			WarnNoLimits: true,
		})
		if err != nil {
			t.Fatalf("Sync: %v", err)
		}
	})
	if !strings.Contains(out, "warning: sync storage is uncapped") {
		t.Fatalf("stderr = %q", out)
	}
}

func historySyncWithTextMessages(chat types.JID, start time.Time, ids ...string) *events.HistorySync {
	msgs := make([]*waHistorySync.HistorySyncMsg, 0, len(ids))
	for i, id := range ids {
		msgs = append(msgs, &waHistorySync.HistorySyncMsg{
			Message: &waWeb.WebMessageInfo{
				Key: &waCommon.MessageKey{
					RemoteJID: proto.String(chat.String()),
					FromMe:    proto.Bool(false),
					ID:        proto.String(id),
				},
				MessageTimestamp: proto.Uint64(uint64(start.Add(time.Duration(i) * time.Second).Unix())),
				Message:          &waProto.Message{Conversation: proto.String("text " + id)},
			},
		})
	}
	return &events.HistorySync{
		Data: &waHistorySync.HistorySync{
			SyncType: waHistorySync.HistorySync_FULL.Enum(),
			Conversations: []*waHistorySync.Conversation{{
				ID:       proto.String(chat.String()),
				Messages: msgs,
			}},
		},
	}
}
