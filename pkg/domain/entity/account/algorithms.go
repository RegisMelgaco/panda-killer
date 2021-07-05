package account

type AccountSecurityAlgorithms interface {
	GenerateSecretFromPassword(string) (string, error)
	// first string argument is secret and second is the password
	CheckSecretAndPassword(string, string) error
}
