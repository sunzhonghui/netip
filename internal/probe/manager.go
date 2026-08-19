package probe

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// ProbeConfigFile structure for probes.yaml.
type ProbeConfigFile struct {
	Probes []NodeConfig `yaml:"probes"`
}

// Manager manages remote probe configurations and lifecycle.
type Manager struct {
	nodes  []NodeConfig
	client *ProbeClient
}

// NewManager loads probe configuration if file exists.
func NewManager(configPath string) *Manager {
	mgr := &Manager{
		client: NewProbeClient(),
	}

	if configPath == "" {
		return mgr
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		slog.Info("No probe config file loaded (optional)", "path", configPath)
		return mgr
	}

	var parsed ProbeConfigFile
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		slog.Warn("Failed to parse probe config YAML", "path", configPath, "err", err.Error())
		return mgr
	}

	mgr.nodes = parsed.Probes
	slog.Info("Loaded probe nodes config", "count", len(mgr.nodes), "path", configPath)
	return mgr
}

// Nodes returns list of configured probe nodes.
func (m *Manager) Nodes() []NodeConfig {
	return m.nodes
}

// Client returns probe RPC client.
func (m *Manager) Client() *ProbeClient {
	return m.client
}
