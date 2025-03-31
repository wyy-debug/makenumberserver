package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"henna-queue/internal/model"
	"henna-queue/internal/repository/mysql"
	"henna-queue/internal/repository/redis"
	"henna-queue/internal/util/wechat"
)

type QueueService struct {
	queueRepo       *mysql.QueueRepository
	userRepo        *mysql.UserRepository
	shopRepo        *mysql.ShopRepository
	serviceRepo     *mysql.ServiceRepository
	statisticsRepo  *mysql.StatisticsRepository
	redisCache      *redis.CacheRepository
}

func NewQueueService() *QueueService {
	return &QueueService{
		queueRepo:       mysql.NewQueueRepository(),
		userRepo:        mysql.NewUserRepository(),
		shopRepo:        mysql.NewShopRepository(),
		serviceRepo:     mysql.NewServiceRepository(),
		statisticsRepo:  mysql.NewStatisticsRepository(),
		redisCache:      redis.NewCacheRepository(),
	}
}

// GetQueueStatus 获取排队状态
func (s *QueueService) GetQueueStatus(userID string, shopID uint) (*model.QueueStatusResponse, error) {
	// 获取店铺信息
	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return nil, errors.New("店铺不存在")
	}
	
	// 检查店铺是否营业
	if shop.Status == 0 {
		return nil, errors.New("店铺暂停营业")
	}
	
	// 获取用户当前排队状态
	queue, err := s.queueRepo.GetActiveByUserID(userID, shopID)
	hasNumber := err == nil && queue != nil
	
	// 获取当前叫号情况
	currentServing, currentWaiting, err := s.redisCache.GetQueueStatus(shopID)
	if err != nil {
		// 如果缓存不存在，从数据库获取
		var servingQueue, waitingQueue *model.Queue
		servingQueue, _ = s.queueRepo.GetByStatus(shopID, 2) // 服务中
		waitingQueue, _ = s.queueRepo.GetByStatus(shopID, 1) // 就位中
		
		if servingQueue != nil {
			currentServing = servingQueue.QueueNumber
		}
		
		if waitingQueue != nil {
			currentWaiting = waitingQueue.QueueNumber
		}
		
		// 更新缓存
		s.redisCache.UpdateQueueStatus(shopID)
	}
	
	// 获取已服务人数
	servedCount, _ := s.queueRepo.GetTodayServedCount(shopID)
	
	// 构造响应
	status := &model.QueueStatusResponse{
		HasNumber:      hasNumber,
		CurrentServing: currentServing,
		CurrentWaiting: currentWaiting,
		TotalServed:    servedCount,
	}
	
	if hasNumber {
		status.QueueNumber = queue.QueueNumber
		
		// 获取前方等待人数
		peopleAhead, _ := s.queueRepo.GetPeopleAhead(shopID, queue.QueueNumber)
		status.PeopleAhead = peopleAhead
		
		// 计算预计等待时间
		if queue.Status == 0 { // 等待中
			waitTime := peopleAhead * queue.Service.Duration
			status.WaitTime = waitTime
		}
	}
	
	return status, nil
}

// GetQueueNumber 取号排队
func (s *QueueService) GetQueueNumber(userID string, shopID, serviceID uint) (*model.Queue, error) {
	// 检查店铺是否营业
	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return nil, errors.New("店铺不存在")
	}
	
	if shop.Status == 0 {
		return nil, errors.New("店铺暂停营业")
	}
	
	// 检查用户是否已在排队
	activeQueue, err := s.queueRepo.GetActiveByUserID(userID, shopID)
	if err == nil && activeQueue != nil {
		return nil, errors.New("您已在排队中")
	}
	
	// 检查服务是否存在
	service, err := s.serviceRepo.GetByID(serviceID)
	if err != nil {
		return nil, errors.New("服务不存在")
	}
	
	if service.ShopID != shopID {
		return nil, errors.New("服务不属于该店铺")
	}
	
	// 生成排队号码
	now := time.Now()
	datePrefix := now.Format("0102") // 月日
	
	// 获取今日队列数量
	count, _ := s.queueRepo.GetDailyCount(shopID, now.Format("20060102"))
	queueNumber := fmt.Sprintf("%s%03d", datePrefix, count+1)
	
	// 创建排队记录
	queue := &model.Queue{
		ShopID:      shopID,
		QueueNumber: queueNumber,
		UserID:      userID,
		ServiceID:   serviceID,
		Status:      0, // 等待中
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	err = s.queueRepo.Create(queue)
	if err != nil {
		return nil, err
	}
	
	// 重新加载，获取关联数据
	queue, _ = s.queueRepo.GetByID(queue.ID)
	
	// 异步发送消息通知
	go func() {
		user, _ := s.userRepo.GetByID(userID)
		if user != nil && user.UnionID != "" {
			wechat.SendSubscribeMessage(user.UnionID, 
				"排队成功通知", 
				fmt.Sprintf("您的排队号码为:%s，当前等待人数:%d", queueNumber, count))
		}
	}()
	
	return queue, nil
}

// CancelQueue 取消排队
func (s *QueueService) CancelQueue(userID string, shopID uint) error {
	// 获取用户当前队列
	queue, err := s.queueRepo.GetActiveByUserID(userID, shopID)
	if err != nil {
		return err
	}
	
	if queue == nil {
		return errors.New("您没有正在进行的排队")
	}
	
	// 更新队列状态为已取消
	queue.Status = 4 // 已取消
	err = s.queueRepo.Update(queue)
	if err != nil {
		return err
	}
	
	// 更新Redis缓存
	s.redisCache.UpdateQueueStatus(shopID)
	
	// 更新统计数据
	s.statisticsRepo.IncrementCancelCount(shopID)
	
	return nil
}

// GetCurrentQueue 获取当前叫号情况
func (s *QueueService) GetCurrentQueue(shopID uint) (*model.CurrentQueueResponse, error) {
	// 获取当前叫号情况
	currentServing, _ := s.queueRepo.GetByStatus(shopID, 2) // 服务中
	currentWaiting, _ := s.queueRepo.GetByStatus(shopID, 1) // 就位中
	
	// 获取今日已服务人数
	todayServed, _ := s.queueRepo.GetTodayServedCount(shopID)
	
	response := &model.CurrentQueueResponse{
		CurrentServing: "",
		CurrentWaiting: "",
		TotalServed:    todayServed,
	}
	
	if currentServing != nil {
		response.CurrentServing = currentServing.QueueNumber
	}
	
	if currentWaiting != nil {
		response.CurrentWaiting = currentWaiting.QueueNumber
	}
	
	return response, nil
}

// CallNextNumber 商户叫号(下一位)
func (s *QueueService) CallNextNumber(shopID uint) (*model.Queue, error) {
	// 查找当前正在服务的队列
	serving, _ := s.queueRepo.GetByStatus(shopID, 2) // 服务中
	if serving != nil {
		// 将当前服务设为已完成
		serving.Status = 3 // 已完成
		err := s.queueRepo.Update(serving)
		if err != nil {
			return nil, err
		}
		
		// 更新统计数据
		s.statisticsRepo.IncrementServedCount(shopID)
	}
	
	// 查找当前等待就位的队列
	waiting, _ := s.queueRepo.GetByStatus(shopID, 1) // 就位中
	if waiting != nil {
		// 将等待就位设为服务中
		waiting.Status = 2 // 服务中
		err := s.queueRepo.Update(waiting)
		if err != nil {
			return nil, err
		}
	}
	
	// 查找下一个等待中的队列
	next, err := s.queueRepo.GetNextWaiting(shopID)
	if err != nil || next == nil {
		return nil, errors.New("没有下一位等待的客人")
	}
	
	// 将下一个等待设为就位中
	next.Status = 1 // 就位中
	err = s.queueRepo.Update(next)
	if err != nil {
		return nil, err
	}
	
	// 更新Redis缓存
	s.redisCache.UpdateQueueStatus(shopID)
	
	// 发送订阅消息通知用户
	go func() {
		err = wechat.SendSubscribeMessage(next.UserID, "就位通知", fmt.Sprintf("您的号码 %s 请就位", next.QueueNumber))
		if err != nil {
			log.Printf("发送订阅消息失败: %v", err)
		}
	}()
	
	return next, nil
}

// UpdateQueueStatus 更新队列状态
func (s *QueueService) UpdateQueueStatus(queueID uint, status int8) error {
	queue, err := s.queueRepo.GetByID(queueID)
	if err != nil {
		return err
	}
	
	if queue == nil {
		return errors.New("队列不存在")
	}
	
	// 更新状态
	queue.Status = status
	err = s.queueRepo.Update(queue)
	if err != nil {
		return err
	}
	
	// 更新Redis缓存
	s.redisCache.UpdateQueueStatus(queue.ShopID)
	
	// 更新统计数据
	if status == 3 { // 如果状态是已完成
		s.statisticsRepo.IncrementServedCount(queue.ShopID)
	} else if status == 4 { // 如果状态是已取消
		s.statisticsRepo.IncrementCancelCount(queue.ShopID)
	}
	
	return nil
}

// ToggleQueuePause 暂停/恢复排队
func (s *QueueService) ToggleQueuePause(shopID uint) (*model.Shop, error) {
	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return nil, err
	}
	
	// 切换状态
	if shop.Status == 1 {
		shop.Status = 0 // 暂停
	} else {
		shop.Status = 1 // 恢复
	}
	
	err = s.shopRepo.Update(shop)
	if err != nil {
		return nil, err
	}
	
	return shop, nil
}

// GetAdminQueue 获取管理员排队列表
func (s *QueueService) GetAdminQueue(shopID uint) ([]*model.Queue, error) {
	// 获取正在排队的列表
	queues, err := s.queueRepo.GetActiveQueues(shopID)
	if err != nil {
		return nil, err
	}
	
	return queues, nil
} 