package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type StoredToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	// Unix seconds
	ExpiresAt int64 `json:"expires_at"`
}

var tokenFileMu sync.Mutex

func TokenFilePath() (string, error) {
	// Debug override: allow storing token into a fixed location.
	// Priority:
	// 1) BAIDUPAN_CLI_TOKEN_FILE (absolute or relative)
	// 2) BAIDUPAN_CLI_TOKEN_DIR  (directory; file name is token.json)
	if p := os.Getenv("BAIDUPAN_CLI_TOKEN_FILE"); p != "" {
		return p, nil
	}
	if d := os.Getenv("BAIDUPAN_CLI_TOKEN_DIR"); d != "" {
		return filepath.Join(d, "token.json"), nil
	}

	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exe)
	return filepath.Join(dir, "token.json"), nil
}

func LoadStoredToken() (*StoredToken, error) {
	p, err := TokenFilePath()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var t StoredToken
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, fmt.Errorf("parse token file %s error: %w", p, err)
	}
	if t.AccessToken == "" {
		return nil, nil
	}
	return &t, nil
}

func SaveStoredToken(t StoredToken) error {
	p, err := TokenFilePath()
	if err != nil {
		return err
	}
	tokenFileMu.Lock()
	defer tokenFileMu.Unlock()

	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}

	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

func (t StoredToken) Expired(skew time.Duration) bool {
	if t.ExpiresAt <= 0 {
		return false
	}
	return time.Now().Add(skew).Unix() >= t.ExpiresAt
}
