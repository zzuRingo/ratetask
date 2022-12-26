package application

import (
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"sync"
	_ "github.com/lib/pq"
)

type config struct {
	Postgres postgresConf
}

type postgresConf struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

var confOnce sync.Once
var connectionString *string

func GetDbConfInstance() *string {
	confOnce.Do(func() {
		conf := &config{}
		if _, err := toml.DecodeFile("./conf/config.toml", conf); err != nil {
			panic(err)
		}
		if conf.Postgres.Host == "" || conf.Postgres.Port == 0 ||
			conf.Postgres.User == "" || conf.Postgres.Password == "" || conf.Postgres.Dbname == "" {
			panic(fmt.Errorf("postgres conf missed some fields"))
		}

		fmt.Printf("host: %s\n", conf.Postgres.Host)
		fmt.Printf("port: %d\n", conf.Postgres.Port)
		fmt.Printf("user: %s\n", conf.Postgres.User)
		fmt.Printf("password: %s\n", conf.Postgres.Password)
		fmt.Printf("dbname: %s\n", conf.Postgres.Dbname)
		cs := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable", conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.User, conf.Postgres.Password, conf.Postgres.Dbname)
		connectionString = &cs
	})
	return connectionString
}

func GetStubCorrespondingPort(stub string) ([]string, error) {
	db, err := sql.Open("postgres", *GetDbConfInstance())
	fmt.Println("*GetDbConfInstance():",*GetDbConfInstance())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	fmt.Println("connect okay")
	// Execute the query
	query := fmt.Sprintf(`(WITH RECURSIVE regional_cte AS (
		select slug FROM regions WHERE parent_slug = '%s'
			UNION
			SELECT t.slug
			FROM regions t
			JOIN regional_cte rt ON rt.slug = t.parent_slug
		)
		select code FROM ports where parent_slug IN (
			SELECT slug FROM regional_cte
			UNION
			SELECT slug FROM regions WHERE slug = '%s')
		)`, stub, stub)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	portCodes := []string{}
	for rows.Next() {
		var code string
		err = rows.Scan(&code)
		if err != nil {
			return nil, err
		}
		portCodes = append(portCodes, code)
	}
	return portCodes, nil
}