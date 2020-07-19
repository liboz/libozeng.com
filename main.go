package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func indexLayoutStart(title string) string {
	return `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>` + title + `</title>
		<style>
			div.paragraph {
				display: flex;
			}
			div.date {
				width: 6em;
				padding-bottom: 1em;
			}
			@media only screen and (max-device-width: 667px) {
				div.date {
				width: 12em;
				}
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
			<style>
			</style>
			<script>
			</script>
		</head>
		<body>
			`
}

func getFile(f string) []byte {
	b, err := ioutil.ReadFile(f)

	if err != nil {
		panic(err)
	}

	return b
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
			I'm Libo Zeng and this is my personal web site. For now this will just have some random translations
			</p>
	`)
	b.WriteString(`
			<p>
			私はリーボと言います。これは私の個人サイトです。とりあえず、雑な翻訳を投稿します。
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
	for id, metaInfo := range postNameMap {
		b.WriteString(`
			<div class="paragraph">
				<div class="date">` + metaInfo.date + `</div>
				<a href="` + id + `">
				` + metaInfo.blurb +
			`</a>
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
	date := fi.Name()[0:9]

	file, err := os.Open("posts/" + id)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	title := scanner.Text()

	blurbString := ""
	start := false
	for scanner.Scan() {
		text := scanner.Text()
		if text == "</blurb>" {
			break
		}
		if start {
			blurbString += text
		}
		if text == "<blurb>" {
			start = true
		}

	}
	metaInfo := postMeta{}
	metaInfo.date = date
	metaInfo.title = title
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
		for scanner.Scan() {
			b.WriteString("\n")
			b.WriteString("			" + scanner.Text())
		}
		b.WriteString(`
			<p>
				Any errata or comments can be made by <a href="https://github.com/liboz/libozeng.com/pulls">sending me a pull request</a>.
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
	os.MkdirAll("public", 0755)
}

func addMiscFiles() {
	misc := getDir("misc")

	for i := 0; i < len(misc); i++ {
		Copy("misc/"+misc[i].Name(), "public/"+misc[i].Name())
	}
}

func main() {
	createFilesAndDirs()
	addMiscFiles()
	postNameMap := make(map[string]postMeta)
	generatePosts(postNameMap)
	generateIndex(postNameMap)
}
