package descriptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	testInfo := Info{
		Name:        "testName",
		Summary:     "testSummary",
		Description: "testDescription",
	}

	expectedDesc := defaultDescriptor
	expectedDesc.Info = testInfo

	actual := Base(testInfo)

	assert.Equal(t, expectedDesc, actual)

	assert.NotNil(t, actual.Routes)
	assert.NotNil(t, actual.Middlwares)
}
