package config

type Conf struct {
	Server *ConfServer
	DB     *ConfDB
	Redis  *ConfRedis
}

func New() *Conf {
	return &Conf{
		Server: NewConfServer(),
		DB:     NewConfDB(),
		Redis:  NewConfRedis(),
	}
}
