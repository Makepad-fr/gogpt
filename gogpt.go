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
	History() ([]ConversationHistoryItem, error)
	AccountInfo() UserAccountInfo
	LoadConversation(uuid string) (*Conversation, error)
	Close() error
	NewChat()
	Session() Session
	Models() ([]ModelInfo, error)
	Debug()
	CreateConversation(message, model string, onResponseCallback conversationResponseConsumer) error
}

type Options struct {
	BrowserContextPath string
	Headless           bool
	Debug              *bool
	TimeZoneOffset     int
	Timeout            *float64
}

func New(options Options) (GoGPT, error) {
	var loadFromBrowserContext bool
	s, err := os.Stat(options.BrowserContextPath)
	if err != nil {
		loadFromBrowserContext = false
	} else {
		loadFromBrowserContext = !s.IsDir()
	}
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	var browserContextPathPtr *string = nil
	if loadFromBrowserContext {
		browserContextPathPtr = &options.BrowserContextPath
	}
	var l *zap.Logger
	if options.Debug != nil && *options.Debug {
		l, err = zap.NewDevelopment()
		options.Headless = false
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}
	logger = l
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{Headless: &options.Headless})
	if err != nil {
		return nil, err
	}
	page, err := browser.NewPage(playwright.BrowserNewContextOptions{
		StorageStatePath: browserContextPathPtr,
	})
	if err != nil {
		return nil, err
	}

	return &gpt{
		browserContextPath:  options.BrowserContextPath,
		browser:             browser,
		page:                page,
		session:             nil,
		popupPassed:         false,
		conversationHistory: newIdBasedSet[ConversationHistoryItem](100),
		availableModels:     []string{},
		timeZoneOffset:      options.TimeZoneOffset,
		timeout:             options.Timeout,
	}, nil
}

// DumpCookie lets session login to the ChatGPT with headless mode disabled and dumps the browser context to the given browserContextPath string passed in parameters
func DumpCookie(browserContextPath string) error {
	// TODO: Complete function definition
	return nil
}
