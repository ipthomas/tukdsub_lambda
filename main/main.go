package main

import (
	"log"
	"net/http"

	"github.com/ipthomas/tukdsub"
	"github.com/ipthomas/tukint"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	DB_URL          = "https://5k2o64mwt5.execute-api.eu-west-1.amazonaws.com/beta/"
	DSUB_BROKER_URL = "http://spirit-test-01.tianispirit.co.uk:8081/SpiritXDSDsub/Dsub"
	PIX_MANAGER_URL = "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient"
	REGIONAL_OID    = "2.16.840.1.113883.2.1.3.31.2.1.1"
	NHS_OID         = "2.16.840.1.113883.2.1.4.1"
)

func main() {
	lambda.Start(Handle_Request)
}
func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	tukint.Set_AWS_Env_Vars(DB_URL, DSUB_BROKER_URL, PIX_MANAGER_URL, NHS_OID, REGIONAL_OID)
	dsubevent := tukdsub.DSUBEvent{Message: req.Body}
	return handle_Response(tukdsub.NewEvent(&dsubevent))
}
func handle_Response(err error) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: tukint.SOAP_XML_Content_Type_EventHeaders, MultiValueHeaders: map[string][]string{}, IsBase64Encoded: false}
	if err == nil {
		resp.StatusCode = http.StatusOK
		resp.Body = tukdsub.DSUB_ACK_TEMPLATE
	} else {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()
		log.Println(err.Error())
	}
	return &resp, err
}
