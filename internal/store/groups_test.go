package store

import (
	"testing"
	"time"
)

func TestGroupsUpsertListAndParticipantsReplace(t *testing.T) {
	db := openTestDB(t)

	gid := "123@g.us"
	created := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	if err := db.UpsertGroup(gid, "Group", "owner@s.whatsapp.net", created); err != nil {
		t.Fatalf("UpsertGroup: %v", err)
	}
	if err := db.ReplaceGroupParticipants(gid, []GroupParticipant{
		{GroupJID: gid, UserJID: "a@s.whatsapp.net", Role: "admin"},
		{GroupJID: gid, UserJID: "b@s.whatsapp.net", Role: ""},
	}); err != nil {
		t.Fatalf("ReplaceGroupParticipants: %v", err)
	}

	gs, err := db.ListGroups("Gro", 10)
	if err != nil {
		t.Fatalf("ListGroups: %v", err)
	}
	if len(gs) != 1 || gs[0].JID != gid {
		t.Fatalf("expected group in list, got %+v", gs)
	}

	admins := countRows(t, db.sql, "SELECT COUNT(*) FROM group_participants WHERE group_jid=? AND role='admin'", gid)
	members := countRows(t, db.sql, "SELECT COUNT(*) FROM group_participants WHERE group_jid=? AND role='member'", gid)
	if admins != 1 || members != 1 {
		t.Fatalf("expected roles admin=1 member=1, got admin=%d member=%d", admins, members)
	}
}
