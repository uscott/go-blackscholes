package blackscholes_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uscott/go-blackscholes"
)

func TestVega(t *testing.T) {

	assert := assert.New(t)

	vega, err := blackscholes.Vega(0, 0, 0, 0, 0, 0, blackscholes.OptionType(' '))
	assert.Error(err)
	assert.True(math.IsNaN(vega))
}
