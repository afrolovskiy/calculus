package calculus

// The zero value for Stack is an empty stack ready to use.
type Stack struct {
	data []interface{}
}

// Push adds x to the top of the stack.
func (s *Stack) Push(x interface{}) {
	s.data = append(s.data, x)
}

// Pop removes and returns the top element of the stack.
// It’s a run-time error to call Pop on an empty stack.
func (s *Stack) Pop() interface{} {
	i := len(s.data) - 1
	res := s.data[i]
	s.data[i] = nil // to avoid memory leak
	s.data = s.data[:i]
	return res
}

// Peek returns the top element of the stack.
// It’s a run-time error to call Peek on an empty stack.
func (s *Stack) Peek() interface{} {
	return s.data[len(s.data)-1]
}

// Len returns the number of elements in the stack.
func (s *Stack) Len() int {
	return len(s.data)
}

// Empty returns true if stack is empty.
func (s *Stack) Empty() bool {
	return len(s.data) == 0
}
