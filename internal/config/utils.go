package config

import (
	"log"
	"strings"
)

const (
	UTF8CodeOfComma  = 44
	UTF8CodeOfDot    = 46
	UTF8CodeOfSlash  = 47
	UTF8CodeOfZero   = 48
	UTF8CodeOfNine   = 57
	UTF8CodeOfColon  = 58
	UTF8CodeOfUpperA = 65
	UTF8CodeOfUpperZ = 90
	UTF8CodeOfLowerA = 97
	UTF8CodeOfLowerZ = 122
)

// strCleanUp removes all the extra characters added by different OSs environments.
func strCleanUp(strToCleanUp string) string {
	var builder strings.Builder
	for _, char := range strToCleanUp {
		switch {
		case char == UTF8CodeOfComma:
			builder.WriteRune(char)
		case char == UTF8CodeOfDot:
			builder.WriteRune(char)
		case char == UTF8CodeOfSlash:
			builder.WriteRune(char)
		case char == UTF8CodeOfColon:
			builder.WriteRune(char)
		case contains(UTF8CodeOfZero, UTF8CodeOfNine, char):
			builder.WriteRune(char)
		case contains(UTF8CodeOfZero, UTF8CodeOfNine, char):
			builder.WriteRune(char)
		case contains(UTF8CodeOfZero, UTF8CodeOfNine, char):
			builder.WriteRune(char)
		case contains(UTF8CodeOfUpperA, UTF8CodeOfUpperZ, char):
			builder.WriteRune(char)
		case contains(UTF8CodeOfLowerA, UTF8CodeOfLowerZ, char):
			builder.WriteRune(char)
		default:
			continue
		}
	}

	return builder.String()
}

func panicOnEmptyEnvVar(name string, value any) {
	switch value.(type) {
	case int, int8, int16, int32, uint, uint8, uint16, uint32:
		if value == 0 {
			log.Panicf("Zero value env var: %s", name)
		}
	case string:
		if value == "" {
			log.Panicf("Zero value env var: %s", name)
		}
	case any:
		if value == nil {
			log.Panicf("Zero value env var: %s", name)
		}
	}
}

func makeRange(min, max int32) []int32 {
	a := make([]int32, max-min+1)
	for i := range a {
		a[i] = min + int32(i)
	}
	return a
}

func contains(start, end, value int32) bool {
	for _, s := range makeRange(start, end) {
		if value == s {
			return true
		}
	}
	return false
}
