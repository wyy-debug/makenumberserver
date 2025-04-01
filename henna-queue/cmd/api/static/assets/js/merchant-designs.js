// 图案管理JS

// 全局变量
let shopId = null;
let currentPage = 1;
let pageSize = 10;
let totalDesigns = 0;
let selectedCategory = '';
let selectedStatus = '';

// 初始化图案管理
function initDesigns() {
    // 获取店铺ID
    shopId = getUrlParameter('id');
    if (!shopId) {
        showToast('错误', '店铺ID不能为空', 'error');
        return;
    }

    // 加载图案列表
    loadDesigns();

    // 绑定事件
    $('#designSearchBtn').on('click', function() {
        currentPage = 1;
        selectedCategory = $('#designCategoryFilter').val();
        selectedStatus = $('#designStatusFilter').val();
        loadDesigns();
    });

    // 图片预览
    $('#designImage').on('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function(e) {
                $('#designPreview').attr('src', e.target.result);
            };
            reader.readAsDataURL(file);
        }
    });

    // 保存图案按钮
    $('#saveDesignBtn').on('click', saveDesign);

    // 重置表单
    $('#designModal').on('hidden.bs.modal', function() {
        resetDesignForm();
    });
}

// 加载图案列表
function loadDesigns() {
    $('#designList').html('<tr><td colspan="7" class="text-center py-3"><div class="spinner-border spinner-border-sm" role="status"></div> 加载中...</td></tr>');

    // 构建查询参数
    const params = {
        shop_id: shopId,
        page: currentPage,
        page_size: pageSize
    };

    if (selectedCategory) params.category = selectedCategory;
    if (selectedStatus) params.status = selectedStatus;

    // 发送请求到管理员API
    $.ajax({
        url: API_BASE_URL + '/api/v1/admin/designs',
        type: 'GET',
        data: params,
        headers: getAuthHeaders(),
        success: function(res) {
            if (res.code === 0) {
                renderDesignList(res.data.designs);
                totalDesigns = res.data.total;
                $('#designTotal').text(totalDesigns);
                renderPagination(totalDesigns, currentPage, pageSize, 'designPagination', function(page) {
                    currentPage = page;
                    loadDesigns();
                });
            } else {
                showToast('错误', res.msg || '加载图案失败', 'error');
                $('#designList').html('<tr><td colspan="7" class="text-center text-danger">加载失败</td></tr>');
            }
        },
        error: function() {
            showToast('错误', '网络错误，请稍后重试', 'error');
            $('#designList').html('<tr><td colspan="7" class="text-center text-danger">加载失败</td></tr>');
        }
    });
}

// 渲染图案列表
function renderDesignList(designs) {
    if (!designs || designs.length === 0) {
        $('#designList').html('<tr><td colspan="7" class="text-center">暂无图案</td></tr>');
        return;
    }

    let html = '';
    designs.forEach(function(design) {
        const categoryText = getCategoryText(design.category);
        const statusBadge = design.status === 1 ? 
            '<span class="badge bg-success">已上架</span>' : 
            '<span class="badge bg-secondary">已下架</span>';
        
        html += `
        <tr>
            <td><img src="${design.image_url}" alt="${design.title}" class="img-thumbnail" style="width: 60px; height: 60px; object-fit: cover;"></td>
            <td>${design.title}</td>
            <td>${categoryText}</td>
            <td>${design.likes_count || 0}</td>
            <td>${statusBadge}</td>
            <td>${formatDateTime(design.created_at)}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary me-1" onclick="editDesign(${design.id})">
                    <i class="bi bi-pencil-square"></i>
                </button>
                <button class="btn btn-sm btn-outline-danger" onclick="deleteDesign(${design.id})">
                    <i class="bi bi-trash"></i>
                </button>
            </td>
        </tr>`;
    });

    $('#designList').html(html);
}

// 保存图案
function saveDesign() {
    const designId = $('#designId').val();
    const isEdit = !!designId;
    
    // 表单验证
    if (!$('#designTitle').val()) {
        showToast('提示', '请输入图案标题', 'warning');
        return;
    }
    
    if (!$('#designCategory').val()) {
        showToast('提示', '请选择图案分类', 'warning');
        return;
    }
    
    // 创建FormData对象
    const formData = new FormData();
    formData.append('title', $('#designTitle').val());
    formData.append('category', $('#designCategory').val());
    formData.append('description', $('#designDescription').val());
    
    // 编辑模式添加状态
    if (isEdit) {
        formData.append('status', $('#designStatus').val());
    }
    
    // 添加图片文件（只有新增或者编辑时上传了新图片才添加）
    const imageFile = $('#designImage')[0].files[0];
    if (imageFile) {
        formData.append('image', imageFile);
    } else if (!isEdit) {
        showToast('提示', '请选择图案图片', 'warning');
        return;
    }
    
    // 显示保存中
    $('#saveDesignBtn').prop('disabled', true).html('<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 保存中...');
    
    // 发送请求
    $.ajax({
        url: API_BASE_URL + `/api/v1/admin/designs${isEdit ? '/' + designId : ''}`,
        type: isEdit ? 'PUT' : 'POST',
        data: formData,
        headers: getAuthHeaders(),
        processData: false,
        contentType: false,
        success: function(res) {
            $('#saveDesignBtn').prop('disabled', false).text('保存');
            if (res.code === 0) {
                $('#designModal').modal('hide');
                showToast('成功', isEdit ? '图案更新成功' : '图案添加成功', 'success');
                loadDesigns();
            } else {
                showToast('错误', res.msg || '操作失败', 'error');
            }
        },
        error: function() {
            $('#saveDesignBtn').prop('disabled', false).text('保存');
            showToast('错误', '网络错误，请稍后重试', 'error');
        }
    });
}

// 编辑图案
function editDesign(id) {
    // 重置表单
    resetDesignForm();
    
    // 设置模态框标题
    $('#designModalLabel').text('编辑图案');
    
    // 显示状态选择
    $('#designStatusGroup').show();
    
    // 获取图案详情
    $.ajax({
        url: API_BASE_URL + `/api/v1/admin/designs/${id}`,
        type: 'GET',
        headers: getAuthHeaders(),
        success: function(res) {
            if (res.code === 0) {
                const design = res.data;
                
                // 填充表单
                $('#designId').val(design.id);
                $('#designTitle').val(design.title);
                $('#designCategory').val(design.category);
                $('#designDescription').val(design.description);
                $('#designStatus').val(design.status);
                
                // 显示图片预览
                $('#designPreview').attr('src', design.image_url);
                
                // 图片上传非必填
                $('#designImage').prop('required', false);
                
                // 显示模态框
                $('#designModal').modal('show');
            } else {
                showToast('错误', res.msg || '获取图案详情失败', 'error');
            }
        },
        error: function() {
            showToast('错误', '网络错误，请稍后重试', 'error');
        }
    });
}

// 删除图案
function deleteDesign(id) {
    if (confirm('确定要删除这个图案吗？此操作不可恢复。')) {
        $.ajax({
            url: API_BASE_URL + `/api/v1/admin/designs/${id}`,
            type: 'DELETE',
            headers: getAuthHeaders(),
            success: function(res) {
                if (res.code === 0) {
                    showToast('成功', '图案已删除', 'success');
                    loadDesigns();
                } else {
                    showToast('错误', res.msg || '删除失败', 'error');
                }
            },
            error: function() {
                showToast('错误', '网络错误，请稍后重试', 'error');
            }
        });
    }
}

// 重置图案表单
function resetDesignForm() {
    $('#designForm')[0].reset();
    $('#designId').val('');
    $('#designModalLabel').text('添加图案');
    $('#designPreview').attr('src', '../assets/img/design-placeholder.png');
    $('#designStatusGroup').hide();
    $('#designImage').prop('required', true);
}

// 获取分类文本
function getCategoryText(category) {
    // 如果已加载分类，使用动态分类
    if (typeof getCategoryTextByCode === 'function') {
        return getCategoryTextByCode(category);
    }
    
    // 如果分类尚未加载，使用静态映射（后备方案）
    const categoryMap = {
        'hand': '手部图案',
        'foot': '脚部图案',
        'arm': '手臂图案',
        'other': '其他图案'
    };
    return categoryMap[category] || '未知分类';
}

// 页面加载完成后初始化
$(function() {
    initDesigns();
}); 