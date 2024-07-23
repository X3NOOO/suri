package routes

import "github.com/X3NOOO/suri/internal/ai"

type RoutingContext struct {
	AI        ai.AI
	MaxMemory int64
}
