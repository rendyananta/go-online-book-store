package cache

import (
	"bytes"
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"testing"
	"time"
)

type dbConnManagerMock struct {
	connectionCallback func(name string) (*sqlx.DB, error)
}

func (d dbConnManagerMock) Connection(name string) (*sqlx.DB, error) {
	if d.connectionCallback != nil {
		return d.connectionCallback(name)
	}

	return sqlx.Open("sqlite3", ":memory:")
}

func TestNewDatabaseDriver(t *testing.T) {
	type args struct {
		config      DriverDatabaseConfig
		connManager dbConnManager
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "can handle success",
			args: args{
				config:      DriverDatabaseConfig{},
				connManager: dbConnManagerMock{},
			},
			wantErr: false,
		},
		{
			name: "can handle err",
			args: args{
				config: DriverDatabaseConfig{},
				connManager: dbConnManagerMock{
					connectionCallback: func(name string) (*sqlx.DB, error) {
						return nil, errors.New("error get connection")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDatabaseDriver(tt.args.config, tt.args.connManager)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDatabaseDriver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDriverDatabase_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name       string
		args       args
		want       []byte
		beforeTest func(d *DriverDatabase, t *testing.T)
		wantErr    bool
	}{
		{
			name: "can get existing cache",
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			beforeTest: func(d *DriverDatabase, t *testing.T) {
				_ = d.Set(context.Background(), "foo", []byte("bar"), 1*time.Minute)
			},
			want:    []byte("bar"),
			wantErr: false,
		},
		{
			name: "can get not existing cache",
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			beforeTest: func(d *DriverDatabase, t *testing.T) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "can get expiring cache",
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			beforeTest: func(d *DriverDatabase, t *testing.T) {
				_ = d.Set(context.Background(), "foo", []byte("bar"), 500*time.Millisecond)
				time.Sleep(500 * time.Millisecond)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, _ := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

			if tt.beforeTest != nil {
				tt.beforeTest(d, t)
			}

			got, err := d.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDriverDatabase_Set(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		val []byte
		ttl time.Duration
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(d *DriverDatabase, t *testing.T)
		wantVal    []byte
		wantErr    bool
	}{
		{
			name: "can set cache",
			args: args{
				ctx: context.Background(),
				key: "foo",
				val: []byte("bar"),
				ttl: 5 * time.Minute,
			},
			wantVal: []byte("bar"),
			wantErr: false,
		},
		{
			name: "can set exists cache",
			args: args{
				ctx: context.Background(),
				key: "foo",
				val: []byte("bar"),
				ttl: 5 * time.Minute,
			},
			beforeTest: func(d *DriverDatabase, t *testing.T) {
				_ = d.Set(context.Background(), "foo", []byte("buzz"), 5*time.Minute)
			},
			wantVal: []byte("bar"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, _ := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

			if tt.beforeTest != nil {
				tt.beforeTest(d, t)
			}

			if err := d.Set(tt.args.ctx, tt.args.key, tt.args.val, tt.args.ttl); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			val, err := d.Get(context.Background(), tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check Val error = %v, wantErr %v", err, tt.wantErr)
			}

			if !bytes.Equal(val, tt.wantVal) {
				t.Errorf("val = %v, wantval %v", string(val), string(tt.wantVal))
			}
		})
	}
}

func TestDriverDatabase_Del(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(d *DriverDatabase, t *testing.T)
		wantErr    bool
	}{
		{
			name: "can delete existing cache",
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			beforeTest: func(d *DriverDatabase, t *testing.T) {
				_ = d.Set(context.Background(), "foo", []byte("buzz"), 5*time.Minute)
			},
			wantErr: false,
		},
		{
			name: "can delete non existing cache",
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, _ := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

			if tt.beforeTest != nil {
				tt.beforeTest(d, t)
			}
			if err := d.Del(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
