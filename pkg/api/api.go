// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Course defines model for Course.
type Course struct {
	Holes []Hole  `dynamodbav:"h" json:"holes"`
	Name  string  `dynamodbav:"n" json:"name"`
	Tees  *string `dynamodbav:"t" json:"tees,omitempty"`
}

// CreatedRound defines model for CreatedRound.
type CreatedRound struct {
	Id RoundID `dynamodbav:"id" json:"id"`
}

// Hole defines model for Hole.
type Hole struct {
	DistanceYards *int       `dynamodbav:"d" json:"distanceYards,omitempty"`
	Hole          HoleNumber `dynamodbav:"n" json:"hole"`
	Par           int        `dynamodbav:"p" json:"par"`
	StrokeIndex   *int       `dynamodbav:"si" json:"strokeIndex,omitempty"`
}

// HoleNumber defines model for HoleNumber.
type HoleNumber = int

// HoleScore defines model for HoleScore.
type HoleScore struct {
	Hole  HoleNumber `dynamodbav:"n" json:"hole"`
	Score int        `dynamodbav:"s" json:"score"`
}

// PlayerData defines model for PlayerData.
type PlayerData struct {
	Name string `dynamodbav:"n" json:"name"`
}

// PlayerGroup defines model for PlayerGroup.
type PlayerGroup = []PlayerData

// PlayerScore defines model for PlayerScore.
type PlayerScore = []HoleScore

// PlayerScoreEvent defines model for PlayerScoreEvent.
type PlayerScoreEvent struct {
	Hole        HoleNumber `dynamodbav:"n" json:"hole"`
	PlayerIndex int        `json:"playerIndex"`
	Score       int        `dynamodbav:"s" json:"score"`
}

// Round defines model for Round.
type Round struct {
	Course  Course       `dynamodbav:"c" json:"course"`
	Id      RoundID      `dynamodbav:"id" json:"id"`
	Players RoundPlayers `dynamodbav:"rp" json:"players"`
	Title   RoundTitle   `dynamodbav:"t" json:"title"`
}

// RoundID defines model for RoundID.
type RoundID = uuid.UUID

// RoundPlayerData defines model for RoundPlayerData.
type RoundPlayerData struct {
	Name   string      `dynamodbav:"n" json:"name"`
	Scores PlayerScore `dynamodbav:"ps" json:"scores"`
}

// RoundPlayers defines model for RoundPlayers.
type RoundPlayers = []RoundPlayerData

// RoundRequest defines model for RoundRequest.
type RoundRequest struct {
	Course  Course      `dynamodbav:"c" json:"course"`
	Players PlayerGroup `dynamodbav:"pg" json:"players"`
	Title   *RoundTitle `dynamodbav:"t" json:"title,omitempty"`
}

// RoundTitle defines model for RoundTitle.
type RoundTitle = string

// CreateRoundJSONRequestBody defines body for CreateRound for application/json ContentType.
type CreateRoundJSONRequestBody = RoundRequest

// SendScoreJSONRequestBody defines body for SendScore for application/json ContentType.
type SendScoreJSONRequestBody = PlayerScoreEvent

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetLatestRoundID request
	GetLatestRoundID(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateRoundWithBody request with any body
	CreateRoundWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateRound(ctx context.Context, body CreateRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetRound request
	GetRound(ctx context.Context, roundID string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SendScoreWithBody request with any body
	SendScoreWithBody(ctx context.Context, roundID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SendScore(ctx context.Context, roundID string, body SendScoreJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetLatestRoundID(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetLatestRoundIDRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateRoundWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateRoundRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateRound(ctx context.Context, body CreateRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateRoundRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetRound(ctx context.Context, roundID string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetRoundRequest(c.Server, roundID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SendScoreWithBody(ctx context.Context, roundID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSendScoreRequestWithBody(c.Server, roundID, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SendScore(ctx context.Context, roundID string, body SendScoreJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSendScoreRequest(c.Server, roundID, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetLatestRoundIDRequest generates requests for GetLatestRoundID
func NewGetLatestRoundIDRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/latest/roundID")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateRoundRequest calls the generic CreateRound builder with application/json body
func NewCreateRoundRequest(server string, body CreateRoundJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateRoundRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateRoundRequestWithBody generates requests for CreateRound with any type of body
func NewCreateRoundRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/round")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetRoundRequest generates requests for GetRound
func NewGetRoundRequest(server string, roundID string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "roundID", runtime.ParamLocationPath, roundID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/round/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewSendScoreRequest calls the generic SendScore builder with application/json body
func NewSendScoreRequest(server string, roundID string, body SendScoreJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSendScoreRequestWithBody(server, roundID, "application/json", bodyReader)
}

// NewSendScoreRequestWithBody generates requests for SendScore with any type of body
func NewSendScoreRequestWithBody(server string, roundID string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "roundID", runtime.ParamLocationPath, roundID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/round/%s/score", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetLatestRoundIDWithResponse request
	GetLatestRoundIDWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetLatestRoundIDResponse, error)

	// CreateRoundWithBodyWithResponse request with any body
	CreateRoundWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateRoundResponse, error)

	CreateRoundWithResponse(ctx context.Context, body CreateRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateRoundResponse, error)

	// GetRoundWithResponse request
	GetRoundWithResponse(ctx context.Context, roundID string, reqEditors ...RequestEditorFn) (*GetRoundResponse, error)

	// SendScoreWithBodyWithResponse request with any body
	SendScoreWithBodyWithResponse(ctx context.Context, roundID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SendScoreResponse, error)

	SendScoreWithResponse(ctx context.Context, roundID string, body SendScoreJSONRequestBody, reqEditors ...RequestEditorFn) (*SendScoreResponse, error)
}

type GetLatestRoundIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *RoundID
}

// Status returns HTTPResponse.Status
func (r GetLatestRoundIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetLatestRoundIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateRoundResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *CreatedRound
}

// Status returns HTTPResponse.Status
func (r CreateRoundResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateRoundResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetRoundResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Round
}

// Status returns HTTPResponse.Status
func (r GetRoundResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetRoundResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SendScoreResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r SendScoreResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SendScoreResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetLatestRoundIDWithResponse request returning *GetLatestRoundIDResponse
func (c *ClientWithResponses) GetLatestRoundIDWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetLatestRoundIDResponse, error) {
	rsp, err := c.GetLatestRoundID(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetLatestRoundIDResponse(rsp)
}

// CreateRoundWithBodyWithResponse request with arbitrary body returning *CreateRoundResponse
func (c *ClientWithResponses) CreateRoundWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateRoundResponse, error) {
	rsp, err := c.CreateRoundWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateRoundResponse(rsp)
}

func (c *ClientWithResponses) CreateRoundWithResponse(ctx context.Context, body CreateRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateRoundResponse, error) {
	rsp, err := c.CreateRound(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateRoundResponse(rsp)
}

// GetRoundWithResponse request returning *GetRoundResponse
func (c *ClientWithResponses) GetRoundWithResponse(ctx context.Context, roundID string, reqEditors ...RequestEditorFn) (*GetRoundResponse, error) {
	rsp, err := c.GetRound(ctx, roundID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetRoundResponse(rsp)
}

// SendScoreWithBodyWithResponse request with arbitrary body returning *SendScoreResponse
func (c *ClientWithResponses) SendScoreWithBodyWithResponse(ctx context.Context, roundID string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SendScoreResponse, error) {
	rsp, err := c.SendScoreWithBody(ctx, roundID, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSendScoreResponse(rsp)
}

func (c *ClientWithResponses) SendScoreWithResponse(ctx context.Context, roundID string, body SendScoreJSONRequestBody, reqEditors ...RequestEditorFn) (*SendScoreResponse, error) {
	rsp, err := c.SendScore(ctx, roundID, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSendScoreResponse(rsp)
}

// ParseGetLatestRoundIDResponse parses an HTTP response from a GetLatestRoundIDWithResponse call
func ParseGetLatestRoundIDResponse(rsp *http.Response) (*GetLatestRoundIDResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetLatestRoundIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest RoundID
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateRoundResponse parses an HTTP response from a CreateRoundWithResponse call
func ParseCreateRoundResponse(rsp *http.Response) (*CreateRoundResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateRoundResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest CreatedRound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseGetRoundResponse parses an HTTP response from a GetRoundWithResponse call
func ParseGetRoundResponse(rsp *http.Response) (*GetRoundResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetRoundResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Round
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseSendScoreResponse parses an HTTP response from a SendScoreWithResponse call
func ParseSendScoreResponse(rsp *http.Response) (*SendScoreResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SendScoreResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Gets the latest round ID
	// (GET /latest/roundID)
	GetLatestRoundID(ctx echo.Context) error
	// Creates a new round
	// (POST /round)
	CreateRound(ctx echo.Context) error
	// Gets round information
	// (GET /round/{roundID})
	GetRound(ctx echo.Context, roundID string) error
	// Send a score event
	// (PUT /round/{roundID}/score)
	SendScore(ctx echo.Context, roundID string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetLatestRoundID converts echo context to params.
func (w *ServerInterfaceWrapper) GetLatestRoundID(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetLatestRoundID(ctx)
	return err
}

// CreateRound converts echo context to params.
func (w *ServerInterfaceWrapper) CreateRound(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateRound(ctx)
	return err
}

// GetRound converts echo context to params.
func (w *ServerInterfaceWrapper) GetRound(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "roundID" -------------
	var roundID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "roundID", runtime.ParamLocationPath, ctx.Param("roundID"), &roundID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter roundID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetRound(ctx, roundID)
	return err
}

// SendScore converts echo context to params.
func (w *ServerInterfaceWrapper) SendScore(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "roundID" -------------
	var roundID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "roundID", runtime.ParamLocationPath, ctx.Param("roundID"), &roundID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter roundID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SendScore(ctx, roundID)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/latest/roundID", wrapper.GetLatestRoundID)
	router.POST(baseURL+"/round", wrapper.CreateRound)
	router.GET(baseURL+"/round/:roundID", wrapper.GetRound)
	router.PUT(baseURL+"/round/:roundID/score", wrapper.SendScore)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
