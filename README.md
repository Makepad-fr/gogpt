# gogpt
GoGPT is a Go library which let users use their ChatGPT accounts through Go.
With GoGPT you can do everything that you can do with your ChatGPT account in the Web interface.

## Installation

To install, you can use the following command:
```bash
go get -u https://github.com/Makepad-fr/gogpt
```

## Usage

To use the ChatGPT from your Go application, you need to create a new instance by passing the file path for browser context, a boolean for headless, a boolean indicating the debug mode and your current timezone offset.

**IMPORTANT:** Please note that, the timezone offset is an integer in minutes. The value is negative for timezones ahead of UTC and positive for timezones behind UTC.

```go
package main
import (
	"github.com/Makepad-fr/gogpt"
	"log"
)
func main() {
	gpt, err := gogpt.New("./gogpt.json", false, true, -120)
	if err != nil {
		log.Fatal(err)
	}
}   
```

### Login

To do any operation on your ChatGPT account, you need to login to your account first.

#### Interactive Login

TBD

#### Headless login

**IMPORTANT:** Headless login is only available if you're using email and password login.

You can simply call the Login method with your email and password on the gpt instance that you've previously created.

```go
package main
import (
	"github.com/Makepad-fr/gogpt"
	"log"
)
func main() {
	gpt, err := gogpt.New("./gogpt.json", false, true, -120)
	if err != nil {
		log.Fatal(err)
	}
	err = gpt.Login("<YOUR_CHATGPT_EMAIL>", "<YOU_CHATGPT_PASSWORD>")
	if err != nil {
		log.Fatal(err)
	}
}
```

### Sending a message (prompt)

#### By creating a new conversation

You can create a new conversation with an initial prompt using the `CreateConversation` method on the `gpt` instance that you previously created

```go
package main
import (
	"github.com/Makepad-fr/gogpt"
	"log"
)
func main() {
	...
	err = gpt.CreateConversation("<PROMPT>", "<MDDEL_SLUG>", func(response gogpt.ConversationResponse) {
		log.Printf("Received response %+v", response)
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

#### To an existing conversation

TBD

### Getting conversation history

You can get your conversation history using the `History` method on the `gpt` instance that you previously created

```go
package main
... 
func main() {
	...
	conversations, err := gpt.History()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Number of conversations in history is %d", len(conversations))
}
```

### Loading a conversation

You can load the conversation details using `LoadConversation` method and the ID of the conversation that you want to load

```go
package main
...
func main() {
	...
	const conversationId = "<CONVERSATION_ID>"
	conversationDetails, err := gpt.LoadConversation(conversationId)
	if err != nil {
		log.Fatal(err)
	}
}
```

### Get available models

You can get the list of the available models for account using the `Models` method.

```go
package main
...
func main() {
	...
	models, err := gpt.Models()
	if err != nil {
		log.Fatal(err)
	}
}
```





