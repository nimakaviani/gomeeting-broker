package models

type AlexaRequest struct {
	Version string       `json:"version"`
	Session AlexaSession `json:"session"`
	Request AlexaReqBody `json:"request"`
}

type AlexaSession struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes struct {
		String map[string]interface{} `json:"string"`
	} `json:"attributes"`
	User struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

type AlexaReqBody struct {
	Type      string      `json:"type"`
	RequestID string      `json:"requestId"`
	Timestamp string      `json:"timestamp"`
	Intent    AlexaIntent `json:"intent,omitempty"`
	Reason    string      `json:"reason,omitempty"`
}

type AlexaIntent struct {
	Name  string               `json:"name"`
	Slots map[string]AlexaSlot `json:"slots"`
}

type AlexaSlot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Response Types

type AlexaResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          AlexaRespBody          `json:"response"`
}

type AlexaRespBody struct {
	OutputSpeech     *AlexaRespPayload `json:"outputSpeech,omitempty"`
	Card             *AlexaRespPayload `json:"card,omitempty"`
	Reprompt         *AlexaReprompt    `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	ShouldEndSession bool              `json:"shouldEndSession"`
}

type AlexaReprompt struct {
	OutputSpeech AlexaRespPayload `json:"outputSpeech,omitempty"`
}

type AlexaRespImage struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

type AlexaRespPayload struct {
	Type    string         `json:"type,omitempty"`
	Title   string         `json:"title,omitempty"`
	Text    string         `json:"text,omitempty"`
	SSML    string         `json:"ssml,omitempty"`
	Content string         `json:"content,omitempty"`
	Image   AlexaRespImage `json:"image,omitempty"`
}
