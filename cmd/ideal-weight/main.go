// Ideal-weight-функція LibraTrack: HTTP POST /ideal-weight через API Gateway.
// Stateless — без бази даних і без сховища ідемпотентності.
package main

import (
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

// Логер на рівні пакета створюється один раз під час init-фази, а не на кожен виклик.
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	lambda.Start(handler)
}
