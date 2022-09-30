// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"

	"github.com/polarismesh/polaris-go"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// AdapterPolaris is the configuration for the polaris plugin
type AdapterPolaris struct {
	Namespace     string           `json:"namespace"`
	FileGroup     string           `json:"fileGroup"`
	FileName      string           `json:"fileName"`
	searchPaths   *garray.StrArray // Searching path array.
	jsonMap       *gmap.StrAnyMap  // The pared JSON objects for configuration files.
	violenceCheck bool             // Whether it does violence check in value index searching. It affects the performance when set true(false in default).
}

// NewAdapterFile creates a new AdapterPolaris
func NewAdapterFile(file ...string) (*AdapterPolaris, error) {
	conf, err := polaris.NewConfigAPI()
	if err != nil {
		return nil, err
	}

	if conf != nil {
		return &AdapterPolaris{
			Namespace:   "",
			FileGroup:   "",
			FileName:    "",
			searchPaths: garray.NewStrArray(true),
			jsonMap:     gmap.NewStrAnyMap(true),
		}, nil
	}
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
	// var usedFileName string = c.FileName
	// if path, _ := c.GetFilePath(usedFileName); path != "" {
	// 	return true
	// }
	// if c.GetContent(usedFileName) != "" {
	// 	return true
	// }
	return false
}

// SetContent sets customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func (c *AdapterPolaris) SetContent(content string, file ...string) {
	name := DefaultConfigFileName
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached `name`.
	localInstances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(name) {
			for _, v := range m {
				if configInstance, ok := v.(*Config); ok {
					if fileConfig, ok := configInstance.GetAdapter().(*AdapterPolaris); ok {
						fileConfig.jsonMap.Remove(name)
					}
				}
			}
		}
		customConfigContentMap.Set(name, content)
	})
}

// GetContent returns customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func (c *AdapterPolaris) GetContent(file ...string) string {
	name := DefaultConfigFileName
	if len(file) > 0 {
		name = file[0]
	}
	return customConfigContentMap.Get(name)
}

// RemoveContent removes the global configuration with specified `file`.
// If `name` is not passed, it removes configuration of the default group name.
func (c *AdapterPolaris) RemoveContent(file ...string) {
	name := DefaultConfigFileName
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached `name`.
	localInstances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(name) {
			for _, v := range m {
				if configInstance, ok := v.(*Config); ok {
					if fileConfig, ok := configInstance.GetAdapter().(*AdapterPolaris); ok {
						fileConfig.jsonMap.Remove(name)
					}
				}
			}
			customConfigContentMap.Remove(name)
		}
	})

	g.Log().Printf(context.TODO(), `RemoveContent: %s`, name)
}

// ClearContent removes all global configuration contents.
func (c *AdapterPolaris) ClearContent() {
	customConfigContentMap.Clear()
	// Clear cache for all instances.
	localInstances.LockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			if configInstance, ok := v.(*Config); ok {
				if fileConfig, ok := configInstance.GetAdapter().(*AdapterPolaris); ok {
					fileConfig.jsonMap.Clear()
				}
			}
		}
	})
	g.Log().Print(context.TODO(), `RemoveConfig`)
}

// getJson returns a *gjson.Json object for the specified `file` content.
// It would print error if file reading fails. It returns nil if any error occurs.
func (c *AdapterPolaris) getJson(fileName ...string) (configJson *gjson.Json, err error) {
	// usedFileName := c.FileName
	// if len(fileName) > 0 && fileName[0] != "" {
	// 	usedFileName = fileName[0]
	// } else {
	// 	usedFileName = c.FileName
	// }
	// // It uses json map to cache specified configuration file content.
	// result := c.jsonMap.GetOrSetFuncLock(usedFileName, func() interface{} {
	// 	var (
	// 		content  string
	// 		filePath string
	// 	)
	// 	// The configured content can be any kind of data type different from its file type.
	// 	isFromConfigContent := true
	// 	if content = c.GetContent(usedFileName); content == "" {
	// 		isFromConfigContent = false
	// 		filePath, err = c.GetFilePath(usedFileName)
	// 		if err != nil {
	// 			return nil
	// 		}
	// 		if filePath == "" {
	// 			return nil
	// 		}
	// 		if file := gres.Get(filePath); file != nil {
	// 			content = string(file.Content())
	// 		} else {
	// 			content = gfile.GetContents(filePath)
	// 		}
	// 	}
	// 	// Note that the underlying configuration json object operations are concurrent safe.
	// 	dataType := gfile.ExtName(filePath)
	// 	if gjson.IsValidDataType(dataType) && !isFromConfigContent {
	// 		configJson, err = gjson.LoadContentType(dataType, content, true)
	// 	} else {
	// 		configJson, err = gjson.LoadContent(content, true)
	// 	}
	// 	if err != nil {
	// 		if filePath != "" {
	// 			err = gerror.Wrapf(err, `load config file "%s" failed`, filePath)
	// 		} else {
	// 			err = gerror.Wrap(err, `load configuration failed`)
	// 		}
	// 		return nil
	// 	}
	// 	configJson.SetViolenceCheck(c.violenceCheck)
	// 	// Add monitor for this configuration file,
	// 	// any changes of this file will refresh its cache in Config object.
	// 	if filePath != "" && !gres.Contains(filePath) {
	// 		_, err = gfsnotify.Add(filePath, func(event *gfsnotify.Event) {
	// 			c.jsonMap.Remove(usedFileName)
	// 		})
	// 		if err != nil {
	// 			return nil
	// 		}
	// 	}
	// 	return configJson
	// })
	// if result != nil {
	// 	return result.(*gjson.Json), err
	// }
	return
}
