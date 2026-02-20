package main

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/trust-me-im-an-engineer/logcheck/analyser"
)

func init() {
	register.Plugin("logcheck", New)
}

// PluginWrapper — обертка, которую требует golangci-lint
type PluginWrapper struct {
	// здесь можно хранить настройки (Bonus Task 1)
}

func New(settings any) (register.LinterPlugin, error) {
	// В будущем здесь можно десериализовать настройки через register.DecodeSettings
	return &PluginWrapper{}, nil
}

func (p *PluginWrapper) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyser.Analyzer}, nil
}

func (p *PluginWrapper) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
