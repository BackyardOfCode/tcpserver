package utils

func GetHelp() string {
	s := `
 
		Commands:
		/help            获取帮助

		/all             获取所有在线用户
		/rename:tom      改用户名为tom
		/msg/tom/hello   给tom发信息
 
		`
	return s
}
