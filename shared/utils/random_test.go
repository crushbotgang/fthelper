package utils_test

import (
	"testing"

	"github.com/kamontat/fthelper/shared/utils"
	"github.com/kamontat/fthelper/shared/xtests"
)

func TestRandom(t *testing.T) {
	var assertion = xtests.New(t)

	assertion.NewName("random should always return different").
		WithExpected(utils.RandString(5)).
		WithActual(utils.RandString(5)).
		MustNotEqual()
}
