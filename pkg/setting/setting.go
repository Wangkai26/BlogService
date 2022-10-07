package setting

import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper
}

func NewSetting() (*Setting,error) {
	vp := viper.New()
	vp.SetConfigName("config")   //文件名
	vp.AddConfigPath("configs/") //文件路径
	vp.SetConfigType("yaml")     //文件类型
	err := vp.ReadInConfig()
	if err != nil{
		return nil,err
	}

	return &Setting{vp},nil
}
