package config

//Secret holds the key/value and images for Drone CI (0.5) secret variables: http://readme.drone.io/0.5/usage/secrets/
type Secret struct {
	Key    string
	Value  string
	Images []string
}
