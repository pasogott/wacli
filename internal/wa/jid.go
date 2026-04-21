package wa

import (
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

func ParseUserOrJID(s string) (types.JID, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return types.JID{}, fmt.Errorf("recipient is required")
	}
	if strings.Contains(s, "@") {
		return types.ParseJID(s)
	}
	s = strings.TrimPrefix(s, "+")
	if s == "" {
		return types.JID{}, fmt.Errorf("recipient is required")
	}
	return types.JID{User: s, Server: types.DefaultUserServer}, nil
}

func IsGroupJID(jid types.JID) bool {
	return jid.Server == types.GroupServer
}
