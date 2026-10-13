package main

import (
	"bytes"
	"conways/input"
	"conways/life"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func request(t *testing.T, handler http.Handler, method, path, body, contentType string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}

func TestParseHandler(t *testing.T) {
	recorder := request(t, newHandler(), http.MethodPost, "/api/parse", "3 2\n#..\n..#\n", "text/plain")
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200: %s", recorder.Code, recorder.Body)
	}
	var response parseResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Width != 3 || response.Height != 2 || response.Grid[0][0] != 1 {
		t.Fatalf("unexpected parse response: %+v", response)
	}
}

func TestParseHandlerInvalidGrid(t *testing.T) {
	recorder := request(t, newHandler(), http.MethodPost, "/api/parse", "2 2\n#x\n..\n", "text/plain")
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; want 400", recorder.Code)
	}
	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Error == "" {
		t.Fatal("expected JSON error message")
	}
}

func TestStepHandlerDeadBorder(t *testing.T) {
	body := `{"grid":[[0,1,0],[0,1,0],[0,1,0]],"wrap":false}`
	recorder := request(t, newHandler(), http.MethodPost, "/api/step", body, "application/json")
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200: %s", recorder.Code, recorder.Body)
	}
	var response gridResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	want := life.Grid{{0, 0, 0}, {1, 1, 1}, {0, 0, 0}}
	if !reflect.DeepEqual(response.Grid, want) {
		t.Fatalf("step = %v; want %v", response.Grid, want)
	}
}

func TestStepHandlerToroidal(t *testing.T) {
	grid, err := input.ReadGrid("../../testdata/t4_wrap.txt")
	if err != nil {
		t.Fatal(err)
	}
	bodyBytes, err := json.Marshal(stepRequest{Grid: grid, Wrap: true})
	if err != nil {
		t.Fatal(err)
	}
	recorder := request(t, newHandler(), http.MethodPost, "/api/step", string(bodyBytes), "application/json")
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200: %s", recorder.Code, recorder.Body)
	}
	var response gridResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	want, err := input.ReadGrid("../../testdata/t4_wrap.gen1.want")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(response.Grid, want) {
		t.Fatalf("step = %v; want %v", response.Grid, want)
	}
}

func TestStepHandlerMalformedRequest(t *testing.T) {
	for _, body := range []string{"not json", `{"grid":[[1],[0,1]],"wrap":false}`, `{"grid":[[2]],"wrap":false}`} {
		recorder := request(t, newHandler(), http.MethodPost, "/api/step", body, "application/json")
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("body %q: status = %d; want 400", body, recorder.Code)
		}
	}
}
