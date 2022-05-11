package ctx

import (
	"context"
	"net/http"
	"net/textproto"
	"strings"

	"gin-boilerplate/comm/context/metadata"
	"gin-boilerplate/comm/trace"

	"github.com/teris-io/shortid"
)

func FromContext(ctx context.Context, patchMd ...metadata.Metadata) context.Context {
	if _, _, ok := trace.FromContext(ctx); !ok {
		ctx = trace.ToContext(ctx, shortid.MustGenerate(), shortid.MustGenerate())
	}

	for i := range patchMd {
		ctx = metadata.MergeContext(ctx, patchMd[i], true)
	}
	return ctx
}

func FromRequest(r *http.Request) context.Context {
	md, ok := metadata.FromContext(r.Context())
	if !ok {
		md = make(metadata.Metadata)
	}

	for k, v := range r.Header {
		md[textproto.CanonicalMIMEHeaderKey(k)] = strings.Join(v, ",")
	}

	md["Host"] = r.Host
	md["Method"] = r.Method
	if r.URL != nil {
		md["URL"] = r.URL.String()
	}

	ctx := FromContext(r.Context(), md)
	traceID, parentSpanID, _ := trace.FromContext(ctx)

	if v := r.Header.Get(trace.TraceIDKey); len(v) <= 0 {
		r.Header.Set(trace.TraceIDKey, traceID)
	}

	if v := r.Header.Get(trace.SpanIDKey); len(v) <= 0 {
		r.Header.Set(trace.SpanIDKey, parentSpanID)
	}

	return ctx
}
