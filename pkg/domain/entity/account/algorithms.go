package account

type AccountSecurityAlgorithms interface {
	GenerateSecretFromPassword(string) (string, error)
}
