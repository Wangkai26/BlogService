package setting

import "time"

type ServerSettingS struct {
	RunMode		 string
	HttpPort	 string
	ReadTimeout	 time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize int
	MaxPageSize		int
	LogSavePath 	string
	LogFileName 	string
	LogFileExt 		string

	// 7part add code
	UploadSavePath  	 string
	UploadServerUrl 	 string
	UploadImageMaxSize   int
	UploadImageAllowExts []string
}

type DatabaseSettingS struct {
	DBType		string
	UserName	string
	Password	string
	Host		string
	DBName		string
	TablePrefix string
	Charset		string
	ParseTime	bool
	MaxIdleConns int
	MaxOpenConns int
}

func (s *Setting) ReadSection(k string,v interface{}) error  {
	err := s.vp.UnmarshalKey(k,v)
	if err != nil{
		return err
	}

	return nil
}

// 8section,设置JWT的一些相关配置
type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}