package pankat

import (
	"fmt"
	"github.com/fatih/color"
	"regexp"
	"strings"
	"time"
)

func ProcessPlugins(_article []byte, article *Article) []byte {
	var _articlePostprocessed []byte

	re := regexp.MustCompile("\\[\\[!(.*?)\\]\\]")
	z := re.FindAllIndex(_article, -1)

	prevPos := 0
	var foundPlugins []string
	for i := 0; i <= len(z); i++ {
		if i == len(z) {
			_articlePostprocessed = append(_articlePostprocessed, _article[prevPos:]...)
			break
		}
		n := z[i]

		// include normal content (not plugin processed)
		if prevPos != n[0] {
			_articlePostprocessed = append(_articlePostprocessed, _article[prevPos:n[0]]...)
		}

		// include plugin processed stuff
		t, name := callPlugin(_article[n[0]:n[1]], article)
		foundPlugins = append(foundPlugins, name)
		_articlePostprocessed = append(_articlePostprocessed, t...)
		prevPos = n[1]
	}
	if Config().Verbose > 1 {
		fmt.Println(article.DstFileName, color.GreenString("plugins:"), foundPlugins)
	}
	return _articlePostprocessed
}

func callPlugin(in []byte, article *Article) ([]byte, string) {
	a := len(in) - 2
	p := string(in[3:a])
	//   fmt.Println(p)
	var output []byte

	f := strings.Fields(p)
	var name string
	if len(f) > 0 {
		name = f[0]
	} else {
		var z []byte
		return z, ""
	}

	//   fmt.Println("\n=========== ", name, " ===========")
	switch strings.ToLower(name) {
	case "specialpage":
		article.SpecialPage = true
	case "draft":
		article.Draft = true
	case "meta":
		re := regexp.MustCompile("[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+")
		z := re.FindIndex(in)
		var t time.Time
		if z != nil {
			s := string(in[z[0]:z[1]])
			//           fmt.Println(s)
			const longForm = "2006-01-02 15:04"
			t, _ = time.Parse(longForm, s)
			article.ModificationDate = t
			//           fmt.Println(t)
		}
		// 	case "warning":
		//     if len(f) > 1 {
		//       o := `<div id="bar">` + strings.Join(f[1:len(f)], " ") + `</div>`
		//       output = []byte(o)
		//     }
	case "series":
		if len(f) > 1 {
			article.Series = strings.Join(f[1:], " ")
		}
	case "tag":
		if len(f) > 1 {
			article.Tags = f[1:]
		}
	case "img":
		b := strings.Join(f[1:], " ")
		//      fmt.Println("\n------------\n", article.SrcDirectoryName)
		//      fmt.Println(f[1])

		o := `<a href="` + f[1] + `"><img src=` + b + `></a>`
		output = []byte(o)

	case "summary":
		article.Summary = strings.Join(f[1:], " ")
	case "title":
		article.Title = strings.Join(f[1:], " ")
	default:
		fmt.Println(article.SrcFileName + ": plugin '" + name + "'" + color.RedString(" NOT supported"))
	}
	return output, name
}
