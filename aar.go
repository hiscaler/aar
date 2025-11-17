package aar

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"time"
)

type AAR struct {
	filename string        // Cached file name
	duration time.Duration // Cached expiration time
}

func New(name string, args ...string) (*AAR, error) {
	aar := &AAR{}
	if len(args) != 0 {
		name = fmt.Sprintf(name, args)
	}
	h := md5.New()
	_, err := io.WriteString(h, name)
	if err != nil {
		return nil, err
	}
	aar.filename = path.Join(os.TempDir(), fmt.Sprintf("%x", h.Sum(nil)))
	return aar, nil
}

// SetDuration set expiration time if you know specific hours
func (aar *AAR) SetDuration(hours int) *AAR {
	aar.duration = time.Duration(hours)
	return aar
}

// SetExpiredTime set expiration time if you only know the expiration time
func (aar *AAR) SetExpiredTime(d time.Time) *AAR {
	aar.duration = time.Until(d)
	return aar
}

// Read reads token from cache file if exists, and will be checked expiration time
func (aar *AAR) Read() (string, error) {
	var finfo fs.FileInfo
	var err error
	if finfo, err = os.Stat(aar.filename); !os.IsNotExist(err) && !finfo.IsDir() && finfo.ModTime().Add(aar.duration*time.Hour).After(time.Now()) {
		var b []byte
		if b, err = os.ReadFile(aar.filename); err == nil {
			return string(b), nil
		}
	}

	if err == nil {
		err = errors.New("aar: read failed or expired")
	}
	return "", err
}

// Write writes token to cache file
func (aar *AAR) Write(content []byte) error {
	return os.WriteFile(aar.filename, content, 0644)
}
