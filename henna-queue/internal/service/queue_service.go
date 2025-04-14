package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/internal/repository"
	"example.com/henna-queue/internal/repository/factory"
	"example.com/henna-queue/internal/repository/redis"
	"example.com/henna-queue/internal/util/wechat"
)

type QueueService struct {
	queueRepo     repository.QueueRepository
	userRepo      repository.UserRepository
	shopRepo      repository.ShopRepository
	serviceRepo   repository.ServiceRepository
	statisticRepo repository.StatisticRepository
	redisCache    *redis.CacheRepository
}

func NewQueueService() *QueueService {
	return &QueueService{
		queueRepo:     factory.NewQueueRepository(),
		userRepo:      factory.NewUserRepository(),
		shopRepo:      factory.NewShopRepository(),
		serviceRepo:   factory.NewServiceRepository(),
		statisticRepo: factory.NewStatisticRepository(),
		redisCache:    redis.NewCacheRepository(),
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
	s.statisticRepo.IncrementCancelCount(shopID)

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
		s.statisticRepo.IncrementServedCount(shopID)
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
func (s *QueueService) UpdateQueueStatus(queueID uint, shopID uint, status int8) (*model.Queue, error) {
	// 查找队列
	queue, err := s.queueRepo.GetByID(queueID)
	if err != nil {
		return nil, err
	}

	// 检查队列是否属于该店铺
	if queue.ShopID != shopID {
		return nil, errors.New("无权操作该队列")
	}

	// 更新状态
	queue.Status = status
	queue.UpdatedAt = time.Now()

	err = s.queueRepo.Update(queue)
	if err != nil {
		return nil, err
	}

	// 更新Redis缓存
	s.redisCache.UpdateQueueStatus(shopID)

	// 更新统计数据
	if status == 3 { // 如果状态是已完成
		s.statisticRepo.IncrementServedCount(queue.ShopID)
	} else if status == 4 { // 如果状态是已取消
		s.statisticRepo.IncrementCancelCount(queue.ShopID)
	}

	return queue, nil
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

// GetQueues 获取队列列表
func (s *QueueService) GetQueues(shopID uint, status string, serviceID uint, date string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	// 添加日志以便于调试
	log.Printf("获取队列列表，店铺ID: %d, 状态: %s, 服务ID: %d, 日期: %s, 页码: %d, 每页数量: %d", 
		shopID, status, serviceID, date, page, pageSize)
	
	if s.queueRepo == nil {
		return nil, 0, errors.New("队列仓库未初始化")
	}
	
	// 尝试获取活跃队列
	queues, err := s.queueRepo.GetActiveQueues(shopID)
	if err != nil {
		log.Printf("获取活跃队列错误: %v", err)
		return nil, 0, err
	}
	
	if queues == nil {
		// 如果返回 nil 且没有错误，返回空结果
		return []map[string]interface{}{}, 0, nil
	}

	// 过滤队列
	var filteredQueues []*model.Queue
	for _, queue := range queues {
		// 跳过无效的队列记录
		if queue == nil {
			continue
		}
		
		// 状态过滤
		if status != "" {
			statusInt, err := strconv.ParseInt(status, 10, 8)
			if err == nil && queue.Status != int8(statusInt) {
				continue
			}
		}

		// 服务ID过滤
		if serviceID > 0 && queue.ServiceID != serviceID {
			continue
		}

		// 日期过滤
		if date != "" {
			startTime, err := time.Parse("20060102", date)
			if err == nil {
				endTime := startTime.AddDate(0, 0, 1)
				if queue.CreatedAt.Before(startTime) || queue.CreatedAt.After(endTime) {
					continue
				}
			}
		}

		filteredQueues = append(filteredQueues, queue)
	}

	// 计算总数
	total := int64(len(filteredQueues))

	// 分页处理
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(filteredQueues) {
		return []map[string]interface{}{}, total, nil
	}
	if end > len(filteredQueues) {
		end = len(filteredQueues)
	}
	pagedQueues := filteredQueues[start:end]

	// 构造返回数据
	result := make([]map[string]interface{}, 0, len(pagedQueues))
	for _, queue := range pagedQueues {
		// 确保关联数据存在
		customerName := ""
		phone := ""
		serviceName := ""
		
		// User 是值类型，不能直接和 nil 比较
		// 使用 UserID 字段判断是否有关联数据
		if queue.UserID != "" {
			customerName = queue.User.Nickname
			phone = queue.User.Phone
		}
		
		// Service 是值类型，不能直接和 nil 比较
		// 使用 ServiceID 字段判断是否有关联数据
		if queue.ServiceID > 0 {
			serviceName = queue.Service.Name
		}
		
		queueData := map[string]interface{}{
			"id":            queue.ID,
			"shop_id":       queue.ShopID,
			"customer_name": customerName,
			"phone":         phone,
			"queue_number":  queue.QueueNumber,
			"service_id":    queue.ServiceID,
			"service_name":  serviceName,
			"status":        queue.Status,
			"created_at":    queue.CreatedAt.Format("2006-01-02 15:04:05"),
			"note":          "", // 如果需要备注字段，可以在Queue模型中添加
		}
		
		result = append(result, queueData)
	}

	return result, total, nil
}

// CreateQueueByAdmin 管理员创建队列
func (s *QueueService) CreateQueueByAdmin(shopID uint, customerName, phone string, serviceID uint, note string, status int8) (map[string]interface{}, error) {
	log.Printf("管理员创建队列，店铺ID: %d, 客户名: %s, 电话: %s, 服务ID: %d, 状态: %d", 
		shopID, customerName, phone, serviceID, status)

	// 移除服务检查
	// 获取服务名称，如果无法获取，使用默认值
	serviceName := "未知服务"
	service, err := s.serviceRepo.GetByID(serviceID)
	log.Printf("服务: %v", service)
	if err == nil && service != nil {
		serviceName = service.Name
	}
	
	// 检查是否已有相同电话的客户
	var userID string
	user, err := s.userRepo.GetByPhone(phone)
	if err != nil || user == nil {
		// 如果用户不存在，创建一个临时用户
		tempUser := &model.User{
			ID:        fmt.Sprintf("temp_%d", time.Now().UnixNano()),
			Nickname:  customerName,
			Phone:     phone,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = s.userRepo.Create(tempUser)
		if err != nil {
			return nil, errors.New("创建临时用户失败: " + err.Error())
		}
		userID = tempUser.ID
	} else {
		userID = user.ID
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
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	err = s.queueRepo.Create(queue)
	if err != nil {
		return nil, errors.New("创建队列失败: " + err.Error())
	}
	
	// 构造返回数据，不再尝试重新加载队列
	result := map[string]interface{}{
		"id":            queue.ID,
		"shop_id":       queue.ShopID,
		"customer_name": customerName,
		"phone":         phone,
		"queue_number":  queue.QueueNumber,
		"service_id":    queue.ServiceID,
		"service_name":  serviceName,
		"status":        queue.Status,
		"created_at":    queue.CreatedAt.Format("2006-01-02 15:04:05"),
		"note":          note,
	}
	
	return result, nil
}
