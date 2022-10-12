package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukdsub"
	"github.com/ipthomas/tukutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handle_Request)
}
func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	dsubEvent := tukdsub.DSUBEvent{
		Action:          req.QueryStringParameters[tukcnst.QUERY_PARAM_ACTION],
		BrokerURL:       os.Getenv(tukcnst.AWS_ENV_DSUB_BROKER_URL),
		ConsumerURL:     os.Getenv(tukcnst.AWS_ENV_DSUB_CONSUMER_URL),
		PDQ_SERVER_URL:  os.Getenv(tukcnst.AWS_ENV_PDQ_SERVER_URL),
		PDQ_SERVER_TYPE: os.Getenv(tukcnst.AWS_ENV_PDQ_SERVER_TYPE),
		REG_OID:         os.Getenv(tukcnst.AWS_ENV_REG_OID),
		NHS_OID:         os.Getenv(tukcnst.AWS_ENV_NHS_OID),
		Pathway:         req.QueryStringParameters[tukcnst.QUERY_PARAM_PATHWAY],
		EventMessage:    req.Body,
		DBConnection:    tukdbint.TukDBConnection{DB_URL: os.Getenv(tukcnst.AWS_ENV_TUK_DB_URL)},
	}

	if req.QueryStringParameters[tukcnst.QUERY_PARAM_EXPRESSION] != "" {
		var exs []string
		if strings.Contains(req.QueryStringParameters[tukcnst.QUERY_PARAM_EXPRESSION], "|") {
			reqexs := strings.Split(req.QueryStringParameters[tukcnst.QUERY_PARAM_EXPRESSION], "|")
			exs = append(exs, reqexs...)
		} else {
			exs = append(exs, req.QueryStringParameters[tukcnst.QUERY_PARAM_EXPRESSION])
		}
		dsubEvent.Expressions = exs
	}
	log.Println("Set Expressions")
	tukutil.Log(dsubEvent.Expressions)
	if err := tukdsub.New_Transaction(&dsubEvent); err != nil {
		log.Println(err.Error())
	}

	apiResp := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(dsubEvent.Response),
	}
	return &apiResp, nil
}
