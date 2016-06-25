package models

import "encoding/json"

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

func NewAlexaResponse() *AlexaResponse {
	er := &AlexaResponse{
		Version: "1.0",
		Response: AlexaRespBody{
			ShouldEndSession: true,
		},
	}

	return er
}

func (this *AlexaResponse) OutputSpeech(text string) *AlexaResponse {
	this.Response.OutputSpeech = &AlexaRespPayload{
		Type: "PlainText",
		Text: text,
	}

	return this
}

func (this *AlexaResponse) OutputSpeechSSML(text string) *AlexaResponse {
	this.Response.OutputSpeech = &AlexaRespPayload{
		Type: "SSML",
		SSML: text,
	}

	return this
}

func (this *AlexaResponse) Card(title string, content string) *AlexaResponse {
	return this.SimpleCard(title, content)
}

func (this *AlexaResponse) SimpleCard(title string, content string) *AlexaResponse {
	this.Response.Card = &AlexaRespPayload{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return this
}

func (this *AlexaResponse) StandardCard(title string, content string, smallImg string, largeImg string) *AlexaResponse {
	this.Response.Card = &AlexaRespPayload{
		Type:    "Standard",
		Title:   title,
		Content: content,
	}

	if smallImg != "" {
		this.Response.Card.Image.SmallImageURL = smallImg
	}

	if largeImg != "" {
		this.Response.Card.Image.LargeImageURL = largeImg
	}

	return this
}

func (this *AlexaResponse) LinkAccountCard() *AlexaResponse {
	this.Response.Card = &AlexaRespPayload{
		Type: "LinkAccount",
	}

	return this
}

func (this *AlexaResponse) Reprompt(text string) *AlexaResponse {
	this.Response.Reprompt = &AlexaReprompt{
		OutputSpeech: AlexaRespPayload{
			Type: "PlainText",
			Text: text,
		},
	}

	return this
}

func (this *AlexaResponse) EndSession(flag bool) *AlexaResponse {
	this.Response.ShouldEndSession = flag

	return this
}

func (this *AlexaResponse) String() ([]byte, error) {
	jsonStr, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
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
