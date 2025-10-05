package ratelimit

import (
	"errors"
	"time"
)

// ================================================
// Policy rate limiting
// ================================================
type Policy struct {
	Name   string        `json:"name"`
	Limit  int           `json:"limit"`
	Window time.Duration `json:"window"`
}

type PolicyInput struct {
	Name   string `json:"name"`
	Limit  int    `json:"limit"`
	Window string `json:"window"`
}

/**
 * PolicytRegistry is an interface for managing and retrieving rate limiting policies.
 */
type IPolicytRegistry interface {
	GetPolicy(name string) (Policy, bool)
}

// ================================================
//
//	Policy manager
//
// ================================================
type PolicyManager struct {
	policies map[string]Policy
}

/**
 * Manager instance Policy manager
 */
var _vIPolicytRegistry IPolicytRegistry

func GetPolicytRegistry() IPolicytRegistry {
	return _vIPolicytRegistry
}

func SetPolicytRegistry(v IPolicytRegistry) error {
	if v == nil {
		return errors.New("data init is nil")
	}
	if _vIPolicytRegistry != nil {
		return errors.New("data init already exists")
	}
	_vIPolicytRegistry = v
	return nil
}

/**
 * Implementation of the GetPolicy method to retrieve a policy by name.
 */
func NewPolicyManager(policies []PolicyInput) (IPolicytRegistry, error) {
	policyMap := make(map[string]Policy)
	for _, cfg := range policies {
		window, err := time.ParseDuration(cfg.Window)
		if err != nil {
			return nil, err
		}
		policyMap[cfg.Name] = Policy{
			Name:   cfg.Name,
			Limit:  cfg.Limit,
			Window: window,
		}
	}
	return &PolicyManager{policies: policyMap}, nil
}

func (pm *PolicyManager) GetPolicy(name string) (Policy, bool) {
	policy, exists := pm.policies[name]
	return policy, exists
}
