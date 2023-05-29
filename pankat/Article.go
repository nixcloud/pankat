package pankat

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
)

type Article struct {
	Title             string
	ArticleMDWNSource []byte
	ModificationDate  time.Time
	Summary           string
	Tags              []string
	Series            string
	SrcFileName       string // /home/user/documents/foo.mdwn
	DstFileName       string // /home/user/documents/foo.html
	SpecialPage       bool   // used for timeline.html, about.html (not added to timeline if true, not added in list of articles)
	Draft             bool
	Anchorjs          bool
	Tocify            bool
	Timeline          bool // generating timeline.html uses this flag in RenderTimeline(..)
	SourceReference   bool // switch for showing the document source mdwn at bottom
	WebsocketSupport  bool // live update support via WS on/off
}

func (a Article) Hash() md5hash {
	bytes, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
	}
	return md5.Sum(bytes)
}
