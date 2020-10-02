package urlshortener

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
)

/*
	StateActive and the other states decribe thr assets lifecycle
	StateArchived = the asset was disabled by its owner
	StateRetired = the asset was not activated for x days
	StateBroken = the assets target does not exist
*/
const (
	StateActive = iota
	StateArchived
	StateRetired
	StateBroken
)

type (
	// AssetRequest is the request body used to create a new asset
	AssetRequest struct {
		// Link is the long form URL
		Link string `json:"link" binding:"required"`
		// Owner identifies the owner of the asset
		Owner string `json:"owner" binding:"required"`
		// Secret is an optional attribute that can be used to 'claim' the asset
		Secret string `json:"secret,omitempty"`
		// ParentID is the id of the category the asset belongs to
		ParentID string `json:"parent,omitempty"`
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
		// AccessToken is used together with 'Secret' in order to claim the asset
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
		Secret      string `json:"secret,omitempty"`
		AccessToken string `json:"token,omitempty"`
		// metadata
		Tags        string `json:"tags,omitempty"`
		ParentID    string `json:"parent,omitempty"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		// status
		State int `json:"state,omitempty"`
		// segmentation
		//Source string `json:"source,omitempty"`
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
		platform.Report(err)
		return nil, err
	}
	return asset.asExternal(), nil
}

// GetURL retrieves the asset
func GetURL(ctx context.Context, uri string) (*AssetResponse, error) {
	var as Asset
	k := assetKey(uri)

	if err := platform.DataStore().Get(ctx, k, &as); err != nil {
		return nil, err
	}
	return as.asExternal(), nil
}

// LogRedirectRequest creates the analytics data for a redirect request
func LogRedirectRequest(ctx context.Context, shortLink string, c *gin.Context) error {

	h := RedirectHistory{
		ShortLink:      shortLink,
		Requester:      "unknown",                 // we don't know as we do not cookie requests
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
		platform.Report(err)
		return err
	}

	return nil
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
	shortLink, _ := util.ShortUUID()
	token, _ := util.ShortUUID()
	now := util.Timestamp()

	return &Asset{
		LongLink:    t.Link,
		ShortLink:   shortLink,
		PreviewLink: fmt.Sprintf("https://lnkk.host/r/%s", shortLink),
		Owner:       t.Owner,
		Secret:      t.Secret,
		AccessToken: token,
		ParentID:    t.ParentID,
		State:       StateActive,
		Created:     now,
		Modified:    now,
	}
}
