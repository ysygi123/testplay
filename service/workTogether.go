package service

import (
	"strconv"
	"testplay/utils"
)

type WT struct {
	MaxUserID  int
	WorkerList []*Worker
}

type MissionType struct {
	MID    int
	MS     string
	ToUser int
}

type Worker struct {
	UserID         int
	NextWorkerTime int64
}

func (w *WT) generateData() (list []*MissionType) {
	list = make([]*MissionType, 10000)
	for i := 0; i < 10000; i++ {
		list[i] = &MissionType{
			MID:    i,
			MS:     "mission : " + strconv.Itoa(i),
			ToUser: utils.RangeRand(1, w.MaxUserID),
		}
	}
	return
}

func (w *WT) getWorkers() {
	w.WorkerList = make([]*Worker, w.MaxUserID)
	for i := 0; i < w.MaxUserID; i++ {
		w.WorkerList[i] = &Worker{
			UserID: i,
		}
	}
	return
}

func (w *WT) Main() {
	w.MaxUserID = 10
	//missionList := w.generateData()
	w.getWorkers()
	//for _, mission := range missionList {
	//
	//}

}
