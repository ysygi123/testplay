package main

/*
#cgo LDFLAGS: -L./C_sdk -lLibWeWorkFinanceSdk_C
#include "./C_sdk/WeWorkFinanceSdk_C.h"
*/

import (
	"context"
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/nsqio/go-nsq"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"sync"
	"testplay/model/es"
	"testplay/model/mysql"
	"testplay/service"
	"testplay/utils"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)
	ps := sync.Pool{
		New: func() interface{} {
			return 0
		},
	}

	go func() {
		for i := 1; i < 100; i++ {
			ps.Put(i)
		}
	}()
	time.Sleep(1 * time.Second)

	ps.Put(100)
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Println("--", ps.Get())
		}()
	}
	time.Sleep(1 * time.Second)
}

func TestArrayColumn() {
	var err error
	type c struct {
		UID  int
		Name string
	}
	//var k []string = make([]string, 0)
	//
	//a := []*c{
	//	&c{
	//		UID: 111,
	//		Name: "1111111",
	//	},
	//	&c{
	//		UID: 222,
	//		Name: "222222",
	//	},
	//}
	//err = common.ArrayColumn(&k, a, "Name", "")
	//fmt.Println(k, err)
	//
	//var csks map[string]string = make(map[string]string)
	//err = common.ArrayColumn(&csks, a, "Name", "Namesadf")
	//fmt.Println(csks, err)

	//var ciki map[int]int = make(map[int]int)
	//err = common.ArrayColumn(&ciki, a, "UID", "UID")
	//fmt.Println(ciki, err)
	//
	//var cikaptr map[int]*c = make(map[int]*c)
	//err = common.ArrayColumn(&cikaptr, a, "", "UID")
	//fmt.Println(cikaptr, err, cikaptr[111], cikaptr[222])
	//
	//var cika map[int]c = make(map[int]c)
	//err = common.ArrayColumn(&cika, a, "", "UID")
	//fmt.Println(cika, err)
	//
	//var ss []string = make([]string, 5)
	//err = common.ArrayColumn(&ss, a, "Name", "")
	//fmt.Println(ss, err)
	//
	x := []map[string]string{
		{"user": "a", "name": "b"}, {"user": "c", "name": "d"},
	}
	var cikaptrmp1 map[string]string = make(map[string]string)
	err = utils.ArrayColumn(&cikaptrmp1, x, "name", "name")
	fmt.Println(cikaptrmp1, err)
	//
	x2 := []*map[string]string{
		{"user": "a", "name": "b"}, {"user": "c", "name": "d"},
	}
	var cikaptrmp2 map[string]string = make(map[string]string)
	err = utils.ArrayColumn(&cikaptrmp2, x2, "user", "name")
	fmt.Println(cikaptrmp2, err)
}

func TestMyChannel() {
	zimu := [26]string{}
	for i := 97; i <= 122; i++ {
		zimu[i-97] = string(rune(i))
	}

	keysize := 1000000

	slicekey := make([]string, keysize)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < keysize; i++ {
		x := ""
		for j := 0; j < 5; j++ {
			x += zimu[25-rand.Intn(25)]
		}
		slicekey[i] = x
	}

	for uid, userID := range slicekey {
		service.GlobalMyMap.Set(&service.UserInfo{
			Uid:    uid,
			UserID: userID,
		})
		service.GlobalNormalMap.Set(&service.UserInfo{
			Uid:    uid,
			UserID: userID,
		})
	}

	randIndexNum := 10000000
	randUserIDs := make([]string, randIndexNum)
	for i := 0; i < randIndexNum; i++ {
		randUserIDs[i] = slicekey[rand.Intn(keysize-1)]
	}

	w := new(sync.WaitGroup)

	w.Add(4)

	go func() {
		for i := 0; i < 100; i++ {
			go func() {
				service.GlobalNormalMap.Set(&service.UserInfo{
					Uid:    100000,
					UserID: "0123123123123",
				})
			}()
		}
		w.Done()
	}()

	go func() {
		for i := 0; i < 100; i++ {
			go func() {
				service.GlobalMyMap.Set(&service.UserInfo{
					Uid:    100000,
					UserID: "0123123123123",
				})
			}()
		}
		w.Done()
	}()

	goNum := int32(50)

	go func() {
		goChan := make(chan struct{}, goNum)
		exitChan := make(chan struct{}, 1)
		t := time.Now()
		for index, key := range randUserIDs {
			goChan <- struct{}{}
			go func(i int, k string) {
				service.GlobalNormalMap.Get(k)
				<-goChan
				//fmt.Println(i, k)
				if i == randIndexNum-1 {
					fmt.Println("全部结束时间 普通map : ", time.Since(t))
					close(goChan)
					exitChan <- struct{}{}
				}
			}(index, key)
		}
		<-exitChan
		w.Done()
	}()

	go func() {
		goChan := make(chan struct{}, goNum)
		exitChan := make(chan struct{}, 1)
		t := time.Now()
		for index, key := range randUserIDs {
			goChan <- struct{}{}
			go func(i int, k string) {
				service.GlobalMyMap.Search(k)
				<-goChan
				//fmt.Println(i, k)
				if i == randIndexNum-1 {
					fmt.Println("全部结束时间 改map : ", time.Since(t))
					close(goChan)
					exitChan <- struct{}{}
				}
			}(index, key)
		}
		<-exitChan
		w.Done()
	}()
	w.Wait()

}

func loopchannel() {
	c := make(chan int, 5)
	i := 0
	for {
		select {
		case c <- 1:
		default:
			fmt.Println("can't write")
			i++
			if i == 5 {
				goto LoopExit
			}
		}
	}
LoopExit:
	fmt.Println(111, c, len(c))
	close(c)

}

func producer1() {
	addr := "192.168.11.98:4150"
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(addr, cfg)
	fmt.Println(err)
	err = producer.Publish("test1", []byte("hello1"))
	fmt.Println(err)
}

func producer2() {
	addr := "192.168.11.98:4152"
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(addr, cfg)
	fmt.Println(err)
	err = producer.Publish("test1", []byte("hello333333"))
	fmt.Println(err)
}

func comsume1() {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 10
	c, err := nsq.NewConsumer("test1", "test1", cfg)
	if err != nil {
		panic(err)
	}
	hand := func(msg *nsq.Message) error {
		fmt.Println(string(msg.Body))
		return nil
	}

	c.AddHandler(nsq.HandlerFunc(hand))

	if err := c.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		fmt.Println(err)
	}

}

func rpcServerFunc() {
	rpcServer := grpc.NewServer()
	pro := new(service.ProdService)
	service.RegisterGetggServer(rpcServer, pro)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("server start !")
	_ = rpcServer.Serve(listener)
}

func arrayUnique() {
	s1 := make([]string, 0)
	s2 := []string{"a", "a", "a"}
	utils.ArrayUnique(&s1, s2)
	fmt.Println(s1)
}

func esAddData() {
	est := utils.NewClient()
	est.SetIndex("user")
	jj := 10000
	s := make([]*es.UserES, jj)
	for i := 0; i < jj; i++ {
		s[i] = es.RandomUserEs(i)
	}
	sss := make([]elastic.BulkableRequest, 0)
	//userInsertMysql := make([]*mysql.User, 0)
	for _, tmpUser := range s {

		pets := make([]*mysql.Pet, 0)
		for _, p := range tmpUser.Pets {
			pt := &mysql.Pet{
				PetName:   p.PetName,
				PetSex:    p.PetSex,
				PetAttack: p.PetAttack,
				PetTag:    p.PetTag,
			}
			pets = append(pets, pt)
		}
		tmpUMysql := &mysql.User{
			//Id:           tmpUser.ID,
			Username:     tmpUser.Username,
			Profession:   tmpUser.Profession,
			Email:        tmpUser.Email,
			TextInfo:     tmpUser.TextInfo,
			RegisterTime: tmpUser.RegisterTime,
			Attribute:    tmpUser.Attribute,
			Pets:         pets,
		}

		_, _ = mysql.TmpAdd(tmpUMysql)

		req := elastic.NewBulkUpdateRequest().
			Index(est.Index).
			Id(strconv.Itoa(tmpUMysql.Id)).
			Doc(tmpUser).
			DocAsUpsert(true)
		sss = append(sss, req)
		//userInsertMysql = append(userInsertMysql, tmpU)
	}
	//fmt.Println(mysql.TmpAddMany(userInsertMysql))
	b, err := est.EsClient.Bulk().Add(sss...).Do(context.Background())
	fmt.Printf("%+v\n%+v\n", b, err)
}

func delEs() {
	est := utils.NewClient()
	est.SetIndex("user")
	fmt.Println(est.Del(0))
	fmt.Println(est.Del(1))
	fmt.Println(est.Del(2))
	fmt.Println(est.Del(3))
	fmt.Println(est.Del(4))
}

func test(kk interface{}) {
	fmt.Println(kk.([]interface{}))
}

// EsOneSearch  query这样都是单个查询
func EsOneSearch() {
	est := utils.NewClient()
	//会被覆盖，只有1个
	tmpQuery := elastic.NewTermsQuery("id", 1)
	tmpQuery2 := elastic.NewRangeQuery("register_time").Gte(1).Lte(100000011001)
	getRes, err := est.EsClient.Search("user").Query(tmpQuery).Query(tmpQuery2).From(0).Size(100).Pretty(true).Do(context.Background())
	if err != nil {
		fmt.Println("出现了很恶心的错误", err)
		return
	}
	fmt.Println(getRes.Hits.Hits)
}

func EsManySearch() {
	q := elastic.NewBoolQuery().Must(elastic.NewTermsQuery("id", 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20))
	q.Must(elastic.NewRangeQuery("register_time").Gte(1).Lte(1589951633))
	q.Must(elastic.NewMatchQuery("pets.pet_name", "裂空座"))
	query := elastic.NewBoolQuery().Filter(q)
	ess, _ := utils.NewClient().EsClient.Search("user").Query(query).Do(context.Background())
	fmt.Println(ess)
}

func RocketMQTest() {
	//s := "当前时间:" + time.Now().String()
	utils.GloableProducerDBManager.SendMessageSync("testTopic", []byte("当前时间:"+time.Now().String()))
	err := utils.GetStartConsumer("testTopic", "notify_consumer_success", "t1", utils.F1)
	if err != nil {
		panic(err)
	}
	err = utils.GetStartConsumer("testTopic", "notify_consumer_success", "t2", utils.F2)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(time.Second * 1) // 运行时长
	var x int
	for x < 10 {
		select {
		case <-ticker.C:
			x++
			utils.GloableProducerDBManager.SendMessageSync("testTopic", []byte("当前时间:"+time.Now().String()))
		}
	}
	ticker.Stop()

	//utils.GloableProducerDBManager.SendMessageSync("testTopic", []byte(s))
	//time.Sleep(1)
}

func TestAddExternalUser() {
	ins := make([]elastic.BulkableRequest, 0)
	est := utils.NewClient()
	est.SetIndex("sc-external-user")
	page := 1
	size := 5000
	for {
		list, err := mysql.GetExternalByLimit(page, size)
		if err != nil {
			panic(err)
		}
		if len(list) == 0 {
			break
		}
		listEla := make([]*es.ScExternalUser, 0, len(list))
		b, _ := jsoniter.Marshal(list)
		_ = jsoniter.Unmarshal(b, &listEla)
		for _, esV := range listEla {
			req := elastic.NewBulkUpdateRequest().Index(est.Index).Id(strconv.Itoa(int(esV.ID))).Doc(esV).DocAsUpsert(true)
			ins = append(ins, req)
		}
		page += 1
		resp, err := est.EsClient.Bulk().Add(ins...).Do(context.Background())
		ins = make([]elastic.BulkableRequest, 0)
		fmt.Printf("%+v\n%+v\n", resp, err)
	}

}

func TestSearchExternal() {
	option := &es.ExternalOption{
		QueryBase: es.QueryBase{
			Page: 0,
			Size: 0,
			Sort: nil,
		},
		Name:     "杨",
		CorpName: "蝉妈妈数据",
	}

	list, _ := es.GetExternalWithLimit(option, context.Background())
	fmt.Println(utils.Data2json(list))
}

func TestPPROF() {
	counter := func() {
		slice := make([]int, 0)
		c := 1
		for i := 0; i < 100000; i++ {
			c = i + 1 + 2 + 3 + 4 + 5
			slice = append(slice, c)
		}
	}
	workOnce := func(wg *sync.WaitGroup) {
		counter()
		wg.Done()
	}

	var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
	var memProfile = flag.String("memprofile", "", "write mem profile to file")

	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go workOnce(&wg)
	}

	wg.Wait()
}

func TestContext() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func(cancelCtx context.Context) {
		for {
			select {
			case <-cancelCtx.Done():
				fmt.Println("1 g end")
				goto END
			default:
				fmt.Println("1")
				time.Sleep(1 * time.Second)
			}
		}
	END:
	}(ctx)

	fmt.Println("赋予新的值")
	ctx, cancel2 := context.WithCancel(ctx)

	go func(cancelCtx context.Context) {
		for {
			select {
			case <-cancelCtx.Done():
				fmt.Println("2 g end")
				goto END
			default:
				fmt.Println("2")
				time.Sleep(1 * time.Second)
			}
		}
	END:
	}(ctx)
	//time.Sleep(5 * time.Second)
	fmt.Println("2 cancel")
	cancel2()
	//time.Sleep(5 * time.Second)
	fmt.Println("1 cancel")
	cancel()
}

func testchannel() {
	var cccccccccc chan int
	close(cccccccccc)
}