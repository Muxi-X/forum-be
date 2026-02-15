package model

import (
	"errors"
	"fmt"
	"forum/pkg/constvar"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	watcher "github.com/casbin/redis-watcher/v2"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Casbin struct {
	Self *casbin.Enforcer
}

var CB *Casbin

func initCasbin(username, password, addr, DBName string) *casbin.Enforcer {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		username,
		password,
		addr,
		DBName)

	a, err := gormadapter.NewAdapter("mysql", config, true) // Your driver and data source.
	if err != nil {
		zap.L().Error("casbin数据库加载失败!", zap.Error(err))
		panic(err)
	}

	text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		g2 = _, _

		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = g(r.sub, p.sub) && g2(r.obj, p.obj) && r.act == p.act
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		zap.L().Error("字符串加载模型失败!", zap.Error(err))
		panic(err)
	}

	// Initialize the watcher.
	// Use the Redis host as parameter.
	w, _ := watcher.NewWatcher(viper.GetString("redis.addr"), watcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: viper.GetString("redis.password"),
		},
		Channel: "/casbin",
		// Only exists in test, generally be true
		IgnoreSelf: false,
	})

	cb, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}

	// Set the watcher for the enforcer.
	if err := cb.SetWatcher(w); err != nil {
		zap.L().Error("SetWatcher", zap.Error(err))
		// panic(err)
	}

	// Set callback to local example
	f := func(msg string) {
		err := cb.LoadPolicy()
		if err != nil {
			zap.L().Error("LoadPolicy", zap.Error(err))
			// panic(err)
		}
	}

	if err := w.SetUpdateCallback(f); err != nil {
		zap.L().Error("SetUpdateCallback", zap.Error(err))
		// panic(err)
	}

	return cb
}

func (c *Casbin) Init() {
	CB = &Casbin{
		Self: initCasbin(viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.addr"),
			viper.GetString("db.name")),
	}

	var err error

	_, err = CB.Self.AddRoleForUser(constvar.MuxiRole, constvar.NormalRole)
	_, err = CB.Self.AddRoleForUser(constvar.NormalAdminRole, constvar.NormalRole)

	_, err = CB.Self.AddRoleForUser(constvar.MuxiAdminRole, constvar.NormalAdminRole)
	_, err = CB.Self.AddRoleForUser(constvar.MuxiAdminRole, constvar.MuxiRole)

	_, err = CB.Self.AddRoleForUser(constvar.SuperAdminRole, constvar.MuxiAdminRole)

	_, err = CB.Self.AddPermissionForUser(constvar.NormalRole, []string{constvar.NormalDomain, constvar.Read}...)
	_, err = CB.Self.AddPermissionForUser(constvar.MuxiRole, []string{constvar.MuxiDomain, constvar.Read}...)

	_, err = CB.Self.AddPermissionForUser(constvar.NormalRole, []string{constvar.CollectionAndLike, constvar.Read}...)
	_, err = CB.Self.AddPermissionForUser(constvar.NormalRole, []string{constvar.Feed, constvar.Read}...)

	_, err = CB.Self.AddPermissionForUser(constvar.NormalAdminRole, []string{constvar.NormalDomain, constvar.Write}...)
	_, err = CB.Self.AddPermissionForUser(constvar.MuxiAdminRole, []string{constvar.MuxiDomain, constvar.Write}...)

	if err != nil {
		panic(err)
	}

	recourseRules := [][]string{
		{constvar.Post + ":" + constvar.MuxiDomain, constvar.MuxiDomain},
		{constvar.Post + ":" + constvar.NormalDomain, constvar.NormalDomain},
	}

	_, err = CB.Self.AddNamedGroupingPolicies("g2", recourseRules)
	if err != nil {
		panic(err)
	}

}

func Enforce(userId uint32, typeName string, data interface{}, act string) (bool, error) {
	object := typeName + ":"

	switch data.(type) {
	case string:
		object += data.(string)
	case uint32:
		object += strconv.Itoa(int(data.(uint32)))
	case int:
		object += strconv.Itoa(data.(int))
	default:
		return false, errors.New("wrong type")
	}

	return CB.Self.Enforce("user:"+strconv.Itoa(int(userId)), object, act)
}

func AddPolicy(userId uint32, typeName string, id uint32, act string) error {
	ok, err := CB.Self.AddPolicy("user:"+strconv.Itoa(int(userId)), typeName+":"+strconv.Itoa(int(id)), act)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("add policy not ok")
	}

	return nil
}

func DeletePermission(userId uint32, typeName string, id uint32, act string) error {
	ok, err := CB.Self.DeletePermissionForUser("user:"+strconv.Itoa(int(userId)), typeName+":"+strconv.Itoa(int(id)), act)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("delete permission not ok")
	}

	return nil
}

func AddRole(typeName string, id uint32, role string) error {
	ok, err := CB.Self.AddRoleForUser(typeName+":"+strconv.Itoa(int(id)), role)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("add role not ok")
	}

	return nil
}

func AddResourceRole(typeName string, id uint32, role string) error {
	ok, err := CB.Self.AddNamedGroupingPolicy("g2", typeName+":"+strconv.Itoa(int(id)), role)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("add g2 not ok")
	}

	return nil
}

func DeleteResourceRole(typeName string, id uint32, role string) error {
	ok, err := CB.Self.RemoveNamedGroupingPolicy("g2", typeName+":"+strconv.Itoa(int(id)), role)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("delete g2 not ok")
	}

	return nil
}

func HasRole(userId uint32, role string) (bool, error) {
	return CB.Self.HasRoleForUser("user:"+strconv.Itoa(int(userId)), role)
}
