// Copyright 2023 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package promql

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBucketQuantile_ForcedMonotonicity(t *testing.T) {
	eps := 1e-12

	for name, tc := range map[string]struct {
		getInput       func() Buckets // The buckets can be modified in-place so return a new one each time.
		expectedForced bool
		expectedFixed  bool
		expectedValues map[float64]float64
	}{
		"simple - monotonic": {
			getInput: func() Buckets {
				return Buckets{
					{
						UpperBound: 10,
						Count:      10,
					}, {
						UpperBound: 15,
						Count:      15,
					}, {
						UpperBound: 20,
						Count:      15,
					}, {
						UpperBound: 30,
						Count:      15,
					}, {
						UpperBound: math.Inf(1),
						Count:      15,
					},
				}
			},
			expectedForced: false,
			expectedFixed:  false,
			expectedValues: map[float64]float64{
				1:    15.,
				0.99: 14.85,
				0.9:  13.5,
				0.5:  7.5,
			},
		},
		"simple - non-monotonic middle": {
			getInput: func() Buckets {
				return Buckets{
					{
						UpperBound: 10,
						Count:      10,
					}, {
						UpperBound: 15,
						Count:      15,
					}, {
						UpperBound: 20,
						Count:      15.00000000001, // Simulate the case there's a small imprecision in float64.
					}, {
						UpperBound: 30,
						Count:      15,
					}, {
						UpperBound: math.Inf(1),
						Count:      15,
					},
				}
			},
			expectedForced: false,
			expectedFixed:  true,
			expectedValues: map[float64]float64{
				1:    15.,
				0.99: 14.85,
				0.9:  13.5,
				0.5:  7.5,
			},
		},
		"real example - monotonic": {
			getInput: func() Buckets {
				return Buckets{
					{
						UpperBound: 1,
						Count:      6454661.3014166197,
					}, {
						UpperBound: 5,
						Count:      8339611.2001912938,
					}, {
						UpperBound: 10,
						Count:      14118319.2444762159,
					}, {
						UpperBound: 25,
						Count:      14130031.5272856522,
					}, {
						UpperBound: 50,
						Count:      46001270.3030008152,
					}, {
						UpperBound: 64,
						Count:      46008473.8585563600,
					}, {
						UpperBound: 80,
						Count:      46008473.8585563600,
					}, {
						UpperBound: 100,
						Count:      46008473.8585563600,
					}, {
						UpperBound: 250,
						Count:      46008473.8585563600,
					}, {
						UpperBound: 1000,
						Count:      46008473.8585563600,
					}, {
						UpperBound: math.Inf(1),
						Count:      46008473.8585563600,
					},
				}
			},
			expectedForced: false,
			expectedFixed:  false,
			expectedValues: map[float64]float64{
				1:    64.,
				0.99: 49.64475715376406,
				0.9:  46.39671690938454,
				0.5:  31.96098248992002,
			},
		},
		"real example - non-monotonic": {
			getInput: func() Buckets {
				return Buckets{
					{
						UpperBound: 1,
						Count:      6454661.3014166225,
					}, {
						UpperBound: 5,
						Count:      8339611.2001912957,
					}, {
						UpperBound: 10,
						Count:      14118319.2444762159,
					}, {
						UpperBound: 25,
						Count:      14130031.5272856504,
					}, {
						UpperBound: 50,
						Count:      46001270.3030008227,
					}, {
						UpperBound: 64,
						Count:      46008473.8585563824,
					}, {
						UpperBound: 80,
						Count:      46008473.8585563898,
					}, {
						UpperBound: 100,
						Count:      46008473.8585563824,
					}, {
						UpperBound: 250,
						Count:      46008473.8585563824,
					}, {
						UpperBound: 1000,
						Count:      46008473.8585563898,
					}, {
						UpperBound: math.Inf(1),
						Count:      46008473.8585563824,
					},
				}
			},
			expectedForced: false,
			expectedFixed:  true,
			expectedValues: map[float64]float64{
				1:    64.,
				0.99: 49.64475715376406,
				0.9:  46.39671690938454,
				0.5:  31.96098248992002,
			},
		},
		"real example 2 - monotonic": {
			getInput: func() Buckets {
				return Buckets{
					{
						UpperBound: 0.005,
						Count:      9.6,
					}, {
						UpperBound: 0.01,
						Count:      9.688888889,
					}, {
						UpperBound: 0.025,
						Count:      9.755555556,
					}, {
						UpperBound: 0.05,
						Count:      9.844444444,
					}, {
						UpperBound: 0.1,
						Count:      9.888888889,
					}, {
						UpperBound: 0.25,
						Count:      9.888888889,
					}, {
						UpperBound: 0.5,
						Count:      9.888888889,
					}, {
						UpperBound: 1,
						Count:      9.888888889,
					}, {
						UpperBound: 2.5,
						Count:      9.888888889,
					}, {
						UpperBound: 5,
						Count:      9.888888889,
					}, {
						UpperBound: 10,
						Count:      9.888888889,
					}, {
						UpperBound: 25,
						Count:      9.888888889,
					}, {
						UpperBound: 50,
						Count:      9.888888889,
					}, {
						UpperBound: 100,
						Count:      9.888888889,
					}, {
						UpperBound: math.Inf(1),
						Count:      9.888888889,
					},
				}
			},
			expectedForced: false,
			expectedFixed:  false,
			expectedValues: map[float64]float64{
				1:    0.1,
				0.99: 0.03468750000281261,
				0.9:  0.00463541666671875,
				0.5:  0.0025752314815104174,
			},
		},
		"real example 2 - non-monotonic": {
			getInput: func() Buckets {
				return Buckets{
					{
						UpperBound: 0.005,
						Count:      9.6,
					}, {
						UpperBound: 0.01,
						Count:      9.688888889,
					}, {
						UpperBound: 0.025,
						Count:      9.755555556,
					}, {
						UpperBound: 0.05,
						Count:      9.844444444,
					}, {
						UpperBound: 0.1,
						Count:      9.888888889,
					}, {
						UpperBound: 0.25,
						Count:      9.888888889,
					}, {
						UpperBound: 0.5,
						Count:      9.888888889,
					}, {
						UpperBound: 1,
						Count:      9.888888889,
					}, {
						UpperBound: 2.5,
						Count:      9.888888889,
					}, {
						UpperBound: 5,
						Count:      9.888888889,
					}, {
						UpperBound: 10,
						Count:      9.888888889001, // Simulate the case there's a small imprecision in float64.
					}, {
						UpperBound: 25,
						Count:      9.888888889,
					}, {
						UpperBound: 50,
						Count:      9.888888888999, // Simulate the case there's a small imprecision in float64.
					}, {
						UpperBound: 100,
						Count:      9.888888889,
					}, {
						UpperBound: math.Inf(1),
						Count:      9.888888889,
					},
				}
			},
			expectedForced: false,
			expectedFixed:  true,
			expectedValues: map[float64]float64{
				1:    0.1,
				0.99: 0.03468750000281261,
				0.9:  0.00463541666671875,
				0.5:  0.0025752314815104174,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			for q, v := range tc.expectedValues {
				res, forced, fixed := BucketQuantile(q, tc.getInput())
				require.Equal(t, tc.expectedForced, forced)
				require.Equal(t, tc.expectedFixed, fixed)
				require.InEpsilon(t, v, res, eps)
			}
		})
	}
}
