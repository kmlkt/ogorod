package sys

import (
	"strconv"
	"time"
)

type ShortID int64

const baseShortID = 1577836800000

var lastShortID int64

func NewShortID() ShortID {
	result := time.Now().UnixMilli() - baseShortID
	if result <= lastShortID {
		result = lastShortID + 1
	}
	lastShortID = result
	return ShortID(result)
}

func (s ShortID) Next() ShortID {
	return ShortID(int64(NewShortID()) - int64(s))
}

func (s ShortID) String() string {
	return strconv.FormatInt(int64(s), 32)
}

func (s *ShortID) UnmarshalText(text []byte) error {
	result, err := strconv.ParseInt(string(text), 32, 64)
	if err != nil {
		return err
	}
	*s = ShortID(result)
	return nil
}

func (s ShortID) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}
