package metadata

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

	expected := defMetadata
	expected.Info = testInfo

	actual := Base(testInfo)

	assert.Equal(t, expected, actual)

	assert.NotNil(t, actual.Routes)
	assert.NotNil(t, actual.Middlwares)
}
