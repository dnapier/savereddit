package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/turnage/graw/reddit"
	"mvdan.cc/xurls/v2"
)

func get() {
	count := 0
	for count < 600 {
		harvest, err := bot.ListingWithParams("/user/"+cfg.App.Username+"/saved/", map[string]string{
			"count": strconv.Itoa(count),
			"limit": "100",
			"after": after,
			"show":  "all",
		})
		if err != nil {
			Log.Error().Err(err).Send()
		}

		count += 100

		last := len(harvest.Posts) - 1
		after = harvest.Posts[last].Name
		posts = append(posts, harvest.Posts...)
		comments = append(comments, harvest.Comments...)
	}
}
func Sort() {
	// Sort by ID field
	sort.Slice(posts, func(p, q int) bool {
		return posts[p].ID < posts[q].ID
	})

	sort.Slice(comments, func(p, q int) bool {
		return comments[p].ID < comments[q].ID
	})
}
func dedupe() {
	// Remove duplicate posts
	var Posts []*reddit.Post
	for _, post := range posts {
		if _, value := keys[post.ID]; !value {
			keys[post.ID] = true
			Posts = append(Posts, post)
		}
	}
	posts = Posts

	// Remove duplicate comments
	var Comments []*reddit.Comment
	for _, comment := range comments {
		if _, value := keys[comment.ID]; !value {
			keys[comment.ID] = true
			Comments = append(Comments, comment)
		}
	}
	comments = Comments

}
func getFullThreads() {
	for _, post := range posts {
		thread, err := bot.Thread(post.Permalink)
		if err != nil {
			Log.Error().Err(err).Send()
		}
		threads = append(threads, thread)
	}
}
func getImages() {
	// Extract urls from text
	var urls []string
	rx := xurls.Relaxed()
	for _, post := range posts {
		urlmao := rx.FindAllString(post.SelfText, -1)
		urls = append(urls, urlmao...)
	}

	sort.Slice(urls, func(p, q int) bool {
		return urls[p] < urls[q]
	})

	// remove duplicates
	keys := make(map[string]bool)
	var u []string
	for _, entry := range urls {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			u = append(u, entry)
		}
	}

	// Get only urls with 'preview' in the hostname, save images in 'images' dir
	if err := os.Mkdir("images", os.ModeDir); err != nil {
		Log.Error().Err(err).Send()
	}

	for _, v := range u {
		previewUrl, err := url.Parse(v)
		if err != nil {
			Log.Debug().Err(err).Send()
		}

		if strings.Contains(previewUrl.Hostname(), "preview") {
			if err := download(v, previewUrl.Path); err != nil {
				Log.Error().Err(err).Send()
			}
		}
	}

	Log.Debug().Int(`length`, len(u)).Send()
}
func download(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		Log.Error().Err(err).Send()
	}
	defer resp.Body.Close()

	out, err := os.Create("images/" + filename)
	if err != nil {
		Log.Error().Err(err).Send()
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
func saveFile() {
	out, err := json.Marshal(posts)
	if err != nil {
		Log.Debug().Err(err).Send()
	}

	_ = ioutil.WriteFile("Posts.json", out, 0644)

	out, err = json.Marshal(comments)
	if err != nil {
		Log.Debug().Err(err).Send()
	}
	_ = ioutil.WriteFile("Comments.json", out, 0644)

	out, err = json.Marshal(threads)
	if err != nil {
		Log.Debug().Err(err).Send()
	}
	_ = ioutil.WriteFile("Threads.json", out, 0644)
}
func readFile() {
	content, err := ioutil.ReadFile("./Posts.json")
	if err != nil {
		Log.Error().Err(err).Send()
	}

	err = json.Unmarshal(content, &posts)
	if err != nil {
		Log.Error().Err(err).Send()
	}

	content, err = ioutil.ReadFile("./Comments.json")
	if err != nil {
		Log.Error().Err(err).Send()
	}

	err = json.Unmarshal(content, &comments)
	if err != nil {
		Log.Error().Err(err).Send()
	}

	content, err = ioutil.ReadFile("./Threads.json")
	if err != nil {
		Log.Error().Err(err).Send()
	}

	err = json.Unmarshal(content, &threads)
	if err != nil {
		Log.Error().Err(err).Send()
	}
}
func insertdb() {
	for _, thread := range threads {
		s := savereddit{Name: thread.Name, Attrs: Attrs(*thread), CreatedUTC: thread.CreatedUTC}
		s.Insert()
	}
}
