package conf

const (
	IkuraId    = 1
	IkuraStore = "/var/ikura/store/"
	IkuraPort  = 9090

	BabyWidth    = 400
	InfantWidth  = 200
	NewbornWidth = 100
	Sperm        = 1

	CacheMaxAge = 30 * 24 * 60 * 60 // 30 days
	Mime        = "image/jpeg"

	Mongodb = "mongodb://admin:12345678@localhost:27017/sa"
)
