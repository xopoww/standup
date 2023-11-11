package commandtypes

type Desc struct {
	Name  string `yaml:"name" validate:"required"`
	Usage string `yaml:"usage"`
	Short string `yaml:"short" validate:"required"`
	Long  string `yaml:"long"`
}
