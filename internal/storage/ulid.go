package storage

import (
	"fmt"
	"io"
	"time"

	"github.com/oklog/ulid/v2"
)

type ulidGenerator struct {
	now     func() time.Time
	entropy io.Reader
}

func NewUlidGenerator(timeNow func() time.Time, entropy io.Reader) *ulidGenerator {
	return &ulidGenerator{now: timeNow, entropy: entropy}
}

func (u *ulidGenerator) Generate() (string, error) {
	t := u.now()
	id, err := ulid.New(ulid.Timestamp(t), u.entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate ulid: %w", err)
	}

	return id.String(), nil
}
