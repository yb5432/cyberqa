// Package openai provides an OpenAI-compatible API interface.
package openai

// ChatCompletionRequest represents a request to the chat completion API.
// It includes all the parameters that can be sent to the API.
type ChatCompletionRequest struct {
	// Model is the name of the model to use for the completion.
	Model string `json:"model"`

	// Messages is a list of messages comprising the conversation so far.
	Messages []Message `json:"messages"`

	// Temperature is a value between 0 and 2 that controls the randomness of the output.
	// Higher values make the output more random, lower values make it more deterministic.
	Temperature float64 `json:"temperature,omitempty"`

	// TopP is an alternative to temperature for controlling randomness.
	// It specifies the cumulative probability threshold for token selection.
	TopP float64 `json:"top_p,omitempty"`

	// N is the number of chat completion choices to generate for each input message.
	N int `json:"n,omitempty"`

	// Stream indicates whether to stream the response back as it's generated.
	Stream bool `json:"stream,omitempty"`

	// Stop is a list of tokens at which to stop generation.
	Stop []string `json:"stop,omitempty"`

	// MaxTokens is the maximum number of tokens to generate in the chat completion.
	MaxTokens int `json:"max_tokens,omitempty"`

	// PresencePenalty is a value between -2.0 and 2.0 that penalizes new tokens based on whether they appear in the text so far.
	PresencePenalty float64 `json:"presence_penalty,omitempty"`

	// FrequencyPenalty is a value between -2.0 and 2.0 that penalizes new tokens based on their existing frequency in the text so far.
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	// LogitBias is a map that modifies the likelihood of specified tokens appearing in the completion.
	LogitBias map[string]int `json:"logit_bias,omitempty"`

	// User is a unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`
}

// Message represents a single message in a chat conversation.
type Message struct {
	// Role is the role of the message sender (e.g., "system", "user", "assistant").
	Role string `json:"role"`

	// Content is the content of the message.
	Content string `json:"content"`
}

// ChatCompletionResponse represents a response from the chat completion API.
type ChatCompletionResponse struct {
	// ID is a unique identifier for the chat completion.
	ID string `json:"id"`

	// Object is the type of object returned (e.g., "chat.completion").
	Object string `json:"object"`

	// Created is the Unix timestamp (in seconds) of when the chat completion was created.
	Created int64 `json:"created"`

	// Model is the name of the model used to generate the completion.
	Model string `json:"model"`

	// Choices is a list of chat completion choices.
	Choices []Choice `json:"choices"`

	// Usage is the token usage statistics for the completion.
	Usage Usage `json:"usage"`
}

// Choice represents a single choice in a chat completion response.
type Choice struct {
	// Index is the index of the choice in the list of choices.
	Index int `json:"index"`

	// Message is the generated message for this choice.
	Message Message `json:"message"`

	// FinishReason is the reason the model stopped generating tokens.
	// This will be "stop" if the model hit a natural stop point or a provided stop sequence,
	// "length" if the maximum number of tokens specified in the request was reached,
	// or "content_filter" if content was omitted due to a flag from our content filters.
	FinishReason string `json:"finish_reason"`
}

// Usage represents the token usage in a chat completion response.
type Usage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`

	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completion_tokens"`

	// TotalTokens is the total number of tokens used (prompt + completion).
	TotalTokens int `json:"total_tokens"`
}
