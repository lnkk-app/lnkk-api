package backend

const (
	// DatastoreAssets collection ASSET
	DatastoreAssets string = "ASSETS"
	// DatastoreMeasurements collection MEASUREMENT
	DatastoreMeasurements string = "MEASUREMENTS"
	// DatastoreGeoLocations collection GEO_LOCATION
	DatastoreGeoLocations string = "GEO_LOCATIONS"
	// DatastoreAuthorizations collection AUTHORIZATION
	DatastoreAuthorizations string = "AUTHORIZATIONS"
	// DatastoreWorkspaces collection WORKSPACES
	DatastoreWorkspaces string = "WORKSPACES"
	// DatastoreUsers collection USERS
	DatastoreUsers string = "USERS"
	// DatastoreChannels collection CHANNELS
	DatastoreChannels string = "CHANNELS"
	// DatastoreMessages collection of messages
	DatastoreMessages string = "MESSAGES"
	// DatastoreAttachments collection of message attachments
	DatastoreAttachments string = "ATTACHMENTS"

	// BackgroundWorkQueue is the default background job queue
	BackgroundWorkQueue string = "background-work"
	// DefaultCacheDuration default time to keep stuff in memory
	DefaultCacheDuration string = "15m"
	// DefaultCrawlerBatchSize number of messages per crawl
	DefaultCrawlerBatchSize int = 50
	// DefaultCrawlerSchedule default time between crawls in seconds
	DefaultCrawlerSchedule int = 7200
	// DefaultUpdateSchedule default time between workspace updates in seconds
	DefaultUpdateSchedule int = 7200
)
