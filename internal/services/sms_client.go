package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"
)

const (
	// ClientVersion is used in User-Agent request header to provide server with API level.
	ClientVersion = "2.0.0"
	// Endpoint points you to Ippanel REST API.
	Endpoint = "https://api2.ippanel.com/api/v1"
	// httpClientTimeout is used to limit http.Client waiting time.
	httpClientTimeout = 30 * time.Second
)

var (
	ErrUnexpectedResponse = errors.New("the Ippanel API is currently unavailable")
	ErrStatusUnauthorized = errors.New("you api key is not valid")
)

// ResponseCode api response code error type
type ResponseCode int

const (
	ErrForbidden           ResponseCode = 403
	ErrNotFound            ResponseCode = 404
	ErrUnprocessableEntity ResponseCode = 422
	ErrInternalServer      ResponseCode = 500
)

// Error general service error type
type Error struct {
	Code    ResponseCode
	Message interface{}
}

// FieldErrs input field level errors
type FieldErrs map[string][]string

// Error implement error interface
func (e Error) Error() string {
	switch e.Message.(type) {
	case string:
		return e.Message.(string)
	case FieldErrs:
		m, _ := json.Marshal(e.Message)
		return string(m)
	}
	return fmt.Sprint(e.Code)
}

// ListParams ...
type ListParams struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

// PaginationInfo ...
type PaginationInfo struct {
	Total int64   `json:"total"`
	Limit int64   `json:"limit"`
	Page  int64   `json:"page"`
	Pages int64   `json:"pages"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}

// BaseResponse base response model
type BaseResponse struct {
	Status       string          `json:"status"`
	Code         ResponseCode    `json:"code"`
	Data         json.RawMessage `json:"data"`
	Meta         *PaginationInfo `json:"meta"`
	ErrorMessage string          `json:"error_message"`
}

// IPPanelClient - Ippanel client based on official SDK
type IPPanelClient struct {
	Apikey  string
	Client  *http.Client
	BaseURL *url.URL
}

// sendPatternReqType send sms with pattern request template
type sendPatternReqType struct {
	Code      string            `json:"code"`
	Sender    string            `json:"sender"`
	Recipient string            `json:"recipient"`
	Variable  map[string]string `json:"variable"`
}

// sendResType response type for send sms
type sendResType struct {
	MessageId int64 `json:"message_id"`
}

// getCreditResType get credit response type
type getCreditResType struct {
	Credit float64 `json:"credit"`
}

// fieldErrsRes field errors response type
type fieldErrsRes struct {
	Errors FieldErrs `json:"error"`
}

// defaultErrsRes default template for errors body
type defaultErrsRes struct {
	Errors string `json:"error"`
}

// NewIPPanelClient create new ippanel sms instance
func NewIPPanelClient(apikey string) *IPPanelClient {
	u, _ := url.Parse(Endpoint)
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   httpClientTimeout,
	}
	return &IPPanelClient{
		Apikey:  apikey,
		Client:  client,
		BaseURL: u,
	}
}

// request preform http request
func (sms IPPanelClient) request(method string, uri string, params map[string]string, data interface{}) (*BaseResponse, error) {
	u := *sms.BaseURL
	// join base url with extra path
	u.Path = path.Join(sms.BaseURL.Path, uri)

	// set query params
	p := url.Values{}
	for key, param := range params {
		p.Add(key, param)
	}
	u.RawQuery = p.Encode()

	marshaledBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	requestBody := bytes.NewBuffer(marshaledBody)
	req, err := http.NewRequest(method, u.String(), requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Apikey", sms.Apikey)
	req.Header.Set("User-Agent", "Ippanel/ApiClient/"+ClientVersion+" Go/"+runtime.Version())

	res, err := sms.Client.Do(req)
	if err != nil || res == nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated:
		_res := &BaseResponse{}
		if err := json.Unmarshal(responseBody, _res); err != nil {
			return nil, fmt.Errorf("could not decode response JSON, %s: %v", string(responseBody), err)
		}
		return _res, nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("IPPanel API internal server error (500). Response: %s", string(responseBody))
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("IPPanel API unauthorized (401) - check your API key. Response: %s", string(responseBody))
	case http.StatusBadRequest:
		return nil, fmt.Errorf("IPPanel API bad request (400). Response: %s", string(responseBody))
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("IPPanel API rate limit exceeded (429). Response: %s", string(responseBody))
	default:
		_res := &BaseResponse{}
		if err := json.Unmarshal(responseBody, _res); err != nil {
			return nil, fmt.Errorf("could not decode response JSON, %s: %v", string(responseBody), err)
		}
		return _res, sms.parseErrors(_res)
	}
}

// get do get request
func (sms IPPanelClient) get(uri string, params map[string]string) (*BaseResponse, error) {
	return sms.request("GET", uri, params, nil)
}

// post do post request
func (sms IPPanelClient) post(uri string, contentType string, data interface{}) (*BaseResponse, error) {
	return sms.request("POST", uri, nil, data)
}

// parseErrors ...
func (sms IPPanelClient) parseErrors(res *BaseResponse) error {
	var err error
	e := Error{Code: res.Code}

	messageFieldErrs := fieldErrsRes{}
	if err = json.Unmarshal(res.Data, &messageFieldErrs); err == nil {
		e.Message = messageFieldErrs.Errors
	} else {
		messageDefaultErrs := defaultErrsRes{}
		if err = json.Unmarshal(res.Data, &messageDefaultErrs); err == nil {
			e.Message = messageDefaultErrs.Errors
		}
	}

	if err != nil {
		return errors.New("cant marshal errors into standard template")
	}
	return e
}

// SendPattern send a message with pattern
func (sms *IPPanelClient) SendPattern(patternCode string, originator string, recipient string, values map[string]string) (int64, error) {
	data := sendPatternReqType{
		Code:      patternCode,
		Sender:    originator,
		Recipient: recipient,
		Variable:  values,
	}

	// Debug logging
	jsonData, _ := json.Marshal(data)
	fmt.Printf("🔍 SMS Request: %s\n", string(jsonData))

	// Retry mechanism (3 attempts)
	maxAttempts := 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("🔄 SMS Attempt %d/%d\n", attempt, maxAttempts)

		_res, err := sms.post("/sms/pattern/normal/send", "application/json", data)
		if err != nil {
			if attempt == maxAttempts {
				return 0, fmt.Errorf("SMS API request failed after %d attempts: %v", maxAttempts, err)
			}
			fmt.Printf("⚠️ SMS Attempt %d failed, retrying... Error: %v\n", attempt, err)
			time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
			continue
		}

		if _res == nil {
			if attempt == maxAttempts {
				return 0, fmt.Errorf("SMS API returned empty response after %d attempts", maxAttempts)
			}
			fmt.Printf("⚠️ SMS Attempt %d returned empty response, retrying...\n", attempt)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Debug logging
		fmt.Printf("🔍 SMS Response Status: %s, Code: %d\n", _res.Status, _res.Code)
		fmt.Printf("🔍 SMS Response Data: %s\n", string(_res.Data))

		if _res.Data == nil {
			if attempt == maxAttempts {
				return 0, fmt.Errorf("SMS API returned null data after %d attempts", maxAttempts)
			}
			fmt.Printf("⚠️ SMS Attempt %d returned null data, retrying...\n", attempt)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		res := sendResType{}
		if err = json.Unmarshal(_res.Data, &res); err != nil {
			if attempt == maxAttempts {
				return 0, fmt.Errorf("failed to parse SMS API response after %d attempts: %v. Raw data: %s", maxAttempts, err, string(_res.Data))
			}
			fmt.Printf("⚠️ SMS Attempt %d failed to parse response, retrying... Error: %v\n", attempt, err)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		if res.MessageId == 0 {
			if attempt == maxAttempts {
				return 0, fmt.Errorf("SMS API returned invalid message ID after %d attempts. Raw response: %s", maxAttempts, string(_res.Data))
			}
			fmt.Printf("⚠️ SMS Attempt %d returned invalid message ID, retrying...\n", attempt)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		fmt.Printf("✅ SMS sent successfully on attempt %d with MessageId: %d\n", attempt, res.MessageId)
		return res.MessageId, nil
	}

	return 0, fmt.Errorf("unexpected error in SMS retry mechanism")
}

// GetCredit get credit for user
func (sms *IPPanelClient) GetCredit() (float64, error) {
	_res, err := sms.get("/sms/accounting/credit/show", nil)
	if err != nil {
		return 0, err
	}

	res := &getCreditResType{}
	if err = json.Unmarshal(_res.Data, res); err != nil {
		return 0, err
	}

	return res.Credit, nil
}
