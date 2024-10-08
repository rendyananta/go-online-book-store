package book

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/rendyananta/example-online-book-store/internal/entity/book"
	"reflect"
	"testing"
)

func TestRepo_FindByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	type args struct {
		ctx context.Context
		id  []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       []book.Book
		wantErr    bool
	}{
		{
			name: "can get book detail by id",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			args: args{
				ctx: context.Background(),
				id:  []string{"1"},
			},
			beforeTest: func() {
				query, queryArgs, _ := sqlx.In(queryGetBookByIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(query).Return(query)

				var bookResult []tableBook
				dbConnMock.EXPECT().SelectContext(context.Background(), &bookResult, query, queryArgs).Return(nil).
					SetArg(1, []tableBook{
						{
							ID:            "1",
							Title:         "Book 1",
							Description:   "Desc",
							Price:         5.6,
							ISBN:          "123124123",
							PublisherID:   "2",
							PublisherName: "Publisher Book 1",
						},
					})

				var result []tableAuthor
				query, queryArgs, _ = sqlx.In(queryGetAuthorsByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(nil).
					SetArg(1, []tableAuthor{
						{
							ID:     "11",
							Name:   "John",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Doe",
							BookID: "1",
						},
					})

				var genreResult []tableGenre
				genreQuery, genreQueryArgs, _ := sqlx.In(queryGetGenresByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(genreQuery).Return(genreQuery)
				dbConnMock.EXPECT().SelectContext(context.Background(), &genreResult, genreQuery, genreQueryArgs).
					Return(nil).
					SetArg(1, []tableGenre{
						{
							ID:     "11",
							Name:   "Action",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Thriller",
							BookID: "1",
						},
					})
			},
			want: []book.Book{
				{
					ID:          "1",
					Title:       "Book 1",
					Description: "Desc",
					Price:       5.6,
					ISBN:        "123124123",
					Publisher: book.Publisher{
						ID:   "2",
						Name: "Publisher Book 1",
					},
					Authors: []book.Author{
						{
							ID:   "11",
							Name: "John",
						},
						{
							ID:   "12",
							Name: "Doe",
						},
					},
					Genres: []book.Genre{
						{
							ID:   "11",
							Name: "Action",
						},
						{
							ID:   "12",
							Name: "Thriller",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "can handle fail get book by id",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			args: args{
				ctx: context.Background(),
				id:  []string{"1"},
			},
			beforeTest: func() {
				query, queryArgs, _ := sqlx.In(queryGetBookByIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(query).Return(query)

				var bookResult []tableBook
				dbConnMock.EXPECT().SelectContext(context.Background(), &bookResult, query, queryArgs).Return(nil).
					SetArg(1, []tableBook{
						{
							ID:            "1",
							Title:         "Book 1",
							Description:   "Desc",
							Price:         5.6,
							ISBN:          "123124123",
							PublisherID:   "2",
							PublisherName: "Publisher Book 1",
						},
					})

				var result []tableAuthor
				query, queryArgs, _ = sqlx.In(queryGetAuthorsByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(nil).
					SetArg(1, []tableAuthor{
						{
							ID:     "11",
							Name:   "John",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Doe",
							BookID: "1",
						},
					})

				var genreResult []tableGenre
				genreQuery, genreQueryArgs, _ := sqlx.In(queryGetGenresByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(genreQuery).Return(genreQuery)
				dbConnMock.EXPECT().SelectContext(context.Background(), &genreResult, genreQuery, genreQueryArgs).
					Return(sql.ErrConnDone)
			},
			want: []book.Book{
				{
					ID:          "1",
					Title:       "Book 1",
					Description: "Desc",
					Price:       5.6,
					ISBN:        "123124123",
					Publisher: book.Publisher{
						ID:   "2",
						Name: "Publisher Book 1",
					},
					Authors: []book.Author{
						{
							ID:   "11",
							Name: "John",
						},
						{
							ID:   "12",
							Name: "Doe",
						},
					},
					Genres: nil,
				},
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

			got, err := r.FindByIDs(tt.args.ctx, tt.args.id)
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

func TestRepo_PaginateAllBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)
	preparedStmtMock := NewMockpreparedQueryGetter(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	type args struct {
		ctx   context.Context
		param book.PaginationParam
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       book.PaginationResult
		wantErr    bool
	}{
		{
			name: "can get all books using pagination",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					bookPagination: preparedStmtMock,
				},
			},
			args: args{
				ctx: context.Background(),
				param: book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				},
			},
			beforeTest: func() {
				var bookResult []tableBook
				preparedStmtMock.EXPECT().SelectContext(context.Background(), &bookResult, []interface{}{"", 1}).Return(nil).
					SetArg(1, []tableBook{
						{
							ID:            "1",
							Title:         "Book 1",
							Description:   "Desc",
							Price:         5.6,
							ISBN:          "123124123",
							PublisherID:   "2",
							PublisherName: "Publisher Book 1",
						},
					})

				var result []tableAuthor
				query, queryArgs, _ := sqlx.In(queryGetAuthorsByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(nil).
					SetArg(1, []tableAuthor{
						{
							ID:     "11",
							Name:   "John",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Doe",
							BookID: "1",
						},
					})

				var genreResult []tableGenre
				genreQuery, genreQueryArgs, _ := sqlx.In(queryGetGenresByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(genreQuery).Return(genreQuery)
				dbConnMock.EXPECT().SelectContext(context.Background(), &genreResult, genreQuery, genreQueryArgs).
					Return(nil).
					SetArg(1, []tableGenre{
						{
							ID:     "11",
							Name:   "Action",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Thriller",
							BookID: "1",
						},
					})
			},
			want: book.PaginationResult{
				Data: []book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5.6,
						ISBN:        "123124123",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher Book 1",
						},
						Authors: []book.Author{
							{
								ID:   "11",
								Name: "John",
							},
							{
								ID:   "12",
								Name: "Doe",
							},
						},
						Genres: []book.Genre{
							{
								ID:   "11",
								Name: "Action",
							},
							{
								ID:   "12",
								Name: "Thriller",
							},
						},
					},
				},
				PerPage: 1,
				LastID:  "1",
			},
			wantErr: false,
		},
		{
			name: "can handle error when get all books using pagination",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					bookPagination: preparedStmtMock,
				},
			},
			args: args{
				ctx: context.Background(),
				param: book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				},
			},
			beforeTest: func() {
				var bookResult []tableBook
				preparedStmtMock.EXPECT().SelectContext(context.Background(), &bookResult, []interface{}{"", 1}).Return(nil).
					SetArg(1, []tableBook{
						{
							ID:            "1",
							Title:         "Book 1",
							Description:   "Desc",
							Price:         5.6,
							ISBN:          "123124123",
							PublisherID:   "2",
							PublisherName: "Publisher Book 1",
						},
					})

				var result []tableAuthor
				query, queryArgs, _ := sqlx.In(queryGetAuthorsByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(nil).
					SetArg(1, []tableAuthor{
						{
							ID:     "11",
							Name:   "John",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Doe",
							BookID: "1",
						},
					})

				var genreResult []tableGenre
				genreQuery, genreQueryArgs, _ := sqlx.In(queryGetGenresByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(genreQuery).Return(genreQuery)
				dbConnMock.EXPECT().SelectContext(context.Background(), &genreResult, genreQuery, genreQueryArgs).
					Return(sql.ErrConnDone)
			},
			want: book.PaginationResult{
				Data: []book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5.6,
						ISBN:        "123124123",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher Book 1",
						},
						Authors: []book.Author{
							{
								ID:   "11",
								Name: "John",
							},
							{
								ID:   "12",
								Name: "Doe",
							},
						},
					},
				},
				PerPage: 1,
				LastID:  "1",
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
			got, err := r.PaginateAllBooks(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaginateAllBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginateAllBooks() got = %v, want %v", got, tt.want)
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
		beforeTest func(r *Repo)
		wantErr    bool
	}{
		{
			name: "boot repository",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			beforeTest: func(r *Repo) {
				dbConnMock.EXPECT().Rebind(queryPaginateAllBooks).Return(queryPaginateAllBooks)
				dbConnMock.EXPECT().Preparex(queryPaginateAllBooks).Return(&sqlx.Stmt{}, nil)

				dbConnMock.EXPECT().Rebind(queryPaginateBooksSearchResult).Return(queryPaginateBooksSearchResult)
				dbConnMock.EXPECT().Preparex(queryPaginateBooksSearchResult).Return(&sqlx.Stmt{}, nil)
			},
			wantErr: false,
		},
		{
			name: "boot repository",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			beforeTest: func(r *Repo) {
				dbConnMock.EXPECT().Rebind(queryPaginateAllBooks).Return(queryPaginateAllBooks)
				dbConnMock.EXPECT().Preparex(queryPaginateAllBooks).Return(&sqlx.Stmt{}, sql.ErrConnDone)
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
				tt.beforeTest(r)
			}

			if err := r.boot(); (err != nil) != tt.wantErr {
				t.Errorf("boot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_getAuthorsByBookIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	type args struct {
		ctx     context.Context
		bookIDs []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       []tableAuthor
		wantErr    bool
	}{
		{
			name: "can handle many book ids",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			args: args{
				ctx:     context.Background(),
				bookIDs: []string{"1", "2"},
			},
			beforeTest: func() {
				var result []tableAuthor
				query, queryArgs, _ := sqlx.In(queryGetAuthorsByBookIDs, []string{"1", "2"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(nil).
					SetArg(1, []tableAuthor{
						{
							ID:     "11",
							Name:   "John",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Doe",
							BookID: "1",
						},
						{
							ID:     "13",
							Name:   "Max",
							BookID: "2",
						},
					})
			},
			want: []tableAuthor{
				{
					ID:     "11",
					Name:   "John",
					BookID: "1",
				},
				{
					ID:     "12",
					Name:   "Doe",
					BookID: "1",
				},
				{
					ID:     "13",
					Name:   "Max",
					BookID: "2",
				},
			},
			wantErr: false,
		},
		{
			name: "can handle error get many book ids",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			args: args{
				ctx:     context.Background(),
				bookIDs: []string{"1", "2"},
			},
			beforeTest: func() {
				var result []tableAuthor
				query, queryArgs, _ := sqlx.In(queryGetAuthorsByBookIDs, []string{"1", "2"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(sql.ErrConnDone)
			},
			want:    nil,
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

			got, err := r.getAuthorsByBookIDs(tt.args.ctx, tt.args.bookIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAuthorsByBookIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAuthorsByBookIDs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_getGenresByBookIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	type args struct {
		ctx     context.Context
		bookIDs []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       []tableGenre
		wantErr    bool
	}{
		{
			name: "can handle many book ids",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			args: args{
				ctx:     context.Background(),
				bookIDs: []string{"1", "2"},
			},
			beforeTest: func() {
				var result []tableGenre
				query, queryArgs, _ := sqlx.In(queryGetGenresByBookIDs, []string{"1", "2"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(nil).
					SetArg(1, []tableGenre{
						{
							ID:     "11",
							Name:   "Action",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Thriller",
							BookID: "1",
						},
						{
							ID:     "13",
							Name:   "History",
							BookID: "2",
						},
					})
			},
			want: []tableGenre{
				{
					ID:     "11",
					Name:   "Action",
					BookID: "1",
				},
				{
					ID:     "12",
					Name:   "Thriller",
					BookID: "1",
				},
				{
					ID:     "13",
					Name:   "History",
					BookID: "2",
				},
			},
			wantErr: false,
		},
		{
			name: "can handle error get many book ids",
			fields: fields{
				cfg:          Config{},
				dbConn:       dbConnMock,
				preparedStmt: preparedStmt{},
			},
			args: args{
				ctx:     context.Background(),
				bookIDs: []string{"1", "2"},
			},
			beforeTest: func() {
				var result []tableGenre
				query, queryArgs, _ := sqlx.In(queryGetGenresByBookIDs, []string{"1", "2"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(sql.ErrConnDone)
			},
			want:    nil,
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

			got, err := r.getGenresByBookIDs(tt.args.ctx, tt.args.bookIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGenresByBookIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getGenresByBookIDs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_PaginateBookSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnMock := NewMockdbConnection(ctrl)
	preparedStmtMock := NewMockpreparedQueryGetter(ctrl)

	type fields struct {
		cfg          Config
		dbConn       dbConnection
		preparedStmt preparedStmt
	}
	type args struct {
		ctx         context.Context
		searchQuery string
		param       book.PaginationParam
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func()
		want       book.PaginationResult
		wantErr    bool
	}{
		{
			name: "can search books using pagination",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					bookSearchPagination: preparedStmtMock,
				},
			},
			args: args{
				ctx:         context.Background(),
				searchQuery: "potter",
				param: book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				},
			},
			beforeTest: func() {
				var bookResult []tableBook
				preparedStmtMock.EXPECT().SelectContext(context.Background(), &bookResult, []interface{}{"%potter%", "", "%potter%", 1}...).Return(nil).
					SetArg(1, []tableBook{
						{
							ID:            "1",
							Title:         "Book 1",
							Description:   "Desc",
							Price:         5.6,
							ISBN:          "123124123",
							PublisherID:   "2",
							PublisherName: "Publisher Book 1",
						},
					})

				var result []tableAuthor
				query, queryArgs, _ := sqlx.In(queryGetAuthorsByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(nil).
					SetArg(1, []tableAuthor{
						{
							ID:     "11",
							Name:   "John",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Doe",
							BookID: "1",
						},
					})

				var genreResult []tableGenre
				genreQuery, genreQueryArgs, _ := sqlx.In(queryGetGenresByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(genreQuery).Return(genreQuery)
				dbConnMock.EXPECT().SelectContext(context.Background(), &genreResult, genreQuery, genreQueryArgs).
					Return(nil).
					SetArg(1, []tableGenre{
						{
							ID:     "11",
							Name:   "Action",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Thriller",
							BookID: "1",
						},
					})
			},
			want: book.PaginationResult{
				Data: []book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5.6,
						ISBN:        "123124123",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher Book 1",
						},
						Authors: []book.Author{
							{
								ID:   "11",
								Name: "John",
							},
							{
								ID:   "12",
								Name: "Doe",
							},
						},
						Genres: []book.Genre{
							{
								ID:   "11",
								Name: "Action",
							},
							{
								ID:   "12",
								Name: "Thriller",
							},
						},
					},
				},
				PerPage: 1,
				LastID:  "1",
			},
			wantErr: false,
		},
		{
			name: "can handle error when search all books using pagination",
			fields: fields{
				cfg:    Config{},
				dbConn: dbConnMock,
				preparedStmt: preparedStmt{
					bookSearchPagination: preparedStmtMock,
				},
			},
			args: args{
				ctx:         context.Background(),
				searchQuery: "potter",
				param: book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				},
			},
			beforeTest: func() {
				var bookResult []tableBook
				preparedStmtMock.EXPECT().SelectContext(context.Background(), &bookResult, []interface{}{"%potter%", "", "%potter%", 1}...).Return(nil).
					SetArg(1, []tableBook{
						{
							ID:            "1",
							Title:         "Book 1",
							Description:   "Desc",
							Price:         5.6,
							ISBN:          "123124123",
							PublisherID:   "2",
							PublisherName: "Publisher Book 1",
						},
					})

				var result []tableAuthor
				query, queryArgs, _ := sqlx.In(queryGetAuthorsByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(query).Return(query)
				dbConnMock.EXPECT().SelectContext(context.Background(), &result, query, queryArgs).
					Return(nil).
					SetArg(1, []tableAuthor{
						{
							ID:     "11",
							Name:   "John",
							BookID: "1",
						},
						{
							ID:     "12",
							Name:   "Doe",
							BookID: "1",
						},
					})

				var genreResult []tableGenre
				genreQuery, genreQueryArgs, _ := sqlx.In(queryGetGenresByBookIDs, []string{"1"})
				dbConnMock.EXPECT().Rebind(genreQuery).Return(genreQuery)
				dbConnMock.EXPECT().SelectContext(context.Background(), &genreResult, genreQuery, genreQueryArgs).
					Return(sql.ErrConnDone)
			},
			want: book.PaginationResult{
				Data: []book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5.6,
						ISBN:        "123124123",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher Book 1",
						},
						Authors: []book.Author{
							{
								ID:   "11",
								Name: "John",
							},
							{
								ID:   "12",
								Name: "Doe",
							},
						},
					},
				},
				PerPage: 1,
				LastID:  "1",
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
			got, err := r.PaginateBookSearch(tt.args.ctx, tt.args.searchQuery, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaginateBookSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginateBookSearch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
