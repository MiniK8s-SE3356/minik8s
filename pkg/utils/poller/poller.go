package poller

import "time"

func PollerStaticPeriod(interval time.Duration, callback func(), if_loop bool) {
	<-time.After(interval)

	if if_loop {
		for {
			callback()
			<-time.After(interval)
		}
	} else {
		callback()
		<-time.After(interval)
	}
}
