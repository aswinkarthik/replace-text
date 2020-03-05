package replacer_test

import (
	"encoding/json"
	"testing"

	"github.com/aswinkarthik/replacer"
	assertions "github.com/stretchr/testify/assert"
)

func TestNode_AddString(t *testing.T) {
	assert := assertions.New(t)
	t.Run("should add one string", func(t *testing.T) {
		node := replacer.NewNode()

		assert.NoError(node.AddString("hello"))
		assert.NoError(node.AddString("help"))

		{
			data, err := json.MarshalIndent(node, "", "  ")
			assert.NoError(err)

			expectedMap := map[string]interface{}{
				"h": map[string]interface{}{
					"e": map[string]interface{}{
						"l": map[string]interface{}{
							"p": map[string]bool{
								"terminal": true,
							},
							"l": map[string]interface{}{
									"o": map[string]bool{
										"terminal": true,
									},
							},
						},
					},
				},
			}
			expected, _ := json.MarshalIndent(expectedMap, "", "  ")
			assert.Equal(string(expected), string(data))
		}

		{
			assert.True(node.Contains("hello"))
			assert.False(node.Contains("world"))
			assert.False(node.Contains("hell"))
		}
	})
}
