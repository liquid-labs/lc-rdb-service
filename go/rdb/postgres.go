package rdb

import (
  "log"
  "os"

  "github.com/go-pg/pg"

  "github.com/Liquid-Labs/env/go/env"
)

// db the cached reference initialized by Connect.
var db *pg.DB
// debug controls debug logging behavior.
const debugKey string = `DEBUG_SQL`
var debug = os.Getenv(debugKey)

func EchoQueries(opt string) {
  if opt == `` {
    opt = `after`
  }
  debug = opt
}

// Note that an earlier version didn't include the 'before' and 'after' prefix
// and would fail to print the 'after' version. I suspect that it was comparing
// the first X bytes, seeing they were the samae and "very close" in time, and
// deciding that it was a duplicaate and suppressing.
func logQuery(prefix string, qe *pg.QueryEvent) {
  if q, err := qe.FormattedQuery(); err != nil {
    log.Printf(`bad query (%s): %s`, prefix, err)
  } else {
    log.Printf(`%s query: %s`, prefix, q)
  }
}

type dbLogger struct { }
func (d dbLogger) BeforeQuery(qe *pg.QueryEvent) {
  if debug == `before` || debug == `all` {
    logQuery(`before`, qe)
  }
}
func (d dbLogger) AfterQuery(qe *pg.QueryEvent) {
  if debug != `before` && debug != `` {
    logQuery(`after`, qe)
  }
}

// Connect initializes the DB connection. The following environment variables must be defined:
// * CLOUDSQL_CONNECTION_NAME
// * CLOUDSQL_DB
// * CLOUDSQL_PASSWORD
// * CLOUDSQL_USER
func Connect() *pg.DB {
  if (db != nil) {
    return db
  } else {
    options := pg.Options{
      User:     env.MustGet("CLOUDSQL_USER"),
      Password: env.MustGet("CLOUDSQL_PASSWORD"), // NOTE: password may NOT be empty
      Database: env.MustGet("CLOUDSQL_DB"),
    }
    // We currently start the proxy as an external process, so there's no need
    // for this.
    /*if env.IsTest() || env.IsDev() {
      options.Dialer = func(network, addr string) (net.Conn, error) {
        return proxy.Dial(env.MustGet(`CLOUDSQL_CONNECTION_NAME`))
      }
    } else {*/
      options.Addr = env.MustGet("CLOUDSQL_CONNECTION_NAME")
    //}

    db = pg.Connect(&options)
    if debug != `` {
      db.AddQueryHook(dbLogger{})
    }
    return db
  }
}
