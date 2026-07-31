package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	"github.com/jusoaresg/gorgon/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger() *slog.Logger {
	return slog.Default()
}

func newTestService(serverURL string, logger *slog.Logger) *QBittorrentService {
	host := strings.TrimPrefix(serverURL, "http://")
	return &QBittorrentService{
		APIService: services.NewAPIService(serverURL, logger),
		Logger:     logger,
		username:   "admin",
		password:   "adminadmin",
		host:       "http://" + host,
		port:       strings.TrimPrefix(host, "127.0.0.1:"),
	}
}

func TestLogin_Success(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v2/auth/login", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		r.ParseForm()
		assert.Equal(t, "admin", r.FormValue("username"))
		assert.Equal(t, "adminadmin", r.FormValue("password"))

		http.SetCookie(w, &http.Cookie{
			Name:  "QBT_SID_9191",
			Value: "test-sid-value",
			Path:  "/",
		})
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	err := svc.Login(&schema.QBittorrentLoginRequest{
		Username: "admin",
		Password: "adminadmin",
	})

	require.NoError(t, err)
	assert.Equal(t, "test-sid-value", svc.sid)
	assert.Equal(t, "QBT_SID_9191", svc.cookieName)
	assert.True(t, svc.loggedIn)
}

func TestLogin_Forbidden(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden."))
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	err := svc.Login(&schema.QBittorrentLoginRequest{
		Username: "admin",
		Password: "wrong",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestLogin_ClassicSID(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "SID",
			Value: "classic-sid",
			Path:  "/",
		})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	err := svc.Login(&schema.QBittorrentLoginRequest{
		Username: "admin",
		Password: "adminadmin",
	})

	require.NoError(t, err)
	assert.Equal(t, "classic-sid", svc.sid)
	assert.Equal(t, "SID", svc.cookieName)
}

func TestSidVerification_AlreadyLoggedIn(t *testing.T) {
	logger := newTestLogger()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		http.SetCookie(w, &http.Cookie{
			Name:  "QBT_SID_9191",
			Value: "test-sid",
			Path:  "/",
		})
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	svc.loggedIn = true

	err := svc.SidVerification()
	require.NoError(t, err)
	assert.Equal(t, 0, requestCount, "should not make any HTTP requests when already logged in")
}

func TestSidVerification_NeedsLogin(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/auth/login" {
			http.SetCookie(w, &http.Cookie{
				Name:  "QBT_SID_9191",
				Value: "fresh-sid",
				Path:  "/",
			})
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	svc.loggedIn = false

	err := svc.SidVerification()
	require.NoError(t, err)
	assert.True(t, svc.loggedIn)
}

func TestCookieHeader(t *testing.T) {
	svc := &QBittorrentService{
		cookieName: "QBT_SID_9191",
		sid:        "my-sid",
	}
	assert.Equal(t, "QBT_SID_9191=my-sid", svc.cookieHeader())

	svc2 := &QBittorrentService{
		cookieName: "SID",
		sid:        "classic",
	}
	assert.Equal(t, "SID=classic", svc2.cookieHeader())

	svc3 := &QBittorrentService{
		sid: "no-name",
	}
	assert.Equal(t, "SID=no-name", svc3.cookieHeader())
}

func TestAddTorrent_Success(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			http.SetCookie(w, &http.Cookie{Name: "QBT_SID_9191", Value: "test-sid", Path: "/"})
			w.WriteHeader(http.StatusNoContent)
		case "/api/v2/torrents/add":
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.Header.Get("Cookie"), "QBT_SID_9191=test-sid")

			r.ParseMultipartForm(10 << 20)
			assert.Equal(t, "magnet:?xt=urn:btih:test123", r.FormValue("urls"))
			assert.Equal(t, "gorgon", r.FormValue("category"))
			assert.Equal(t, "false", r.FormValue("paused"))

			resp := addTorrentResponse{
				AddedTorrentIds: []string{"test123"},
				FailureCount:    0,
				SuccessCount:    1,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	err := svc.AddTorrent("magnet:?xt=urn:btih:test123")
	require.NoError(t, err)
}

func TestAddTorrent_Failure(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			http.SetCookie(w, &http.Cookie{Name: "QBT_SID_9191", Value: "test-sid", Path: "/"})
			w.WriteHeader(http.StatusNoContent)
		case "/api/v2/torrents/add":
			resp := addTorrentResponse{
				AddedTorrentIds: []string{},
				FailureCount:    1,
				SuccessCount:    0,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	err := svc.AddTorrent("magnet:?xt=urn:btih:bad")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failures=1")
}

func TestDeleteTorrent_Success(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			http.SetCookie(w, &http.Cookie{Name: "QBT_SID_9191", Value: "test-sid", Path: "/"})
			w.WriteHeader(http.StatusNoContent)
		case "/api/v2/torrents/delete":
			assert.Equal(t, "POST", r.Method)
			r.ParseForm()
			assert.Equal(t, "abc123", r.FormValue("hashes"))
			assert.Equal(t, "true", r.FormValue("deleteFiles"))
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	err := svc.DeleteTorrent("abc123", true)
	require.NoError(t, err)
}

func TestDeleteTorrent_WrongMethod(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			http.SetCookie(w, &http.Cookie{Name: "QBT_SID_9191", Value: "test-sid", Path: "/"})
			w.WriteHeader(http.StatusNoContent)
		case "/api/v2/torrents/delete":
			if r.Method == "GET" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method Not Allowed"))
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	err := svc.DeleteTorrent("abc123", false)
	require.NoError(t, err)
}

func TestCheckTorrents_Success(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			http.SetCookie(w, &http.Cookie{Name: "QBT_SID_9191", Value: "test-sid", Path: "/"})
			w.WriteHeader(http.StatusNoContent)
		case "/api/v2/torrents/info":
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "downloading", r.URL.Query().Get("filter"))
			assert.Contains(t, r.Header.Get("Cookie"), "QBT_SID_9191=test-sid")

			torrents := []schema.CheckTorrentResponse{
				{
					Name:     "test-torrent",
					Progress: 0.5,
					State:    "downloading",
					Hash:     "abc123",
					DlSpeed:  1024000,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(torrents)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	var result []schema.CheckTorrentResponse
	err := svc.CheckTorrents("downloading", &result)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "test-torrent", result[0].Name)
	assert.Equal(t, float32(0.5), result[0].Progress)
}

func TestCheckTorrentsWithHash_Success(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			http.SetCookie(w, &http.Cookie{Name: "QBT_SID_9191", Value: "test-sid", Path: "/"})
			w.WriteHeader(http.StatusNoContent)
		case "/api/v2/torrents/info":
			assert.Equal(t, "all", r.URL.Query().Get("filter"))
			assert.Equal(t, "abc123hash", r.URL.Query().Get("hashes"))

			torrents := []schema.CheckTorrentResponse{
				{Name: "specific-torrent", Hash: "abc123hash", Progress: 1.0},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(torrents)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	var result []schema.CheckTorrentResponse
	err := svc.CheckTorrentsWithHash("all", "abc123hash", &result)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "specific-torrent", result[0].Name)
}

func TestCheckConnection_Success(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "QBT_SID_9191", Value: "test-sid", Path: "/"})
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	err := svc.CheckConnection()
	require.NoError(t, err)
}

func TestCheckConnection_Failure(t *testing.T) {
	logger := newTestLogger()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden."))
	}))
	defer server.Close()

	svc := newTestService(server.URL, logger)
	err := svc.CheckConnection()
	require.Error(t, err)
}
