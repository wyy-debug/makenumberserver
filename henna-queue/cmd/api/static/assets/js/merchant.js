$(document).ready(function() {
    // 检查管理员权限
    checkSuperAdminAuth();
    
    // 初始化页面
    let currentPage = 1;
    const pageSize = 10;
    let keyword = '';
    
    // 加载商户列表
    loadMerchants(currentPage, pageSize, keyword);
    
    // 搜索功能
    $('#searchBtn').on('click', function() {
        keyword = $('#searchInput').val().trim();
        currentPage = 1;
        loadMerchants(currentPage, pageSize, keyword);
    });
    
    // 回车搜索
    $('#searchInput').on('keypress', function(e) {
        if (e.which === 13) {
            keyword = $(this).val().trim();
            currentPage = 1;
            loadMerchants(currentPage, pageSize, keyword);
        }
    });
    
    // 删除商户
    let merchantIdToDelete;
    $(document).on('click', '.delete-btn', function() {
        merchantIdToDelete = $(this).data('id');
        $('#deleteModal').modal('show');
    });
    
    $('#confirmDeleteBtn').on('click', function() {
        if (merchantIdToDelete) {
            deleteMerchant(merchantIdToDelete);
        }
    });
    
    // 注销功能
    $('#logoutBtn').on('click', function() {
        logout();
    });
});

// 加载商户列表
function loadMerchants(page, pageSize, keyword) {
    const token = localStorage.getItem('token');
    
    // 构建查询参数
    const params = new URLSearchParams();
    params.append('page', page);
    params.append('page_size', pageSize);
    if (keyword) {
        params.append('keyword', keyword);
    }
    
    $.ajax({
        url: API_BASE_URL + '/api/v1/admin/shops?' + params.toString(),
        type: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                renderMerchantList(res.data.list, res.data.total, page, pageSize);
            } else {
                showToast(res.message || '加载商户列表失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('加载商户列表失败', xhr);
            showToast('加载商户列表失败', 'error');
            
            if (xhr.status === 401) {
                window.location.href = '../login.html';
            }
        }
    });
}

// 渲染商户列表
function renderMerchantList(merchants, total, currentPage, pageSize) {
    const $tableBody = $('#merchantList');
    $tableBody.empty();
    
    if (merchants && merchants.length > 0) {
        merchants.forEach(function(merchant) {
            const statusBadge = merchant.status === 1 
                ? '<span class="badge bg-success">营业中</span>' 
                : '<span class="badge bg-secondary">暂停营业</span>';
            
            const createdAt = new Date(merchant.created_at).toLocaleString();
            
            const row = `
                <tr>
                    <td>${merchant.id}</td>
                    <td>${merchant.name}</td>
                    <td>${merchant.phone}</td>
                    <td>${merchant.address}</td>
                    <td>${statusBadge}</td>
                    <td>${createdAt}</td>
                    <td>
                        <div class="btn-group btn-group-sm" role="group">
                            <a href="./detail.html?id=${merchant.id}" class="btn btn-outline-primary">
                                <i class="bi bi-eye"></i>
                            </a>
                            <a href="./edit.html?id=${merchant.id}" class="btn btn-outline-secondary">
                                <i class="bi bi-pencil"></i>
                            </a>
                            <button class="btn btn-outline-danger delete-btn" data-id="${merchant.id}">
                                <i class="bi bi-trash"></i>
                            </button>
                        </div>
                    </td>
                </tr>
            `;
            
            $tableBody.append(row);
        });
    } else {
        $tableBody.html('<tr><td colspan="7" class="text-center">暂无商户数据</td></tr>');
    }
    
    // 生成分页
    generatePagination(total, pageSize, currentPage, function(page) {
        loadMerchants(page, pageSize, $('#searchInput').val().trim());
    });
}

// 删除商户
function deleteMerchant(merchantId) {
    const token = localStorage.getItem('token');
    
    $.ajax({
        url: API_BASE_URL + `/api/v1/admin/shops/${merchantId}`,
        type: 'DELETE',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                showToast('商户删除成功');
                $('#deleteModal').modal('hide');
                
                // 重新加载当前页
                const currentPage = parseInt($('.pagination .active .page-link').data('page')) || 1;
                const keyword = $('#searchInput').val().trim();
                loadMerchants(currentPage, 10, keyword);
            } else {
                showToast(res.message || '删除商户失败', 'error');
            }
        },
        error: function(xhr) {
            console.error('删除商户失败', xhr);
            showToast(xhr.responseJSON?.message || '服务器错误，请稍后再试', 'error');
            
            if (xhr.status === 401) {
                window.location.href = '../login.html';
            }
        }
    });
} 