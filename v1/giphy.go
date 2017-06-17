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

package giphy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/orijtech/otils"
)

const (
	baseURL = "https://api.giphy.com/v1"

	NoThrottle = int64(-1)
)

type Request struct {
	Query string `json:"query"`

	Rating Rating `json:"rating"`
	Format Format `json:"format"`

	MaxPageNumber uint64    `json:"max_page_number"`
	LimitPerPage  uint64    `json:"limit_per_page"`
	Language      Language  `json:"lang"`
	SortBy        SortOrder `json:"sort_by"`
	Tag           string    `json:"tag"`

	ThrottleDurationMs int64 `json:"throttle_duration_ms"`
}

type Rating string

const (
	RatingPG      Rating = "pg"
	RatingPG13    Rating = "pg-13"
	RatingR       Rating = "r"
	RatingGeneral Rating = "g"
	RatingYouth   Rating = "y"
)

type Format string

const (
	FormatHTMl Format = "html"
	FormatJSON Format = "json"
)

type GIF struct {
	Type     string `json:"type,omitempty"`
	ID       string `json:"id,omitempty"`
	URL      string `json:"url"`
	Width    int    `json:"width,string,omitempty"`
	Height   int    `json:"height,string,omitempty"`
	Size     int64  `json:"size,string,omitempty"`
	MP4      string `json:"mp4,omitempty"`
	MP4Size  int64  `json:"mp4_size,string,omitempty"`
	Webp     string `json:"webp,omitempty"`
	WebpSize int64  `json:"webp_size,string,omitempty"`
}

type Giph struct {
	Type        string `json:"type,omitempty"`
	ID          string `json:"id,omitempty"`
	Slug        string `json:"slug,omitempty"`
	BitlyURL    string `json:"bitly_url,omitempty"`
	BitlyGIFURL string `json:"bitly_gif_url,omitempty"`
	EmbedURL    string `json:"embed_url,omitempty"`
	Owner       string `json:"username,omitempty"`
	Source      string `json:"source,omitempty"`
	Rating      string `json:"rating,omitempty"`
	Caption     string `json:"caption,omitempty"`
	ContentURL  string `json:"content_url,omitempty"`

	SourceTopLevelDomain string `json:"source_tld,omitempty"`
	SourcePostURL        string `json:"source_post_url,omitempty"`

	ImportDate   *GiphyTime `json:"import_datetime,omitempty"`
	TrendingDate *GiphyTime `json:"trending_datetime,omitempty"`

	Sizes map[string]*GIF `json:"images"`

	ImageOriginalURL string `json:"image_original_url,omitempty"`
	ImageURL         string `json:"image_url,omitempty"`
	FrameCount       uint   `json:"image_frames,string,omitempty"`
	ImageWidth       int    `json:"image_width,string,omitempty"`
	ImageHeight      int    `json:"image_height,string,omitempty"`

	FixedHeightDownsampledURL    string `json:"fixed_height_downsampled_url,omitempty"`
	FixedHeightDownsampledHeight int    `json:"fixed_height_downsampled_height,string,omitempty"`
	FixedHeightDownsampledWidth  int    `json:"fixed_height_downsampled_width,string,omitempty"`

	FixedHeightSmallURL    string `json:"fixed_height_small_url,omitempty"`
	FixedHeightSmallHeight int    `json:"fixed_height_small_height,string,omitempty"`
	FixedHeightSmallWidth  int    `json:"fixed_height_small_width,string,omitempty"`

	FixedHeightSmallStillURL    string `json:"fixed_height_small_still_url,omitempty"`
	FixedHeightSmallStillHeight int    `json:"fixed_height_small_still_height,string,omitempty"`
	FixedHeightSmallStillWidth  int    `json:"fixed_height_small_still_width,string,omitempty"`

	FixedWidthDownsampledURL    string `json:"fixed_width_downsampled_url,omitempty"`
	FixedWidthDownsampledHeight int    `json:"fixed_width_downsampled_height,string,omitempty"`
	FixedWidthDownsampledWidth  int    `json:"fixed_width_downsampled_width,string,omitempty"`

	FixedWidthSmallURL    string `json:"fixed_width_small_url,omitempty"`
	FixedWidthSmallHeight int    `json:"fixed_width_small_height,string,omitempty"`
	FixedWidthSmallWidth  int    `json:"fixed_width_small_width,string,omitempty"`

	FixedWidthSmallStillURL    string `json:"fixed_width_small_still_url,omitempty"`
	FixedWidthSmallStillHeight int    `json:"fixed_width_small_still_height,string,omitempty"`
	FixedWidthSmallStillWidth  int    `json:"fixed_width_small_still_width,string,omitempty"`
}

// GiphyTime sends time back in the format:
//    2015-08-22 15:23:22 and that trips out the default
// JSON unmarshaling, so make a custom unmarshaling.
type GiphyTime time.Time

const (
	giphyTimeFormat   = "2006-01-02 15:04:05"
	blankGiphyTimeStr = "0000-00-00 00:00:00"
)

func blankTimeStr(str string) bool {
	return str == "" || str == blankGiphyTimeStr
}

func (gt *GiphyTime) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	if blankTimeStr(unquoted) {
		return nil
	}
	t, err := time.Parse(giphyTimeFormat, unquoted)
	if err != nil {
		return err
	}
	*gt = GiphyTime(t)
	return nil
}

type Pagination struct {
	TotalCount uint64 `json:"total_count,omitempty"`
	Offset     uint64 `json:"offset,omitempty"`
	Count      uint64 `json:"count,omitempty"`
}

type Response struct {
	Giphs      []*Giph                `json:"data,omitempty"`
	Pagination *Pagination            `json:"pagination,omitempty"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
}

type Client struct {
	sync.RWMutex

	rt     http.RoundTripper
	apiKey string
}

func (c *Client) httpClient() *http.Client {
	c.RLock()
	defer c.RUnlock()

	rt := c.rt
	return &http.Client{Transport: rt}
}

func (c *Client) _apiKey() string {
	c.RLock()
	defer c.RUnlock()

	return c.apiKey
}

func (c *Client) SetAPIKey(key string) {
	c.Lock()
	defer c.Unlock()

	c.apiKey = key
}

func (c *Client) SetHTTPRoundTripper(rt http.RoundTripper) {
	c.Lock()
	defer c.Unlock()

	c.rt = rt
}

type ResponsePager struct {
	Cancel func() error `json:"-"`
	Pages  <-chan *Page `json:"-"`
}

type Page struct {
	Giphs []*Giph `json:"giphs"`
	Err   error   `json:"error"`

	PageNumber uint64 `json:"page_number"`
}

var errAlreadyClosed = errors.New("already closed")

func makeCanceler() (chan bool, func() error) {
	ch := make(chan bool)
	var closeOnce sync.Once
	closeFn := func() error {
		var err error = errAlreadyClosed
		closeOnce.Do(func() {
			close(ch)
			err = nil
		})
		return err
	}

	return ch, closeFn
}

type pager struct {
	Query    string    `json:"q"`
	Limit    uint64    `json:"limit"`
	Rating   Rating    `json:"rating"`
	Format   Format    `json:"fmt"`
	Offset   uint64    `json:"offset"`
	Language Language  `json:"lang"`
	SortBy   SortOrder `json:"sort"`
	Tag      string    `json:"tag"`
}

var errEmptyResponse = errors.New("could not parse the response from the server")
var blankGiph Giph

type giphWrap struct {
	Giph *Giph `json:"data"`
}

func (c *Client) RandomSticker(req *Request) (*Giph, error) {
	return c.randomGIF(req, "/stickers/random")
}

func (c *Client) RandomGIF(req *Request) (*Giph, error) {
	return c.randomGIF(req, "/gifs/random")
}

func (c *Client) randomGIF(req *Request, route string) (*Giph, error) {
	if req == nil {
		req = new(Request)
	}
	pager := &pager{
		Rating: req.Rating,
		Format: req.Format,
		Tag:    req.Tag,
	}
	qv, err := otils.ToURLValues(pager)
	if err != nil {
		return nil, err
	}
	qv.Set("api_key", c._apiKey())
	theURL := fmt.Sprintf("%s%s?%s", baseURL, route, qv.Encode())
	return c.fetchGIF(theURL)
}

func (c *Client) fetchGIF(theURL string) (*Giph, error) {
	httpReq, err := http.NewRequest("GET", theURL, nil)
	if err != nil {
		return nil, err
	}
	slurp, _, err := c.doHTTPReq(httpReq)
	if err != nil {
		return nil, err
	}
	gWrap := new(giphWrap)
	if err := json.Unmarshal(slurp, gWrap); err != nil {
		return nil, err
	}
	if reflect.DeepEqual(*gWrap.Giph, blankGiph) {
		return nil, errEmptyResponse
	}
	return gWrap.Giph, nil
}

func (c *Client) GIFByID(id string) (*Giph, error) {
	qv := make(url.Values)
	qv.Set("api_key", c._apiKey())
	theURL := fmt.Sprintf("%s/gifs/%s?%s", baseURL, id, qv.Encode())
	return c.fetchGIF(theURL)
}

func (c *Client) Trending(req *Request) (*ResponsePager, error) {
	return c.fetch(req, "/gifs/trending")
}

func (c *Client) TrendingStickers(req *Request) (*ResponsePager, error) {
	return c.fetch(req, "/stickers/trending")
}

func (c *Client) SearchStickers(req *Request) (*ResponsePager, error) {
	return c.fetch(req, "/stickers/search")
}

func (c *Client) Search(req *Request) (*ResponsePager, error) {
	return c.fetch(req, "/gifs/search")
}

func (c *Client) fetch(req *Request, route string) (*ResponsePager, error) {
	if req == nil {
		req = new(Request)
	}

	maxPage := req.MaxPageNumber
	pageExceeds := func(page uint64) bool {
		if maxPage <= 0 {
			return false
		}
		return page >= maxPage
	}

	cancelChan, cancelFn := makeCanceler()
	pagesChan := make(chan *Page, 1)

	throttleDuration := 150 * time.Millisecond
	if req.ThrottleDurationMs == NoThrottle {
		throttleDuration = 0
	} else if req.ThrottleDurationMs > 0 {
		throttleDuration = time.Duration(req.ThrottleDurationMs) * time.Millisecond
	}

	go func() {
		defer close(pagesChan)

		pageNumber := uint64(0)
		offset := uint64(0)

		for {
			pager := &pager{
				Limit:    req.LimitPerPage,
				Rating:   req.Rating,
				Format:   req.Format,
				Offset:   offset,
				Query:    req.Query,
				SortBy:   req.SortBy,
				Language: req.Language,
			}

			page := &Page{PageNumber: pageNumber}
			qv, err := otils.ToURLValues(pager)
			if err != nil {
				page.Err = err
				pagesChan <- page
				return
			}
			qv.Set("api_key", c._apiKey())

			theURL := fmt.Sprintf("%s%s?%s", baseURL, route, qv.Encode())
			req, err := http.NewRequest("GET", theURL, nil)
			if err != nil {
				page.Err = err
				pagesChan <- page
				return
			}
			slurp, _, err := c.doHTTPReq(req)
			if err != nil {
				page.Err = err
				pagesChan <- page
				return
			}
			res := new(Response)
			if err := json.Unmarshal(slurp, res); err != nil {
				page.Err = err
				pagesChan <- page
				return
			}

			if len(res.Giphs) == 0 {
				// No more results here
				return
			}
			page.Giphs = res.Giphs
			pagesChan <- page

			pageNumber += 1
			if pageExceeds(pageNumber) {
				return
			}

			select {
			case <-cancelChan:
				return
			case <-time.After(throttleDuration):
			}

			if res.Pagination != nil {
				offset += res.Pagination.Count
			}
		}
	}()

	return &ResponsePager{Pages: pagesChan, Cancel: cancelFn}, nil
}

func (c *Client) doHTTPReq(req *http.Request) ([]byte, http.Header, error) {
	res, err := c.httpClient().Do(req)
	if err != nil {
		return nil, nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	if !otils.StatusOK(res.StatusCode) {
		// TODO: Perhaps read the body
		// in case they sent content in there.
		return nil, res.Header, errors.New(res.Status)
	}

	slurp, err := ioutil.ReadAll(res.Body)
	return slurp, res.Header, err
}

const publicAPIKey = "dc6zaTOxFJmzC"

func NewClientFromEnvOrDefault() (*Client, error) {
	apiKey := strings.TrimSpace(os.Getenv("GIPHY_API_KEY"))
	if apiKey == "" {
		apiKey = publicAPIKey
	}
	return &Client{apiKey: apiKey}, nil
}

var errBlankAPIKey = errors.New("expecting a non-blank API key")

func NewClient(apiKey string) (*Client, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, errBlankAPIKey
	}
	return &Client{apiKey: apiKey}, nil
}

// Sort order
type SortOrder string

const (
	SortRecent   SortOrder = "recent"
	SortRelevant SortOrder = "relevant"
)
