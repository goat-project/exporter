package utils

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

// StrToF64 converts string to float64.
func StrToF64(s string) float64 {
	if Null(s) {
		return 0
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 0
}

// StrToUint64 converts string to *uint64.
func StrToUint64(s string) *uint64 {
	return strToUint(s, 64)
}

// StrToUint32 converts string to uint32.
func StrToUint32(s string) uint32 {
	c := strToUint(s, 32)
	if c != nil {
		return uint32(*c)
	}

	return 0
}

// StrToFloat32 converts string to *float32.
func StrToFloat32(s string) *float32 {
	if Null(s) {
		return nil
	}

	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "string": s}).Error("unable to parse string to float")
		return nil
	}
	ff := float32(f)

	return &ff
}

// String converts string to *string.
func String(s string) *string {
	if Null(s) {
		return nil
	}

	return &s
}

// Null checks if string is one of NULL values.
func Null(s string) bool {
	return s == "" || s == "NULL" || s == "null" || s == "nil" || s == "<nil>"
}

func strToUint(s string, bitSize int) *uint64 {
	if Null(s) {
		return nil
	}

	u, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "string": s}).Error("unable to parse string to uint")
		return nil
	}

	return &u
}
