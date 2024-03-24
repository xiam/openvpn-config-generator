package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	v []configValue
	r string
}

var testCases = []testCase{
	{
		v: []configValue{
			{
				Name:   "verb",
				Type:   configTypeString,
				String: []string{"5"},
			},
		},
		r: `verb "5"`,
	},
	{
		v: []configValue{
			{
				Name:   "server",
				Type:   configTypeString,
				String: []string{"10.9.0.1", "255.255.0.0"},
			},
			{
				Name:   "keep-alive",
				Type:   configTypeString,
				String: []string{"10", "20"},
			},
		},
		r: "server \"10.9.0.1\" \"255.255.0.0\"\nkeep-alive \"10\" \"20\"",
	},
	{
		v: []configValue{
			{
				Name:   "verb",
				Type:   configTypeString,
				String: []string{"5"},
			},
			{
				Name:   "server",
				Type:   configTypeString,
				String: []string{"10.9.0.1", "255.255.0.0"},
			},
			{
				Name:   "keep-alive",
				Type:   configTypeString,
				String: []string{"10", "20"},
			},
		},
		r: "verb \"5\"\nserver \"10.9.0.1\" \"255.255.0.0\"\nkeep-alive \"10\" \"20\"",
	},
	{
		v: []configValue{
			{
				Name:   "verb",
				Type:   configTypeString,
				String: []string{"5"},
			},
			{
				Name: "client-to-client",
			},
		},
		r: "verb \"5\"\nclient-to-client",
	},
	{
		v: []configValue{
			{
				Name: "client-to-client",
			},
		},
		r: "client-to-client",
	},
	{
		v: []configValue{
			{
				Name: "client-to-client",
			},
			{
				Name: "persist-key",
			},
		},
		r: "client-to-client\npersist-key",
	},
	{
		v: []configValue{
			{
				Name: "client-to-client",
			},
			{
				Name:   "verb",
				Type:   configTypeString,
				String: []string{"5"},
			},
		},
		r: "client-to-client\nverb \"5\"",
	},
	{
		v: []configValue{
			{
				Name: "client-to-client",
			},
			{
				Name: "persist-key",
			},
			{
				Name:   "verb",
				Type:   configTypeString,
				String: []string{"5"},
			},
			{
				Name: "persist-tun",
			},
		},
		r: "client-to-client\npersist-key\nverb \"5\"\npersist-tun",
	},
	{
		v: []configValue{
			{
				Name: "client-to-client",
			},
			{
				Name: "persist-key",
			},
			{
				Name: "persist-tun",
			},
		},
		r: "client-to-client\npersist-key\npersist-tun",
	},
	{
		v: []configValue{
			{
				Name:  "tls-crypt",
				Type:  configTypeEmbed,
				Embed: []byte("foo"),
			},
		},
		r: "<tls-crypt>\nfoo\n</tls-crypt>",
	},
	{
		v: []configValue{
			{
				Name:  "tls-crypt",
				Type:  configTypeEmbed,
				Embed: []byte("foo\nbar\nbaz"),
			},
		},
		r: "<tls-crypt>\nfoo\nbar\nbaz\n</tls-crypt>",
	},
	{
		v: []configValue{
			{
				Name:  "tls-crypt",
				Type:  configTypeEmbed,
				Embed: []byte("foo\nbar\nbaz"),
			},
			{
				Name: "persist-tun",
			},
			{
				Name: "client-to-client",
			},
			{
				Name: "persist-key",
			},
		},
		r: "<tls-crypt>\nfoo\nbar\nbaz\n</tls-crypt>\npersist-tun\nclient-to-client\npersist-key",
	},
	{
		v: []configValue{
			{
				Name:  "tls-crypt",
				Type:  configTypeEmbed,
				Embed: []byte("foo\nbar\nbaz"),
			},
			{
				Name: "persist-tun",
			},
		},
		r: "<tls-crypt>\nfoo\nbar\nbaz\n</tls-crypt>\npersist-tun",
	},
	{
		v: []configValue{
			{
				Name:   "verb",
				Type:   configTypeString,
				String: []string{"5"},
			},
			{
				Name:  "tls-crypt",
				Type:  configTypeEmbed,
				Embed: []byte("foo\nbar\nbaz"),
			},
		},
		r: "verb \"5\"\n<tls-crypt>\nfoo\nbar\nbaz\n</tls-crypt>",
	},
	{
		v: []configValue{
			{
				Name:   "verb",
				Type:   configTypeString,
				String: []string{"5"},
			},
			{
				Name:  "tls-crypt",
				Type:  configTypeEmbed,
				Embed: []byte("foo\nbar\nbaz"),
			},
			{
				Name:   "keep-alive",
				Type:   configTypeString,
				String: []string{"10", "20"},
			},
			{
				Name: "client-to-client",
			},
		},
		r: "verb \"5\"\n<tls-crypt>\nfoo\nbar\nbaz\n</tls-crypt>\nkeep-alive \"10\" \"20\"\nclient-to-client",
	},
}

func TestTemplate(t *testing.T) {
	for _, testCase := range testCases {
		buf, err := compile(testCase.v)
		assert.NoError(t, err)
		assert.Equal(t, testCase.r, string(buf))
	}
}
