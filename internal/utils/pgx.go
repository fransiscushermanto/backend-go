package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ToPgUUID converts google/uuid.UUID to pgtype.UUID
func ToPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// FromPgUUID converts pgtype.UUID to google/uuid.UUID
func FromPgUUID(id pgtype.UUID) uuid.UUID {
	return id.Bytes
}

// ToPgTimestamp converts time.Time to pgtype.Timestamptz
func ToPgTimestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// ToPgTimestampPtr converts *time.Time to pgtype.Timestamptz
func ToPgTimestampPtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// FromPgTimestamp converts pgtype.Timestamptz to time.Time
func FromPgTimestamp(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// FromPgTimestampPtr converts pgtype.Timestamptz to *time.Time
func FromPgTimestampPtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
