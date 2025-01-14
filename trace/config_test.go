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

package trace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestNewSpanConfig(t *testing.T) {
	k1v1 := attribute.String("key1", "value1")
	k1v2 := attribute.String("key1", "value2")
	k2v2 := attribute.String("key2", "value2")

	timestamp0 := time.Unix(0, 0)
	timestamp1 := time.Unix(0, 0)

	link1 := Link{
		SpanContext: SpanContext{traceID: TraceID([16]byte{1, 1}), spanID: SpanID{3}},
		Attributes:  []attribute.KeyValue{k1v1},
	}
	link2 := Link{
		SpanContext: SpanContext{traceID: TraceID([16]byte{1, 1}), spanID: SpanID{3}},
		Attributes:  []attribute.KeyValue{k1v2, k2v2},
	}

	tests := []struct {
		options  []SpanStartOption
		expected *SpanConfig
	}{
		{
			// No non-zero-values should be set.
			[]SpanStartOption{},
			new(SpanConfig),
		},
		{
			[]SpanStartOption{
				WithAttributes(k1v1),
			},
			&SpanConfig{
				attributes: []attribute.KeyValue{k1v1},
			},
		},
		{
			// Multiple calls should append not overwrite.
			[]SpanStartOption{
				WithAttributes(k1v1),
				WithAttributes(k1v2),
				WithAttributes(k2v2),
			},
			&SpanConfig{
				// No uniqueness is guaranteed by the API.
				attributes: []attribute.KeyValue{k1v1, k1v2, k2v2},
			},
		},
		{
			[]SpanStartOption{
				WithAttributes(k1v1, k1v2, k2v2),
			},
			&SpanConfig{
				// No uniqueness is guaranteed by the API.
				attributes: []attribute.KeyValue{k1v1, k1v2, k2v2},
			},
		},
		{
			[]SpanStartOption{
				WithTimestamp(timestamp0),
			},
			&SpanConfig{
				timestamp: timestamp0,
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls overwrites with last-one-wins.
				WithTimestamp(timestamp0),
				WithTimestamp(timestamp1),
			},
			&SpanConfig{
				timestamp: timestamp1,
			},
		},
		{
			[]SpanStartOption{
				WithLinks(link1),
			},
			&SpanConfig{
				links: []Link{link1},
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls should append not overwrite.
				WithLinks(link1),
				WithLinks(link1, link2),
			},
			&SpanConfig{
				// No uniqueness is guaranteed by the API.
				links: []Link{link1, link1, link2},
			},
		},
		{
			[]SpanStartOption{
				WithNewRoot(),
			},
			&SpanConfig{
				newRoot: true,
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls should not change NewRoot state.
				WithNewRoot(),
				WithNewRoot(),
			},
			&SpanConfig{
				newRoot: true,
			},
		},
		{
			[]SpanStartOption{
				WithSpanKind(SpanKindConsumer),
			},
			&SpanConfig{
				spanKind: SpanKindConsumer,
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls overwrites with last-one-wins.
				WithSpanKind(SpanKindClient),
				WithSpanKind(SpanKindConsumer),
			},
			&SpanConfig{
				spanKind: SpanKindConsumer,
			},
		},
		{
			// Everything should work together.
			[]SpanStartOption{
				WithAttributes(k1v1),
				WithTimestamp(timestamp0),
				WithLinks(link1, link2),
				WithNewRoot(),
				WithSpanKind(SpanKindConsumer),
			},
			&SpanConfig{
				attributes: []attribute.KeyValue{k1v1},
				timestamp:  timestamp0,
				links:      []Link{link1, link2},
				newRoot:    true,
				spanKind:   SpanKindConsumer,
			},
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, NewSpanStartConfig(test.options...))
	}
}

func TestEndSpanConfig(t *testing.T) {
	timestamp := time.Unix(0, 0)

	tests := []struct {
		options  []SpanEndOption
		expected *SpanConfig
	}{
		{
			[]SpanEndOption{},
			new(SpanConfig),
		},
		{
			[]SpanEndOption{
				WithStackTrace(true),
			},
			&SpanConfig{
				stackTrace: true,
			},
		},
		{
			[]SpanEndOption{
				WithTimestamp(timestamp),
			},
			&SpanConfig{
				timestamp: timestamp,
			},
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, NewSpanEndConfig(test.options...))
	}
}

func TestTracerConfig(t *testing.T) {
	v1 := "semver:0.0.1"
	v2 := "semver:1.0.0"
	schemaURL := "https://opentelemetry.io/schemas/1.2.0"
	tests := []struct {
		options  []TracerOption
		expected *TracerConfig
	}{
		{
			// No non-zero-values should be set.
			[]TracerOption{},
			new(TracerConfig),
		},
		{
			[]TracerOption{
				WithInstrumentationVersion(v1),
			},
			&TracerConfig{
				instrumentationVersion: v1,
			},
		},
		{
			[]TracerOption{
				// Multiple calls should overwrite.
				WithInstrumentationVersion(v1),
				WithInstrumentationVersion(v2),
			},
			&TracerConfig{
				instrumentationVersion: v2,
			},
		},
		{
			[]TracerOption{
				WithSchemaURL(schemaURL),
			},
			&TracerConfig{
				schemaURL: schemaURL,
			},
		},
	}
	for _, test := range tests {
		config := NewTracerConfig(test.options...)
		assert.Equal(t, test.expected, config)
	}
}
