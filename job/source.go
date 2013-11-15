package job

type Source struct {
	name string
}

const defaultSource = "HTTP"
const overflowSource = "SQS"

func (s *Source) SetDefault() {
	s.name = defaultSource
}

func (s *Source) SetOverflow() {
	s.name = overflowSource
}

func (s *Source) IsDefault() bool {
	return s.name == defaultSource
}

func (s *Source) IsOverflow() bool {
	return s.name == overflowSource
}
