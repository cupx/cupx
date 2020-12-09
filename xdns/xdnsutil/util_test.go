// Copyright 2020 The CupX Authors.
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

package xdnsutil

import "testing"

func TestTrimSubDomain(t *testing.T) {
	type args struct {
		name string
		n    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				name: "q.w.c.cupx.net",
				n:    4,
			},
			want: "net",
		},
		{
			name: "t2",
			args: args{
				name: "q.w.c.cupx.net",
				n:    5,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimSubDomain(tt.args.name, tt.args.n); got != tt.want {
				t.Errorf("GetSubDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
