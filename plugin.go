package logcheck

import (
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
	return &PluginWrapper{}, nil
}

func (p *PluginWrapper) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyser.Analyzer}, nil
}

func (p *PluginWrapper) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
