package replacer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReplacer(t *testing.T) {
	t.Run("should create replacer with state machines initialized", func(t *testing.T) {
		replacement := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		r, err := NewReplacer(replacement)

		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.NotNil(t, r.root)
	})
}

func TestReplacer_Run(t *testing.T) {
	t.Run("should find text from reader and replace it in writer", func(t *testing.T) {
		replacement := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		r, err := NewReplacer(replacement)
		assert.NoError(t, err)

		input := "this is text with key1 and key2 in their contents. Other keys need not be replaced. This is key1"
		expectedOutput := "this is text with value1 and value2 in their contents. Other keys need not be replaced. This is value1"

		reader := strings.NewReader(input)
		writer := &bytes.Buffer{}

		{
			err := r.Replace(reader, writer)
			assert.NoError(t, err)
		}

		assert.Equal(t, expectedOutput, writer.String())
	})

	t.Run("should work with input text larger than buffer size", func(t *testing.T) {
		replacement := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		r, err := NewReplacer(replacement)
		assert.NoError(t, err)

		input := "this is text with key1 and key2 in their contents. Other keys need not be replaced. This is key1"
		expectedOutput := "this is text with value1 and value2 in their contents. Other keys need not be replaced. This is value1"

		reader := strings.NewReader(input)
		writer := &bytes.Buffer{}

		{
			err := r.run(10, reader, writer)
			assert.NoError(t, err)
		}

		assert.Equal(t, expectedOutput, writer.String())
	})
}

func TestReplacer_ReplaceString(t *testing.T) {
	t.Run("should replace given string and return string when matches are found", func(t *testing.T) {
		replacement := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		r, err := NewReplacer(replacement)
		assert.NoError(t, err)

		input := "key1 need to be replaced by key2"
		expectedOutput := "value1 need to be replaced by value2"

		{
			output, err := r.ReplaceString(input)

			assert.NoError(t, err)
			assert.Equal(t, expectedOutput, output)
		}
	})

	t.Run("should return source string if there are no matches", func(t *testing.T) {
		replacement := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		r, err := NewReplacer(replacement)
		assert.NoError(t, err)

		input := "key3 need not be replaced by key4"

		{
			output, err := r.ReplaceString(input)

			assert.Equal(t, ErrNoMatchesFound, err)
			assert.Equal(t, input, output)
		}
	})
}
