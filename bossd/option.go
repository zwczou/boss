package bossd

import (
	"os"
	"reflect"

	"github.com/flosch/pongo2"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Static struct {
	Dir  string
	Path string
}

type Template struct {
	Dir     string
	Context map[string]map[string]interface{}
}

func (t *Template) toPongoCtx() pongo2.Context {
	out := pongo2.Context{}
	for k, v := range t.Context {
		out[k] = v
	}
	return out
}

type Database struct {
	DataType     string `yaml:"data_type"`
	DataSource   string `yaml:"data_source"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type Redis struct {
	Addr        string
	DB          int
	MaxIdle     int `yaml:"max_idle"`
	IdleTimeout int `yaml:"idle_timeout"`
}

type Logger struct {
	Level string
	Skip  int
}

type option struct {
	Verbose  bool
	Secret   string
	HTTPAddr string `yaml:"http_addr"`
	Logger   Logger
	Database Database
	Redis    Redis
	Static   Static
	Template Template
}

func NewOption() *option {
	return &option{
		Verbose:  true,
		Secret:   "zwczou",
		HTTPAddr: ":8050",
		Logger: Logger{
			Level: "debug",
			Skip:  5,
		},
		Database: Database{
			DataType:     "mysql",
			DataSource:   "root:root@tcp(127.0.0.1:3306)/bossd?charset=utf8mb4&parseTime=true&loc=Local",
			MaxIdleConns: 20,
			MaxOpenConns: 200,
		},
		Redis: Redis{
			Addr:        "127.0.0.1:6379",
			DB:          1,
			MaxIdle:     3,
			IdleTimeout: 60,
		},
		Static: Static{
			Dir:  "./web/static",
			Path: "/static",
		},
		Template: Template{
			Dir: "./web/templates",
		},
	}
}

func (opts *option) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	err = yaml.NewDecoder(file).Decode(opts)
	if err != nil {
		file.Close()
		return err
	}
	file.Close()

	level, err := log.ParseLevel(opts.Logger.Level)
	if err != nil {
		return err
	}
	log.SetLevel(level)
	filenameHook := filename.NewHook()
	filenameHook.Field = "source"
	filenameHook.Skip = opts.Logger.Skip
	log.AddHook(filenameHook)
	return nil
}

func (opt *option) Print() {
	s := reflect.ValueOf(opt).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		log.WithField(typeOfT.Field(i).Name, f.Interface()).Info("option")
	}
}
