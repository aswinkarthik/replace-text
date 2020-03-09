package replacer_test

import (
	"encoding/json"
	"testing"

	"github.com/aswinkarthik/replacer"
	assertions "github.com/stretchr/testify/assert"
)

func TestNode_AddString(t *testing.T) {
	assert := assertions.New(t)
	t.Run("should add strings and form a trie", func(t *testing.T) {
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
	})

	t.Run("should throw error if a string that is a prefix of the given string is already present in trie", func(t *testing.T) {
		node := replacer.NewNode()

		assert.NoError(node.AddString("hell"))
		assert.EqualError(node.AddString("hello"), replacer.ErrPrefixConflict.Error())
	})

	t.Run("should throw error if the given string is a prefix of an existing string in trie", func(t *testing.T) {
		node := replacer.NewNode()

		assert.NoError(node.AddString("hello"))
		assert.EqualError(node.AddString("hell"), replacer.ErrContainsConflict.Error())
	})
}

func TestNode_Contains(t *testing.T) {
	assert := assertions.New(t)
	t.Run("should return true for words in trie", func(t *testing.T) {
		node := replacer.NewNode()

		{
			assert.NoError(node.AddString("hello"))
			assert.NoError(node.AddString("help"))
		}

		{
			assert.True(node.Contains("hello"))
			assert.False(node.Contains("world"))
			assert.False(node.Contains("hell"))
		}
	})
}

func TestNode_AddReplacement(t *testing.T) {
	assert := assertions.New(t)
	t.Run("should add replacement text for given string", func(t *testing.T) {
		node := replacer.NewNode()

		{
			assert.NoError(node.AddReplacement("hello","world"))
			assert.Error(node.AddReplacement("hello","world"))
		}
	})
}

func TestNode_Next(t *testing.T) {
	t.Run("should return error if next node cannot be found", func(t *testing.T) {
		assert := assertions.New(t)
		node := replacer.NewNode()
		assert.NoError(node.AddString("test"))

		nextNode, err := node.Next(byte('a'))
		assert.Nil(nextNode)
		assert.EqualError(err, replacer.ErrNodeNotFound.Error())
	})

	t.Run("should return next node when found", func(t *testing.T) {
		assert := assertions.New(t)
		node := replacer.NewNode()
		assert.NoError(node.AddString("test"))

		nextNode, err := node.Next(byte('t'))

		assert.NotNil(nextNode)
		assert.NoError(err)
	})
}