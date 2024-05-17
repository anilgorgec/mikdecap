package tzsp

import "errors"

var (
	hLen                     = 4
	ErrHeaderTooShort        = errors.New("header too short")
	ErrUnsupportedVersion    = errors.New("unsupported version")
	ErrUnsupportedPacketType = errors.New("unsupported packet type")
	ErrTruncatedTag          = errors.New("truncated tag")
	ErrMissingEndTag         = errors.New("packet truncated")
)
