// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import "testing"

func TestCodeAndCase(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name:  "rune exists (lowercase)",
			args:  args{r: 'a'},
			want:  30,
			want1: false,
		},
		{
			name:  "rune exists (uppercase)",
			args:  args{r: 'A'},
			want:  30,
			want1: true,
		},
		{
			name:  "rune does not exist",
			args:  args{r: '🦍'},
			want:  0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CodeAndCase(tt.args.r)
			if got != tt.want {
				t.Errorf("CodeAndCase() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CodeAndCase() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
