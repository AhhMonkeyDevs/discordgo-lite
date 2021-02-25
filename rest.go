package discordgo

import (
	"bytes"
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

const apiVersion = 8

var apiBase = fmt.Sprintf("%s%d", "https://discord.com/api/v", apiVersion)

var buckets = make(map[string]*Bucket)

type request struct {
	url             string
	bucket          string
	callbackChannel chan []byte
	token           string
	method          string
	body            []byte
	contentType     string
}

func NewRestRequest() *request {
	return &request{
		url:    apiBase,
		method: "GET",
	}
}

func (r *request) major(path string) {
	r.url += "/" + path
	r.bucket += "/" + path
}

func (r *request) minor(path string) {
	r.url += "/" + path
	r.bucket += "/:minor"
}

func (r *request) getBucket() *Bucket {
	if val, ok := buckets[r.bucket]; ok {
		return val
	}

	newBucket := Bucket{}
	newBucket.queue = list.New()
	buckets[r.bucket] = &newBucket
	return &newBucket
}

func (r *request) Method(method string) *request {
	r.method = method
	return r
}

func (r *request) Token(token string) *request {
	r.token = token
	return r
}

func (r *request) Route(path string) *request {
	r.major(path)
	return r
}

func (r *request) Query(query string) *request{
	r.url += "?" + query
	return r
}

func (r *request) Guild(id string) *request {
	r.major(id)
	return r
}

func (r *request) Channel(id string) *request {
	r.major(id)
	return r
}

func (r *request) Id(id string) *request {
	r.minor(id)
	return r
}

func (r *request) Callback(callback chan []byte) *request {
	r.callbackChannel = callback
	return r
}
func (r *request) Body(contentType string, content []byte) *request {
	r.contentType = contentType
	r.body = content
	return r
}

func (r *request) Enqueue() {
	r.getBucket().enqueue(r)
}

func (r *request) execute() (*http.Response, []byte, error) {
	client := &http.Client{}

	var reqBody io.Reader
	if r.body != nil {
		reqBody = bytes.NewReader(r.body)
	}

	req, err := http.NewRequest(r.method, r.url, reqBody)
	req.Header.Add("Authorization", "Bot "+r.token)

	if r.contentType != "" {
		req.Header.Add("Content-Type", r.contentType)
	}

	req.Header.Add("User-Agent", "DiscordBot (https://andrewwilsonwebdesign.com, 1)")

	resp, err := client.Do(req)

	if err != nil {
		return resp, nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	r.callbackChannel <- body

	return resp, body, err
}

type Bucket struct {
	running     bool
	initialized bool
	limit       int
	remaining   int
	reset       time.Duration
	queue       *list.List
}

func (b *Bucket) enqueue(r *request) {

	b.queue.PushBack(r)

	if !b.running {
		go b.run()
	}
}

func (b *Bucket) run() {
	b.running = true
	for b.queue.Front() != nil {

		element := b.queue.Front()
		request := element.Value.(*request)

		if b.remaining == 0 && b.initialized {

			time.Sleep(b.reset)
			b.remaining = b.limit
		}

		b.remaining--
		success := b.execute(request)

		if success {
			b.queue.Remove(element)
		}

	}
	b.running = false
}

func (b *Bucket) execute(request *request) bool {

	resp, _, execErr := request.execute()

	if execErr != nil{
		return false
	}

	limit := resp.Header.Get("X-RateLimit-Limit")
	reset := resp.Header.Get("X-RateLimit-Reset-After")
	remaining := resp.Header.Get("X-RateLimit-Remaining")

	if limitNumber, err := strconv.Atoi(limit); err == nil {
		b.limit = limitNumber
	}

	if resetNumber, err := strconv.ParseFloat(reset, 64); err == nil {
		newTime := time.Unix(0, int64(resetNumber*math.Pow(10, 9)))
		b.reset = newTime.Sub(time.Unix(0, 0))
	}

	if remainingNumber, err := strconv.Atoi(remaining); err == nil {
		b.initialized = true
		b.remaining = remainingNumber
	}

	if resp.StatusCode == 429 {
		fmt.Println("Rate limited")
		return false
	}

	return true

}
