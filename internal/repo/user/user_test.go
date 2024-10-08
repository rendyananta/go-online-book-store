package user

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/rendyananta/example-online-book-store/internal/entity/user"
	"reflect"
	"testing"
	"time"
)

func TestRepo_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)
	preparedStmtMock := NewMockqueryGetter(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	type args struct {
		ctx   context.Context
		param user.User
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       user.User
		wantErr    bool
	}{
		{
			name: "can create user",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					findByEmailStmt: preparedStmtMock,
					findByIDStmt:    preparedStmtMock,
				},
			},
			args: args{
				ctx: context.Background(),
				param: user.User{
					Name:     "example",
					Email:    "example@email.com",
					Password: "hashed-password",
				},
			},
			beforeTest: func() {
				dbConnMock.EXPECT().
					ExecContext(context.Background(), queryInsertUser, gomock.AssignableToTypeOf(""), "example", "example@email.com", "hashed-password", gomock.AssignableToTypeOf(time.Time{}), gomock.AssignableToTypeOf(time.Time{})).
					Return(&sqlite3.SQLiteResult{}, nil)
			},
			want: user.User{
				ID:       "123",
				Name:     "example",
				Email:    "example@email.com",
				Password: "hashed-password",
			},
			wantErr: false,
		},
		{
			name: "can handle error when creating user",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					findByEmailStmt: preparedStmtMock,
					findByIDStmt:    preparedStmtMock,
				},
			},
			args: args{
				ctx: context.Background(),
				param: user.User{
					Name:     "example",
					Email:    "example@email.com",
					Password: "hashed-password",
				},
			},
			beforeTest: func() {
				dbConnMock.EXPECT().
					ExecContext(context.Background(), queryInsertUser, gomock.AssignableToTypeOf(""), "example", "example@email.com", "hashed-password", gomock.AssignableToTypeOf(time.Time{}), gomock.AssignableToTypeOf(time.Time{})).
					Return(&sqlite3.SQLiteResult{}, sql.ErrConnDone)
			},
			want:    user.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repo{
				cfg:          tt.fields.cfg,
				dbConn:       tt.fields.dbConn,
				preparedStmt: tt.fields.preparedStmt,
			}

			if tt.beforeTest != nil {
				tt.beforeTest()
			}

			got, err := r.Create(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Name != tt.want.Name {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}

			if got.Email != tt.want.Email {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}

			if got.Password != tt.want.Password {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_FindByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)
	preparedStmtMock := NewMockqueryGetter(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       user.User
		wantErr    bool
	}{
		{
			name: "can handle get by email",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					findByEmailStmt: preparedStmtMock,
					findByIDStmt:    preparedStmtMock,
				},
			},
			args: args{
				ctx:   context.Background(),
				email: "example@email.com",
			},
			beforeTest: func() {
				var res tableUser
				preparedStmtMock.EXPECT().
					GetContext(context.Background(), &res, "example@email.com").
					Return(nil).
					SetArg(1, tableUser{
						ID:       "123",
						Name:     "Example",
						Email:    "example@email.com",
						Password: "hashed-password",
					})
			},
			want: user.User{
				ID:       "123",
				Name:     "Example",
				Email:    "example@email.com",
				Password: "hashed-password",
			},
			wantErr: false,
		},
		{
			name: "can handle error when getting user by email",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					findByEmailStmt: preparedStmtMock,
					findByIDStmt:    preparedStmtMock,
				},
			},
			args: args{
				ctx:   context.Background(),
				email: "example@email.com",
			},
			beforeTest: func() {
				var res tableUser
				preparedStmtMock.EXPECT().
					GetContext(context.Background(), &res, "example@email.com").
					Return(sql.ErrConnDone)
			},
			want:    user.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repo{
				cfg:          tt.fields.cfg,
				dbConn:       tt.fields.dbConn,
				preparedStmt: tt.fields.preparedStmt,
			}

			if tt.beforeTest != nil {
				tt.beforeTest()
			}

			got, err := r.FindByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_boot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	tests := []struct {
		name       string
		fields     fields
		beforeTest func()
		wantErr    bool
	}{
		{
			name: "can prepare stmt",
			fields: fields{
				dbConn: dbConnMock,
			},
			beforeTest: func() {
				dbConnMock.EXPECT().Rebind(queryGetUserByEmail).Return(queryGetUserByEmail)
				dbConnMock.EXPECT().Preparex(queryGetUserByEmail).Return(&sqlx.Stmt{}, nil)

				dbConnMock.EXPECT().Rebind(queryGetUserByID).Return(queryGetUserByID)
				dbConnMock.EXPECT().Preparex(queryGetUserByID).Return(&sqlx.Stmt{}, nil)
			},
			wantErr: false,
		},
		{
			name: "can handle error when preparing stmt",
			fields: fields{
				dbConn: dbConnMock,
			},
			beforeTest: func() {
				dbConnMock.EXPECT().Rebind(queryGetUserByEmail).Return(queryGetUserByEmail)
				dbConnMock.EXPECT().Preparex(queryGetUserByEmail).Return(&sqlx.Stmt{}, sql.ErrConnDone)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				cfg:          tt.fields.cfg,
				dbConn:       tt.fields.dbConn,
				preparedStmt: tt.fields.preparedStmt,
			}

			if tt.beforeTest != nil {
				tt.beforeTest()
			}

			if err := r.boot(); (err != nil) != tt.wantErr {
				t.Errorf("boot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)
	preparedStmtMock := NewMockqueryGetter(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       user.User
		wantErr    bool
	}{
		{
			name: "can handle get by id",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					findByEmailStmt: preparedStmtMock,
					findByIDStmt:    preparedStmtMock,
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "1234",
			},
			beforeTest: func() {
				var res tableUser
				preparedStmtMock.EXPECT().
					GetContext(context.Background(), &res, "1234").
					Return(nil).
					SetArg(1, tableUser{
						ID:       "123",
						Name:     "Example",
						Email:    "example@email.com",
						Password: "hashed-password",
					})
			},
			want: user.User{
				ID:       "123",
				Name:     "Example",
				Email:    "example@email.com",
				Password: "hashed-password",
			},
			wantErr: false,
		},
		{
			name: "can handle error when getting user by id",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					findByEmailStmt: preparedStmtMock,
					findByIDStmt:    preparedStmtMock,
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "1234",
			},
			beforeTest: func() {
				var res tableUser
				preparedStmtMock.EXPECT().
					GetContext(context.Background(), &res, "1234").
					Return(sql.ErrConnDone)
			},
			want:    user.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				cfg:          tt.fields.cfg,
				dbConn:       tt.fields.dbConn,
				preparedStmt: tt.fields.preparedStmt,
			}

			if tt.beforeTest != nil {
				tt.beforeTest()
			}

			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByIDs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
