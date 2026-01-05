package utils

import "github.com/google/uuid"

func StringPointer(s string) *string {
	return &s
}

func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
