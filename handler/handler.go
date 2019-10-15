package handler

import "sync"

type Reporter interface {
	ReportPower(on bool)
	ReportVolume(volume int)
}

type multiplexer []Reporter

func Multiplex(reporters ...Reporter) Reporter {
	return (multiplexer)(reporters)
}

func (mux multiplexer) fanOut(inner func(Reporter)) {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(len(mux))
	for _, muxReporter := range mux {
		reporter := muxReporter
		go func() {
			defer wg.Done()
			inner(reporter)
		}()
	}
}

func (mux multiplexer) ReportPower(on bool) {
	mux.fanOut(func(reporter Reporter) {
		reporter.ReportPower(on)
	})
}

func (mux multiplexer) ReportVolume(volume int) {
	mux.fanOut(func(reporter Reporter) {
		reporter.ReportVolume(volume)
	})
}