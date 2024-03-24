package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetVerb(t *testing.T) {
	config := New()

	{
		err := config.Set("verb")
		assert.Error(t, err, "no value was provided")
	}

	{
		err := config.Remove("verb")
		assert.Error(t, err, "setting does not exist")
	}

	{
		err := config.Set("verb", 6)
		assert.NoError(t, err)
	}

	{
		err := config.Set("verb", 5)
		assert.NoError(t, err, "should have been overwritten")
	}

	{
		err := config.Remove("verb")
		assert.NoError(t, err)
	}

	{
		err := config.Remove("verb")
		assert.Error(t, err, "key does not exist anymore")
	}

	{
		err := config.Set("verb", 7)
		assert.NoError(t, err)
	}

	{
		err := config.Set("verb", 5)
		assert.NoError(t, err)
	}

	buf, err := config.Compile()
	assert.NoError(t, err)

	assert.Equal(t, `verb "5"`, string(buf))
}

func TestMultipleRemote(t *testing.T) {
	config := New()

	{
		err := config.Add("remote", "server1.mydomain")
		assert.NoError(t, err)
	}
	{
		err := config.Add("remote", "server2.mydomain")
		assert.NoError(t, err)
	}
	{
		err := config.Add("remote", "server3.mydomain")
		assert.NoError(t, err)
	}

	buf, err := config.Compile()
	assert.NoError(t, err)

	assert.Equal(t, "remote \"server1.mydomain\"\nremote \"server2.mydomain\"\nremote \"server3.mydomain\"", string(buf))
}

func TestKey(t *testing.T) {
	config := New()

	value := []byte(`-----BEGIN OpenVPN Static key V1-----
e5e4d6af39289d53
171ecc237a8f996a
97743d146661405e
c724d5913c550a0c
30a48e52dfbeceb6
e2e7bd4a8357df78
4609fe35bbe99c32
bdf974952ade8fb9
71c204aaf4f256ba
eeda7aed4822ff98
fd66da2efa9bf8c5
e70996353e0f96a9
c94c9f9afb17637b
283da25cc99b37bf
6f7e15b38aedc3e8
e6adb40fca5c5463
-----END OpenVPN Static key V1-----`)

	{
		err := config.Embed("key", value)
		assert.NoError(t, err)
	}

	{
		err := config.Add("remote", "server2.mydomain")
		assert.NoError(t, err)
	}
	{
		err := config.Add("remote", "server3.mydomain")
		assert.NoError(t, err)
	}

	buf, err := config.Compile()
	assert.NoError(t, err)

	assert.Equal(t, "<key>\n"+string(value)+"\n</key>\nremote \"server2.mydomain\"\nremote \"server3.mydomain\"", string(buf))
}
