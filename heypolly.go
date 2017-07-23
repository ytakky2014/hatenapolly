package main

import (
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"io/ioutil"
	"os/exec"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {
	// AWSセッション作成
	sess := session.Must(session.NewSession())

	// Pollyクライアントを作成
	svc := polly.New(sess, aws.NewConfig().WithRegion("us-east-1"))

	// tagは可変させるべきだけど一旦これで
	tag := "docker"
	titles := getHatenaTitle("http://b.hatena.ne.jp/search/tag?safe=off&q=" + tag +"&users=10")

	for _, title := range titles {

		text := "タイトル : " + title
		// SynthesizeSpeechに渡すパラメータを設定
		params := &polly.SynthesizeSpeechInput{
			OutputFormat: aws.String("mp3"),
			Text:         aws.String(text),
			VoiceId:      aws.String("Mizuki"),
		}

		// polly.SynthesizeSpeechを実行
		resp, err := svc.SynthesizeSpeech(params)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// 結果をmp3ファイルで出力
		content, err := ioutil.ReadAll(resp.AudioStream)
		ioutil.WriteFile("/tmp/gopolly.mp3", content, os.ModePerm)

		// mp3ファイルを再生
		exerr := exec.Command("afplay", "/tmp/gopolly.mp3").Run()
		if exerr != nil {
			fmt.Println(exerr.Error())
			return
		}
	}
}

func getHatenaTitle(url string) []string{
	doc, _ := goquery.NewDocument(url)
	titles := []string{}
	doc.Find("li.search-result span.users a").Each(func(_ int, s *goquery.Selection) {
		title, _ := s.Attr("title")
		title = strings.Replace(title, "はてなブックマーク - ", "", 1)
		titles = append(titles, title)
	})
	return titles
}
