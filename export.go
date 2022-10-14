package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

func initAnaconda() *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials(
		os.Getenv("ACCESS_TOKEN"),
		os.Getenv("ACCESS_TOKEN_SECRET"),
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_KEY_SECRET"))
}

func main() {
	api := initAnaconda()

	params := url.Values{}
	params.Set("count", "200")
	params.Set("skip_status", "false")
	params.Set("include_user_entities", "false")
	params.Set("exclude_replies", "true")
	params.Set("include_rts", "false")

	tweets, err := api.GetUserTimeline(params)
	if err != nil {
		log.Fatal(err)
	}

	for _, tweet := range tweets {

		filename := fmt.Sprintf("%s/%s.md", os.Getenv("HUGO_POST_PATH"), tweet.IdStr)
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var content []string

		content = append(content, "---")

		date := fmt.Sprintf("date: \"%s\"", tweet.CreatedAt)

		content = append(content, date)

		fullText := tweet.FullText

		// replace all " with '
		fullText = strings.ReplaceAll(fullText, "\"", "'")

		// replace all \n with empty a whitespace
		fullText = strings.ReplaceAll(fullText, "\n", " ")

		description := fmt.Sprintf("description: \"%s\"", fullText)
		title := fmt.Sprintf("title: \"%s\"", fullText)

		if len(fullText) > 60 {

			// substring the full text to a length of 60
			fullText = fullText[0:59]

			// trim the full text
			fullText = strings.TrimSpace(fullText)

			/* overwrite description and title */
			description = fmt.Sprintf("description: \"%s\"", fullText)
			title = fmt.Sprintf("title: \"%s\"", fullText)
		}

		content = append(content, description)
		content = append(content, title)

		// I think it's a good idea to save the status id
		twitterStatusId := fmt.Sprintf("statusId: \"%s\"", tweet.IdStr)
		content = append(content, twitterStatusId)

		// save all urls
		var urls []string
		for _, value := range tweet.Entities.Urls {
			urls = append(urls, fmt.Sprintf("\"%s\"", value.Expanded_url))
		}
		twitterUrls := fmt.Sprintf("twitterUrls: [%s]", strings.Join(urls, ","))
		content = append(content, twitterUrls)

		// save all hashtags
		var twitterHashtags []string
		for _, value := range tweet.Entities.Hashtags {
			twitterHashtags = append(twitterHashtags, fmt.Sprintf("\"%s\"", value.Text))
		}
		hashtags := fmt.Sprintf("tags: [%s]", strings.Join(twitterHashtags, ","))
		content = append(content, hashtags)

		// Download all images and save it to the static directory
		var images []string
		for _, value := range tweet.Entities.Media {
			_, filename := path.Split(value.Media_url_https)
			downloadImage(value.Media_url_https, fmt.Sprintf("%s/%s", os.Getenv("HUGO_PATH_STATIC"), filename))

			images = append(images, fmt.Sprintf("\"%s\"", filename))
		}
		twitterImages := fmt.Sprintf("images: [%s]", strings.Join(images, ","))
		content = append(content, twitterImages)

		content = append(content, "---")
		content = append(content, "\n")
		content = append(content, tweet.FullText)

		post := strings.Join(content, "\n")
		_, err = fmt.Fprintln(file, post)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// downloadImage save the file from url to filename
func downloadImage(url, filename string) error {
	// Get the response bytes from the url
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Response code != 200")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
