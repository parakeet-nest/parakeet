package wasm

import (
	"context"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
)

// one context for all plugins 🤔
var ctx = context.Background()

type Plugin struct {
	extismPlugin *extism.Plugin
}

func (p *Plugin) Call(functionName string, param []byte) ([]byte, error) {
	_, out, err := p.extismPlugin.Call(functionName, param)
	return out, err
}

func NewPlugin(path string, allowedPath map[string]string) (*Plugin, error) {
	config := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
		// LogLevel:     extism.LogLevelInfo, // Removed as it does not exist in extism.PluginConfig
	}
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path,
			},
		},
		AllowedHosts: []string{"*"}, // 🤔
		AllowedPaths: allowedPath,
		//AllowedPaths: map[string]string{
		//	"data": ".",
		//},
		Config: map[string]string{},
	}

	pluginInst, err := extism.NewPlugin(
		ctx,
		manifest,
		config,
		[]extism.HostFunction{},
	) // new

	if err != nil {
		return nil, err
	}

	return &Plugin{extismPlugin: pluginInst}, nil
}
