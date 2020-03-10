package replacer

// StateMachine represents a single fsm
// It will hold a cursor in Trie based on what
// characters have been passed through the fsm.
type StateMachine struct {
	StartPosition int
	EndPosition   int
	Terminated    bool
	ReplaceWith   string
	Node          *Node
}

// StateMachines holds a collection of machines.
// transitMachines are fsm still in transit nodes.
// TerminalMachines are fsm that reached the terminal nodes.
type StateMachines struct {
	transitMachines  []*StateMachine
	TerminalMachines []*StateMachine
	Root             *Node
}

// NewStateMachines is a constructor for creating StateMachines
func NewStateMachines(root *Node) *StateMachines {
	return &StateMachines{
		transitMachines:  make([]*StateMachine, 0),
		TerminalMachines: make([]*StateMachine, 0),
		Root:             root,
	}
}

// Accept will take in a single character and its position.
// It will pass through all the transit machines and transition it to
// the next state. If there is no state, the transit machine is discarded
// TerminalMachines hold all fsm that reached the terminal nodes.
func (s *StateMachines) Accept(ch byte, pos int) {
	resultMachines := make([]*StateMachine, 0, len(s.transitMachines))

	for _, m := range s.transitMachines {
		nextNode, err := m.Node.Next(ch)
		if err != nil {
			continue
		}

		m.Node = nextNode
		if nextNode.Terminates() {
			m.EndPosition = pos
			m.Terminated = true
			m.ReplaceWith = nextNode.value

			s.TerminalMachines = append(s.TerminalMachines, m)
		} else {
			resultMachines = append(resultMachines, m)
		}
	}

	nextNode, err := s.Root.Next(ch)
	if err == nil {
		m := &StateMachine{
			StartPosition: pos,
			Node:          nextNode,
		}
		if nextNode.Terminates() {
			m.Terminated = true
			m.EndPosition = pos
			m.ReplaceWith = nextNode.value

			s.TerminalMachines = append(s.TerminalMachines, m)
		} else {
			resultMachines = append(resultMachines, m)
		}
	}

	s.transitMachines = resultMachines
}
