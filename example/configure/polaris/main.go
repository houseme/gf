// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"fmt"

	goPolaris "github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/gogf/gf/contrib/configure/polaris/v2"
)

func main() {
	if err := api.SetLoggersDir("/tmp/polaris/log"); err != nil {
		fmt.Println("fail to set logger example.", err)
		return
	}
	configAPI, err := goPolaris.NewConfigAPIByFile(
		"/Users/qun/Documents/golang/houseme/gf/example/configure/polaris/polaris.yaml")

	if err != nil {
		fmt.Println("fail to start example.", err)
		return
	}
	// 获取远程的配置文件
	namespace := "default"
	fileGroup := "polaris-config-example"
	fileName := "example.yaml"

	configFile, err := configAPI.GetConfigFile(namespace, fileGroup, fileName)
	if err != nil {
		fmt.Println("fail to get config.", err)
		return
	}

	// 打印配置文件内容
	fmt.Println(configFile.GetContent())

	// 方式一：添加监听器
	configFile.AddChangeListener(changeListener)

	// 方式二：添加监听器
	changeChan := make(chan model.ConfigFileChangeEvent)
	configFile.AddChangeListenerWithChannel(changeChan)

	for {
		select {
		case event := <-changeChan:
			fmt.Println(fmt.Sprintf("received change event by channel. %+v", event))
		}
	}
}

func changeListener(event model.ConfigFileChangeEvent) {
	fmt.Println(fmt.Sprintf("received change event. %+v", event))
}

func demo() {
	polarisConfig, _ := polaris.NewAdapterFile()
	if polarisConfig != nil {

	}
}
