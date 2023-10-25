package nats

// stream name for the nats event
const (
	MMSStreamName = "mms-stream"

	// subject name for the nats event
	// this stream is using pull based
	// read this subject name stream like this:
	// service-name.event-name
	// that mean the event-name is the event that sent by the service-name
	StreamLaporanServiceAccountStatisticRefreshDataNeeded = "stream.laporan-service.account-statistic-refresh-data-needed"
	StreamAccountConnectedToFacebook                      = "stream.account-api.connected-to-facebook"
	StreamCrawlerServiceFetchedAccountData                = "stream.crawler-service.fetched-account-data"
	StreamCrawlerServiceFetchedMetaAccountData            = "stream.crawler-service.fetched-meta-account-data"

	// durable name for the nats event
	// durable name is the name for the durable consumer
	// durable consumer is a consumer that can receive the event that was sent before it was created
	// so that mean if the consumer is down, it will receive the event that was sent before it was down
	// durable name must be unique
	CrawlerDurableName = "crawler-service-durable"
	LaporanDurableName = "laporan-service-durable"

	// request reply nats
	// request reply nats is a nats that can send request and receive the response
	// request reply nats is used for the synchronous communication
	AccountGetAccountByIdReqRep          = "account-api.get-account-by-id-req-rep"
	AccountGetAssetMonitoredReqRep       = "account-api.get-asset-monitored-req-rep"
	AccountGetAccountIdsByBranchIdReqRep = "account-api.get-account-ids-by-branch-id-req-rep"
	AccountGetAllParentBranchIdReqRep    = "account-api.get-all-parent-branch-id-req-rep"
	AccountGetAllChildBranchIdReqRep     = "account-api.get-all-child-branch-id-req-rep"
)

func GetAllSubjectsForJetstream() []string {
	return []string{
		"stream.crawler-service.*",
		"stream.laporan-service.*",
		"stream.account-api.*",
	}
}
