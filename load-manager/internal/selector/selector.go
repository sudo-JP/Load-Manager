package selector

import (
	"github.com/sudo-JP/Load-Manager/load-manager/internal/registry"
)

type Selector interface {
	SelectNode(nodes []*registry.BackendNode) *registry.BackendNode
}
