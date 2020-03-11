package replacer_test

import (
	"github.com/aswinkarthik/replace-text/replacer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStateMachines_Accept(t *testing.T) {
	t.Run("should accept the given character and move a state machine", func(t *testing.T) {
		node := replacer.NewNode()
		{
			assert.NoError(t, node.AddString("hello"))
			assert.NoError(t, node.AddString("help"))
		}

		fsm := replacer.NewStateMachines(node)
		longText := "This is a long text, hello. After we continue there may be some help, and hello again. Even if hell is present, it does not matter."
		{
			for i := 0; i < len(longText); i++ {
				fsm.Accept(longText[i], int64(i))
			}
		}

		//data, _ := json.MarshalIndent(fsm.TerminalMachines, "", "  ")
		//fmt.Println(string(data))

		actualResults := make([]string, 0)
		actualStartPositions := make([]int64, 0)
		actualEndPositions := make([]int64, 0)
		for _, m := range fsm.TerminalMachines {
			actualResults = append(actualResults, longText[m.StartPosition:m.EndPosition+1])
			actualStartPositions = append(actualStartPositions, m.StartPosition)
			actualEndPositions = append(actualEndPositions, m.EndPosition)
		}

		expectedResults := []string{"hello", "help", "hello"}
		expectedStartPositions := []int64{21, 64, 74}
		expectedEndPositions := []int64{25, 67, 78}

		assert.Equal(t, expectedResults, actualResults)
		assert.Equal(t, expectedStartPositions, actualStartPositions)
		assert.Equal(t, expectedEndPositions, actualEndPositions)
	})
}
