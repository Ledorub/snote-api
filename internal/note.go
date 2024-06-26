package internal

import (
	"errors"
	"fmt"
	"github.com/ledorub/snote-api/internal/datetime"
	"github.com/ledorub/snote-api/internal/validator"
	"time"
)

type Note struct {
	ID                string
	Content           *string
	CreatedAt         time.Time
	ExpiresIn         time.Duration
	ExpiresAt         time.Time
	ExpiresAtTimeZone *time.Location
	KeyHash           string
}

func (n *Note) CheckErrors() error {
	v := validator.Validator{}

	v.Check(validator.ValidateHyphenatedB58String(n.ID), "id should consist of latin letters and/or digits and hyphens")
	v.Check(len(*n.Content) != 0, "content should be provided")
	v.Check(len(*n.Content) <= 1_048_576, "content should not exceed 1 MB")
	v.Check(
		validator.ValidateTimeInRange(n.CreatedAt, time.Now().Add(-1*time.Minute), time.Now()),
		"time of the creation should be in range (now - 1 min, now]",
	)
	v.Check(len(n.KeyHash) == 44, "key hash must be exactly 44 bytes long")

	isExpiresInSet := n.ExpiresIn != 0
	isExpiresAtSet := !n.ExpiresAt.IsZero() && n.ExpiresAtTimeZone != nil
	hasConflict := isExpiresInSet && isExpiresAtSet || !(isExpiresInSet || isExpiresAtSet)
	v.Check(!hasConflict, "either expiration date and time zone or expiration timeout should be provided")
	if n.ExpiresIn != 0 {
		year := 24 * 60 * 365 * time.Minute
		v.Check(
			n.ExpiresIn >= 10*time.Minute && n.ExpiresIn <= 1*year,
			"expiration timeout should be in range [10 min, 365 days]",
		)
	} else {
		localCreatedAt := n.CreatedAt.In(n.ExpiresAtTimeZone)
		expiresAtLowerBound := time.Date(
			localCreatedAt.Year(), localCreatedAt.Month(), localCreatedAt.Day(),
			localCreatedAt.Hour(), localCreatedAt.Minute()+9, localCreatedAt.Second(), localCreatedAt.Nanosecond(),
			localCreatedAt.Location(),
		)
		expiresAtUpperBound := time.Date(
			localCreatedAt.Year()+1, localCreatedAt.Month(), localCreatedAt.Day(),
			localCreatedAt.Hour(), localCreatedAt.Minute(), localCreatedAt.Second(), localCreatedAt.Nanosecond(),
			localCreatedAt.Location(),
		)
		v.Check(
			validator.ValidateTimeInRange(n.ExpiresAt, expiresAtLowerBound, expiresAtUpperBound),
			"expiration date should be in range (local time + 9 min, local time + 1 year]",
		)
	}

	var validationErrors []error
	for _, err := range v.GetErrors() {
		validationErrors = append(validationErrors, err)
	}
	return errors.Join(validationErrors...)
}

func NewNote(
	content *string,
	expiresIn time.Duration,
	expiresAt time.Time,
	expiresAtTimeZone string,
	keyHash string,
) (*Note, error) {
	tz, err := time.LoadLocation(expiresAtTimeZone)
	if err != nil && expiresIn == 0 {
		return &Note{}, fmt.Errorf("unable to load time zone %v: %w", tz, err)
	}
	expiresAt = datetime.TimeAsLocalTime(expiresAt, tz)

	note := &Note{
		Content:           content,
		CreatedAt:         time.Now(),
		ExpiresIn:         expiresIn,
		ExpiresAt:         expiresAt,
		ExpiresAtTimeZone: tz,
		KeyHash:           keyHash,
	}
	return note, nil
}

type NoteModel struct {
	ID                uint64
	Content           *string
	CreatedAt         time.Time
	ExpiresAt         time.Time
	ExpiresAtTimeZone string
	KeyHash           []byte
}
