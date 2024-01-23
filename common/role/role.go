package role

// 将身份代号转为名称
//  0-非在册队员 1-申请队员 2-岗前培训 3-见习队员 4-正式队员 5-督导老师 6-树洞之友 40-普通队员 41-核心队员 42-区域负责人 43-组委会成员 44-组委会主任
func GetRoleName(role int64) string {
	switch role {
	case 0:
		return "非在册队员"
	case 1:
		return "申请队员"
	case 2:
		return "岗前培训"
	case 3:
		return "见习队员"
	case 4:
		return "正式队员"
	case 5:
		return "督导老师"
	case 6:
		return "树洞之友"
	case 40:
		return "普通队员"
	case 41:
		return "核心队员"
	case 42:
		return "区域负责人"
	case 43:
		return "组委会成员"
	case 44:
		return "组委会主任"
	default:
		return "非在册队员"
	}
}

// 将身份名称转为代号
func GetRole(rolename string) int64 {
	switch rolename {
	case "申请队员":
		return 0
	case "岗前培训":
		return 1
	case "见习队员":
		return 2
	case "正式队员":
		return 3
	case "督导老师":
		return 4
	case "普通队员":
		return 30
	case "核心队员":
		return 31
	case "区域负责人":
		return 32
	case "组委会成员":
		return 33
	case "组委会主任":
		return 34
	default:
		return 0
	}
}
