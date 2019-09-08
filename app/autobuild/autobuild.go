package autobuild

type AutoBuildConfig struct {
	Buildpack string `json:"buildpack"`
	ImageName string `json:"image-name"`
}
