package book

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/rendyananta/example-online-book-store/internal/entity/book"
	"reflect"
	"testing"
)

func TestNewQueryUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := NewMockbookRepo(ctrl)

	type args struct {
		repo bookRepo
	}
	tests := []struct {
		name    string
		args    args
		want    *QueriesUseCase
		wantErr bool
	}{
		{
			name: "can init",
			args: args{repo: repoMock},
			want: &QueriesUseCase{
				repo: repoMock,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQueryUseCase(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewQueryUseCase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryUseCase() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueriesUseCase_DetailByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := NewMockbookRepo(ctrl)

	type fields struct {
		repo bookRepo
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
		want       book.Book
		wantErr    bool
	}{
		{
			name: "can get detail by id",
			fields: fields{
				repo: repoMock,
			},
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			beforeTest: func() {
				repoMock.EXPECT().FindByIDs(context.Background(), []string{"1"}).Return([]book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5,
						ISBN:        "12345",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher",
						},
						Authors: []book.Author{
							{
								ID:   "1",
								Name: "Author A",
							},
						},
						Genres: []book.Genre{
							{
								ID:   "2",
								Name: "Action",
							},
						},
					},
				}, nil)
			},
			want: book.Book{
				ID:          "1",
				Title:       "Book 1",
				Description: "Desc",
				Price:       5,
				ISBN:        "12345",
				Publisher: book.Publisher{
					ID:   "2",
					Name: "Publisher",
				},
				Authors: []book.Author{
					{
						ID:   "1",
						Name: "Author A",
					},
				},
				Genres: []book.Genre{
					{
						ID:   "2",
						Name: "Action",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "can handle error get detail by id",
			fields: fields{
				repo: repoMock,
			},
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			beforeTest: func() {
				repoMock.EXPECT().FindByIDs(context.Background(), []string{"1"}).Return([]book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5,
						ISBN:        "12345",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher",
						},
						Authors: []book.Author{
							{
								ID:   "1",
								Name: "Author A",
							},
						},
						Genres: nil,
					},
				}, sql.ErrConnDone)
			},
			want:    book.Book{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := QueriesUseCase{
				repo: tt.fields.repo,
			}
			if tt.beforeTest != nil {
				tt.beforeTest()
			}
			got, err := q.DetailByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetailByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetailByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueriesUseCase_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := NewMockbookRepo(ctrl)

	type fields struct {
		repo bookRepo
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
			name:   "can get all with pagination",
			fields: fields{repo: repoMock},
			args: args{
				ctx: context.Background(),
				param: book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				},
			},
			beforeTest: func() {
				repoMock.EXPECT().PaginateAllBooks(context.Background(), book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				}).Return(book.PaginationResult{
					Data: []book.Book{
						{
							ID:          "1",
							Title:       "Book 1",
							Description: "Desc",
							Price:       5,
							ISBN:        "12345",
							Publisher: book.Publisher{
								ID:   "2",
								Name: "Publisher",
							},
							Authors: []book.Author{
								{
									ID:   "1",
									Name: "Author A",
								},
							},
							Genres: []book.Genre{
								{
									ID:   "2",
									Name: "Action",
								},
							},
						},
					},
					PerPage: 1,
					LastID:  "1",
				}, nil)
			},
			want: book.PaginationResult{
				Data: []book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5,
						ISBN:        "12345",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher",
						},
						Authors: []book.Author{
							{
								ID:   "1",
								Name: "Author A",
							},
						},
						Genres: []book.Genre{
							{
								ID:   "2",
								Name: "Action",
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
			name:   "can get all with pagination",
			fields: fields{repo: repoMock},
			args: args{
				ctx: context.Background(),
				param: book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				},
			},
			beforeTest: func() {
				repoMock.EXPECT().PaginateAllBooks(context.Background(), book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				}).Return(book.PaginationResult{
					Data: []book.Book{
						{
							ID:          "1",
							Title:       "Book 1",
							Description: "Desc",
							Price:       5,
							ISBN:        "12345",
							Publisher: book.Publisher{
								ID:   "2",
								Name: "Publisher",
							},
							Authors: []book.Author{
								{
									ID:   "1",
									Name: "Author A",
								},
							},
						},
					},
					PerPage: 1,
					LastID:  "1",
				}, sql.ErrConnDone)
			},
			want: book.PaginationResult{
				Data: []book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5,
						ISBN:        "12345",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher",
						},
						Authors: []book.Author{
							{
								ID:   "1",
								Name: "Author A",
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
			q := QueriesUseCase{
				repo: tt.fields.repo,
			}

			if tt.beforeTest != nil {
				tt.beforeTest()
			}

			got, err := q.GetAll(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueriesUseCase_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := NewMockbookRepo(ctrl)

	type fields struct {
		repo bookRepo
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
			name:   "can search with pagination",
			fields: fields{repo: repoMock},
			args: args{
				ctx:         context.Background(),
				searchQuery: "potter",
				param: book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				},
			},
			beforeTest: func() {
				repoMock.EXPECT().PaginateBookSearch(context.Background(), "potter", book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				}).Return(book.PaginationResult{
					Data: []book.Book{
						{
							ID:          "1",
							Title:       "Book 1",
							Description: "Desc",
							Price:       5,
							ISBN:        "12345",
							Publisher: book.Publisher{
								ID:   "2",
								Name: "Publisher",
							},
							Authors: []book.Author{
								{
									ID:   "1",
									Name: "Author A",
								},
							},
							Genres: []book.Genre{
								{
									ID:   "2",
									Name: "Action",
								},
							},
						},
					},
					PerPage: 1,
					LastID:  "1",
				}, nil)
			},
			want: book.PaginationResult{
				Data: []book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5,
						ISBN:        "12345",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher",
						},
						Authors: []book.Author{
							{
								ID:   "1",
								Name: "Author A",
							},
						},
						Genres: []book.Genre{
							{
								ID:   "2",
								Name: "Action",
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
			name:   "can handle error search with pagination",
			fields: fields{repo: repoMock},
			args: args{
				ctx:         context.Background(),
				searchQuery: "potter",
				param: book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				},
			},
			beforeTest: func() {
				repoMock.EXPECT().PaginateBookSearch(context.Background(), "potter", book.PaginationParam{
					PerPage: 1,
					LastID:  "",
				}).Return(book.PaginationResult{
					Data: []book.Book{
						{
							ID:          "1",
							Title:       "Book 1",
							Description: "Desc",
							Price:       5,
							ISBN:        "12345",
							Publisher: book.Publisher{
								ID:   "2",
								Name: "Publisher",
							},
							Authors: []book.Author{
								{
									ID:   "1",
									Name: "Author A",
								},
							},
						},
					},
					PerPage: 1,
					LastID:  "1",
				}, sql.ErrConnDone)
			},
			want: book.PaginationResult{
				Data: []book.Book{
					{
						ID:          "1",
						Title:       "Book 1",
						Description: "Desc",
						Price:       5,
						ISBN:        "12345",
						Publisher: book.Publisher{
							ID:   "2",
							Name: "Publisher",
						},
						Authors: []book.Author{
							{
								ID:   "1",
								Name: "Author A",
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
			q := QueriesUseCase{
				repo: tt.fields.repo,
			}
			if tt.beforeTest != nil {
				tt.beforeTest()
			}
			got, err := q.Search(tt.args.ctx, tt.args.searchQuery, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
		})
	}
}
