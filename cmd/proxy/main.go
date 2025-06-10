// Ligolo-ng
// Copyright (C) 2025 Nicolas Chatelain (nicocha30)

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/nicocha30/ligolo-ng/cmd/proxy/config"
	"github.com/nicocha30/ligolo-ng/pkg/tlsutils"
	"github.com/nicocha30/ligolo-ng/types"

	"github.com/desertbit/grumble"
	"github.com/hashicorp/yamux"
	"github.com/nicocha30/ligolo-ng/cmd/proxy/app"
	"github.com/nicocha30/ligolo-ng/pkg/controller"
	"github.com/sirupsen/logrus"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func Proxy(flags *types.ProxyFlags) {
	var allowDomains []string

	flag.Usage = func() {
		fmt.Printf("Ligolo-ng %s / %s / %s\n", version, commit, date)
		fmt.Println("Made in France with love by @Nicocha30!")
		fmt.Println("https://github.com/nicocha30/ligolo-ng")
		fmt.Printf("\nUsage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flags.CPUProfile != "" {
		f, err := os.Create(*flags.CPUProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	config.InitConfig(*flags.ConfigFile)

	if *flags.VersionFlag {
		fmt.Printf("Ligolo-ng %s / %s / %s\n", version, commit, date)
		return
	}

	if *flags.VersionFlag {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// If verbose is set, include method and line in log messages
	logrus.SetReportCaller(*flags.VersionFlag)

	if *flags.DomainWhitelist != "" {
		allowDomains = strings.Split(*flags.DomainWhitelist, ",")
	}

	if !*flags.HideBanner && !*flags.DaemonMode {
		app.App.SetPrintASCIILogo(func(a *grumble.App) {
			a.Println("    __    _             __                       ")
			a.Println("   / /   (_)___ _____  / /___        ____  ____ _")
			a.Println("  / /   / / __ `/ __ \\/ / __ \\______/ __ \\/ __ `/")
			a.Println(" / /___/ / /_/ / /_/ / / /_/ /_____/ / / / /_/ / ")
			a.Println("/_____/_/\\__, /\\____/_/\\____/     /_/ /_/\\__, /  ")
			a.Println("        /____/                          /____/   ")
			a.Println("\n  Made in France ♥            by @Nicocha30!")
			a.Printf("  Version: %s\n\n", version)
		})
	}

	if *flags.EnableSelfcert && *flags.SelfcertDomain == "ligolo" {
		logrus.Warning("Using default selfcert domain 'ligolo', beware of CTI, SOC and IoC!")
	}

	app.Run()

	proxyController := controller.New(controller.ControllerConfig{
		Address: *flags.ListenInterface,
		CertManagerConfig: &tlsutils.CertManagerConfig{
			SelfCertCache:   "ligolo-selfcerts",
			Certfile:        *flags.CertFile,
			Keyfile:         *flags.KeyFile,
			DomainWhitelist: allowDomains,
			SelfcertDomain:  *flags.SelfcertDomain,
			EnableAutocert:  *flags.EnableAutocert,
			EnableSelfcert:  *flags.EnableSelfcert,
		},
	})

	app.ProxyController = &proxyController

	go proxyController.ListenAndServe()

	// Wait for listener
	if err := proxyController.WaitForReady(); err != nil {
		logrus.Fatal(err)
	}

	// Agent registration goroutine
	go func() {
		for {
			remoteConn := <-proxyController.Connection

			yamuxConn, err := yamux.Client(remoteConn, nil)
			if err != nil {
				logrus.Errorf("could not create yamux client, error: %v", err)
				continue
			}

			agent, err := controller.NewAgent(yamuxConn)
			if err != nil {
				logrus.Errorf("could not register agent, error: %v", err)
				continue
			}

			logrus.WithFields(logrus.Fields{"remote": remoteConn.RemoteAddr(), "name": agent.Name, "id": agent.SessionID}).Info("Agent joined.")

			if err := app.RegisterAgent(agent); err != nil {
				logrus.Errorf("could not register agent: %s", err.Error())
			}

			go func() {
				// Check agent status
				for {
					select {
					case <-agent.Session.CloseChan(): // Agent closed
						logrus.WithFields(logrus.Fields{"remote": remoteConn.RemoteAddr(), "name": agent.Name, "id": agent.SessionID}).Warnf("Agent dropped.")
						return
					}
				}
			}()

		}
	}()

	if *flags.DaemonMode && !config.Config.GetBool("web.enabled") {
		logrus.Warning("daemon mode enabled but web.enabled is false!")
	}

	if config.Config.GetBool("web.enabled") {
		logrus.Infof("Starting Ligolo-ng Web, API URL is set to: %s", app.GetAPIUrl())
		go app.StartLigoloApi()
	}

	if *flags.DaemonMode {
		proxyController.WaitForFinished()
	} else {
		// Grumble doesn't like cli args
		os.Args = []string{}
		grumble.Main(app.App)
	}

	if *flags.MemProfile != "" {
		f, err := os.Create(*flags.MemProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

}
