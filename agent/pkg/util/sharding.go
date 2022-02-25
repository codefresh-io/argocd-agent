package util

type Sharding struct {
	numberOfShard int
	replicas      int
}

func NewSharding(numberOfShard, replicas int) *Sharding {
	return &Sharding{
		numberOfShard,
		replicas,
	}
}

func (sh *Sharding) applicationsRange(amountOfApps int) (int, int) {
	appsPerShard := sh.replicas % amountOfApps
	from := appsPerShard * sh.numberOfShard
	return from, from + appsPerShard
}

func (sh *Sharding) PickApplications(applications []interface{}) []interface{} {
	from, to := sh.applicationsRange(len(applications))
	return applications[from:to]
}
