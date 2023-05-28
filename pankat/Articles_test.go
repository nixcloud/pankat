package pankat

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"
	"time"
)

type watcherMock struct {
	mock.Mock
}

func (o *watcherMock) fileWrite(filename string) {
	fmt.Println("fileWrite: ", filename)
	o.Called(filename)
}

func (o *watcherMock) fileRemove(filename string) {
	fmt.Println("fileRemove: ", filename)
	o.Called(filename)
}

func (o *watcherMock) fileCreate(filename string) {
	fmt.Println("fileCreate: ", filename)
	o.Called(filename)
}

func (o *watcherMock) dirCreate(filename string) {
	fmt.Println("dirCreate: ", filename)
	o.Called(filename)
}

func (o *watcherMock) dirRemove(filename string) {
	fmt.Println("dirRemove: ", filename)
	o.Called(filename)
}

const longForm = "2006-01-02 15:04"

func TestArticlesSortingOrder(t *testing.T) {
	var articles Articles
	// assign a value to a variable
	time1, _ := time.Parse(longForm, "2019-01-01 00:00")
	article1 := Article{Title: "foo", ModificationDate: time1, Summary: "foo summary", Tags: []string{"Linux", "Go"},
		SrcFileName: "/home/user/documents/foo.mdwn", DstFileName: "/home/user/documents/foo.html"}
	time2, _ := time.Parse(longForm, "2022-01-01 00:00")
	article2 := Article{Title: "bar", ModificationDate: time2, Summary: "bar summary", Tags: []string{"Linux", "Go"},
		SrcFileName: "/home/user/documents/bar.mdwn", DstFileName: "/home/user/documents/bar.html"}
	time3, _ := time.Parse(longForm, "2010-01-01 00:00")
	article3 := Article{Title: "batz", ModificationDate: time3, Summary: "batz summary", Tags: []string{"Linux", "Go"},
		SrcFileName: "/home/user/documents/batz.mdwn", DstFileName: "/home/user/documents/batz.html"}
	time4, _ := time.Parse(longForm, "2024-01-01 00:00")
	article4 := Article{Draft: true, Title: "draft", ModificationDate: time4, Summary: "draft summary", Tags: []string{"Linux", "Go"},
		SrcFileName: "/home/user/documents/draft.mdwn", DstFileName: "/home/user/documents/draft.html"}
	articles = append(articles, &article1)
	articles = append(articles, &article2)
	articles = append(articles, &article3)
	articles = append(articles, &article4)
	assert.Equal(t, len(articles), 4)

	filteredArticles := articles.FilterOutDrafts()
	sort.Sort(filteredArticles)
	assert.Equal(t, len(filteredArticles), 3)
	assert.Equal(t, filteredArticles[0].Title, "bar")
	assert.Equal(t, filteredArticles[1].Title, "foo")
	assert.Equal(t, filteredArticles[2].Title, "batz")

	assert.True(t, filteredArticles.NextArticle(filteredArticles[0]) == nil)
	assert.True(t, filteredArticles.PrevArticle(filteredArticles[0]) == filteredArticles[1])

	assert.Equal(t, filteredArticles.NextArticle(filteredArticles[1]), filteredArticles[0])
	assert.True(t, filteredArticles.PrevArticle(filteredArticles[1]) == filteredArticles[2])

	assert.True(t, filteredArticles.NextArticle(filteredArticles[2]) == filteredArticles[1])
	assert.True(t, filteredArticles.PrevArticle(filteredArticles[2]) == nil)
}

func fsNotifyWatchDocumentsDirectory(directory string, mockedWatcher *watcherMock) {
	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:
				//fmt.Println("xxxxxxxxxxxxxxxxx " + event.Name() + " xxxxxxxxxxxxxxxxx")
				if event.FileInfo.IsDir() == false {
					documentsPath, err := os.Getwd()
					if err != nil {
						log.Println(err)
					}
					eventRelFileName, _ := filepath.Rel(documentsPath, event.Path)
					// random int
					//a := rand.Intn(1010)
					//time.Sleep(time.Millisecond * time.Duration(a))
					if event.Op == watcher.Remove {
						fmt.Println("file removed:", event.Name())
						mockedWatcher.fileRemove(eventRelFileName)
					}
					if event.Op == watcher.Write {
						fmt.Println("File write in: ", eventRelFileName)
						mockedWatcher.fileWrite(eventRelFileName)
					}
					if event.Op == watcher.Create {
						fmt.Println("create detected in: ", eventRelFileName)
						mockedWatcher.fileCreate(eventRelFileName)
					}
					//for v := 0; v < 10000; v++ {
					//	os.Create("./test_data/foo" + strconv.Itoa(v) + ".mdwn")
					//}
				} else {
					if event.Op == watcher.Remove {
						fmt.Println("dir removed:", event.Name())
						w.Remove(event.Path)
						mockedWatcher.dirRemove(event.Name())
					}
					if event.Op == watcher.Create {
						fmt.Println("dir created:", event.Name())
						w.Add(event.Path)
						mockedWatcher.dirCreate(event.Name())
					}
				}
			case err := <-w.Error:
				fmt.Println("ERROR in fswatcher", err)
				if err == watcher.ErrWatchedFileDeleted {
					continue
				}
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	walkFunc := watchDir(w)
	fmt.Println("watching dir: ", directory)
	if err := filepath.Walk(directory, walkFunc); err != nil {
		fmt.Println("ERROR", err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func watchDir(w *watcher.Watcher) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			fmt.Println("watching dir: ", path)
			return w.Add(path)
		}
		return nil
	}
}

func createDeleteFileTest(basedir string) {
	dir := filepath.Join(basedir, "createDeleteFileTest")
	os.Mkdir(dir, 0755)
	time.Sleep(1 * time.Second)
	err := os.Remove(dir)
	if err != nil {
		fmt.Println(err)
	}
}

func createWriteDeleteFile(basedir string, filename string) {
	f := filepath.Join(basedir, filename)
	fmt.Println("xxxxx: " + f)
	// create directory
	os.Mkdir(basedir, 0755)
	// create file
	os.Create(f)
	// write string into file
	file, _ := os.OpenFile(f, os.O_RDWR, 0644)
	file.WriteString("asdf")
	file.Sync()
	file.Close()
	time.Sleep(1 * time.Second)
	err := os.Remove(f)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)
}

func createFileAndDir(basedir string, filename string) {
	for i := 0; i < 10; i++ {
		dir := filepath.Join(basedir, filename+strconv.Itoa(i))
		f := filepath.Join(dir, filename+".mdwn")
		// create directory
		os.Mkdir(dir, 0755)
		os.Mkdir(dir+"/asdf"+strconv.Itoa(i), 0755)
		// create file
		os.Create(f)
		// write string into file
		file, _ := os.OpenFile(f, os.O_RDWR, 0644)
		defer file.Close()
		file.WriteString(filename)
		time.Sleep(1 * time.Second)
		//for v := 0; v < 10000; v++ {
		//	os.Create("./test_data/foo" + strconv.Itoa(v) + ".mdwn")
		//}
	}
}

func TestMockWatcherActivity(t *testing.T) {
	mockWatcherActivity := new(watcherMock)
	mockWatcherActivity.On("dirCreate", "test_data").Once()
	mockWatcherActivity.On("dirRemove", "test_data").Twice()

	mockWatcherActivity.dirCreate("test_data")
	mockWatcherActivity.dirRemove("test_data")
	mockWatcherActivity.dirRemove("test_data")

	mockWatcherActivity.On("fileCreate", "test_data.mdwn").Once()
	mockWatcherActivity.On("fileWrite", "test_data.mdwn").Once()
	mockWatcherActivity.On("fileRemove", "test_data.mdwn").Once()

	mockWatcherActivity.fileCreate("test_data.mdwn")
	mockWatcherActivity.fileWrite("test_data.mdwn")
	mockWatcherActivity.fileRemove("test_data.mdwn")
}

func TestCreateAndDeleteDir(t *testing.T) {
	dir := "./test_data"
	defer os.Remove(dir)
	os.Remove(dir)
	os.Mkdir(dir, 0755)
	mockWatcherActivity := new(watcherMock)

	mockWatcherActivity.On("dirCreate", "createDeleteFileTest").Once()
	//mockWatcherActivity.On("dirRemove1", "createDeleteFileTest1").Once()

	go fsNotifyWatchDocumentsDirectory(dir, mockWatcherActivity)
	time.Sleep(1 * time.Second)
	createDeleteFileTest(dir)
	mockWatcherActivity.AssertExpectations(t)

}

func TestCreateWriteDeleteFile(t *testing.T) {
	dir := "./test_data"
	defer os.Remove(dir)
	os.RemoveAll(dir)
	time.Sleep(1 * time.Second)

	os.Mkdir(dir, 0755)
	mockWatcherActivity := new(watcherMock)
	go fsNotifyWatchDocumentsDirectory(dir, mockWatcherActivity)
	mockWatcherActivity.On("fileCreate", "test_data\\testfileCreateWriteDelete.mdwn").Once()
	mockWatcherActivity.On("fileWrite", "test_data\\testfileCreateWriteDelete.mdwn").Once()
	mockWatcherActivity.On("fileRemove", "test_data\\testfileCreateWriteDelete.mdwn").Once()
	time.Sleep(1 * time.Second)

	go createWriteDeleteFile(dir, "testfileCreateWriteDelete.mdwn")

	time.Sleep(2 * time.Second)
	mockWatcherActivity.AssertExpectations(t)
}

func TestCrash(t *testing.T) {
	dir := "./test_data"
	defer os.Remove(dir)
	os.Remove(dir)
	os.Mkdir(dir, 0755)
	mockWatcherActivity := new(watcherMock)
	go fsNotifyWatchDocumentsDirectory(dir, mockWatcherActivity)
	go createFileAndDir(dir, "foo")
	// make sure it is called at least once or several times
	//mockWatcherActivity.On("fileWrite", "test_data\\foo1\\foo.mdwn").Once()
	//mockWatcherActivity.On("fileWrite", "test_data\\foo2\\foo.mdwn").Maybe()
	//mockWatcherActivity.On("fileWrite", "test_data\\foo3\\foo.mdwn").Maybe()
	//mockWatcherActivity.On("fileWrite", "test_data\\foo4\\foo.mdwn").Maybe()
	//mockWatcherActivity.On("fileWrite", "test_data\\foo5\\foo.mdwn").Maybe()
	//mockWatcherActivity.On("fileWrite", "test_data\\foo6\\foo.mdwn").Maybe()
	//mockWatcherActivity.On("fileWrite", "test_data\\foo7\\foo.mdwn").Maybe()
	//mockWatcherActivity.On("fileWrite", "test_data\\foo8\\foo.mdwn").Maybe()
	//mockWatcherActivity.On("fileWrite", "test_data\\foo9\\foo.mdwn").Maybe()

	go createFileAndDir(dir, "bar")
	//go createFileAndDir(dir, "batz")
	//go createFileAndDir(dir, "batz11")
	//go createFileAndDir(dir, "batz22")
	//go createFileAndDir(dir, "batz33")
	//go createFileAndDir(dir, "batz44")
	time.Sleep(10 * time.Second)
}
