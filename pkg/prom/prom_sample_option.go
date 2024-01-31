package prom

type SampleOption func(*Sample)

func withLabel(name, value string) SampleOption {
	return func(s *Sample) {
		s.Labels[name] = value
	}
}
