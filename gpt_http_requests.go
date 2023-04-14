package gogpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

// initCookieJarAndHttpClient initialises the autoFillingCookieJar and http.Client instances inside the current *gpt instance
func (g *gpt) initCookieJarAndHttpClient() error {
	if g.cookieJar == nil {
		cookieJar, err := createNewAutoFillingCookieJar(baseURL, g.getUserCookiesSupplier(baseURL))
		if err != nil {
			return err
		}
		g.cookieJar = cookieJar
		g.httpClient = &http.Client{
			Jar: g.cookieJar,
		}
		return nil
	}
	if g.httpClient == nil {
		g.httpClient = &http.Client{
			Jar: g.cookieJar,
		}
	}
	return nil
}

// initSession initializes the session of the current gpt instance.
// It returns an error if something goes wrong while unmarshalling the api response
func (g *gpt) initSession() error {
	err := g.cookieJar.setExpiredCookies()
	if err != nil {
		return err
	}
	resp, err := g.httpClient.Get("https://chat.openai.com/api/auth/session")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Debug("session token response", zap.String("session-request-response", string(body)))
	s, err := unmarshalGPTSessionResponseJSON(body)
	if err != nil {
		logger.Error("Error while unmarshalling session response")
		return err
	}
	g.session = s

	return nil
}

// refreshSession verifies if there's a session exists. If there's no session exists, creates one using initSession
// if there's an existing session verifies if the session is expired using isExpired function. If the session is expired
// recreates the session using initSession
func (g *gpt) refreshSession() error {
	if g.session == nil {
		return g.initSession()
	}
	isExpired, err := g.session.isExpired()
	if err != nil {
		logger.Error("Error while checking if the existing session is expired", zap.String("expiration-date-string", g.session.Expires))
		return g.initSession()
	}
	if isExpired {
		return g.initSession()
	}
	return nil
}

// prepareRequest prepares the cookies and the user session to use in each http request. This function should be called
// before each http request to ensure that the request will not be blocked
func (g *gpt) prepareRequest() error {
	err := g.initCookieJarAndHttpClient()
	if err != nil {
		return err
	}
	err = g.cookieJar.setExpiredCookies()
	if err != nil {
		return err
	}
	err = g.refreshSession()
	if err != nil {
		return err
	}
	return nil
}

// createAPIURL creates the API url for the given endpoint
func createAPIURL(endpoint string) string {
	return fmt.Sprintf("%s/backend-api/%s", baseURL, endpoint)
}

// createRequest creates a new http.Request using given method, endpoint and body.
func (g *gpt) createRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {
	err := g.prepareRequest()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, createAPIURL(endpoint), body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("authorization", fmt.Sprintf("Bearer %s", g.session.AccessToken))
	request.Header.Set("content-type", "application/json")
	return request, nil
}

// runAPIRequest makes an HTTP request with given method on the given endpoint with the given requestBody as io.Reader.
// It handles the response as JSON and unmarshal it to the parameterized type
func runAPIRequest[T any](g *gpt, method, endpoint string, requestBody io.Reader) (*T, error) {
	request, err := g.createRequest(method, endpoint, requestBody)
	if err != nil {
		return nil, err
	}
	resp, err := g.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("run http %s request on  %s with %+v failed with status code %d", method, endpoint, body, resp.StatusCode)
	}
	var response T
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// getConversationHistory returns the history of conversation using given offset and limit as ConversationHistoryResponse
func (g *gpt) getConversationHistory(offset, limit uint) (*ConversationHistoryResponse, error) {
	return runAPIRequest[ConversationHistoryResponse](g, "GET", fmt.Sprintf("conversations?offset=%d&limit=%d", offset, limit), nil)
}

// getAccountInfo returns the additional information about the user's account as UserAccountInfo pointer
func (g *gpt) getAccountInfo() (*UserAccountInfo, error) {
	return runAPIRequest[UserAccountInfo](g, "GET", "accounts/check", nil)
}

// getConversation get the details of a conversation by its uuid as Conversation pointer
func (g *gpt) getConversation(uuid string) (*Conversation, error) {
	return runAPIRequest[Conversation](g, "GET", fmt.Sprintf("conversation/%s", uuid), nil)
}

// getModels returns the available models as ModelsResponse
func (g *gpt) getModels() (*ModelsResponse, error) {
	return runAPIRequest[ModelsResponse](g, "GET", "models", nil)
}

// sendMessageToNewConversation creates a new conversation by sending the given message and using the given model.
// for each response event it calls onResponse function to handle the response as ConversationResponse
func (g *gpt) sendMessageToNewConversation(message, model string, onResponse conversationResponseConsumer) ([]byte, error) {
	messageRequest, err := createMessageRequestForNewConversation(message, model, g.timeZoneOffset)
	if err != nil {
		return nil, err
	}
	requestBody, err := json.Marshal(*messageRequest)
	if err != nil {
		return nil, err
	}

	request, err := g.createRequest("POST", "conversation", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("DNT", "1")
	request.Header.Set("Origin", "https://chat.openai.com")
	request.Header.Set("Referer", "https://chat.openai.com/")
	request.Header.Set("Sec-Fetch-Dest", "empty")
	request.Header.Set("Sec-Fetch-Mode", "cors")
	request.Header.Set("Sec-Fetch-Site", "same-site")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.146 Safari/537.36")

	resp, err := g.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		reader, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		logger.Error("Send message to new conversation is failed", zap.Int("status-code", resp.StatusCode),
			zap.String("body", string(reader)), zap.String("url", request.URL.String()))
		return nil, fmt.Errorf("send message to new conversation is failed. Status code %d", resp.StatusCode)
	}
	// Read and process the events
	reader := bufio.NewReader(resp.Body)
	conversationId, err := g.handleConversationResponseEvent(reader, onResponse)
	if err != nil {
		return nil, err
	}
	return conversationId, nil
}

// handleConversationResponseEvent handles the conversation response events as *bufio.Reader using the given conversationResponseConsumer function
func (g *gpt) handleConversationResponseEvent(reader *bufio.Reader, onResponse conversationResponseConsumer) ([]byte, error) {
	var conversationId string = ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				// Process the event
				logger.Debug("EOF detected", zap.String("current-line", line))
				break
			}
			logger.Error("Error while handling event", zap.Any("error", err))
			return nil, err
		}
		if isEmpty(line) {
			logger.Debug("Empty line received", zap.String("line", line))
			// If the line is empty do nothing just continue
			continue
		}
		if isEndOfEventStream(line) {
			logger.Debug("End of the event stream received", zap.String("line", line))
			// If the line is the end of the event stream quit the loop
			break
		}
		if strings.HasPrefix(line, "data: ") {

			trimmedLine := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
			var response ConversationResponse
			err := json.Unmarshal([]byte(trimmedLine), &response)
			if err != nil {
				logger.Error("Error while parsing conversation response to ConversationRequestEvent",
					zap.String("line", line))
				return nil, err
			}
			if isEmpty(conversationId) {
				conversationId = response.ConversationID
			} else {
				if conversationId != response.ConversationID {
					logger.Warn("THe conversation id is different then the current one", zap.String("current-conversation-id", response.ConversationID), zap.String("existing-conversation-id", conversationId))
				}
			}
			if response.Message.Author.Role == "user" {
				title, err := g.GenerateTitle(response.ConversationID, response.Message.ID)
				if err != nil {
					logger.Error("Error while generating title")
					return nil, err
				}
				logger.Info("Title generated for the new conversation", zap.ByteString("title", title))
				continue
			}
			if response.Message.Author.Role == "assistant" {
				onResponse(response)
				if response.Message.EndTurn != nil && *response.Message.EndTurn {
					logger.Debug("Received the last message", zap.Any("response", response))
					// If the response indicates the end, quit the loop
					break
				}
			}

		}
		// Process the event
		logger.Debug("Received event", zap.String("current-line", line))
	}
	return []byte(conversationId), nil
}

// GenerateTitle generates the title for the given conversation and given message. It returns the generated title as
// []byte
func (g *gpt) GenerateTitle(conversationId, messageId string) ([]byte, error) {
	// https://chat.openai.com/backend-api/conversation/gen_title/b3a28cf1-7f6e-450e-b08a-dab473549383
	requestBody, err := json.Marshal(GenerateConversationTitleRequestBody{MessageId: messageId})
	if err != nil {
		return nil, err
	}
	endPoint := fmt.Sprintf("conversation/gen_title/%s", conversationId)
	response, err := runAPIRequest[GenerateConversationTitleResponse](g, http.MethodPost, endPoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	return []byte(response.Title), nil
}
