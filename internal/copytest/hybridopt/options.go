package hybridopt

type Options struct {
	CollectFlags bool
}

func defaultOptions() Options {
	return Options{CollectFlags: true}
}
