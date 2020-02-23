package actions

import (
	"context"

	"github.com/lnkk-ai/lnkk/pkg/platform"
)

// StoreActionCorrelation is a helper to mange correlation keys
func StoreActionCorrelation(ctx context.Context, action, viewID, teamID string) error {
	err := platform.Set(ctx, correlationKey(viewID, teamID), action, 1800)
	if err != nil {
		platform.Report(err)
	}
	return err
}

// LookupActionCorrelation is a helper to mange correlation keys
func LookupActionCorrelation(ctx context.Context, viewID, teamID string) string {
	v, err := platform.Get(ctx, correlationKey(viewID, teamID))
	if err != nil {
		platform.Report(err)
		return ""
	}
	return v
}

func correlationKey(viewID, teamID string) string {
	return viewID + "." + teamID
}
