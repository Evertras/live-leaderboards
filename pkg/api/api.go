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
	Name   string       `dynamodbav:"n" json:"name"`
	Scores *PlayerScore `dynamodbav:"ps" json:"scores,omitempty"`
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

// PostRoundJSONRequestBody defines body for PostRound for application/json ContentType.
type PostRoundJSONRequestBody = RoundRequest

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
	// PostRoundWithBody request with any body
	PostRoundWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostRound(ctx context.Context, body PostRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetRoundRoundID request
	GetRoundRoundID(ctx context.Context, roundID string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostRoundWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostRoundRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostRound(ctx context.Context, body PostRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostRoundRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetRoundRoundID(ctx context.Context, roundID string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetRoundRoundIDRequest(c.Server, roundID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostRoundRequest calls the generic PostRound builder with application/json body
func NewPostRoundRequest(server string, body PostRoundJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostRoundRequestWithBody(server, "application/json", bodyReader)
}

// NewPostRoundRequestWithBody generates requests for PostRound with any type of body
func NewPostRoundRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
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

// NewGetRoundRoundIDRequest generates requests for GetRoundRoundID
func NewGetRoundRoundIDRequest(server string, roundID string) (*http.Request, error) {
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
	// PostRoundWithBodyWithResponse request with any body
	PostRoundWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostRoundResponse, error)

	PostRoundWithResponse(ctx context.Context, body PostRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*PostRoundResponse, error)

	// GetRoundRoundIDWithResponse request
	GetRoundRoundIDWithResponse(ctx context.Context, roundID string, reqEditors ...RequestEditorFn) (*GetRoundRoundIDResponse, error)
}

type PostRoundResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *CreatedRound
}

// Status returns HTTPResponse.Status
func (r PostRoundResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostRoundResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetRoundRoundIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Round
}

// Status returns HTTPResponse.Status
func (r GetRoundRoundIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetRoundRoundIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostRoundWithBodyWithResponse request with arbitrary body returning *PostRoundResponse
func (c *ClientWithResponses) PostRoundWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostRoundResponse, error) {
	rsp, err := c.PostRoundWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostRoundResponse(rsp)
}

func (c *ClientWithResponses) PostRoundWithResponse(ctx context.Context, body PostRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*PostRoundResponse, error) {
	rsp, err := c.PostRound(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostRoundResponse(rsp)
}

// GetRoundRoundIDWithResponse request returning *GetRoundRoundIDResponse
func (c *ClientWithResponses) GetRoundRoundIDWithResponse(ctx context.Context, roundID string, reqEditors ...RequestEditorFn) (*GetRoundRoundIDResponse, error) {
	rsp, err := c.GetRoundRoundID(ctx, roundID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetRoundRoundIDResponse(rsp)
}

// ParsePostRoundResponse parses an HTTP response from a PostRoundWithResponse call
func ParsePostRoundResponse(rsp *http.Response) (*PostRoundResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostRoundResponse{
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

// ParseGetRoundRoundIDResponse parses an HTTP response from a GetRoundRoundIDWithResponse call
func ParseGetRoundRoundIDResponse(rsp *http.Response) (*GetRoundRoundIDResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetRoundRoundIDResponse{
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

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /round)
	PostRound(ctx echo.Context) error

	// (GET /round/{roundID})
	GetRoundRoundID(ctx echo.Context, roundID string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostRound converts echo context to params.
func (w *ServerInterfaceWrapper) PostRound(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostRound(ctx)
	return err
}

// GetRoundRoundID converts echo context to params.
func (w *ServerInterfaceWrapper) GetRoundRoundID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "roundID" -------------
	var roundID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "roundID", runtime.ParamLocationPath, ctx.Param("roundID"), &roundID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter roundID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetRoundRoundID(ctx, roundID)
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

	router.POST(baseURL+"/round", wrapper.PostRound)
	router.GET(baseURL+"/round/:roundID", wrapper.GetRoundRoundID)

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