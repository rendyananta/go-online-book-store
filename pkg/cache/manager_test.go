package cache

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestManager_Del(t *testing.T) {
	type fields struct {
		driverSet map[DriverName]Driver
		config    Config
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(m Manager, t *testing.T)
		wantErr    bool
	}{
		{
			name: "can delete existing key",
			fields: fields{
				driverSet: map[DriverName]Driver{
					DrvNameDatabase: func() Driver {
						drv, err := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

						if err != nil {
							t.Errorf("error creating database driver, err: %s", err)
						}

						return drv
					}(),
				},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			beforeTest: func(m Manager, t *testing.T) {
				_ = m.Set(context.Background(), "foo", []byte("bar"), 1*time.Minute)
			},
			wantErr: false,
		},
		{
			name: "can delete non existent key",
			fields: fields{
				driverSet: map[DriverName]Driver{
					DrvNameDatabase: func() Driver {
						drv, err := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

						if err != nil {
							t.Errorf("error creating database driver, err: %s", err)
						}

						return drv
					}(),
				},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manager{
				driverSet: tt.fields.driverSet,
				config:    tt.fields.config,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(m, t)
			}

			if err := m.Del(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Get(t *testing.T) {
	type fields struct {
		driverSet map[DriverName]Driver
		config    Config
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(m Manager, t *testing.T)
		want       []byte
		wantErr    bool
	}{
		{
			name: "can get existing data",
			fields: fields{
				driverSet: map[DriverName]Driver{
					DrvNameDatabase: func() Driver {
						drv, err := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

						if err != nil {
							t.Errorf("error creating database driver, err: %s", err)
						}

						return drv
					}(),
				},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			beforeTest: func(m Manager, t *testing.T) {
				_ = m.Set(context.Background(), "foo", []byte("bar"), 1*time.Minute)
			},
			want:    []byte("bar"),
			wantErr: false,
		},
		{
			name: "can get non existent data",
			fields: fields{
				driverSet: map[DriverName]Driver{
					DrvNameDatabase: func() Driver {
						drv, err := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

						if err != nil {
							t.Errorf("error creating database driver, err: %s", err)
						}

						return drv
					}(),
				},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "can't get expired data",
			fields: fields{
				driverSet: map[DriverName]Driver{
					DrvNameDatabase: func() Driver {
						drv, err := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

						if err != nil {
							t.Errorf("error creating database driver, err: %s", err)
						}

						return drv
					}(),
				},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
			args: args{
				ctx: context.Background(),
				key: "foo",
			},
			beforeTest: func(m Manager, t *testing.T) {
				_ = m.Set(context.Background(), "foo", []byte("bar"), 500*time.Millisecond)
				time.Sleep(500 * time.Millisecond)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manager{
				driverSet: tt.fields.driverSet,
				config:    tt.fields.config,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(m, t)
			}

			got, err := m.Get(tt.args.ctx, tt.args.key)
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

func TestManager_Register(t *testing.T) {
	drv, err := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

	if err != nil {
		t.Errorf("error creating database driver, err: %s", err)
	}

	type fields struct {
		driverSet map[DriverName]Driver
		config    Config
	}
	type args struct {
		name   DriverName
		driver Driver
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Driver
	}{
		{
			name: "register driver",
			fields: fields{
				driverSet: map[DriverName]Driver{},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
			args: args{
				name:   DrvNameDatabase,
				driver: drv,
			},
			want: drv,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manager{
				driverSet: tt.fields.driverSet,
				config:    tt.fields.config,
			}
			m.Register(tt.args.name, tt.args.driver)

			got := m.driverSet[tt.args.name]
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Registerd driver got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Set(t *testing.T) {
	type fields struct {
		driverSet map[DriverName]Driver
		config    Config
	}
	type args struct {
		ctx context.Context
		key string
		val []byte
		ttl time.Duration
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(m Manager, t *testing.T)
		wantVal    []byte
		wantErr    bool
	}{
		{
			name: "can handle error unregistered driver",
			fields: fields{
				driverSet: map[DriverName]Driver{},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
			args: args{
				ctx: context.Background(),
				key: "foo",
				val: []byte("bar"),
				ttl: 1 * time.Second,
			},
			wantVal: nil,
			wantErr: true,
		},
		{
			name: "can set cache key",
			fields: fields{
				driverSet: map[DriverName]Driver{
					DrvNameDatabase: func() Driver {
						drv, err := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

						if err != nil {
							t.Errorf("error creating database driver, err: %s", err)
						}

						return drv
					}(),
				},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
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
			name: "can overwrite existing cache key",
			fields: fields{
				driverSet: map[DriverName]Driver{
					DrvNameDatabase: func() Driver {
						drv, err := NewDatabaseDriver(DriverDatabaseConfig{}, dbConnManagerMock{})

						if err != nil {
							t.Errorf("error creating database driver, err: %s", err)
						}

						return drv
					}(),
				},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
			args: args{
				ctx: context.Background(),
				key: "foo",
				val: []byte("bar"),
				ttl: 5 * time.Minute,
			},
			beforeTest: func(m Manager, t *testing.T) {
				_ = m.Set(context.Background(), "foo", []byte("buzz"), 1*time.Minute)
			},
			wantVal: []byte("bar"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manager{
				driverSet: tt.fields.driverSet,
				config:    tt.fields.config,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(m, t)
			}

			if err := m.Set(tt.args.ctx, tt.args.key, tt.args.val, tt.args.ttl); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestNewManager(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args args
		want Manager
	}{
		{
			name: "new cache manager",
			args: args{},
			want: Manager{
				driverSet: map[DriverName]Driver{},
				config: Config{
					DefaultDriver: DrvNameDatabase,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewManager(tt.args.cfg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() got = %v, want %v", got, tt.want)
			}
		})
	}
}
