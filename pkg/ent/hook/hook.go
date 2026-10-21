package hook

import "github.com/Toflex/directory_v2/ent"

// RegisterHooks attaches validation hooks to the client. Call this once
// at client construction time (e.g. in main.go or your db setup).
func RegisterHooks(client *ent.Client) {
	client.FeeTier.Use(validateTierRangesHook)
}
