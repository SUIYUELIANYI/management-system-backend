package github_image

//定义serve的映射关系
var serveMap = map[string]RepoInterface{
	"github": &GithubServe{},
}

func RepoCreate() RepoInterface {
	return serveMap[PLATFORM]
}
