package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

const layout = "2006-1-2"

type FileInfo struct {
	fileName string
	date     time.Time
}

func findEndOfDate(fileName string) string {
	count := 0
	accumulatingStr := ""
	for _, char := range fileName {
		strChar := string(char)
		if strChar == "-" {
			count += 1
			if count == 3 {
				return accumulatingStr
			}
		}
		accumulatingStr += strChar
	}
	panic("didn't find date like thing")
}

func indexLayoutStart(title string) string {
	return `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>` + title + `</title>
		<link rel="stylesheet" href="sakura-dark.css" type="text/css">
		<style>
			div.paragraph {
				display: flex;
			}
			div.paragraph div.article {
				margin-bottom: 1em;
				flex: 3;
			}
			div.date {
				flex: 1;
			}
		</style>
		<script>
		</script>
	</head>
	<body>
		<div class="main">`
}

func layoutEnd() string {
	return `
		</div>
		<!-- The core Firebase JS SDK is always required and must be listed first -->
		<script src="https://www.gstatic.com/firebasejs/7.17.1/firebase-app.js"></script>

		<!-- TODO: Add SDKs for Firebase products that you want to use
			https://firebase.google.com/docs/web/setup#available-libraries -->
		<script src="https://www.gstatic.com/firebasejs/7.17.1/firebase-analytics.js"></script>

		<script>
		// Your web app's Firebase configuration
		var firebaseConfig = {
			apiKey: "AIzaSyDTRHph0ykvK6jbqH1RQCm0ReoRrS4jXe4",
			authDomain: "libo-283700.firebaseapp.com",
			databaseURL: "https://libo-283700.firebaseio.com",
			projectId: "libo-283700",
			storageBucket: "libo-283700.appspot.com",
			messagingSenderId: "125580517524",
			appId: "1:125580517524:web:7f401d13a30631537ed481",
			measurementId: "G-6YBSVRX3BC"
		};
		// Initialize Firebase
		firebase.initializeApp(firebaseConfig);
		firebase.analytics();
		</script>
	</body>
</html>`
}

func postLayoutStart(title string) string {
	return `<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<title>` + title + `</title>
			<link rel="stylesheet" href="sakura-dark.css" type="text/css">
			<style>
			</style>
			<script>
			</script>
		</head>
		<body>
			`
}

func getDir(dir string) []os.FileInfo {
	p, err := ioutil.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	return p
}

func writeFile(fileName string, b bytes.Buffer) {
	err := ioutil.WriteFile("public/"+fileName, b.Bytes(), 0644)

	if err != nil {
		panic(err)
	}
}

func writeIntro(b *bytes.Buffer) {
	b.WriteString(`
			<p>
			I'm Libo Zeng and this is my personal web site. I might post random translations as well as miscellaneous tech related stuff. Opinions here are my own.
			</p>
	`)
	b.WriteString(`
			<p>
			私はリーボと言います。これは私の個人サイトです。雑な翻訳やテクノロジーに関する物を投稿するつもりです。あくまで個人的な意見です。
			</p>
	`)
}

func generateIndex(postNameMap map[string]postMeta) {
	var b bytes.Buffer
	b.WriteString(indexLayoutStart("Libo Zeng"))
	writeIntro(&b)
	writePostsSection(&b, postNameMap)
	b.WriteString(layoutEnd())
	writeFile("index.html", b)
}

func writePostsSection(b *bytes.Buffer, postNameMap map[string]postMeta) {
	fileNames := make([]FileInfo, 0)
	for id := range postNameMap {
		dateBase := findEndOfDate(id)
		date, err := time.Parse(layout, dateBase)
		if err != nil {
			panic("Time parse error")
		}
		fileNames = append(fileNames, FileInfo{fileName: id, date: date})
	}

	sort.Slice(fileNames, func(i, j int) bool {
		return fileNames[i].date.Before(fileNames[j].date)
	})

	for i := len(fileNames) - 1; i >= 0; i-- {
		id := fileNames[i].fileName
		metaInfo := postNameMap[id]
		b.WriteString(`
			<div class="paragraph">
				<div class="date">` + metaInfo.date + `</div>
				<div class="article"><a href="` + id + `">
				` + metaInfo.blurb +
			`</a></div>
			</div>
	  `)
	}
}

type postMeta struct {
	date  string
	title string
	blurb string
}

func getPostMeta(fi os.FileInfo) (string, postMeta, *os.File, *bufio.Scanner) {
	id := fi.Name()
	re := regexp.MustCompile(`(\d+-)+`)
	match := re.FindStringSubmatch(fi.Name())
	date := match[0][0 : len(match[0])-1]

	file, err := os.Open("posts/" + id)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	blurbString := ""
	for scanner.Scan() {
		text := scanner.Text()
		blurbString += text
		if strings.Contains(blurbString, "</blurb>") {
			break
		}
	}
	blurbString = strings.Trim(blurbString[7:len(blurbString)-8], " ") // remove <blurb> form the start and </blurb> from the end
	metaInfo := postMeta{}
	metaInfo.date = date
	// len(id) - 5 to remove the .html

	splitStrings := strings.Split(id[len(date):len(id)-5], "-")
	titleString := ""
	for _, word := range splitStrings {
		titleString += strings.Title(word) + " "
	}
	titleString += date
	metaInfo.title = titleString
	metaInfo.blurb = blurbString

	return id, metaInfo, file, scanner
}

func generatePosts(postNameMap map[string]postMeta) {
	posts := getDir("posts")

	for i := 0; i < len(posts); i++ {
		id, metaInfo, file, scanner := getPostMeta(posts[i])
		postNameMap[id] = metaInfo

		var b bytes.Buffer
		b.WriteString(postLayoutStart(metaInfo.title))
		b.WriteString("<p>" + metaInfo.date + "</p>")
		b.WriteString("<h2>" + metaInfo.blurb + "</h2>")
		preEncountered := false
		for scanner.Scan() {
			text := scanner.Text()
			if strings.Contains(text, "<pre>") {
				preEncountered = true
			}
			spacing := "			"
			if preEncountered {
				spacing = ""
			}
			b.WriteString("\n")
			b.WriteString(spacing + scanner.Text())
			if strings.Contains(text, "</pre>") {
				preEncountered = false
			}
		}
		b.WriteString("<hr></hr>")
		b.WriteString(`
			<p>
				Any error corrections or comments can be made by <a href="https://github.com/liboz/libozeng.com/pulls">sending me a pull request</a>.
			</p>
		`)
		b.WriteString("<p><a href=\"index.html\">Back</a></p>")
		b.WriteString(layoutEnd())

		writeFile(id, b)
		file.Close()

	}
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func createFilesAndDirs() {
	os.MkdirAll("public/assets", 0755)
}

func addMiscFiles() {
	misc := getDir("misc")

	for i := 0; i < len(misc); i++ {
		Copy("misc/"+misc[i].Name(), "public/"+misc[i].Name())
	}
}

func addAssets() {
	assets := getDir("assets")

	for i := 0; i < len(assets); i++ {
		Copy("assets/"+assets[i].Name(), "public/assets/"+assets[i].Name())
	}
}

func main() {
	createFilesAndDirs()
	addMiscFiles()
	addAssets()
	postNameMap := make(map[string]postMeta)
	generatePosts(postNameMap)
	generateIndex(postNameMap)
}
