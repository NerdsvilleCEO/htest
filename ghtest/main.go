/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 03-26-2016

* Last Modified : Tue Mar 29 18:00:48 2016

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"flag"
	"fmt"
	"github.com/kiyor/golib"
	"github.com/kiyor/htest"
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
	"runtime"
)

var (
	flagConfig        *string = flag.String("c", "./config", "config file or path")
	flagTemplate      *string = flag.String("t", "./template", "template file or path")
	flagExampleConfig *bool   = flag.Bool("example", false, "generate example config")
	flagIp            *string = flag.String("ip", "", "testing ip")
	flagVerbose       *bool   = flag.Bool("vv", false, "verbose output")
	flagRaw           *bool   = flag.Bool("raw", false, "raw response output")
	flagMethod        *bool   = flag.Bool("method", false, "check support compare method")
	flagTimeout       *string = flag.String("timeout", "", "overwrite timeout value")
	flagCheckOnly     *bool   = flag.Bool("check", false, "check config file only")
	flagNewConfig     *string = flag.String("new", "http://a.com/b", "create new config")
	flagCurl          *bool   = flag.Bool("curl", false, "output curl command")
)

func init() {
	flag.Parse()
	htest.Logger = golib.NewLogger(&golib.LogOptions{
		Name:      "htest",
		ShowErr:   true,
		ShowDebug: *flagVerbose,
		ShowColor: true,
	})
	htest.Verbose = *flagVerbose
	htest.RawResp = *flagRaw
	htest.ShowCurl = *flagCurl

	if *flagMethod {
		fmt.Println(htest.SupportMethod())
		os.Exit(0)
	}
	if *flagExampleConfig {
		configs := htest.ExampleConfig()
		d, _ := yaml.Marshal(configs)
		fmt.Println(string(d))
		os.Exit(0)
	}
	if *flagNewConfig != "http://a.com/b" {
		configs := htest.ExampleConfig()
		u, _ := url.Parse(*flagNewConfig)
		configs[0].Hash = *flagNewConfig
		configs[0].Request.Scheme = u.Scheme
		configs[0].Request.Uri = u.RequestURI()
		configs[0].Request.Hostname = u.Host
		d, _ := yaml.Marshal(configs)
		fmt.Println(string(d))
		os.Exit(0)
	}
	if *flagCheckOnly {
		htest.VerifyYaml(*flagTemplate, true)
		htest.VerifyYaml(*flagConfig, false)
		os.Exit(0)
	}

	htest.OverTimeout = *flagTimeout
	runtime.GOMAXPROCS(runtime.NumCPU())

	htest.LoadTemplate(*flagTemplate)
}

func main() {
	results := make(chan *htest.Result)
	go func() {
		for {
			select {
			case r := <-results:
				fmt.Println(r.String())
			}
		}
	}()
	htest.DoCheck(*flagConfig, results, *flagIp)
}
