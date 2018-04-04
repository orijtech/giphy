// Copyright 2017 orijtech. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package giphy_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/orijtech/giphy/v1"
)

func TestTrending(t *testing.T) {
	client, err := giphy.NewClient(testAPIKey1)
	if err != nil {
		t.Fatal(err)
	}
	tr := &transport{route: trendingRoute}
	client.SetHTTPRoundTripper(tr)

	tests := [...]struct {
		req     *giphy.Request
		wantErr bool
	}{
		0: {
			req: &giphy.Request{},
		},
		1: {
			req: nil,
		},
	}

	for i, tt := range tests {
		res, err := client.Trending(context.Background(), tt.req)
		if tt.wantErr {
			if err == nil {
				t.Errorf("#%d: want non-nil error", i)
			}
			continue
		}

		if res == nil {
			t.Errorf("#%d: expected non-nil response", i)
			continue
		}
		for page := range res.Pages {
			if page.Err != nil {
				t.Errorf("Page #%d err: %v", page.PageNumber, page.Err)
				continue
			}
		}
	}
}

func TestTrendingStickers(t *testing.T) {
	client, err := giphy.NewClient(testAPIKey1)
	if err != nil {
		t.Fatal(err)
	}
	tr := &transport{route: trendingStickersRoute}
	client.SetHTTPRoundTripper(tr)

	tests := [...]struct {
		req      *giphy.Request
		wantErr  bool
		maxPages int
	}{
		0: {
			req: &giphy.Request{},
		},
		1: {
			req: nil,
		},
		2: {
			req:      &giphy.Request{MaxPageNumber: 2},
			maxPages: 2,
		},
	}

	for i, tt := range tests {
		res, err := client.TrendingStickers(context.Background(), tt.req)
		if tt.wantErr {
			if err == nil {
				t.Errorf("#%d: want non-nil error", i)
			}
			continue
		}

		if res == nil {
			t.Errorf("#%d: expected non-nil response", i)
			continue
		}
		pageCount := 0
		for page := range res.Pages {
			if page.Err != nil {
				t.Errorf("Page #%d err: %v", page.PageNumber, page.Err)
				continue
			}
			if len(page.Giphs) > 0 {
				pageCount += 1
			}
		}

		if tt.maxPages > 0 && tt.maxPages != pageCount {
			t.Errorf("#%d gotPageCount: %d wantPageCount: %d", i, tt.maxPages, pageCount)
		}
	}
}

func TestSearch(t *testing.T) {
	t.Errorf("Unimplemented")
}

func TestSearchStickers(t *testing.T) {
	t.Errorf("Unimplemented")
}

func TestGIFByID(t *testing.T) {
	t.Errorf("Unimplemented")
}

func TestRandomGIF(t *testing.T) {
	t.Errorf("Unimplemented")
}

func TestRandomSticker(t *testing.T) {
	t.Errorf("Unimplemented")
}

type transport struct {
	route string
}

var _ http.RoundTripper = (*transport)(nil)

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch t.route {
	case trendingRoute:
		return t.trendingRoundTrip(req)
	case trendingStickersRoute:
		return t.trendingStickersRoundTrip(req)
	case gifByIDRoute:
		return t.gifByIDRoundTrip(req)
	case randomGIFRoute:
		return t.randomGIFRoundTrip(req)
	case randomStickerRoute:
		return t.randomStickerRoundTrip(req)
	default:
		return nil, errUnimplemented
	}
}

const (
	testAPIKey1 = "test-api-key1"
	testAPIKey2 = "test-api-key2"

	trendingRoute         = "/trending"
	trendingStickersRoute = "/trending-stickers"
	gifByIDRoute          = "/gif-by-id"
	randomGIFRoute        = "/random-gif"
	randomStickerRoute    = "/random-sticker"
)

var errUnimplemented = errors.New("unimplemented")
var (
	blankReqResp = makeResp("blank request", http.StatusBadRequest, nil)
)

func makeResp(status string, code int, body io.ReadCloser) *http.Response {
	res := &http.Response{
		Header:     make(http.Header),
		StatusCode: code,
		Status:     status,
		Body:       body,
	}
	return res
}

func checkBadAuthAndMethod(req *http.Request, wantMethod string) (*http.Response, error) {
	if req == nil {
		return blankReqResp, nil
	}
	if got := req.Method; got != wantMethod {
		return makeResp(fmt.Sprintf("only %q supported not %q", wantMethod, got), http.StatusBadRequest, nil), nil
	}
	query := req.URL.Query()
	if strings.TrimSpace(query.Get("api_key")) == "" {
		return makeResp(`expecting "api_key"`, http.StatusUnauthorized, nil), nil
	}
	return nil, nil
}

func (t *transport) trendingRoundTrip(req *http.Request) (*http.Response, error) {
	if badAuthResp, err := checkBadAuthAndMethod(req, "GET"); badAuthResp != nil || err != nil {
		return badAuthResp, err
	}
	query := req.URL.Query()
	offsetStr := query.Get("offset")
	var offset int
	if offsetStr != "" {
		var err error
		if offset, err = strconv.Atoi(offsetStr); err != nil {
			return makeResp(fmt.Sprintf(`%q could not be parsed as "offset"`, offsetStr), http.StatusBadRequest, nil), nil
		}
	}
	offset /= 20
	srcPath := trendingPathByOffset(offset)
	f, err := os.Open(srcPath)
	if err != nil {
		return makeResp(err.Error(), http.StatusBadRequest, nil), nil
	}
	return makeResp("200", http.StatusOK, f), nil
}

func trendingPathByOffset(offset int) string {
	return fmt.Sprintf("./testdata/trending-%d.json", offset)
}

func trendingStickersPathByOffset(offset int) string {
	return fmt.Sprintf("./testdata/trending-stickers-%d.json", offset)
}

func (t *transport) trendingStickersRoundTrip(req *http.Request) (*http.Response, error) {
	if badAuthResp, err := checkBadAuthAndMethod(req, "GET"); badAuthResp != nil || err != nil {
		return badAuthResp, err
	}
	query := req.URL.Query()
	offsetStr := query.Get("offset")
	var offset int
	if offsetStr != "" {
		var err error
		if offset, err = strconv.Atoi(offsetStr); err != nil {
			return makeResp(fmt.Sprintf(`%q could not be parsed as "offset"`, offsetStr), http.StatusBadRequest, nil), nil
		}
	}
	offset /= 20
	srcPath := trendingPathByOffset(offset)
	f, err := os.Open(srcPath)
	if err != nil {
		return makeResp(err.Error(), http.StatusBadRequest, nil), nil
	}
	return makeResp("200", http.StatusOK, f), nil
}

func (t *transport) gifByIDRoundTrip(req *http.Request) (*http.Response, error) {
	if badAuthResp, err := checkBadAuthAndMethod(req, "GET"); badAuthResp != nil || err != nil {
		return badAuthResp, err
	}
	return nil, errUnimplemented
}

func (t *transport) randomGIFRoundTrip(req *http.Request) (*http.Response, error) {
	if badAuthResp, err := checkBadAuthAndMethod(req, "GET"); badAuthResp != nil || err != nil {
		return badAuthResp, err
	}
	return nil, errUnimplemented
}

func (t *transport) randomStickerRoundTrip(req *http.Request) (*http.Response, error) {
	if badAuthResp, err := checkBadAuthAndMethod(req, "GET"); badAuthResp != nil || err != nil {
		return badAuthResp, err
	}
	return nil, errUnimplemented
}
