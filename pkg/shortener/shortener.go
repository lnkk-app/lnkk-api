package shortener

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"
)

const (
	// DatastoreAssets collection ASSETS
	DatastoreAssets string = "ASSETS"
	// DatastoreRedirectHistory collection MEASUREMENT
	DatastoreRedirectHistory string = "REDIRECT_HISTORY"

	// StateActive and the other states decribe thr assets lifecycle
	StateActive = iota
	// StateArchived = the asset was disabled by its owner
	StateArchived
	// StateExpired = the asset was not activated for x days
	StateExpired
	// StateBroken = the asset's target does not exist
	StateBroken

	// LastAccessThreshold is the time we allow to pass before updating the LastAccess attribute, again
	LastAccessThreshold = 3600 * 6 // 6h
	// ExpireAfter defines the age of an asset before it expires (in seconds).
	ExpireAfter int64 = 30 // days

	// DailyExpiration is a metric to track how many assets were expired
	DailyExpiration = "DAILY_EXPIRATION"
)

type (
	// AssetRequest is the request body used to create a new asset
	AssetRequest struct {
		// Link is the long form URL
		Link string `json:"link" binding:"required"`
		// Owner identifies the owner of the asset
		Owner string `json:"owner" binding:"required"`
		// ParentID is the id of the category the asset belongs to
		ParentID string `json:"parent,omitempty"`
		// Source identiefies the client who created the request
		Source string `json:"source,omitempty"`
	}

	// AssetResponse contains the relevant attributes after creating a new asset
	AssetResponse struct {
		// Link is the long form URL
		Link string `json:"link" binding:"required"`
		// ShortLink is the ID or suffix
		ShortLink string `json:"short_link,omitempty"`
		// PreviewLink is not use for now. Defaults to the canonical short link for now
		PreviewLink string `json:"preview_link,omitempty"`
		// Owner identifies the owner of the asset
		Owner string `json:"owner" binding:"required"`
		// AccessToken is used as a 'Secret' in order to claim or access the asset
		AccessToken string `json:"token,omitempty"`
	}

	//
	// INTERNAL
	//

	// Asset is the interal datastore structure used to store assets
	Asset struct {
		LongLink    string `json:"long_link" binding:"required"`
		ShortLink   string `json:"short_link" binding:"required"`
		PreviewLink string `json:"preview_link,omitempty"`
		// ownership etc
		Owner       string `json:"owner,omitempty"`
		AccessToken string `json:"token,omitempty"`
		// metadata
		Tags        string `json:"tags,omitempty"`
		ParentID    string `json:"parent,omitempty"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		// status
		State      int   `json:"state,omitempty"`
		LastAccess int64 `json:"last_access"`

		// segmentation
		Source string `json:"source,omitempty"`
		//Client string `json:"client,omitempty"`

		// internal
		Created  int64 `json:"-"`
		Modified int64 `json:"-"`
	}

	// RedirectHistory records redirect events
	RedirectHistory struct {
		ShortLink string `json:"short_link" binding:"required"`
		// requester metadata
		Requester string `json:"requester" binding:"required"`
		IP        string `json:"ip,omitempty"`
		Owner     string `json:"owner,omitempty"`
		// browser metadata
		UserAgent      string `json:"user_agent,omitempty"`
		AcceptLanguage string `json:"accept_language,omitempty"`
		// campaign metadata, see https://support.google.com/analytics/answer/1033863
		Source   string `json:"utm_source,omitempty"`
		Medium   string `json:"utm_medium,omitempty"`
		Campaign string `json:"utm_campaign,omitempty"`
		Content  string `json:"utm_content,omitempty"`
		// internal metadata
		Created int64 `json:"-"`
	}
)

// CreateURL creates a new asset
func CreateURL(ctx context.Context, as *AssetRequest) (*AssetResponse, error) {

	asset := as.asInternal()

	k := assetKey(asset.ShortLink)
	if _, err := platform.DataStore().Put(ctx, k, asset); err != nil {
		platform.ReportError(err)
		return nil, err
	}
	return asset.asExternal(), nil
}

// GetURL retrieves the asset
func GetURL(ctx context.Context, shortLink string, touch bool) (*Asset, error) {
	var asset Asset
	k := assetKey(shortLink)

	// FIXME: add caching at some point ...
	if err := platform.DataStore().Get(ctx, k, &asset); err != nil {
		return nil, err
	}

	if touch {
		// not every retrieval needs updating ...
		now := util.Timestamp()
		if asset.LastAccess < now-LastAccessThreshold {
			asset.LastAccess = now

			if _, err := platform.DataStore().Put(ctx, k, &asset); err != nil {
				platform.ReportError(err)
			}
		}
	}
	return &asset, nil
}

// LogRedirectRequest creates the analytics data for a redirect request
func LogRedirectRequest(ctx context.Context, asset *Asset, c *gin.Context) error {

	h := RedirectHistory{
		ShortLink:      asset.ShortLink,
		Requester:      "unknown", // we don't know as we do not cookie requests
		Owner:          asset.Owner,
		IP:             anonimizeIP(c.ClientIP()), // anonimize the IP to be GDPR compliant
		UserAgent:      strings.ToLower(c.GetHeader("User-Agent")),
		AcceptLanguage: strings.ToLower(c.GetHeader("Accept-Language")),
		Source:         c.Query("mtu_source"),
		Medium:         c.Query("mtu_medium"),
		Campaign:       c.Query("mtu_campaign"),
		Content:        c.Query("mtu_content"),
		Created:        util.Timestamp(),
	}

	// find the aproximation of the IPs real world location
	CreateGeoLocation(ctx, h.IP)

	k := datastore.IncompleteKey(DatastoreRedirectHistory, nil)
	if _, err := platform.DataStore().Put(ctx, k, &h); err != nil {
		platform.ReportError(err)
		return err
	}

	return nil
}

// GetAssets returns an array of count assets for owner, starting at page
func GetAssets(ctx context.Context, owner string, count, page int) *datastore.Iterator {
	q := datastore.NewQuery(DatastoreAssets).Filter("Owner =", owner).Limit(count).Order("-Created")
	return platform.DataStore().Run(ctx, q)
}

func assetKey(uri string) *datastore.Key {
	return datastore.NameKey(DatastoreAssets, uri, nil)
}

func (t *Asset) asExternal() *AssetResponse {
	return &AssetResponse{
		Link:        t.LongLink,
		ShortLink:   t.ShortLink,
		PreviewLink: t.PreviewLink,
		Owner:       t.Owner,
		AccessToken: t.AccessToken,
	}
}

func (t *AssetRequest) asInternal() *Asset {
	shortLink, _ := util.ShortUUID() // FIXME: SHORT or UNGUESSABLE
	token, _ := util.ShortUUID()
	now := util.Timestamp()

	return &Asset{
		LongLink:    t.Link,
		ShortLink:   shortLink,
		PreviewLink: fmt.Sprintf("https://lnkk.host/r/%s", shortLink), // FIXME: remove the static declaration of the host!
		Owner:       t.Owner,
		AccessToken: token,
		ParentID:    t.ParentID,
		State:       StateActive,
		LastAccess:  0, // easy to select assets that were never used ...
		Source:      t.Source,
		Created:     now,
		Modified:    now,
	}
}
