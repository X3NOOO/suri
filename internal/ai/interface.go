package ai

type AI interface {
	Query(string) (string, error)
}
