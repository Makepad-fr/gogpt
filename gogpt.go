package gogpt

import (
	"github.com/playwright-community/playwright-go"
	"go.uber.org/zap"
	"log"
	"os"
)

var logger *zap.Logger

func init() {
	err := playwright.Install()
	if err != nil {
		log.Fatal(err)
	}
}

type GoGPT interface {
	Login(username, password string) error
	Ask(question string, version Version)
	History()
	OpenFromHistory(index uint)
	NewChat()
	Debug()
}




func New(browserContextPath string, headless, debug bool) (GoGPT,error) {
	var loadFromBrowserContext bool
	s, err := os.Stat(browserContextPath)
	if err != nil {
		loadFromBrowserContext = false
	} else {
		loadFromBrowserContext = !s.IsDir()
	}
	pw,err := playwright.Run()
	if err != nil {
		return nil, err
	}
	var browserContextPathPtr *string = nil
	if loadFromBrowserContext {
		browserContextPathPtr = &browserContextPath
	}
	var l *zap.Logger
	if debug {
		l, err = zap.NewDevelopment()
		headless = false
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}
	logger = l
	browser,err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless})
	if err != nil {
		return nil, err
	}
	page,err := browser.NewPage(playwright.BrowserNewContextOptions{
		StorageStatePath:  browserContextPathPtr,
	})
	if err != nil {
		return nil, err
	}


	return &gpt{
		browserContextPath: browserContextPath,
		browser: browser,
		page: page,
	}, nil
}

// DumpCookie lets user login to the ChatGPT with headless mode disabled and dumps the browser context to the given browserContextPath string passed in parameters
func DumpCookie(browserContextPath string) error {
	// TODO: Complete function definition
	return nil
}
