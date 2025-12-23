package clientcmd

import "k8s.io/client-go/rest"

const RecommendedHomeFile = ""

func BuildConfigFromFlags(_ string, _ string) (*rest.Config, error) {
	return &rest.Config{}, nil
}
