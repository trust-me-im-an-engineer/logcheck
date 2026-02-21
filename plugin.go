package logcheck

import (
	"fmt"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/trust-me-im-an-engineer/logcheck/analyser"
)

func init() {
	register.Plugin("logcheck", New)
}

type PluginWrapper struct {
}

func New(settings any) (register.LinterPlugin, error) {
	s, ok := settings.(map[string]any)
	if !ok {
		return &PluginWrapper{}, nil
	}

	if val, ok := s["sensitive-keywords"].(string); ok {
		if err := analyser.Analyzer.Flags.Set("sensitive-keywords", val); err != nil {
			return nil, fmt.Errorf("failed to set sensitive-keywords: %w", err)
		}
	}

	if s, ok := s["watched-logs"].(string); ok {
		if err := analyser.Analyzer.Flags.Set("watched-logs", s); err != nil {
			return nil, fmt.Errorf("failed to set watched-logs: %w", err)
		}
	}

	return &PluginWrapper{}, nil
}

func (p *PluginWrapper) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyser.Analyzer}, nil
}

func (p *PluginWrapper) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
