// The rdb package provides a basic bootstrap connection to a relational (SQL) database for Liquid Dev projects.
//
// rdb supports logging queries to stderr. This can activated via `LogQueries` called from a test file or by setting the 'DEBUG_SQL' environment variable. A value of 'before' will log the query prior to processing DEFAULTs. 'all' will show both prior and post-DEFAULT processing, and any other non-blank value will show post-processing values.
package rdb
