// API基础URL，根据您的实际部署环境修改
const API_BASE_URL = 'http://localhost:8080';

// Toast通知
function showToast(message, type = 'success') {
    // 如果已存在toast，先移除
    $('.toast').remove();
    
    // 创建toast元素
    const toastHtml = `
        <div class="toast position-fixed top-0 end-0 m-4" role="alert" aria-live="assertive" aria-atomic="true">
            <div class="toast-header ${type === 'success' ? 'bg-success' : 'bg-danger'} text-white">
                <strong class="me-auto">${type === 'success' ? '成功' : '错误'}</strong>
                <button type="button" class="btn-close btn-close-white" data-bs-dismiss="toast" aria-label="Close"></button>
            </div>
            <div class="toast-body">
                ${message}
            </div>
        </div>
    `;
    
    // 添加到body
    $('body').append(toastHtml);
    
    // 显示toast
    const toastElement = new bootstrap.Toast($('.toast'), { 
        delay: 3000,
        autohide: true
    });
    
    toastElement.show();
}

// 退出登录
function logout() {
    localStorage.removeItem('token');
    window.location.href = '../login.html';
}

// 检查认证状态
function checkAuth() {
    const token = localStorage.getItem('token');
    
    if (!token) {
        window.location.href = '../login.html';
        return false;
    }
    
    return true;
}

// 检查超级管理员权限
function checkSuperAdminAuth() {
    const token = localStorage.getItem('token');
    
    if (!token) {
        window.location.href = '../login.html';
        return;
    }
    
    $.ajax({
        url: API_BASE_URL + '/api/admin/profile',
        type: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 200) {
                if (res.data.role !== 2) { // 超级管理员角色
                    showToast('您没有权限访问此页面', 'error');
                    setTimeout(function() {
                        window.location.href = '../dashboard/index.html';
                    }, 2000);
                }
            }
        },
        error: function(xhr) {
            console.error('获取管理员信息失败', xhr);
            if (xhr.status === 401) {
                window.location.href = '../login.html';
            }
        }
    });
}

// 上传图片到服务器
async function uploadImage(file) {
    const token = localStorage.getItem('token');
    const formData = new FormData();
    formData.append('file', file);
    
    return new Promise((resolve, reject) => {
        $.ajax({
            url: API_BASE_URL + '/api/upload',
            type: 'POST',
            data: formData,
            processData: false,
            contentType: false,
            headers: {
                'Authorization': 'Bearer ' + token
            },
            success: function(res) {
                if (res.code === 200) {
                    resolve(res.data.url);
                } else {
                    reject(res.message || '上传失败');
                }
            },
            error: function(xhr) {
                reject('上传失败: ' + xhr.statusText);
            }
        });
    });
}

// 获取URL查询参数
function getUrlParam(param) {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(param);
}

// 生成分页
function generatePagination(total, pageSize, currentPage, onPageClick) {
    const totalPages = Math.ceil(total / pageSize);
    const $pagination = $('#pagination');
    $pagination.empty();
    
    // 如果总页数小于等于1，不显示分页
    if (totalPages <= 1) {
        return;
    }
    
    // 上一页
    const $prevLi = $(`
        <li class="page-item ${currentPage === 1 ? 'disabled' : ''}">
            <a class="page-link" href="#" data-page="${currentPage - 1}">上一页</a>
        </li>
    `);
    $pagination.append($prevLi);
    
    // 页码
    let startPage = Math.max(1, currentPage - 2);
    let endPage = Math.min(startPage + 4, totalPages);
    
    if (endPage - startPage < 4 && startPage > 1) {
        startPage = Math.max(1, endPage - 4);
    }
    
    for (let i = startPage; i <= endPage; i++) {
        const $pageLi = $(`
            <li class="page-item ${i === currentPage ? 'active' : ''}">
                <a class="page-link" href="#" data-page="${i}">${i}</a>
            </li>
        `);
        $pagination.append($pageLi);
    }
    
    // 下一页
    const $nextLi = $(`
        <li class="page-item ${currentPage === totalPages ? 'disabled' : ''}">
            <a class="page-link" href="#" data-page="${currentPage + 1}">下一页</a>
        </li>
    `);
    $pagination.append($nextLi);
    
    // 绑定点击事件
    $('.page-link').on('click', function(e) {
        e.preventDefault();
        if ($(this).parent().hasClass('disabled')) {
            return;
        }
        
        const page = parseInt($(this).data('page'));
        if (onPageClick && typeof onPageClick === 'function') {
            onPageClick(page);
        }
    });
} 