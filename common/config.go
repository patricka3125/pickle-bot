package common

type Config struct {
	OpenAPI    OpenAPI    `yaml:"openAPI"`
	PickleBall Pickleball `yaml:"pickleball"`
}

type OpenAPI struct {
	AppID     string `yaml:"appId"`
	AppKey    string `yaml:"appKey"`
	ReceiveID string `yaml:"receiveId"`
}

type Pickleball struct {
	DocumentID string `yaml:"documentId"`
}
