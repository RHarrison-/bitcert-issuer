package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Config ...
type Config struct {
	OutDir  *string  `json:"outDir"`
	Wallets []Wallet `json:"wallets"`
}

// Vars
var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	testnet = flag.Bool("testnet", false, "enables use of bitcoin testnet")
	w       *astilectron.Window
)

func main() {
	// Init
	flag.Parse()
	astilog.FlagInit()

	if *testnet {
		fmt.Println("testnet flag active")
	}

	var config = Config{}

	conf, _ := json.Marshal(config)
	tempDir := (os.TempDir() + "/bitcerts")
	// Build temp files if not existing
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, os.ModePerm)
		err := ioutil.WriteFile(tempDir+"/config", []byte(conf), 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

	if _, err := os.Stat(tempDir + "/config.json"); os.IsNotExist(err) {
		err := ioutil.WriteFile(tempDir+"/config", []byte(conf), 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

	configData, _ := ioutil.ReadFile(tempDir + "/config.json")

	conf, _ = json.Marshal(configData)
	json.Unmarshal(conf, &config)

	

	if len(config.Wallets) > 0 {
		fmt.Println("Importing wallets")
		importWallets(config.Wallets)
	} else {
		// only for test application
		var defaults = []Wallet{}
		defaults = append(defaults, Wallet{
			address: "mzq1guetvcVB1Cqpu7UNXb7Y3Kcme47z7x",
			key: "cN6Ui3NpUcqPRZZ4FWEb8LRcDv2Mms1i2NZ8RE1VNdV8vqhxArx2",
		})
		defaults = append(defaults, Wallet{
			address: "mzto4erzfnExogdNMzG4mUmJ1dwFxAMuQb",
			key: "cRQnhsQ2wDYmwo2w4fq2pbYb7hdTnf3MW5hasAr2uGrFKhdkAvw4",
		})

		fmt.Println("Import Default Wallets")

		importWallets(defaults)
	}

	fmt.Println("-----------------")
	fmt.Println(Wallets)
	fmt.Println("-----------------")


	net.connect(*testnet)

	// run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset: Asset,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:          *debug,
		Homepage:       "index.html",
		MessageHandler: handleMessages, // message handler
		RestoreAssets:  RestoreAssets,
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#333"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(600),
			Width:           astilectron.PtrInt(1000),
			Frame:           astilectron.PtrBool(false),
		},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

// MenuOptions: []*astilectron.MenuItemOptions{{
// 	Label: astilectron.PtrStr("File"),
// 	SubMenu: []*astilectron.MenuItemOptions{
// 		{
// 			Label: astilectron.PtrStr("About"),
// 			OnClick: func(e astilectron.Event) (deleteListener bool) {
// 				if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
// 					// Unmarshal payload
// 					var s string
// 					if err := json.Unmarshal(m.Payload, &s); err != nil {
// 						astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
// 						return
// 					}
// 					astilog.Infof("About modal has been displayed and payload is %s!", s)
// 				}); err != nil {
// 					astilog.Error(errors.Wrap(err, "sending about event failed"))
// 				}
// 				return
// 			},
// 		},
// 		{Role: astilectron.MenuItemRoleClose},
// 	},
// }},
// OnWait: func(_ *astilectron.Astilectron, iw *astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
// 	w = iw
// 	go func() {
// 		time.Sleep(5 * time.Second)
// 		if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
// 			astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
// 		}
// 	}()
// 	return nil
// },
