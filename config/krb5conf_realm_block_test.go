package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromString_Clockskew(t *testing.T) {
	c, err := NewFromString(`
[libdefaults]
  default_realm = EXAMPLE.COM
  clockskew = 900
`)
	require.NoError(t, err)
	assert.Equal(t, 900*time.Second, c.LibDefaults.Clockskew)
}

func TestNewFromString_AuthToLocalNamesSkipped(t *testing.T) {
	testCases := map[string]string{
		"block after kdc": `
[realms]
  EXAMPLE.COM = {
    kdc = kdc.example.com
    admin_server = kdc.example.com
    auth_to_local_names = {
      alice = local
    }
  }
`,
		"block before kdc": `
[realms]
  EXAMPLE.COM = {
    auth_to_local_names = {
      alice = local
    }
    kdc = kdc.example.com
    admin_server = kdc.example.com
  }
`,
		"one line after kdc": `
[realms]
  EXAMPLE.COM = {
    kdc = kdc.example.com
    admin_server = kdc.example.com
    auth_to_local_names = { alice = local }
  }
`,
		"one line before kdc": `
[realms]
  EXAMPLE.COM = {
    auth_to_local_names = { alice = local }
    kdc = kdc.example.com
    admin_server = kdc.example.com
  }
`,
	}

	for name, conf := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := NewFromString(conf)
			require.NoError(t, err)
			require.Len(t, c.Realms, 1)
			assert.Equal(t, "EXAMPLE.COM", c.Realms[0].Realm)
			assert.Equal(t, []string{"kdc.example.com:88"}, c.Realms[0].KDC)
			assert.Equal(t, []string{"kdc.example.com"}, c.Realms[0].AdminServer)
		})
	}
}
