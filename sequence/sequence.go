package sequence

// SequenceFunc is a function which can be used for generating a sequence of strings
type SequenceFunc func(min, max int) []string
