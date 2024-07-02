package main

type T struct {
	Model             string   `json:"model"`
	Messages          []MSG    `json:"messages"`
	Stream            *bool    `json:"stream"`
	FrequencyPenalty  *float32 `json:"frequency_penalty"`
	MaxTokens         *float32 `json:"max_tokens"`
	MinP              *float32 `json:"min_p"`
	PresencePenalty   *float32 `json:"presence_penalty"`
	RepetitionPenalty *float32 `json:"repetition_penalty"`
	Temperature       *float32 `json:"temperature"`
	TopA              *float32 `json:"top_a"`
	TopK              *float32 `json:"top_k"`
	TopP              *float32 `json:"top_p"`
	Tools             *[]Tool  `json:"tools"`
}

type Tool struct {
	Type     string      `json:"type"`
	Function llmFunction `json:"function"`
}

type llmFunction struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Parameters  llmFunctionParam `json:"parameters"`
}

type llmFunctionParam struct {
	Type       string                `json:"type"`
	Properties llmFunctionProperties `json:"properties"`
	Required   []string              `json:"required"`
}

type llmFunctionProperties struct {
	Query `json:"query"`
}

type Query struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type AIResponse struct {
	Id      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   TokUsage `json:"usage"`
}

type Choice struct {
	Index        int    `json:"index"`
	Message      *MSG   `json:"message,omitempty"`
	Delta        *Delta `json:"delta,omitempty"`
	FinishReason string `json:"finish_reason"`
}

type MSG struct {
	Role      string       `json:"role"`
	Name      *string      `json:"name,omitempty"`
	Content   interface{}  `json:"content"`
	ToolCalls *[]Tool_Call `json:"tool_calls,omitempty"`
	ToolID    *string      `json:"tool_call_id,omitempty"`
}

type Tool_Call struct {
	Id       string        `json:"id"`
	Type     string        `json:"type"`
	Function Function_Call `json:"function"`
}

type Function_Call struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type TokUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type Delta struct {
	Role      string       `json:"role"`
	Content   interface{}  `json:"content"`
	ToolCalls *[]Tool_Call `json:"tool_calls,omitempty"`
}

type Chunk struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
