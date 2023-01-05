package pankat

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"os"
)

type md5hash [md5.Size]byte

type ArticlesCacheList struct {
	Hash    md5hash
	Article string
}

type ArticlesCache struct {
	Store map[md5hash]string
}

func (s ArticlesCache) computeHash(a Article) md5hash {
	bytes, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
	}
	return md5.Sum(bytes)
}

// load hashes and articles via json from disk
func (s ArticlesCache) load() {
	var v = []ArticlesCacheList{}
	if GetConfig().Force == 1 {
		fmt.Println(color.MagentaString("Forcing reevaluation, ignoring ArticlesCache"))
	} else {
		b, errReadFile := os.ReadFile(GetConfig().MyMd5HashMapJson)
		if errReadFile != nil {
			fmt.Println(errReadFile)
		} else {
			jBuff := bytes.NewBufferString(string(b))
			dec := json.NewDecoder(jBuff)
			if err := dec.Decode(&v); err != nil {
				fmt.Println(err)
			}
		}
	}
	for i := range v {
		//fmt.Println(v[i].Hash)
		s.Store[v[i].Hash] = v[i].Article
	}
}

// store hashes and articles (hash set) as list via json to disk
func (s ArticlesCache) save() {
	var v = []ArticlesCacheList{}
	for key, value := range s.Store {
		//fmt.Println("Key:", key, "Value:", value)
		var q = ArticlesCacheList{
			key,
			value,
		}
		v = append(v, q)
	}
	jsonBuff := bytes.NewBufferString("")
	enc := json.NewEncoder(jsonBuff)
	if errEnc := enc.Encode(&v); errEnc != nil {
		fmt.Println(errEnc)
	}
	//fmt.Println(string(jsonBuff.Bytes()))
	errn := os.WriteFile(GetConfig().MyMd5HashMapJson, jsonBuff.Bytes(), 0644)
	if errn != nil {
		panic(errn)
	}
}

// query the local cache for the article
func (s ArticlesCache) Get(a Article) string {
	// FIXME add error handling!
	hash := s.computeHash(a)
	return s.Store[hash]
}

// update the local cache for a given article
func (s ArticlesCache) Set(a Article, text string) {

	hash := s.computeHash(a)
	s.Store[hash] = text
	s.save()
}
