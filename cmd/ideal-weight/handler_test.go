package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func call(t *testing.T, body string) (int, map[string]any) {
	t.Helper()
	resp, err := handler(context.Background(), events.APIGatewayProxyRequest{Body: body})
	if err != nil {
		t.Fatalf("handler must not return a Go error (would be 502): %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(resp.Body), &out); err != nil {
		t.Fatalf("response body is not JSON: %v (%q)", err, resp.Body)
	}
	return resp.StatusCode, out
}

func TestBadRequest(t *testing.T) {
	tests := []struct{ name, body, wantErr string }{
		{"empty body", ``, "invalid JSON"},
		{"malformed json", `{"age": 30,`, "invalid JSON"},
		{"json array", `[1,2]`, "invalid JSON"},
		{"missing height", `{"age":30,"sex":"male"}`, "missing required field(s): height_cm"},
		{"missing all", `{}`, "missing required field(s): age, height_cm, sex"},
		{"age as string", `{"age":"30","height_cm":180,"sex":"male"}`, "age has the wrong type"},
		{"age fractional", `{"age":30.5,"height_cm":180,"sex":"male"}`, "age has the wrong type"},
		{"height as string", `{"age":30,"height_cm":"180","sex":"male"}`, "height_cm has the wrong type"},
		{"sex as number", `{"age":30,"height_cm":180,"sex":1}`, "sex has the wrong type"},
		{"age 17", `{"age":17,"height_cm":180,"sex":"male"}`, "age must be between 18 and 120"},
		{"age 121", `{"age":121,"height_cm":180,"sex":"male"}`, "age must be between 18 and 120"},
		{"height 119", `{"age":30,"height_cm":119,"sex":"male"}`, "height_cm must be between 120 and 230"},
		{"height 231", `{"age":30,"height_cm":231,"sex":"male"}`, "height_cm must be between 120 and 230"},
		{"sex other", `{"age":30,"height_cm":180,"sex":"other"}`, `sex must be "male" or "female"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, out := call(t, tt.body)
			if status != 400 {
				t.Fatalf("status = %d, want 400", status)
			}
			if msg, _ := out["error"].(string); !strings.Contains(msg, tt.wantErr) {
				t.Errorf("error = %q, want it to contain %q", msg, tt.wantErr)
			}
		})
	}
}

func TestBoundariesAccepted(t *testing.T) {
	for _, body := range []string{
		`{"age":18,"height_cm":180,"sex":"male"}`,
		`{"age":120,"height_cm":180,"sex":"male"}`,
		`{"age":30,"height_cm":120,"sex":"female"}`,
		`{"age":30,"height_cm":230,"sex":"female"}`,
	} {
		if status, _ := call(t, body); status != 200 {
			t.Errorf("%s: status = %d, want 200", body, status)
		}
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		body    string
		average float64
	}{
		{`{"age":30,"height_cm":180,"sex":"male"}`, 74.1},
		{`{"age":45,"height_cm":165,"sex":"female"}`, 57.7},
	}
	for _, tt := range tests {
		status, out := call(t, tt.body)
		if status != 200 {
			t.Fatalf("%s: status = %d, want 200", tt.body, status)
		}
		if out["average_kg"] != tt.average {
			t.Errorf("%s: average_kg = %v, want %v", tt.body, out["average_kg"], tt.average)
		}
		if note, _ := out["note"].(string); !strings.HasPrefix(note, "Not medical advice") {
			t.Errorf("%s: note = %q", tt.body, note)
		}
	}
}
