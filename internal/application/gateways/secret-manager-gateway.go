package gateways

type ISecretManagerGateway interface {
	Get(key string) (string, error)
}
