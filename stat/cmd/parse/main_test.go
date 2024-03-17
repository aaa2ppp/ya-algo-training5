package main

import "testing"

func Test_parseCell(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"1",
			args{[]byte(`
			<div class="table__data table__data_mood_pos">+2</div>
			<div class="table__data table__data_type_time">06:58</div>
			`)},
			`+2 06:58`,
			false,
		},
		{
			"2",
			args{[]byte(`
			<div class="table__data table__data_mood_pos">+2</div><div class="table__data table__data_type_time">06:58</div>`)},
			`+2 06:58`,
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCell(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCell() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseCell() = %q, want %q", got, tt.want)
			}
		})
	}
}
