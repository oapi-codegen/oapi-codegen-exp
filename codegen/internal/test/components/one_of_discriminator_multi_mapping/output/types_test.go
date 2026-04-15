package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOneOfDiscriminatorMultiMapping tests that multiple discriminator mapping
// entries pointing to the same schema all resolve correctly via
// ValueByDiscriminator.
func TestOneOfDiscriminatorMultiMapping(t *testing.T) {
	httpConfigTypes := []string{
		"another_server",
		"apache_server",
		"web_server",
	}

	for _, configType := range httpConfigTypes {
		t.Run("http-"+configType, func(t *testing.T) {
			// Build the JSON directly to set the discriminator
			raw, err := json.Marshal(ConfigHTTP{
				ConfigType: configType,
				Host:       "example.com",
				Port:       8080,
			})
			require.NoError(t, err)

			var saveReq ConfigSaveReq
			err = saveReq.UnmarshalJSON(raw)
			require.NoError(t, err)

			cfg, err := saveReq.AsConfigHTTP()
			require.NoError(t, err)
			assert.Equal(t, configType, cfg.ConfigType)
			assert.Equal(t, "example.com", cfg.Host)

			got, err := saveReq.ValueByDiscriminator()
			require.NoError(t, err)
			assert.Equal(t, cfg, got)
		})
	}

	t.Run("ssh", func(t *testing.T) {
		raw, err := json.Marshal(ConfigSSH{
			ConfigType: "ssh_server",
		})
		require.NoError(t, err)

		var saveReq ConfigSaveReq
		err = saveReq.UnmarshalJSON(raw)
		require.NoError(t, err)

		cfg, err := saveReq.AsConfigSSH()
		require.NoError(t, err)
		assert.Equal(t, "ssh_server", cfg.ConfigType)

		got, err := saveReq.ValueByDiscriminator()
		require.NoError(t, err)
		assert.Equal(t, cfg, got)
	})
}

// TestFromConfigHTTPSetsDiscriminator verifies that FromConfigHTTP sets the
// discriminator value. V3 codegen forces the lexicographically first mapping
// value ("another_server") regardless of what the caller sets.
func TestFromConfigHTTPSetsDiscriminator(t *testing.T) {
	var saveReq ConfigSaveReq
	err := saveReq.FromConfigHTTP(ConfigHTTP{
		ConfigType: "web_server", // caller wants web_server
		Host:       "example.com",
		Port:       443,
	})
	require.NoError(t, err)

	disc, err := saveReq.Discriminator()
	require.NoError(t, err)
	// V3 codegen forces "another_server" (sorted first among the HTTP mappings).
	// If/when this is fixed to preserve the caller's value, change this assertion.
	assert.Equal(t, "another_server", disc)
}

func TestApplyDefaults(t *testing.T) {
	h := &ConfigHTTP{}
	h.ApplyDefaults()
	s := &ConfigSSH{}
	s.ApplyDefaults()
	r := &ConfigSaveReq{}
	r.ApplyDefaults()
}
