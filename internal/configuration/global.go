package configuration

type File struct {
	CurrentContext CurrentContext `yaml:"currentContext"`
}

type CurrentContext struct {
	AwsProfile string `yaml:"awsProfile"`
}
