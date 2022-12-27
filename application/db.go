package application

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
)

type config struct {
	Postgres postgresConf
}

type postgresConf struct {
	Host         string
	Port         int
	User         string
	Password     string
	Dbname       string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifeTime  int
}

type priceCons struct {
	Date  string
	Price sql.NullInt32
}

var confOnce sync.Once
var dbConnectionString *string
var dbConnectionConf *postgresConf

var dbMutex sync.Mutex
var db *sql.DB

func GetDbConfInstance() (*string, *postgresConf) {
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
		fmt.Printf("maxOpenConns: %d\n", conf.Postgres.MaxOpenConns)
		fmt.Printf("maxIdleConns: %d\n", conf.Postgres.MaxIdleConns)
		fmt.Printf("maxLifeTime: %d\n", conf.Postgres.MaxLifeTime)
		cs := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable", conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.User, conf.Postgres.Password, conf.Postgres.Dbname)
		dbConnectionString = &cs
		dbConnectionConf = &conf.Postgres
	})
	return dbConnectionString, dbConnectionConf
}

func GetDbInstance() (*sql.DB, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	if db != nil {
		return db, nil
	}
	var err error = nil
	connectionStr, conf := GetDbConfInstance()
	db, err = sql.Open("postgres", *connectionStr)
	if err != nil {
		fmt.Printf("db open connection error :%+v", db)
		return nil, err
	}
	if conf.MaxOpenConns != 0 {
		db.SetMaxOpenConns(conf.MaxOpenConns)
	}
	if conf.MaxIdleConns != 0 {
		db.SetMaxIdleConns(conf.MaxIdleConns)
	}
	if conf.MaxLifeTime != 0 {
		db.SetConnMaxLifetime(time.Duration(conf.MaxLifeTime) * time.Minute)
	}
	return db, nil
}

func GetStubCorrespondingPort(stub string) ([]string, error) {
	db, e := GetDbInstance()
	if e != nil {
		return nil, e
	}
	fmt.Println("connect okay")

	query := `(WITH RECURSIVE regional_cte AS (
		select slug FROM regions WHERE parent_slug = $1
			UNION
			SELECT t.slug
			FROM regions t
			JOIN regional_cte rt ON rt.slug = t.parent_slug
		)
		select code FROM ports where parent_slug IN (
			SELECT slug FROM regional_cte
			UNION
			SELECT slug FROM regions WHERE slug = $2)
		)`
	rows, err := db.Query(query, stub, stub)
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

func GetDailyAvgPrice(oriPortCodes []string, destPortCodes []string,
	startDate string, endDate string) ([]priceCons, error) {
	db, err := GetDbInstance()
	if err != nil {
		return nil, err
	}
	fmt.Println("connect okay")

	// Execute the query
	query := fmt.Sprintf(`SELECT day,
		case when count(*)>=3 then round(avg(price),0)
			else null
		end as price
		FROM prices
   		where day BETWEEN $1 AND $2 and
		orig_code in ('%s') and
		dest_code in ('%s')
   		group by day
   		order by day asc`, strings.Join(oriPortCodes, `','`),
		strings.Join(destPortCodes, `','`))
	fmt.Println(query)
	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	cons := []priceCons{}
	for rows.Next() {
		p := priceCons{}
		err = rows.Scan(&p.Date, &p.Price)
		if err != nil {
			return nil, err
		}
		if len(p.Date) > 10 {
			p.Date = p.Date[0:10]
		}
		fmt.Println(p)
		cons = append(cons, p)
	}
	return cons, nil
}
