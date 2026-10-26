package main

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Вказівники, щоб відрізнити відсутнє поле від нульового значення.
type request struct {
	Age      *int     `json:"age"`
	HeightCm *float64 `json:"height_cm"`
	Sex      *string  `json:"sex"`
}

// Помилки клієнта — звичайна відповідь 400 і nil-error: non-nil error дав би 502 від API Gateway.
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestID := requestIDFrom(ctx, req)
	log := logger.With("request_id", requestID)

	in, msg := validate(req.Body)
	if msg != "" {
		log.WarnContext(ctx, "bad request", "error", msg)
		return respond(requestID, 400, map[string]string{"error": msg}), nil
	}

	res := Calculate(*in.HeightCm, *in.Sex)
	log.InfoContext(ctx, "ideal weight calculated", "height_cm", *in.HeightCm, "sex", *in.Sex, "average_kg", res.AverageKg)
	return respond(requestID, 200, res), nil
}

// validate повертає запит або зрозуміле повідомлення про помилку.
func validate(body string) (request, string) {
	var in request
	if err := json.Unmarshal([]byte(body), &in); err != nil {
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) && typeErr.Field != "" {
			return in, typeErr.Field + " has the wrong type (expected " + expectedType(typeErr.Field) + ")"
		}
		return in, "invalid JSON: body must be an object with age, height_cm and sex"
	}

	var missing []string
	if in.Age == nil {
		missing = append(missing, "age")
	}
	if in.HeightCm == nil {
		missing = append(missing, "height_cm")
	}
	if in.Sex == nil {
		missing = append(missing, "sex")
	}
	if len(missing) > 0 {
		return in, "missing required field(s): " + strings.Join(missing, ", ")
	}

	if *in.Age < 18 || *in.Age > 120 {
		return in, "age must be between 18 and 120"
	}
	if *in.HeightCm < 120 || *in.HeightCm > 230 {
		return in, "height_cm must be between 120 and 230"
	}
	if *in.Sex != "male" && *in.Sex != "female" {
		return in, `sex must be "male" or "female"`
	}
	return in, ""
}

func expectedType(field string) string {
	switch field {
	case "age":
		return "whole number"
	case "height_cm":
		return "number"
	default:
		return "string"
	}
}

func requestIDFrom(ctx context.Context, req events.APIGatewayProxyRequest) string {
	if lc, ok := lambdacontext.FromContext(ctx); ok && lc.AwsRequestID != "" {
		return lc.AwsRequestID
	}
	if req.RequestContext.RequestID != "" {
		return req.RequestContext.RequestID
	}
	return "local"
}

func respond(requestID string, status int, v any) events.APIGatewayProxyResponse {
	b, _ := json.Marshal(v)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Request-Id": requestID,
		},
		Body: string(b),
	}
}
