package jsonpath

// Integer Stack

type intStack struct {
	values []int
}

func newIntStack() *intStack {
	return &intStack{
		values: make([]int, 0, 100),
	}
}

func (s *intStack) len() int {
	return len(s.values)
}

func (s *intStack) push(r int) {
	s.values = append(s.values, r)
}

func (s *intStack) pop() (int, bool) {
	if s.len() == 0 {
		return 0, false
	}

	v, _ := s.peek()
	s.values = s.values[:len(s.values)-1]

	return v, true
}

func (s *intStack) peek() (int, bool) {
	if s.len() == 0 {
		return 0, false
	}

	v := s.values[len(s.values)-1]

	return v, true
}

// Interface Stack

type stack struct {
	values []interface{}
}

func newStack() *stack {
	return &stack{
		values: make([]interface{}, 0, 100),
	}
}

func (s *stack) len() int {
	return len(s.values)
}

func (s *stack) push(r interface{}) {
	s.values = append(s.values, r)
}

func (s *stack) pop() (interface{}, bool) {
	if s.len() == 0 {
		return nil, false
	}

	v, _ := s.peek()
	s.values = s.values[:len(s.values)-1]

	return v, true
}

func (s *stack) peek() (interface{}, bool) {
	if s.len() == 0 {
		return nil, false
	}

	v := s.values[len(s.values)-1]

	return v, true
}

func (s *stack) clone() *stack {
	d := stack{
		values: make([]interface{}, s.len()),
	}
	copy(d.values, s.values)

	return &d
}

func (s *stack) toArray() []interface{} {
	return s.values
}
