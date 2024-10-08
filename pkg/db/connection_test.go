package db

import (
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
)

func TestManager_Connection(t *testing.T) {
	conn1, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("error opening connection")
	}

	type fields struct {
		connections map[string]*sqlx.DB
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sqlx.DB
		wantErr bool
	}{
		{
			name: "can get initialized connection",
			fields: fields{connections: map[string]*sqlx.DB{
				"conn1": conn1,
			}},
			args:    args{name: "conn1"},
			want:    conn1,
			wantErr: false,
		},
		{
			name: "can't get uninitialized connection",
			fields: fields{connections: map[string]*sqlx.DB{
				"conn1": conn1,
			}},
			args:    args{name: "conn2"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ConnManager{
				connections: tt.fields.connections,
			}
			got, err := m.Connection(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connection() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConnectionManager(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name       string
		args       args
		wantNotNil bool
		wantErr    bool
	}{
		{
			name: "can init db connection manager",
			args: args{
				cfg: Config{
					Connections: map[string]ConnectionConfig{},
				},
			},
			wantNotNil: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnectionManager(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnectionManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != nil) != tt.wantNotNil {
				t.Errorf("NewConnectionManager() got = %v, want not nil %v", got, tt.wantNotNil)
			}
		})
	}
}
