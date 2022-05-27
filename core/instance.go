package core

import "flag"

func New() *Garden {
	service := Garden{}

	var bootstrap string
	var env string
	flag.StringVar(&bootstrap, "bootstrap", "boots", "bootstrap yml files path")
	flag.StringVar(&env, "bootstrap", "dev", "bootstrap yml files path")
	flag.Parse()
	service.bootstrap(bootstrap, env)
	return &service
}
