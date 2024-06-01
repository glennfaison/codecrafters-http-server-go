package myhttp

import (
	"reflect"
	"testing"
)

func TestParseRequestFromString(t *testing.T) {
	type args struct {
		requestString string
	}
	tests := []struct {
		name    string
		args    args
		want    Request
		wantErr bool
	}{
		{
			name: "GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			args: args{
				requestString: "GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			},
			want: Request{
				method:   "GET",
				path:     "/index.html",
				protocol: "HTTP/1.1",
				headers: map[string]string{
					"host":       "localhost:4221",
					"user-agent": "curl/7.64.1",
					"accept":     "*/*",
				},
			},
			wantErr: false,
		},
		{
			name: "GET /echo/abc HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			args: args{
				requestString: "GET /echo/abc HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			},
			want: Request{
				method:   "GET",
				path:     "/echo/abc",
				protocol: "HTTP/1.1",
				headers: map[string]string{
					"host":       "localhost:4221",
					"user-agent": "curl/7.64.1",
					"accept":     "*/*",
				},
			},
			wantErr: false,
		},
		{
			name: "GET /user-agent HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			args: args{
				requestString: "GET /user-agent HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: foobar/1.2.3\r\nAccept: */*\r\n\r\n",
			},
			want: Request{
				method:   "GET",
				path:     "/user-agent",
				protocol: "HTTP/1.1",
				headers: map[string]string{
					"host":       "localhost:4221",
					"user-agent": "foobar/1.2.3",
					"accept":     "*/*",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequestFromString(tt.args.requestString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRequestFromString()\n%v\n\n\nwant\n%v", got, tt.want)
			}
		})
	}
}
