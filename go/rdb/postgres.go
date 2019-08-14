package rdb

import (
  "log"
  "os"
  // "net"

  "github.com/go-pg/pg"
  // "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/proxy"

  "github.com/Liquid-Labs/env/go/env"
)

// db is the package local reference initialized by Connect.
var db *pg.DB
/*
type postModelInitHook func(db *pg.DB) error

// postModelInitHooks stores the hooks to be executed after model
// initialization.
var postModelInitHooks = make([]postModelInitHook, 0, 8)

// RegisterPostModelInitHook allows packages defining models to include routines
// to be run after the model-based schema has been created and before the DB
// accepts user connections. This may be used to insert procedures, indexes,
// or any other action which is not reflected in the models themselves.
//
// Any failure should result in a panic. Once all hooks have been successfully
// executed, the DB is considered fully initialized and ready for general use.
func RegisterPostModelInitHook(hook postModelInitHook) {
  postModelInitHooks = append(postModelInitHooks, hook)
}

// InitializeDB creates the model schema and runs any post-model init
// initialization hooks.
func InitializeDB(modelDefs ...interface{}) {
  Connect()
  db.AddQueryHook(InitLogger{})

  createOptions := pg.CreateTableOptions{
    FKConstraints : true,
  }

  for _, modelDef := range modelDefs {
    if err := db.CreateTable(modeDef, &createOptions); err != nil {
      // TODO: can we get this to print the struct def? If not, just name using
      // reflect?
      log.Panicf("Could not create table for %v", modelDef)
    }
  }

  for _, hook := range postModelInitHooks {
    if err := hook(db); err != nil {
      log.Panicf(`Could not execute post-model init hook: %+v; %s`, hook, err)
    }
  }
}
*/

var debug = os.Getenv(`DEBUG_SQL`)

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
