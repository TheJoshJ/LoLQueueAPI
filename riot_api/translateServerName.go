package riot_api

func translateServerName(server string) string {

	var serverString string

	switch server {
	case "BR":
		serverString = "BR1"
	case "EUNE":
		serverString = "EUN1"
	case "EUW":
		serverString = "EUW1"
	case "JP":
		serverString = "JP1"
	case "KR":
		serverString = "KR1"
	case "LAN":
		serverString = "LA1"
	case "LAS":
		serverString = "LA2"
	case "NA":
		serverString = "NA1"
	case "OCE":
		serverString = "OCE1"
	case "RU":
		serverString = "RU1"
	case "TR":
		serverString = "TR1"
	default:
		serverString = "NA1"
	}

	return serverString
}
