package gogpt

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"go.uber.org/zap"
	"log"
	"strings"
)

type gpt struct {
	GoGPT
	browserContextPath string
	browser            playwright.Browser
	page playwright.Page
	user Session
}


//getChallenge returns  a playwright.ElementHandle related to the challenge and an error if there's an error returned by navigate
func (g *gpt) getChallenge() (playwright.ElementHandle, error) {
	err := g.navigate()
	if err != nil {
		return nil,err
	}
	selector, err := g.page.WaitForSelector(challengeDivSelector)
	if err != nil {
		return nil, nil
	}
	return selector, nil
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
	iFrameElementHandle, err := challengeElementHandle.WaitForSelector(iframeSelector)
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

//userNeedsToLogin returns true if the user needs to be logged in by navigating to the default url of ChatGPT
func (g *gpt) userNeedsToLogin() (bool, error) {
	err := g.navigate()
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
	}
	_, err = g.page.WaitForSelector(loginPageTextSelector)
	if err != nil {
		// TODO: Verify if there's no other error -> service unavailable or already logged in
		logger.Debug("Not on the login page. No need to login")
		return false, nil
	}
	return true, nil
}

//Login let you log in to your ChatGPT account using given username and password
func (g *gpt) Login(username, password string) error {
	needLogin, err := g.userNeedsToLogin()
	if err != nil {
		return err
	}
	if needLogin {
		logger.Debug("Session needs to login")
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
			logger.Error("Error while waiting the url changes to logged in URL")
			return err
		}
		g.saveBrowserContexts()
	}
	err = g.passPopupDialog()
	if err != nil {
		return err
	}
	cookies, err := g.page.Context().Cookies(baseURL)
	if err != nil {
		return err
	}
	logger.Debug("Number of cookies ", zap.Int("number-of-cookies", len(cookies)))
	g.user.Cookies = playwrightCookiesToHttpCookies(cookies)
	// TODO: Do not update the session now. Put in a variable then get the session details and create the user at once
	return nil
}

//Session returns the information about the current session
func (g *gpt) Session() Session {
	return g.user
}

//Ask let you ask a new question with the given Version
func (*gpt) Ask(question string, version Version) {

}

//History returns the history of conversations
func (*gpt) History() {

}

//OpenFromHistory let you select a chat from the history with the given index
func (*gpt) OpenFromHistory(index uint) {

}

//NewChat creates a new chat
func (*gpt) NewChat() {

}

//Debug function is only used for debugging purposes, it disables the default behavior of playwright which is closing
//browser and page once the execution is done
func (g *gpt) Debug() {
	g.page.WaitForTimeout(100000000000)
}

//saveBrowserContexts saves the browser contexts of the *gpt to the browserContextPath
func (g *gpt) saveBrowserContexts() error{
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

//getPopupDialog returns the playwright.ElementHandle related to the popupDialog selected by popupDialogSelector
//if there's no pop-up dialog it returns nil
func (g *gpt) getPopupDialog() playwright.ElementHandle{
	elementHandle, err := g.page.WaitForSelector(popupDialogSelector)
	if err != nil {
		logger.Debug("Popup selector does not exists")
		return nil
	}
	logger.Debug("Popup exists returning the element handle")
	return elementHandle
}

//passPopupDialog closes the pop-up dialog if there's any. To avoid that it happens again and again it updates the
//browserContext identified by browserContextPath. If something getc
func (g *gpt) passPopupDialog() error {
	popupDialogElementHandler := g.getPopupDialog()
	if popupDialogElementHandler == nil {
		logger.Debug("There's no popup to pass")
		// If there's nothing to pass, just return
		return nil
	}
	for popupDialogElementHandler != nil {
		last := false
		logger.Debug("Waiting for next button")
		buttonHandle, err := popupDialogElementHandler.WaitForSelector(nextButtonSelector)
		if err != nil {
			logger.Debug("Next button does not exists in popup. Waiting for done button")
			buttonHandle, err = popupDialogElementHandler.WaitForSelector(doneButtonSelector)
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
			break
		}
		logger.Debug("Updating popup element handler")
		popupDialogElementHandler = g.getPopupDialog()
	}
	logger.Debug("Updating the browser context after popup")
	// Update the browser context once the dialog is closed
	return g.saveBrowserContexts()
}


