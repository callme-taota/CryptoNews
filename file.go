package main

import (
	"encoding/csv"
	"io"
	"reflect"
	"strconv"
)

func PreDayOnly(list []Feed) []Feed {
	var feedList []Feed
	for _, item := range list {
		if TimeCaseCmp(JSUnixTimeToTime(int(item.PublishTime)), GetPreDay()) {
			feedList = append(feedList, item)
		}
	}
	return feedList
}

func CurrDayOnly(list []Feed) []Feed {
	var feedList []Feed
	for _, item := range list {
		if TimeCaseCmp(JSUnixTimeToTime(int(item.PublishTime)), GetCurrentDay()) {
			feedList = append(feedList, item)
		}
	}
	return feedList
}

func structToCsv(records []Feed, w io.Writer) error {
	writer := csv.NewWriter(w)

	// Write header
	header := []string{
		"ID",
		"Title",
		"Content",
		"TranslatedTitle",
		"TranslatedContent",
		"PublishTime",
		"Important",
		"SourceURL",
		"RelatedFeeds",
		"Nickname",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write records
	for _, record := range records {
		var values []string
		s := reflect.ValueOf(&record).Elem()
		for i := 0; i < s.NumField(); i++ {
			field := s.Field(i)
			var value string
			switch field.Interface().(type) {
			case int, int64:
				value = strconv.FormatInt(field.Int(), 10)
			case bool:
				value = strconv.FormatBool(field.Interface().(bool))
			default:
				value = field.Interface().(string)
			}
			values = append(values, value)
		}
		if err := writer.Write(values); err != nil {
			return err
		}
	}

	writer.Flush()
	return nil
}
