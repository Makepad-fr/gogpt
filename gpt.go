package gogpt

import (
	"errors"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"go.uber.org/zap"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

type gpt struct {
	GoGPT
	browserContextPath  string
	browser             playwright.Browser
	page                playwright.Page
	session             *Session
	httpClient          *http.Client
	cookieJar           *autoFillingCookieJar
	accountInfo         *UserAccountInfo
	username, password  *string
	popupPassed         bool
	conversationHistory *idBasedSet[ConversationHistoryItem]
	availableModels     []string
	timeZoneOffset      int
	timeout             *float64
}

// getChallenge returns  a playwright.ElementHandle related to the challenge and an error if there's an error returned by navigate
func (g *gpt) getChallenge() (playwright.ElementHandle, error) {
	err := g.navigate()
	if err != nil {
		return nil, err
	}
	selector, err := g.page.WaitForSelector(challengeDivSelector, playwright.PageWaitForSelectorOptions{Timeout: g.timeout})
	if err != nil {
		return nil, nil
	}
	return selector, nil
}

// Close closes the open page and the browser window
func (g *gpt) Close() error {
	err := g.page.Close()
	if err != nil {
		return err
	}
	return g.browser.Close()
}

// navigate goes to the baseURL
func (g *gpt) navigate() error {
	if strings.HasPrefix(g.page.URL(), baseURL) {
		logger.Debug("No need to navigate", zap.String("current-url", g.page.URL()))
		return nil
	}
	logger.Debug("Navigating to the base url")
	_, err := g.page.Goto(baseURL)
	if err != nil {
		logger.Error("Error while navigating to the default URL")
		return err
	}
	return nil
}

// solveChallenge solves the challenge in the login screen
func (g *gpt) solveChallenge(challengeElementHandle playwright.ElementHandle) error {
	iFrameElementHandle, err := challengeElementHandle.WaitForSelector(iframeSelector, playwright.ElementHandleWaitForSelectorOptions{Timeout: g.timeout})
	if err != nil {
		logger.Error("iframeSelector does not exists in challengeElementHandle")
		return err
	}

	iFrame, err := iFrameElementHandle.ContentFrame()
	if err != nil {
		return err
	}
	locator, err := iFrame.Locator(checkboxSelector)
	if err != nil {
		return err
	}
	rto := randomTimeOut()
	log.Printf("Will wait for %fms\n", rto)
	g.page.WaitForTimeout(rto)
	err = locator.Click()
	if err != nil {
		return err
	}
	return nil
}

// userNeedsToLogin returns true if the user needs to be logged in by navigating to the default url of ChatGPT
func (g *gpt) userNeedsToLogin() (bool, error) {
	err := g.navigate()
	if g.page.URL() == fmt.Sprintf("%s/chat", baseURL) {
		logger.Debug("Already on the application page by the URL. No need to login")
		return false, nil
	}
	if err != nil {
		return true, err
	}
	challengeElement, err := g.getChallenge()
	if err != nil {
		logger.Error("Error while getting challenge element")
		return true, err
	}
	if challengeElement != nil {
		err := g.solveChallenge(challengeElement)
		if err != nil {
			logger.Error("Error while solving challenge")
			return true, err
		}
		err = g.saveBrowserContexts()
		if err != nil {
			return true, err
		}
	}
	_, err = g.page.WaitForSelector(loginPageTextSelector, playwright.PageWaitForSelectorOptions{Timeout: g.timeout})
	if err != nil {
		// TODO: Verify if there's no other error -> service unavailable or already logged in
		logger.Debug("Not on the login page. No need to login")
		return false, nil
	}
	return true, nil
}

// Login let you log in to your ChatGPT account using given username and password
func (g *gpt) Login(username, password string) error {
	err := g.internalLogin(username, password)
	if err != nil {
		return err
	}
	err = g.initCookieJarAndHttpClient()
	if err != nil {
		return err
	}
	err = g.initSession()
	if err != nil {
		return err
	}
	err = g.initUserAccountInfo()
	if err != nil {
		return err
	}
	err = g.initAvailableModels()
	if err != nil {
		return err
	}
	return nil
}

// initAvailableModels initialises availableModels in the current gpt instance.
func (g *gpt) initAvailableModels() error {
	modelInfo, err := g.Models()
	if err != nil {
		return err
	}
	for _, mi := range modelInfo {
		g.availableModels = append(g.availableModels, mi.Slug)
	}
	return nil
}

// Session returns the information about the current session
func (g *gpt) Session() Session {
	return *g.session
}

// initUserAccountInfo initialises the account information for the current user
func (g *gpt) initUserAccountInfo() error {
	accountInfo, err := g.getAccountInfo()
	if err != nil {
		return err
	}
	g.accountInfo = accountInfo
	return nil
}

// Ask let you ask a new question with the given Version
func (*gpt) Ask(question string, version Version) {

}

// updateConversationHistory updates the conversationHistory attribute of the current gpt instance
func (g *gpt) updateConversationHistory(maxRetry, limit uint) error {
	logger.Debug("Make a first request to get the total number of conversations")
	response, err := g.getConversationHistory(0, limit)
	if err != nil {
		logger.Error("Error while getting user's conversations")
		return err
	}
	g.conversationHistory.addAll(response.Items)
	logger.Debug("Items added ", zap.Int("number-of-items", g.conversationHistory.size()))
	var attempts uint = 0
	for g.conversationHistory.size() < response.Total && attempts < maxRetry {
		// Wait for a random timeout
		randomTimeOut := randomTimeOut()
		logger.Debug("Waiting for ", zap.Float64("timeout", randomTimeOut))
		time.Sleep(time.Duration(randomTimeOut) * time.Millisecond)
		response, err = g.getConversationHistory(uint(g.conversationHistory.size()), uint(math.Min(float64(limit), float64(response.Total-g.conversationHistory.size()))))
		before := g.conversationHistory.size()
		if err != nil {
			return err
		}
		g.conversationHistory.addAll(response.Items)
		logger.Debug("Items added ", zap.Int("number-of-items", g.conversationHistory.size()))
		if before == g.conversationHistory.size() {
			logger.Warn("The call was not bring any result", zap.Int("size", g.conversationHistory.size()), zap.Uint("attempts-count", attempts))
			/* For some reason the total number does not always match with the reel number of conversations.
			While this solution needs to be investigated more carefully, for now we are counting a number of unsuccessful
			attempts and stop on maxAttempts number of attempts
			*/
			attempts++
		}

	}
	return nil
}

// History returns the history of conversations as a slice of ConversationHistoryItem
func (g *gpt) History() ([]ConversationHistoryItem, error) {
	err := g.updateConversationHistory(5, 100)
	if err != nil {
		return nil, err
	}
	// Return the created items
	return g.conversationHistory.content, nil
}

// LoadConversation loads a conversation from the chat history using the conversation uuid
func (g *gpt) LoadConversation(uuid string) (*Conversation, error) {
	element := g.conversationHistory.find(uuid)
	if element == nil {
		logger.Warn("Can not find conversation, we'll try to reload the conversation history and re-try", zap.String("conversation-uuid", uuid))
		err := g.updateConversationHistory(5, 100)
		if err != nil {
			return nil, err
		}
		// once the conversation history is loaded, re-try to find the related element
		element = g.conversationHistory.find(uuid)
		if element == nil {
			return nil, fmt.Errorf("can not find conversation with uuid %s", uuid)
		}
	}
	return g.getConversation(uuid)

}

// NewChat creates a new chat
func (*gpt) NewChat() {

}

// Debug function is only used for debugging purposes, it disables the default behavior of playwright which is closing
// browser and page once the execution is done
func (g *gpt) Debug() {
	g.page.WaitForTimeout(100000000000)
}

// saveBrowserContexts saves the browser contexts of the *gpt to the browserContextPath
func (g *gpt) saveBrowserContexts() error {
	contexts := g.browser.Contexts()
	if len(contexts) > 1 {
		logger.Fatal("Multiple contexts contexts detected", zap.Int("length", len(contexts)))
	}
	logger.Debug("Updating browser context", zap.String("path", g.browserContextPath))
	_, err := contexts[0].StorageState(g.browserContextPath)
	if err != nil {
		logger.Error("Something went wrong when saving the browser context", zap.String("path", g.browserContextPath))
		return err
	}
	logger.Debug("Browser context updated", zap.String("path", g.browserContextPath))
	return nil
}

// getPopupDialog returns the playwright.ElementHandle related to the popupDialog selected by popupDialogSelector
// if there's no pop-up dialog it returns nil
func (g *gpt) getPopupDialog() playwright.ElementHandle {
	elementHandle, err := g.page.WaitForSelector(popupDialogSelector, playwright.PageWaitForSelectorOptions{Timeout: g.timeout})
	if err != nil {
		logger.Debug("Popup selector does not exists")
		return nil
	}
	logger.Debug("Popup exists returning the element handle")
	return elementHandle
}

// passPopupDialog closes the pop-up dialog if there's any. To avoid that it happens again and again it updates the
// browserContext identified by browserContextPath. If something getc
func (g *gpt) passPopupDialog() error {
	if g.popupPassed {
		logger.Debug("Pop-up already passed no need to repass")
		return nil
	}
	popupDialogElementHandler := g.getPopupDialog()
	if popupDialogElementHandler == nil {
		logger.Debug("There's no popup to pass")
		g.popupPassed = true
		// If there's nothing to pass, just return
		return nil
	}
	for popupDialogElementHandler != nil {
		last := false
		logger.Debug("Waiting for next button")
		buttonHandle, err := popupDialogElementHandler.WaitForSelector(nextButtonSelector, playwright.ElementHandleWaitForSelectorOptions{Timeout: g.timeout})
		if err != nil {
			logger.Debug("Next button does not exists in popup. Waiting for done button")
			buttonHandle, err = popupDialogElementHandler.WaitForSelector(doneButtonSelector, playwright.ElementHandleWaitForSelectorOptions{Timeout: g.timeout})
			if err != nil {
				logger.Error("Either next button or done button should appear. Something is probably wrong")
				return err
			}
			last = true
			logger.Debug("Done button appeared")
		}
		logger.Debug("Will click on button")
		err = buttonHandle.Click()
		if err != nil {
			logger.Error("Something went wrong while clicking on button inside of the pop-up")
			return err
		}
		if last {
			logger.Debug("Dialog passed")
			g.popupPassed = true
			break
		}
		logger.Debug("Updating popup element handler")
		popupDialogElementHandler = g.getPopupDialog()
	}
	logger.Debug("Updating the browser context after popup")
	// Update the browser context once the dialog is closed
	return g.saveBrowserContexts()
}

// getUserCookiesSupplier creates a httpCookieSupplier for the given url string passed in parameters
func (g *gpt) getUserCookiesSupplier(u string) httpCookieSupplier {
	return func() ([]*http.Cookie, error) {
		loginNeeded, err := g.userNeedsToLogin()
		if err != nil {
			return nil, err
		}
		if loginNeeded {
			// If user needs to log in
			// Check if both username and password are provided
			if g.username == nil || g.password == nil {
				return nil, errors.New("can generate cookies as the username or password is not provided and user needs to be logged in")
			}
			err = g.internalLogin(*g.username, *g.password)
			if err != nil {
				return nil, err
			}
		} else {
			// If user does not need to log in pass te pop-up dialog if applicable
			err = g.passPopupDialog()
			if err != nil {
				return nil, err
			}
		}
		// Get cookies for the given url string
		cookies, err := g.page.Context().Cookies(u)
		if err != nil {
			return nil, err
		}
		// Convert playwright.BrowserContextCookiesResult to http.Cookie
		httpCookies := playwrightCookiesToHttpCookies(cookies)
		// Return them
		return httpCookies, nil
	}
}

// internalLogin just handles the login with the given username and password without any side effects
func (g *gpt) internalLogin(username, password string) error {
	needLogin, err := g.userNeedsToLogin()
	if err != nil {
		return err
	}
	if needLogin {
		logger.Debug("User needs to login")
		err := g.page.Click(loginButtonSelector)
		if err != nil {
			logger.Error("Error while clicking on login button selector")
			return err
		}
		err = g.page.Fill(usernameInputSelector, username)
		if err != nil {
			logger.Error("Error while filling username input")
			return err
		}
		err = g.page.Click(continueButtonSelector)
		if err != nil {
			logger.Error("Error while clicking on continue button")
			return err
		}
		err = g.page.Fill(passwordInputSelector, password)
		if err != nil {
			logger.Error("Error while filling the password input")
			return err
		}
		err = g.page.Click(continueButtonSelector)
		if err != nil {
			logger.Error("Error while clicking on continue button")
			return err
		}
		err = g.page.WaitForURL(fmt.Sprintf("%s/chat", baseURL))
		if err != nil {
			logger.Error("Error while waiting the u changes to logged in URL")
			return err
		}
		err = g.saveBrowserContexts()
		if err != nil {
			return err
		}
		// Login successful save login information
		g.username = &username
		g.password = &password
	}
	err = g.passPopupDialog()
	if err != nil {
		return err
	}
	return nil
}

// AccountInfo returns the UserAccountInfo instance related to the current user account
func (g *gpt) AccountInfo() UserAccountInfo {
	return *g.accountInfo
}

// Models returns the list of available models for the userr
func (g *gpt) Models() ([]ModelInfo, error) {
	modelResponses, err := g.getModels()
	if err != nil {
		return nil, err
	}
	return modelResponses.Models, err
}

func (g *gpt) isModelExists(model string) bool {
	for _, availableModelSlug := range g.availableModels {
		if model == availableModelSlug {
			return true
		}
	}
	return false
}

// CreateConversation creates a new conversation by sending the given message and using the given model. For each response
// received by the ChatGPT, it calls the onResponse callback with the received response as CovnersationResponse.
func (g *gpt) CreateConversation(message, model string, onResponse conversationResponseConsumer) error {
	// TODO: Check if the given model is in the models array (to add in gpt structure)
	if !g.isModelExists(model) {
		return fmt.Errorf("%s is not a valid model", model)
	}
	// TODO: Add callbacks to handle the event
	return g.sendMessageToNewConversation(message, model, onResponse)

}
