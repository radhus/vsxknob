package vsx

var sources = map[string]string{
	"00": "PHONO",
	"01": "CD",
	"02": "TUNER",
	"03": "CD-R/TAPE",
	"04": "DVD",
	"05": "TV/SAT",
	"10": "VIDEO 1",
	"12": "MULTI CH IN",
	"14": "VIDEO 2",
	"15": "DVR/BDR",
	"17": "iPod/USB",
	"18": "XM RADIO",
	"19": "HDMI 1",
	"20": "HDMI 2",
	"21": "HDMI 3",
	"22": "HDMI 4",
	"23": "HDMI 5",
	"25": "BD",
	"26": "Internet Radio",
	"27": "SIRIUS",
	"33": "ADAPTER PORT",
}

var sourceIDs = map[string]string{}

func init() {
	for id, value := range sources {
		sourceIDs[value] = id
	}
}
