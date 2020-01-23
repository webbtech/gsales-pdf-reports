package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pulpfree/gsales-pdf-reports/config"
	"github.com/pulpfree/gsales-pdf-reports/model"
	"github.com/pulpfree/gsales-pdf-reports/pkgerrors"
	"github.com/pulpfree/gsales-pdf-reports/report"
	"github.com/pulpfree/gsales-pdf-reports/validate"
	log "github.com/sirupsen/logrus"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
)

// Response data format
type Response struct {
	Code      int         `json:"code"`      // HTTP status code
	Data      interface{} `json:"data"`      // Data payload
	Message   string      `json:"message"`   // Error or status message
	Status    string      `json:"status"`    // Status code (error|fail|success)
	Timestamp int64       `json:"timestamp"` // Machine-readable UTC timestamp in nanoseconds since EPOCH
}

// SignedURL struct
type SignedURL struct {
	URL string `json:"url"`
}

var (
	cfg      *config.Config
	stdError *pkgerrors.StdError
)

func init() {
	cfg = &config.Config{}
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
}

// HandleRequest function
// NOTE: strange, the error parameter cannot be used or removed... would be good to dig into
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	hdrs["Access-Control-Allow-Origin"] = "*"
	hdrs["Access-Control-Allow-Methods"] = "GET,OPTIONS,POST,PUT"
	hdrs["Access-Control-Allow-Headers"] = "Authorization,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"

	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{Body: string("null"), Headers: hdrs, StatusCode: 200}, nil
	}

	t := time.Now()

	// If this is a ping test, intercept and return
	if req.HTTPMethod == "GET" {
		log.Info("Ping test in handleRequest")
		return gatewayResponse(Response{
			Code:      200,
			Data:      "pong",
			Status:    "success",
			Timestamp: t.Unix(),
		}, hdrs, nil), nil
	}

	var r *model.RequestInput
	json.Unmarshal([]byte(req.Body), &r)

	// validate input
	reportRequest, err := validate.SetRequest(r)
	if err != nil {
		fmt.Printf("err in validate %+v\n", err)
		return gatewayResponse(Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	rpt, err := report.New(reportRequest, cfg)
	if err != nil {
		return gatewayResponse(Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	url, err := rpt.CreateSignedURL()
	if err != nil {
		return gatewayResponse(Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	urlStr := url[0:100]
	log.Infof("signed url created %s", urlStr)

	return gatewayResponse(Response{
		Code:      201,
		Data:      SignedURL{URL: url},
		Status:    "success",
		Timestamp: t.Unix(),
	}, hdrs, nil), nil

}

func main() {
	lambda.Start(thundra.Wrap(HandleRequest))
}

func gatewayResponse(resp Response, hdrs map[string]string, err error) events.APIGatewayProxyResponse {

	if err != nil {
		resp.Code = 500
		resp.Status = "error"
		log.Error(err)
		// send friendly error to client
		if ok := errors.As(err, &stdError); ok {
			resp.Message = stdError.Msg
		} else {
			resp.Message = err.Error()
		}
	}
	body, _ := json.Marshal(&resp)

	return events.APIGatewayProxyResponse{Body: string(body), Headers: hdrs, StatusCode: resp.Code}
}
