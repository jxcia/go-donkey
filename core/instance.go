package core

import "flag"

func New() *Garden {
	service := Garden{}

	var bootstrap string
	var env string
	flag.StringVar(&bootstrap, "bootstrap", "bootstrap.yml", "bootstrap yml files path")
	flag.StringVar(&env, "env", "dev", "bootstrap yml files path")
	flag.Parse()
	service.bootstrap(bootstrap, env)
	return &service
}
