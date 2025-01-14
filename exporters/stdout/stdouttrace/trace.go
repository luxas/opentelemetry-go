// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stdouttrace // import "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var zeroTime time.Time

// Exporter is an implementation of trace.SpanSyncer that writes spans to stdout.
type traceExporter struct {
	config config

	stoppedMu sync.RWMutex
	stopped   bool
}

// ExportSpans writes spans in json format to stdout.
func (e *traceExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped {
		return nil
	}

	if len(spans) == 0 {
		return nil
	}
	stubs := tracetest.SpanStubsFromReadOnlySpans(spans)
	if !e.config.Timestamps {
		for i := range stubs {
			stub := &stubs[i]
			stub.StartTime = zeroTime
			stub.EndTime = zeroTime
			for j := range stub.Events {
				ev := &stub.Events[j]
				ev.Time = zeroTime
			}
		}
	}

	out, err := e.marshal(stubs)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(e.config.Writer, string(out))
	return err
}

// Shutdown is called to stop the exporter, it preforms no action.
func (e *traceExporter) Shutdown(ctx context.Context) error {
	e.stoppedMu.Lock()
	e.stopped = true
	e.stoppedMu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

// marshal v with approriate indentation.
func (e *traceExporter) marshal(v interface{}) ([]byte, error) {
	if e.config.PrettyPrint {
		return json.MarshalIndent(v, "", "\t")
	}
	return json.Marshal(v)
}
