$(document).ready(function() {
    // 检查登录状态
    if (!checkAuth()) {
        return; // checkAuth函数内部已处理重定向
    }
    
    // 检查超级管理员权限
    checkSuperAdminAuth();
    
    // 图片预览
    $('#coverImage').on('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function(e) {
                $('#coverPreview').attr('src', e.target.result).show();
            };
            reader.readAsDataURL(file);
        }
    });
    
    // 表单提交
    $('#merchantForm').on('submit', async function(e) {
        e.preventDefault();
        
        // 表单验证
        if (!this.checkValidity()) {
            e.stopPropagation();
            $(this).addClass('was-validated');
            return;
        }
        
        try {
            // 禁用提交按钮，防止重复提交
            const $submitBtn = $(this).find('button[type="submit"]');
            $submitBtn.prop('disabled', true).html('<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 保存中...');
            
            // 准备表单数据
            const formData = {
                name: $('#shopName').val(),
                phone: $('#phone').val(),
                address: $('#address').val(),
                latitude: parseFloat($('#latitude').val()) || 0,
                longitude: parseFloat($('#longitude').val()) || 0,
                business_hours: $('#businessHours').val(),
                description: $('#description').val()
            };
            
            // 处理图片上传
            const coverImageFile = $('#coverImage')[0].files[0];
            if (coverImageFile) {
                formData.cover_image = await uploadImage(coverImageFile);
            }
            
            // 发送API请求
            await createMerchant(formData);
            
            // 提示成功并返回列表页
            showToast('商户添加成功');
            setTimeout(() => {
                window.location.href = './index.html';
            }, 1500);
            
        } catch (error) {
            showToast(error, 'error');
            // 恢复提交按钮
            const $submitBtn = $(this).find('button[type="submit"]');
            $submitBtn.prop('disabled', false).text('保存');
        }
    });
    
    // 选择位置按钮
    $('#selectLocationBtn').on('click', function() {
        // 这里可以调用地图选择器API，如高德地图、百度地图等
        // 简单示例，实际中应该使用地图API
        alert('请在实际应用中集成地图选择器API');
    });
    
    // 注销功能
    $('#logoutBtn').on('click', function() {
        logout();
    });
});

// 创建商户
async function createMerchant(data) {
    const token = localStorage.getItem('token');
    
    return new Promise((resolve, reject) => {
        $.ajax({
            url: API_BASE_URL + '/api/admin/shops',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(data),
            headers: {
                'Authorization': 'Bearer ' + token
            },
            success: function(res) {
                if (res.code === 200) {
                    resolve(res.data);
                } else {
                    reject(res.message || '创建商户失败');
                }
            },
            error: function(xhr) {
                console.error('创建商户失败', xhr);
                reject(xhr.responseJSON?.message || '服务器错误，请稍后再试');
                
                if (xhr.status === 401) {
                    window.location.href = '../login.html';
                }
            }
        });
    });
} 