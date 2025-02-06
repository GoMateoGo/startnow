package db

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

type TableEngine interface {
	DatabaseName() string
}

type LinkInfo struct {
	User    string       `mapstructure:"user" json:"user" yaml:"user"`
	Home    string       `mapstructure:"home" json:"home" yaml:"home"`
	Port    uint32       `mapstructure:"port" json:"port" yaml:"port"`
	Name    string       `mapstructure:"name" json:"name" yaml:"name"`
	Pass    string       `mapstructure:"pass" json:"pass" yaml:"pass"`
	Show    bool         `mapstructure:"show" json:"show" yaml:"show"`
	Idle    int          `mapstructure:"idle" json:"idle" yaml:"idle"`
	Open    int          `mapstructure:"open" json:"open" yaml:"open"`
	CId     int          `mapstructure:"c_id" json:"c_id" yaml:"c_id"`
	XEngine *xorm.Engine `json:"-"`
}

var EngineMap map[string]*LinkInfo

func NewEngine(companyId int64, table TableEngine) (*xorm.Engine, error) {
	var name string
	name = table.DatabaseName()
	if companyId > 0 {
		name = fmt.Sprintf("%s_%d", name, companyId)
	}
	engine, ok := EngineMap[name]
	if !ok {
		fmt.Println(companyId, name)
		return nil, errors.New("未找到对应的数据库链接信息，请联系技术")
	}
	return engine.XEngine, nil
}

func InitEngine(val []LinkInfo) error {
	EngineMap = make(map[string]*LinkInfo)
	for i := range val {
		if err := linkEngine(&val[i]); nil != err {
			return err
		}
		if 0 != val[i].CId {
			fmt.Println("初始化数据库名:", val[i].Name)
			EngineMap[val[i].Name] = &val[i]
		} else {
			EngineMap[val[i].Name] = &val[i]
		}
	}
	return nil
}

func linkEngine(val *LinkInfo) error {
	str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", val.User, val.Pass, val.Home, val.Port, val.Name)
	dbHand, err := xorm.NewEngine("mysql", str)
	if err != nil {
		fmt.Println("Web database link failed", err)
		return err
	}
	if err := dbHand.Ping(); err != nil {
		fmt.Println("Test link failed", err)
		return err
	}
	dbHand.SetTableMapper(core.SameMapper{})   // core.SameMapper{}
	dbHand.SetColumnMapper(core.SnakeMapper{}) // SnakeMapper  SameMapper
	dbHand.ShowSQL(val.Show)
	dbHand.SetMaxIdleConns(val.Idle)
	dbHand.SetMaxOpenConns(val.Open)
	val.XEngine = dbHand
	return nil
}
