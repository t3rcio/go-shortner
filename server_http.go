package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/t3rcio/go-shortner/settings"
)

const _VERSION = "1.0"
const DEFAULT_LOG_LEVEL = 0
const LOG_INFO_LEVEL = 1
const LOG_WARNING_LEVEL = 2
const LOG_ERROR_LEVEL = 3
const STDOUT = ""
const WEB_SERVER_PORT = "8000"
const URL_SHORT_PATTERN = "/s/"
const REDIS_HOST = "localhost"
const REDIS_PORT = "6379"
const REDIS_DEFAULT_DB = 0
const REDIS_EXPIRATION_KEY_TIME = 0 // no expiration time
const REDIS_KEY_LENGHT = 8
const CHARSET = "abcdefghijklmnopqrstuvwxyz"

var (
	port         string
	redis_server string
)

func logger(msg string, level int, file string) {
	/*
	* A simple logger
	 */
	if level > DEFAULT_LOG_LEVEL {
		/*
		* TODO: add the file operation here
		 */

		fmt.Println(msg)
	}
}

func GenerateARandonKey(lenght int) string {
	var result []byte
	for i := 0; i < lenght; i++ {
		j := rand.Int31n(26)
		result = append(result, CHARSET[j])
	}
	return decodeBytes(result)
}

func getURLFromKey(conn *redis.Client, ctx *context.Context, path string) string {
	value, err := conn.Get(*ctx, path).Result()
	if err != nil {
		fmt.Println(err)
		value = ""
	}
	return value
}

func getRedisKeys(conn *redis.Client, ctx *context.Context) []string {
	/*
	* Return the keys on redis server
	 */
	keys, err := conn.Keys(*ctx, "*").Result()
	if err != nil {
		panic(err)
	}
	return keys
}

func setRedisKey(conn *redis.Client, ctx *context.Context, key string, value string) bool {
	/*
	* Storages a url
	 */
	err := conn.Set(*ctx, key, value, REDIS_EXPIRATION_KEY_TIME).Err()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func decodeBytes(bytesSet []byte) string {
	/*
	* Converts a bytes slice into a string
	 */
	var result []string
	for _, value := range bytesSet {
		result = append(result, strings.ToValidUTF8(string(value), ""))
	}
	return strings.Join(result, "")
}

func loadHtmlFile(htmlFile string) string {
	/*
	* Loads the htmlFile content into a string
	 */
	_file := settings.TEMPLATES_ROOT + "/" + htmlFile
	content, err := os.ReadFile(_file)
	if err != nil {
		logger("Err: HTML File not found", LOG_WARNING_LEVEL, STDOUT)
	}
	_content := decodeBytes(content)
	return _content
}

func main() {
	flag.StringVar(&port, "port", WEB_SERVER_PORT, "Port listen to")
	flag.StringVar(&redis_server, "redisserver", REDIS_HOST, "Redis server")
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_server + ":" + REDIS_PORT,
		Password: "",
		DB:       REDIS_DEFAULT_DB,
	})
	ctx := context.Background()

	// Root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
	})

	// Home
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		var htmlContent string
		w.Header().Add("version", _VERSION)
		htmlContent = loadHtmlFile(settings.INDEX_FILE)
		fmt.Fprint(w, htmlContent)
	})

	// POST a URL
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				panic(err)
			}
			url := r.PostForm.Get("url")
			if url != "" {
				_key := GenerateARandonKey(REDIS_KEY_LENGHT)
				result := setRedisKey(rdb, &ctx, _key, url)
				if result {
					fmt.Fprint(w, _key)
					return
				}
				fmt.Fprint(w, "500")
			}
			fmt.Fprint(w, "404")
		}
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		keys := getRedisKeys(rdb, &ctx)
		fmt.Fprint(w, keys)
	})

	// Redirect shorts
	http.HandleFunc(URL_SHORT_PATTERN, func(w http.ResponseWriter, r *http.Request) {
		links := getRedisKeys(rdb, &ctx)
		link := strings.Split(r.URL.Path, URL_SHORT_PATTERN)[1]
		_url := getURLFromKey(rdb, &ctx, link)

		logger("Request: "+r.URL.Path, LOG_INFO_LEVEL, STDOUT)
		logger("Links: "+strings.Join(links, ", "), LOG_INFO_LEVEL, STDOUT)
		logger("URL from KEY: "+link+": "+_url, LOG_INFO_LEVEL, STDOUT)

		if _url != "" {
			http.Redirect(w, r, _url, http.StatusMovedPermanently)
		}
		fmt.Fprint(w, "404")
	})

	httpErr := http.ListenAndServe(":"+port, nil)
	if httpErr != nil {
		fmt.Print(httpErr)
	}

}
