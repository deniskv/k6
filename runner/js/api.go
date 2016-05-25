package js

import (
	"errors"
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"gopkg.in/olebedev/go-duktape.v2"
	"time"
)

type apiFunc func(r *Runner, c *duktape.Context, ch chan<- runner.Result) int

func apiHTTPDo(r *Runner, c *duktape.Context, ch chan<- runner.Result) int {
	method := argString(c, 0)
	if method == "" {
		ch <- runner.Result{Error: errors.New("Missing method in http call")}
		return 0
	}

	url := argString(c, 1)
	if url == "" {
		ch <- runner.Result{Error: errors.New("Missing URL in http call")}
		return 0
	}

	args := struct {
		Report bool `json:"report"`
	}{}
	if err := argJSON(c, 2, &args); err != nil {
		ch <- runner.Result{Error: errors.New("Invalid arguments to http call")}
		return 0
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	startTime := time.Now()
	err := r.Client.Do(req, res)
	duration := time.Since(startTime)

	if args.Report {
		ch <- runner.Result{Error: err, Time: duration}
	}

	return 0
}
