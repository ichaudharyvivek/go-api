package config

type Conf struct {
	Server *ConfServer
	DB     *ConfDB
}

func New() *Conf {
	return &Conf{
		Server: NewConfServer(),
		DB:     NewConfDB(),
	}
}
