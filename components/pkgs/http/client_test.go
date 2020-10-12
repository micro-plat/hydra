package http

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func getTestTlsByCert() *tls.Config {
	cert, _ := tls.LoadX509KeyPair("client_test_crt.txt", "client_test_key.txt")
	return &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		Rand:               rand.Reader,
	}
}
func getTestTlsByCa() *tls.Config {
	caData, _ := ioutil.ReadFile("client_test_crs.txt")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	return &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            pool,
		Rand:               rand.Reader,
	}
}

func Test_getCert(t *testing.T) {
	type args struct {
		c *conf
	}
	tests := []struct {
		name    string
		args    args
		want    *tls.Config
		wantErr bool
	}{
		{name: "1", args: args{c: &conf{}}, want: &tls.Config{InsecureSkipVerify: true}, wantErr: false},
		{name: "2", args: args{c: &conf{Certs: []string{""}}}, want: &tls.Config{InsecureSkipVerify: true}, wantErr: false},
		{name: "3", args: args{c: &conf{Certs: []string{"client_test_crt.txt", "client_test_key.txt"}}}, want: getTestTlsByCert(), wantErr: false},
		{name: "4", args: args{c: &conf{Ca: "client_test_crs.txt"}}, want: getTestTlsByCa(), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCert(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getProxy(t *testing.T) {
	type args struct {
		c *conf
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1", args: args{c: &conf{Proxy: "http://127.0.0.1:6547"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getProxy(tt.args.c)
			url, err := got(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("getProxy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			s, _ := url.Parse(tt.args.c.Proxy)
			if !reflect.DeepEqual(url, s) {
				t.Errorf("getProxy() = %v, want %v", url, s)
			}
		})
	}
}
