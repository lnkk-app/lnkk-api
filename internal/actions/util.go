package actions

import (
	"context"
	"log"

	"github.com/lnkk-ai/lnkk/pkg/platform"
)

// StoreActionCorrelation is a helper to mange correlation keys
func StoreActionCorrelation(ctx context.Context, action, viewID, teamID string) error {
	// FIXME remove this
	log.Printf("ctx-> %v", ctx)
	log.Printf("store-> %s,%s,%s", action, viewID, teamID)

	return platform.Set(ctx, correlationKey(viewID, teamID), action, "30")
}

// LookupActionCorrelation is a helper to mange correlation keys
func LookupActionCorrelation(ctx context.Context, viewID, teamID string) string {
	// FIXME remove this
	log.Printf("ctx-> %v", ctx)
	log.Printf("retrieve-> %s,%s", viewID, teamID)

	v, err := platform.Get(ctx, correlationKey(viewID, teamID))
	if err != nil {
		platform.Report(err)
		return ""
	}
	return v
}

func correlationKey(viewID, teamID string) string {
	return viewID + teamID
}
