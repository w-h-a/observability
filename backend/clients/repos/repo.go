package repos

type Repo interface {
	Options() RepoOptions
	Read(dest interface{}, str string, additional ...interface{}) error
}
