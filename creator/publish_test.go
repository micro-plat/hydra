package creator

import (
	"testing"

	"github.com/micro-plat/hydra/services"
)

func Test_conf_Pub(t *testing.T) {
	type fields struct {
		data         map[string]iCustomerBuilder
		vars         map[string]map[string]interface{}
		routerLoader func(string) *services.ORouter
	}
	type args struct {
		data         map[string]iCustomerBuilder
		platName     string
		systemName   string
		clusterName  string
		registryAddr string
		cover        bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &conf{
				data:         tt.fields.data,
				vars:         tt.fields.vars,
				routerLoader: tt.fields.routerLoader,
			}
			if err := c.Pub(tt.args.data, tt.args.platName, tt.args.systemName, tt.args.clusterName, tt.args.registryAddr, tt.args.cover); (err != nil) != tt.wantErr {
				t.Errorf("conf.Pub() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
