package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/rendyananta/example-online-book-store/internal/config"
	"github.com/rendyananta/example-online-book-store/pkg/db"
	"github.com/rendyananta/example-online-book-store/pkg/log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "01/02/06"

func main() {
	appCfg := config.LoadAppConfig()

	log.SetUp(appCfg.Global.Log)

	dbManager, err := db.NewConnectionManager(appCfg.Global.DB)
	if err != nil {
		slog.Error("cannot initialize db connection manager", slog.String("err", err.Error()))
		panic(err)
	}

	defaultConn, err := dbManager.Connection(db.ConnDefault)
	if err != nil {
		slog.Error("cannot get connection")
	}

	var filename string

	args := flag.Args()
	fmt.Println("args: ", args)
	flag.StringVar(&filename, "file", "", "book seeder input file")
	flag.Parse()

	if filename == "" {
		slog.Error("filename arg is required")
		os.Exit(1)
	}

	file, err := os.Open(filename)
	if err != nil {
		slog.Error("file can't be opened", slog.String("error", err.Error()))
		os.Exit(1)
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	// prevent duplicated row
	genreIDByName := make(map[string]string)
	authorIDByName := make(map[string]string)
	publisherIDByName := make(map[string]string)
	for i := 1; i < 500; i++ {
		record := records[i]

		publisher := record[13]
		var publisherID string
		existingPublisherID, ok := publisherIDByName[publisher]
		if !ok {
			publisherUUID, err := uuid.NewV7()
			if err != nil {
				fmt.Printf("err generate publisher uuid %s\n", err)
			}

			publisherID = publisherUUID.String()

			_, err = defaultConn.Exec(
				defaultConn.Rebind(`insert into publishers (id, name) values (?, ?)`),
				publisherUUID.String(),
				publisher,
			)

			publisherIDByName[publisher] = publisherID
		} else {
			publisherID = existingPublisherID
		}

		var authors = strings.Split(record[3], ",")

		authorIDs := make([]string, 0, len(authors))
		for _, author := range authors {
			authorID, err := uuid.NewV7()
			if err != nil {
				continue
			}

			authorUUID, ok := authorIDByName[author]

			if !ok {
				authorIDByName[author] = authorID.String()
				_, err = defaultConn.Exec(
					defaultConn.Rebind(`insert into authors (id, name) values (?, ?)`),
					authorID.String(),
					author,
				)

				if err != nil {
					slog.Error("can't insert author", slog.String("error", err.Error()))
					continue
				}

				authorUUID = authorID.String()
			}

			authorIDs = append(authorIDs, authorUUID)
		}

		genreRawStr := strings.ReplaceAll(record[8], "'", "\"")

		var genres []string
		err = json.Unmarshal([]byte(genreRawStr), &genres)
		if err != nil {
			slog.Error("can't unmarshal genres", slog.String("error", err.Error()), slog.String("raw", genreRawStr))
		}

		genreIDs := make([]string, 0, len(genres))
		for _, genre := range genres {
			genreID, err := uuid.NewV7()
			if err != nil {
				continue
			}

			genreUUID, ok := genreIDByName[genre]

			if !ok {
				genreIDByName[genre] = genreID.String()
				_, err = defaultConn.Exec(
					defaultConn.Rebind(`insert into genres (id, name) values (?, ?)`),
					genreID.String(),
					genre,
				)

				if err != nil {
					slog.Error("can't insert author", slog.String("error", err.Error()))
					continue
				}

				genreUUID = genreID.String()
			}

			genreIDs = append(genreIDs, genreUUID)
		}

		id, err := uuid.NewV7()
		if err != nil {
			continue
		}

		title := record[1]
		rawRating := record[4]

		rating, err := strconv.ParseFloat(rawRating, 64)
		if err != nil {
			slog.Error("can't parse rating", slog.String("error", err.Error()), slog.String("price", rawRating))
		}

		description := record[5]
		language := record[6]
		isbn := record[7]
		rawPrice := record[24]
		if rawPrice == "" {
			continue
		}
		price, err := strconv.ParseFloat(rawPrice, 64)
		if err != nil {
			slog.Error("can't parse price", slog.String("error", err.Error()), slog.String("price", rawPrice))
		}

		edition := record[11]

		rawPages := record[12]
		pages, err := strconv.ParseInt(rawPages, 10, 64)
		if err != nil {
			slog.Error("can't parse pages", slog.String("error", err.Error()))
		}

		var publishDate string
		publishDateRaw := record[14]

		if publishDateRaw != "" {
			t, err := time.Parse(dateLayout, publishDateRaw)
			if err != nil {
				slog.Error("can't parse date", slog.String("error", err.Error()))
			}

			publishDate = t.Format(time.DateTime)
		}

		var firstPublishDate string
		firstPublishDateRaw := record[15]
		if firstPublishDateRaw != "" {
			t, err := time.Parse(dateLayout, publishDateRaw)
			if err != nil {
				slog.Error("can't parse date", slog.String("error", err.Error()))
			}

			firstPublishDate = t.Format(time.DateTime)
		}

		coverImg := record[21]

		_, err = defaultConn.Exec(
			defaultConn.Rebind(`insert into books (id, title, description, price, isbn, language, edition, pages, publisher_id, published_at, first_published_at, cover_img, rating) 
				values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`),
			id,
			title,
			description,
			price,
			isbn,
			language,
			edition,
			pages,
			publisherID,
			publishDate,
			firstPublishDate,
			coverImg,
			rating,
		)
		if err != nil {
			continue
		}

		for _, authorID := range authorIDs {
			bookAuthorID, err := uuid.NewV7()
			if err != nil {
				continue
			}

			_, err = defaultConn.Exec(
				defaultConn.Rebind(`insert into books_authors (id, book_id, author_id) VALUES (?, ?, ?)`),
				bookAuthorID,
				id,
				authorID,
			)
			if err != nil {
				continue
			}
		}

		for _, genreID := range genreIDs {
			bookGenreID, err := uuid.NewV7()
			if err != nil {
				continue
			}

			_, err = defaultConn.Exec(
				defaultConn.Rebind(`insert into books_genres (id, book_id, genre_id) VALUES (?, ?, ?)`),
				bookGenreID,
				id,
				genreID,
			)
			if err != nil {
				continue
			}
		}

		//_, err = defaultConn.Exec(
		//	defaultConn.Rebind(`insert into book_search_index (id, title, description, genres, authors, publisher)
		//		values (?, ?, ?, ?, ?, ?)`),
		//	id,
		//	title,
		//	description,
		//	strings.Join(genres, ","),
		//	strings.Join(authors, ","),
		//	publisher,
		//)
		//if err != nil {
		//	slog.Error("error insert into virtual table", slog.String("error", err.Error()))
		//}
	}
}
