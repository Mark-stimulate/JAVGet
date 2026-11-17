package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAuthenticatedRoutes(t *testing.T) {
    apiKey = "test-only-key"
    handler := enableCORS(requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNoContent)
    })))
    for _, tc := range []struct { name, method, path, token, cookie string; want int }{
        {"missing", "GET", "/api/videos", "", "", 401},
        {"wrong", "GET", "/api/videos", "wrong", "", 401},
        {"valid", "GET", "/api/videos", "test-only-key", "", 204},
        {"preflight", "OPTIONS", "/api/videos", "", "", 200},
        {"file cookie", "GET", "/file/ABC-123/poster.jpg", "", "test-only-key", 204},
        {"cookie cannot add", "GET", "/api/addvideo/ABC-123", "", "test-only-key", 401},
    } {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest(tc.method, tc.path, nil)
            if tc.token != "" { req.Header.Set("Authorization", "Bearer "+tc.token) }
            if tc.cookie != "" { req.AddCookie(&http.Cookie{Name: "javget_key", Value: tc.cookie}) }
            response := httptest.NewRecorder()
            handler.ServeHTTP(response, req)
            if response.Code != tc.want { t.Fatalf("got %d, want %d", response.Code, tc.want) }
        })
    }
}

func TestVideoIDs(t *testing.T) {
    for _, value := range []string{"ABC-123", "ABC-DEF-123", "AB12-345"} {
        if !validVideoID.MatchString(value) { t.Errorf("valid ID rejected: %s", value) }
    }
    for _, value := range []string{"../ABC-123", "ABC-123;echo", "ABC", "ABC-123/secret"} {
        if validVideoID.MatchString(value) { t.Errorf("invalid ID accepted: %s", value) }
    }
}

func TestEnvironmentCompatibility(t *testing.T) {
    t.Setenv("JAVGET_TEST", "")
    t.Setenv("AVGRABBER_TEST", "")
    t.Setenv("NASSAV_TEST", "legacy")
    if value := envCompat("JAVGET_TEST", "AVGRABBER_TEST", "NASSAV_TEST", "default"); value != "legacy" { t.Fatalf("oldest setting ignored: %s", value) }
    t.Setenv("AVGRABBER_TEST", "previous")
    if value := envCompat("JAVGET_TEST", "AVGRABBER_TEST", "NASSAV_TEST", "default"); value != "previous" { t.Fatalf("previous setting ignored: %s", value) }
    t.Setenv("JAVGET_TEST", "current")
    if value := envCompat("JAVGET_TEST", "AVGRABBER_TEST", "NASSAV_TEST", "default"); value != "current" { t.Fatalf("new setting should take precedence: %s", value) }
}
