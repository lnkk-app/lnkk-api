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

	// BackgroundWorkQueue is the default background job queue
	BackgroundWorkQueue string = "background-work"
	// DefaultCacheDuration default time to keep stuff in memory
	DefaultCacheDuration string = "10m"
	// DefaultCrawlerBatchSize number of messages per crawl
	DefaultCrawlerBatchSize int = 20
	// DefaultCrawlerSchedule default time between crawls in seconds
	DefaultCrawlerSchedule int = 600
	// DefaultUpdateSchedule default time between workspace updates in seconds
	DefaultUpdateSchedule int = 6000
)
