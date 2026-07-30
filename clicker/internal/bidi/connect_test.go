package bidi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// fakeSessionServer answers session.new and session.status the way the
// handler decides, reusing the helpers from session_test.go.
func fakeSessionServer(t *testing.T, handler func(*testServerConn, testCommand)) *httptest.Server {
	t.Helper()
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer ws.Close()
		s := &testServerConn{ws: ws}
		for {
			cmd, err := s.readCommand()
			if err != nil {
				return
			}
			handler(s, cmd)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

// A chromedriver per-session endpoint already has a live session and rejects
// session.new. ConnectRemote must attach to the existing session instead of
// failing.
func TestConnectRemoteAttachesWhenSessionExists(t *testing.T) {
	srv := fakeSessionServer(t, func(s *testServerConn, cmd testCommand) {
		switch cmd.Method {
		case "session.new":
			s.writeJSON(map[string]interface{}{
				"type": "error", "id": cmd.ID,
				"error": "session not created", "message": "session not created",
			})
		case "session.status":
			s.respond(cmd.ID, map[string]interface{}{"ready": false, "message": "already connected"})
		default:
			s.respond(cmd.ID, map[string]interface{}{})
		}
	})

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/session/abc123"
	_, client, sessionID, err := ConnectRemote(url, nil)
	if err != nil {
		t.Fatalf("ConnectRemote should attach to the existing session, got: %v", err)
	}
	defer client.Close()
	if sessionID != "abc123" {
		t.Errorf("sessionID: got %q, want %q", sessionID, "abc123")
	}
}

// When both session.new and session.status fail, ConnectRemote must return
// the original session.new error.
func TestConnectRemoteFailsWhenNoSessionUsable(t *testing.T) {
	srv := fakeSessionServer(t, func(s *testServerConn, cmd testCommand) {
		s.writeJSON(map[string]interface{}{
			"type": "error", "id": cmd.ID,
			"error": "session not created", "message": "session not created",
		})
	})

	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	_, _, _, err := ConnectRemote(url, nil)
	if err == nil {
		t.Fatal("expected error when no session can be created or attached")
	}
	if !strings.Contains(err.Error(), "session not created") {
		t.Errorf("error should carry the session.new failure, got: %v", err)
	}
}

func TestSessionIDFromURL(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"ws://127.0.0.1:9515/session/abc123", "abc123"},
		{"ws://host/session/", ""},
		{"ws://host:9222", ""},
		{"ws://host/session/a/b", ""},
	}
	for _, c := range cases {
		if got := sessionIDFromURL(c.in); got != c.want {
			t.Errorf("sessionIDFromURL(%q): got %q, want %q", c.in, got, c.want)
		}
	}
}
