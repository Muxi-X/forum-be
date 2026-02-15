package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"sync"
)

var categories = []string{
	constvar.DailyLife, constvar.Study, constvar.Project, constvar.Emotion, constvar.Campus, constvar.Trade,
}

func (s *PostService) GetUnReadPostNum(_ context.Context, req *pb.Request, resp *pb.UnReadPostNumResponse) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(categories))

	for _, c := range categories {
		go func(category string) {
			defer wg.Done()

			clickModel := &dao.LastClickModel{
				UserId:      req.UserId,
				Category:    category,
				LastClickAt: "2022-10-20T13:34:58+08:00",
			}
			err := clickModel.GetOrCreateLastRead()
			if err != nil {
				logger.Error("GetUnReadPostNum GetOrCreateLastRead", logger.String(err.Error()))
			}
			count, err := s.Dao.CountPostByTime(clickModel.LastClickAt, clickModel.Category)
			if err != nil {
				logger.Error("GetUnReadPostNum CountPostByTime", logger.String(err.Error()))
			}
			unReadPostNum := &pb.UnReadPostNum{
				Category: category,
				Num:      uint32(count),
			}
			mu.Lock()
			resp.UnReadNum = append(resp.UnReadNum, unReadPostNum)
			mu.Unlock()
		}(c)
	}

	wg.Wait()
	return nil
}
