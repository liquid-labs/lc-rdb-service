package rdb_test

import (
  "os"
  "strconv"
  "testing"

  "github.com/stretchr/testify/require"

  // the package we're testing
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
)

func TestPostgresDBIntegration(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  }

  t.Run(`DBSetup`, testDBSetup)
}

type stringResult struct {
  Result string
}

func testDBSetup(t *testing.T) {
  rdb := rdb.Connect()
  defer rdb.Close()

  ver := &stringResult{}
  _, err := rdb.QueryOne(ver, "SELECT current_setting('server_version_num') AS result")
  require.NoError(t, err, `Unexpected error testing DB connection.`)
  verNum, err := strconv.Atoi(ver.Result)
  require.NoErrorf(t, err, "Unexpected DB version format: %s", ver.Result)
  require.NotEmpty(t, ver.Result, "Version result is empty.")
  // Not yet released:
  // require.Greaterf(t, verNum, 0, "Unexpected version number: %i", verNum)
  require.Truef(t, verNum > 0, "Unexpected version number: %i", verNum)
}
