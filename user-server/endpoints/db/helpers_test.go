package db

import (
    "testing"
    "time"
)

func Test_getUpdates_behaviour(t *testing.T) {
    now := time.Now()
    url := &UrlData{Key: "k", Url: "u", Env: "e", UpdatedAt: &now}
    upd := getUpdates(url)
    // should include $set and the provided fields
    if len(upd) == 0 {
        t.Fatalf("expected updates non-empty")
    }
}
