package github_image

type RepoInterface interface {
	Push(PATH, filename, content string) (string, string, string)
	Del(filepath, sha string) string
}