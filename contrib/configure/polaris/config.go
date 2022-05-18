package polaris

import (
	"context"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
)

// AdapterPolaris is the configuration for the polaris plugin
type AdapterPolaris struct {
	Namespace   string           `json:"namespace"`
	FileGroup   string           `json:"fileGroup"`
	FileName    string           `json:"fileName"`
	searchPaths *garray.StrArray // Searching path array.
	jsonMap     *gmap.StrAnyMap  // The pared JSON objects for configuration files.
}

var (
	supportedFileTypes     = []string{"toml", "yaml", "yml", "json", "ini", "xml"} // All supported file types suffixes.
	localInstances         = gmap.NewStrAnyMap(true)                               // Instances map containing configuration instances.
	customConfigContentMap = gmap.NewStrStrMap(true)                               // Customized configuration content.

	// Prefix array for trying searching in resource manager.
	resourceTryFolders = []string{
		"", "/", "config/", "config", "/config", "/config/",
		"manifest/config/", "manifest/config", "/manifest/config", "/manifest/config/",
	}

	// Prefix array for trying searching in local system.
	localSystemTryFolders = []string{"", "config/", "manifest/config"}
)

// NewAdapterPolaris creates a new AdapterPolaris
func NewAdapterPolaris(file ...string) (*AdapterPolaris, error) {
	return nil, nil
}

// Get returns the configuration value for the given key.
func (c *AdapterPolaris) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	j, err := c.getJson()
	if err != nil {
		return nil, err
	}
	if j != nil {
		return j.Get(pattern).Val(), nil
	}
	return nil, nil
}

// Data retrieves and returns all configuration data as map type.
func (c *AdapterPolaris) Data(ctx context.Context) (data map[string]interface{}, err error) {
	j, err := c.getJson()
	if err != nil {
		return nil, err
	}
	if j != nil {
		return j.Var().Map(), nil
	}
	return nil, nil
}

// Available checks and returns whether configuration of given `file` is available.
func (c *AdapterPolaris) Available(ctx context.Context, fileName string) bool {
	var usedFileName string = c.FileName
	if path, _ := c.GetFilePath(usedFileName); path != "" {
		return true
	}
	if c.GetContent(usedFileName) != "" {
		return true
	}
	return false
}
