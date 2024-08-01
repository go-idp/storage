package config

type Config struct {
	Port int

	// BaseDir is the base directory of the server
	//	- default: ""
	BaseDir string

	OSSAccessKeyID     string
	OSSAccessKeySecret string
	OSSBucket          string
	OSSEndpoint        string
	// OSSBaseDir is the base directory of the OSS
	//	- default: ""
	OSSBaseDir string
}
