package config

// defaults struct
type defaults struct {
	AWSRegion string `yaml:"AWSRegion"`
	DBHost    string `yaml:"DBHost"`
	DBName    string `yaml:"DBName"`
	S3Bucket  string `yaml:"S3Bucket"`
	SsmPath   string `yaml:"SsmPath"`
	Stage     string `yaml:"Stage"`
}

type config struct {
	AWSRegion    string
	DBConnectURL string
	DBName       string
	S3Bucket     string
	Stage        StageEnvironment
}
