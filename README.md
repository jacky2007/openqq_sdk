openqq_sdk
==========

OpenQQ SDK for golang
OpenQQ SDK for golang
腾讯开放平台V3版OpenAPI的golang SDK
=============================example=========================================

	parmas := make(map[string]string, 10)
	parmas["openid"] = "11111111111111111"
	parmas["openkey"] = "2222222222222222"
	parmas["pf"] = "pengyou"
	re, err := openqq.API("/v3/user/get_info", parmas, "http")
	if err == nil {
		fmt.Printf("%v\n", re)
		fmt.Println(re["nickname"])
		fmt.Println(re["ret"])
	}
