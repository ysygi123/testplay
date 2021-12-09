package utils

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"strconv"
)

var E *elastic.Client

type ES struct {
	EsClient *elastic.Client
	Index    string
}

type QueryColumn struct {
	key string
	val interface{}
}

func NewClient() *ES {
	es := new(ES)
	if E != nil {
		es.EsClient = E
	} else {
		client, err := elastic.NewClient(
			elastic.SetSniff(false),
			elastic.SetURL("http://127.0.0.1:9200"),
			elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		)
		if err != nil {
			log.Fatalln("failedto create -------" + err.Error())
		}
		es.EsClient = client
	}
	return es
}

func (e *ES) SetIndex(index string) *ES {
	e.Index = index
	return e
}

func (e *ES) Upsert(id string, data interface{}) (b *elastic.BulkResponse, err error) {
	if e.Index == "" {
		err = fmt.Errorf("没设置index")
		return
	}
	req := elastic.NewBulkUpdateRequest().
		Index(e.Index).
		Id(id).
		Doc(data).
		DocAsUpsert(true)

	b, err = e.EsClient.Bulk().Add(req).
		Do(context.Background())

	if err != nil {
		return
	}

	return
}

// Del 单个删除
func (e *ES) Del(id int) (esc *elastic.DeleteResponse, err error) {
	if e.Index == "" {
		err = fmt.Errorf("没设置index")
		return
	}
	return e.EsClient.Delete().Index(e.Index).Id(strconv.Itoa(id)).Do(context.Background())
}

