package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeConverter(t *testing.T) {
	t.Parallel()
	testCases := []string{
		"12:00:00AM",
		"12:00:00PM",
		"01:15:00AM",
		"01:15:00PM",
		"07:45:23AM",
		"08:23:45PM",
		"11:45:41PM",
		"10:12:07AM",
		"00:45:56AM",
		"12:23:53PM",
	}
	OKanswers := []string{
		"00:00:00",
		"12:00:00",
		"01:15:00",
		"13:15:00",
		"07:45:23",
		"20:23:45",
		"23:45:41",
		"10:12:07",
		"00:45:56",
		"12:23:53",
	}
	for i, clock := range testCases {
		clock := clock
		myCase := TimeConverter(clock)
		assert.Equal(t, myCase, OKanswers[i])

	}
}
