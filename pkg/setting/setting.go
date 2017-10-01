package setting

import ini "gopkg.in/ini.v1"

func NewContext() error {
	iniFile, err := ini.Load("conf/app.ini")
	if err != nil {
		return err
	}
	iniFile.NameMapper = ini.TitleUnderscore
	err = iniFile.MapTo(&App)
	return err
}

var (
	App struct {
		Telegram struct {
			Token string
			Admin int64
		}

		Secret string
	}
)
