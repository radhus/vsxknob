package handler

type Reporter interface {
	ReportPower(on bool)
	ReportVolume(volume int)
}
