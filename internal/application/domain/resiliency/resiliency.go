package resiliency

const (
	OK                uint32 = 0
	CANCELLED         uint32 = 1
	UNKNOWN           uint32 = 2
	INVALIDARGUMENT   uint32 = 3
	DEADLINEEXCEEDED  uint32 = 4
	NOTFOUND          uint32 = 5
	ALREADYEXISTS     uint32 = 6
	PERMISSIONDENIED  uint32 = 7
	RESOURCEEXHAUSTED uint32 = 8
)
