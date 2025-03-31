$(document).ready(function() {
    // 检查登录状态
    if (!checkAuth()) {
        return; // checkAuth函数内部已处理重定向
    }
    
    // 检查超级管理员权限
    checkSuperAdminAuth();
    
    // 获取商户ID
    const shopId = getUrlParam('id');
    if (!shopId) {
        showToast('商户ID不存在', 'error');
        setTimeout(() => {
            window.location.href = './index.html';
        }, 1500);
        return;
    }
    
    // 加载商户信息
    loadMerchantDetail(shopId);
    
    // 加载商户服务
    loadShopServices(shopId);
    
    // 加载商户管理员
    loadShopAdmins(shopId);
    
    // 排队管理相关
    let queueStatus = 1; // 1: 运行中, 0: 已暂停
    let currentCustomer = null;
    let waitingList = [];
    let serviceTimer = null;
    let serviceStartTime = null;

    // 加载排队数据
    loadQueueData(shopId);
    
    // 编辑按钮
    $('#editBtn').on('click', function() {
        window.location.href = `./edit.html?id=${shopId}`;
    });
    
    // 添加服务按钮
    $('#addServiceBtn').on('click', function() {
        // 重置表单
        $('#serviceForm')[0].reset();
        $('#serviceId').val('');
        $('#serviceModalLabel').text('添加服务');
        $('#statusGroup').hide(); // 新增时不显示状态
        $('#serviceModal').modal('show');
    });
    
    // 保存服务
    $('#saveServiceBtn').on('click', function() {
        const serviceForm = $('#serviceForm')[0];
        
        // 表单验证
        if (!serviceForm.checkValidity()) {
            $(serviceForm).addClass('was-validated');
            return;
        }
        
        // 获取表单数据
        const serviceData = {
            name: $('#serviceName').val(),
            duration: parseInt($('#duration').val()),
            description: $('#serviceDescription').val(),
            sort_order: parseInt($('#sortOrder').val() || 0)
        };
        
        const serviceId = $('#serviceId').val();
        
        if (serviceId) {
            // 更新服务
            serviceData.status = parseInt($('#serviceStatus').val());
            updateService(shopId, serviceId, serviceData);
        } else {
            // 创建服务
            createService(shopId, serviceData);
        }
    });
    
    // 添加管理员按钮
    $('#addAdminBtn').on('click', function() {
        // 重置表单
        $('#adminForm')[0].reset();
        $('#adminModal').modal('show');
    });
    
    // 保存管理员
    $('#saveAdminBtn').on('click', function() {
        const adminForm = $('#adminForm')[0];
        
        // 表单验证
        if (!adminForm.checkValidity()) {
            $(adminForm).addClass('was-validated');
            return;
        }
        
        // 获取表单数据
        const adminData = {
            username: $('#adminUsername').val(),
            password: $('#adminPassword').val(),
            role: parseInt($('#adminRole').val()),
            shop_id: parseInt(shopId)
        };
        
        createAdmin(adminData);
    });
    
    // 服务操作：编辑、删除
    $(document).on('click', '.edit-service', function() {
        const serviceId = $(this).data('id');
        const service = window.servicesData.find(s => s.id === serviceId);
        
        if (service) {
            // 填充表单
            $('#serviceId').val(service.id);
            $('#serviceName').val(service.name);
            $('#duration').val(service.duration);
            $('#serviceDescription').val(service.description);
            $('#sortOrder').val(service.sort_order);
            $('#serviceStatus').val(service.status.toString());
            
            // 显示状态字段
            $('#statusGroup').show();
            
            // 更新标题
            $('#serviceModalLabel').text('编辑服务');
            
            // 显示模态框
            $('#serviceModal').modal('show');
        }
    });
    
    $(document).on('click', '.delete-service', function() {
        const serviceId = $(this).data('id');
        if (confirm('确定要删除这个服务吗？')) {
            deleteService(shopId, serviceId);
        }
    });
    
    // 管理员操作：删除
    $(document).on('click', '.delete-admin', function() {
        if (confirm('确定要删除此管理员吗？删除后无法恢复。')) {
            const adminId = $(this).data('id');
            deleteAdmin(adminId);
        }
    });
    
    // 注销功能
    $('#logoutBtn').on('click', function() {
        logout();
    });

    // 暂停/恢复取号
    $('#pauseQueueBtn').on('click', function() {
        toggleQueueStatus(shopId, 0);
    });

    $('#resumeQueueBtn').on('click', function() {
        toggleQueueStatus(shopId, 1);
    });

    // 完成服务
    $('#completeServiceBtn').on('click', function() {
        if (currentCustomer) {
            completeService(shopId, currentCustomer.id);
        }
    });

    // 取消服务
    $('#cancelServiceBtn').on('click', function() {
        if (currentCustomer) {
            if (confirm('确定要取消当前客户的服务吗？')) {
                cancelService(shopId, currentCustomer.id);
            }
        }
    });

    // 呼叫下一位
    $('#callNextBtn').on('click', function() {
        callNextCustomer(shopId);
    });
});

// 加载商户详情
function loadMerchantDetail(shopId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}`,
        type: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                displayMerchantDetail(res.data);
                // 加载统计数据
                loadStatistics(shopId);
            } else {
                showToast(res.message || '加载商户信息失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('加载商户信息失败', xhr);
            showToast('加载商户信息失败', 'error');
            
            if (xhr.status === 401) {
                window.location.href = '../login.html';
            } else if (xhr.status === 404) {
                showToast('商户不存在', 'error');
                setTimeout(() => {
                    window.location.href = './index.html';
                }, 1500);
            }
        }
    });
}

// 显示商户详情
function displayMerchantDetail(shop) {
    // 设置页面标题
    document.title = `${shop.name} - 商户详情`;
    
    // 填充基本信息
    $('#shopName').text(shop.name);
    $('#phone').text(shop.phone || '-');
    $('#address').text(shop.address || '-');
    $('#businessHours').text(shop.business_hours || '-');
    $('#description').text(shop.description || '-');
    $('#createdAt').text(new Date(shop.created_at).toLocaleString());
    
    // 状态标签
    if (shop.status === 1) {
        $('#statusBadge').removeClass('bg-secondary').addClass('bg-success').text('营业中');
    } else {
        $('#statusBadge').removeClass('bg-success').addClass('bg-secondary').text('暂停营业');
    }
    
    // 封面图
    if (shop.cover_image) {
        $('#coverImage').attr('src', shop.cover_image);
    }
}

// 加载统计数据
function loadStatistics(shopId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/statistics`,
        type: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                // 更新统计数据
                $('#totalServed').text(res.data.today.completed || 0);
                $('#totalWaiting').text(res.data.today.waiting || 0);
                $('#totalServices').text(res.data.today.total || 0);
            }
        },
        error: function(xhr) {
            console.error('加载统计数据失败', xhr);
        }
    });
}

// 显示统计数据
function displayStatistics(stats) {
    if (stats.today) {
        $('#todayQueueCount').text(stats.today.queue_count || 0);
        $('#todayServedCount').text(stats.today.served_count || 0);
        $('#avgWaitTime').text((stats.today.avg_wait_time || 0) + '分钟');
    }
    
    if (stats.comparison) {
        const queuePercentage = stats.comparison.queue_percentage;
        const servedPercentage = stats.comparison.served_percentage;
        const waitPercentage = stats.comparison.wait_percentage;
        
        $('#queueCompare').text((queuePercentage > 0 ? '+' : '') + queuePercentage.toFixed(1) + '%');
        $('#queueCompare').addClass(queuePercentage > 0 ? 'text-success' : 'text-danger');
        
        $('#servedCompare').text((servedPercentage > 0 ? '+' : '') + servedPercentage.toFixed(1) + '%');
        $('#servedCompare').addClass(servedPercentage > 0 ? 'text-success' : 'text-danger');
        
        $('#waitCompare').text((waitPercentage > 0 ? '+' : '') + waitPercentage.toFixed(1) + '%');
        // 注意等待时间减少是好事
        $('#waitCompare').addClass(waitPercentage < 0 ? 'text-success' : 'text-danger');
    }
}

// 加载商户服务
function loadShopServices(shopId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/services`,
        type: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                // 保存服务数据到全局变量，方便编辑时使用
                window.servicesData = res.data;
                displayServices(res.data);
            } else {
                showToast(res.message || '加载服务列表失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('加载服务列表失败', xhr);
            $('#serviceList').html('<tr><td colspan="5" class="text-center">加载服务失败</td></tr>');
        }
    });
}

// 显示服务列表
function displayServices(services) {
    const $serviceList = $('#serviceList');
    $serviceList.empty();
    
    if (services && services.length > 0) {
        services.forEach(function(service) {
            const statusBadge = service.status === 1 
                ? '<span class="badge bg-success">启用</span>' 
                : '<span class="badge bg-secondary">禁用</span>';
            
            const row = `
                <tr>
                    <td>${service.name}</td>
                    <td>${service.duration}分钟</td>
                    <td>${service.description || '-'}</td>
                    <td>${statusBadge}</td>
                    <td>
                        <div class="btn-group btn-group-sm" role="group">
                            <button class="btn btn-outline-secondary edit-service" data-id="${service.id}">
                                <i class="bi bi-pencil"></i>
                            </button>
                            <button class="btn btn-outline-danger delete-service" data-id="${service.id}">
                                <i class="bi bi-trash"></i>
                            </button>
                        </div>
                    </td>
                </tr>
            `;
            
            $serviceList.append(row);
        });
    } else {
        $serviceList.html('<tr><td colspan="5" class="text-center">暂无服务</td></tr>');
    }
}

// 创建服务
function createService(shopId, data) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/services`,
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(data),
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('服务创建成功');
                $('#serviceModal').modal('hide');
                // 重新加载服务列表
                loadShopServices(shopId);
            } else {
                showToast(res.message || '创建服务失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('创建服务失败', xhr);
            showToast(xhr.responseJSON?.message || '服务器错误，请稍后再试', 'error');
        }
    });
}

// 更新服务
function updateService(shopId, serviceId, data) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/services/${serviceId}`,
        type: 'PUT',
        contentType: 'application/json',
        data: JSON.stringify(data),
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('服务更新成功');
                $('#serviceModal').modal('hide');
                // 重新加载服务列表
                loadShopServices(shopId);
            } else {
                showToast(res.message || '更新服务失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('更新服务失败', xhr);
            showToast(xhr.responseJSON?.message || '服务器错误，请稍后再试', 'error');
        }
    });
}

// 删除服务
function deleteService(shopId, serviceId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/services/${serviceId}`,
        type: 'DELETE',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('服务删除成功');
                // 重新加载服务列表
                loadShopServices(shopId);
            } else {
                showToast(res.message || '删除服务失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('删除服务失败', xhr);
            showToast(xhr.responseJSON?.message || '服务器错误，请稍后再试', 'error');
        }
    });
}

// 加载商户管理员
function loadShopAdmins(shopId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/admins`,
        type: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                displayAdmins(res.data);
            } else {
                showToast(res.message || '加载管理员列表失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('加载管理员列表失败', xhr);
            $('#adminList').html('<tr><td colspan="5" class="text-center">加载管理员失败</td></tr>');
        }
    });
}

// 显示管理员列表
function displayAdmins(admins) {
    const $adminList = $('#adminList');
    $adminList.empty();
    
    if (admins && admins.length > 0) {
        admins.forEach(function(admin) {
            const role = admin.role === 1 ? '管理员' : '店员';
            const lastLogin = admin.last_login ? new Date(admin.last_login).toLocaleString() : '-';
            const createdAt = new Date(admin.created_at).toLocaleString();
            
            const row = `
                <tr>
                    <td>${admin.username}</td>
                    <td>${role}</td>
                    <td>${lastLogin}</td>
                    <td>${createdAt}</td>
                    <td>
                        <button class="btn btn-sm btn-outline-danger delete-admin" data-id="${admin.id}">
                            <i class="bi bi-trash"></i>
                        </button>
                    </td>
                </tr>
            `;
            
            $adminList.append(row);
        });
    } else {
        $adminList.html('<tr><td colspan="5" class="text-center">暂无管理员</td></tr>');
    }
}

// 创建管理员
function createAdmin(data) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + '/api/admin/admins',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(data),
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('管理员创建成功');
                $('#adminModal').modal('hide');
                // 重新加载管理员列表
                loadShopAdmins(data.shop_id);
            } else {
                showToast(res.message || '创建管理员失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('创建管理员失败', xhr);
            showToast(xhr.responseJSON?.message || '服务器错误，请稍后再试', 'error');
        }
    });
}

// 删除管理员
function deleteAdmin(adminId) {
    const token = localStorage.getItem('token');
    const shopId = getUrlParam('id');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/admins/${adminId}`,
        type: 'DELETE',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('管理员删除成功');
                // 重新加载管理员列表
                loadShopAdmins(shopId);
            } else {
                showToast(res.message || '删除管理员失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('删除管理员失败', xhr);
            showToast(xhr.responseJSON?.message || '服务器错误，请稍后再试', 'error');
        }
    });
}

// 加载排队数据
function loadQueueData(shopId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/queue`,
        type: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                // 更新店铺排队状态
                queueStatus = res.data.status;
                updateQueueStatusUI();
                
                // 更新当前服务客户
                if (res.data.current) {
                    currentCustomer = res.data.current;
                    showCurrentCustomer();
                    startServiceTimer(new Date(currentCustomer.start_time));
                } else {
                    hideCurrentCustomer();
                }
                
                // 更新等待队列
                waitingList = res.data.waiting || [];
                updateWaitingList();
            } else {
                showToast(res.message || '加载排队数据失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('加载排队数据失败', xhr);
            showToast('加载排队数据失败', 'error');
        }
    });
}

// 更新队列状态UI
function updateQueueStatusUI() {
    if (queueStatus === 1) {
        $('#pauseQueueBtn').removeClass('d-none');
        $('#resumeQueueBtn').addClass('d-none');
    } else {
        $('#pauseQueueBtn').addClass('d-none');
        $('#resumeQueueBtn').removeClass('d-none');
    }
}

// 显示当前服务客户
function showCurrentCustomer() {
    if (!currentCustomer) return;
    
    $('#currentNumber').text(currentCustomer.queue_number);
    $('#currentName').text(currentCustomer.customer_name || '顾客');
    $('#currentService').text(`(${currentCustomer.service_name})`);
    
    $('#noCurrentCustomer').addClass('d-none');
    $('#currentCustomer').removeClass('d-none');
}

// 隐藏当前服务客户
function hideCurrentCustomer() {
    $('#noCurrentCustomer').removeClass('d-none');
    $('#currentCustomer').addClass('d-none');
    
    // 停止计时器
    if (serviceTimer) {
        clearInterval(serviceTimer);
        serviceTimer = null;
    }
}

// 开始服务计时器
function startServiceTimer(startTime) {
    // 停止之前的计时器
    if (serviceTimer) {
        clearInterval(serviceTimer);
    }
    
    serviceStartTime = startTime;
    
    // 更新计时显示
    updateServiceTimer();
    
    // 设置计时器，每秒更新一次
    serviceTimer = setInterval(updateServiceTimer, 1000);
}

// 更新服务计时器显示
function updateServiceTimer() {
    if (!serviceStartTime) return;
    
    const now = new Date();
    const elapsed = now - serviceStartTime;
    
    // 计算时、分、秒
    const hours = Math.floor(elapsed / 3600000);
    const minutes = Math.floor((elapsed % 3600000) / 60000);
    const seconds = Math.floor((elapsed % 60000) / 1000);
    
    // 格式化显示
    const timeString = 
        String(hours).padStart(2, '0') + ':' +
        String(minutes).padStart(2, '0') + ':' +
        String(seconds).padStart(2, '0');
    
    $('#serviceTimer').text(timeString);
}

// 更新等待队列
function updateWaitingList() {
    const $queueItems = $('#queueItems');
    $queueItems.empty();
    
    if (waitingList && waitingList.length > 0) {
        $('#emptyQueue').addClass('d-none');
        $('#waitingList').removeClass('d-none');
        
        waitingList.forEach(function(item) {
            const waitingSince = new Date(item.created_at);
            const now = new Date();
            const waitTime = Math.floor((now - waitingSince) / 60000); // 等待分钟数
            
            const row = `
                <tr>
                    <td><span class="badge bg-secondary">${item.queue_number}</span></td>
                    <td>${item.customer_name || '顾客'}</td>
                    <td>${item.service_name}</td>
                    <td>${waitTime} 分钟</td>
                    <td>
                        <button class="btn btn-sm btn-outline-primary call-customer" data-id="${item.id}">
                            呼叫
                        </button>
                    </td>
                </tr>
            `;
            
            $queueItems.append(row);
        });
        
        // 绑定呼叫按钮事件
        $('.call-customer').on('click', function() {
            const queueId = $(this).data('id');
            callCustomer(shopId, queueId);
        });
    } else {
        $('#emptyQueue').removeClass('d-none');
        $('#waitingList').addClass('d-none');
    }
}

// 切换队列状态
function toggleQueueStatus(shopId, status) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/queue/status`,
        type: 'PUT',
        contentType: 'application/json',
        data: JSON.stringify({ status: status }),
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                queueStatus = status;
                updateQueueStatusUI();
                
                if (status === 1) {
                    showToast('取号已恢复');
                } else {
                    showToast('取号已暂停');
                }
            } else {
                showToast(res.message || '操作失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('操作失败', xhr);
            showToast('操作失败，请稍后再试', 'error');
        }
    });
}

// 完成服务
function completeService(shopId, queueId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/queue/${queueId}/complete`,
        type: 'PUT',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('服务已完成');
                
                // 重新加载排队数据
                loadQueueData(shopId);
            } else {
                showToast(res.message || '操作失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('操作失败', xhr);
            showToast('操作失败，请稍后再试', 'error');
        }
    });
}

// 取消服务
function cancelService(shopId, queueId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/queue/${queueId}/cancel`,
        type: 'PUT',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('服务已取消');
                
                // 重新加载排队数据
                loadQueueData(shopId);
            } else {
                showToast(res.message || '操作失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('操作失败', xhr);
            showToast('操作失败，请稍后再试', 'error');
        }
    });
}

// 呼叫下一位
function callNextCustomer(shopId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/queue/next`,
        type: 'PUT',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('已呼叫下一位客户');
                
                // 重新加载排队数据
                loadQueueData(shopId);
            } else {
                showToast(res.message || '操作失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('操作失败', xhr);
            showToast('操作失败，请稍后再试', 'error');
        }
    });
}

// 呼叫特定客户
function callCustomer(shopId, queueId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/admin/shops/${shopId}/queue/${queueId}/call`,
        type: 'PUT',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('已呼叫客户');
                
                // 重新加载排队数据
                loadQueueData(shopId);
            } else {
                showToast(res.message || '操作失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('操作失败', xhr);
            showToast('操作失败，请稍后再试', 'error');
        }
    });
}

// 定时刷新排队数据（每30秒）
setInterval(function() {
    const shopId = getUrlParam('id');
    if (shopId) {
        loadQueueData(shopId);
    }
}, 30000); 