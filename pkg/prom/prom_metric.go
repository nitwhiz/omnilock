package prom

import (
	"encoding/json"
	"fmt"
	"github.com/nitwhiz/omnilock/pkg/server"
	"strings"
)

type SamplerFunc func(s *server.Server) *Sample

type Metric struct {
	Name    string
	Type    string
	Help    string
	Sampler SamplerFunc
}

func (m *Metric) helpLine() string {
	return fmt.Sprintf("# HELP %s %s", m.Name, m.Help)
}

func (m *Metric) typeLine() string {
	return fmt.Sprintf("# TYPE %s", m.Type)
}

func (m *Metric) value(s *server.Server) *Sample {
	return m.Sampler(s)
}

func (m *Metric) string(s *server.Server) string {
	str := strings.Builder{}

	if _, err := str.WriteString(m.helpLine()); err != nil {
		return ""
	}

	if err := str.WriteByte('\n'); err != nil {
		return ""
	}

	if _, err := str.WriteString(m.typeLine()); err != nil {
		return ""
	}

	if err := str.WriteByte('\n'); err != nil {
		return ""
	}

	sample := m.value(s)

	if _, err := str.WriteString(m.Name); err != nil {
		return ""
	}

	if len(sample.Labels) > 0 {
		bs, err := json.Marshal(sample.Labels)

		if err == nil {
			if _, err := str.Write(bs); err != nil {
				return ""
			}
		}
	}

	if err := str.WriteByte(' '); err != nil {
		return ""
	}

	if _, err := fmt.Fprintf(&str, "%f", sample.Value); err != nil {
		return ""
	}

	if err := str.WriteByte(' '); err != nil {
		return ""
	}

	if _, err := fmt.Fprintf(&str, "%d", sample.Timestamp); err != nil {
		return ""
	}

	if _, err := str.WriteString("\n\n"); err != nil {
		return ""
	}

	return str.String()
}
