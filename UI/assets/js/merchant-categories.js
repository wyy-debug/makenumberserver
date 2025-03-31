// 分类管理JS

// 全局变量
let categoryShopId = null;
let categories = [];

// 初始化分类管理
function initCategories() {
    // 获取店铺ID（与设计共用）
    categoryShopId = shopId || getUrlParameter('id');
    if (!categoryShopId) {
        showToast('错误', '店铺ID不能为空', 'error');
        return;
    }

    // 加载分类列表
    loadCategories();

    // 绑定事件
    $('#addCategoryBtn').on('click', function() {
        resetCategoryForm();
    });

    // 保存分类按钮
    $('#saveCategoryBtn').on('click', saveCategory);

    // 重置表单
    $('#categoryModal').on('hidden.bs.modal', function() {
        resetCategoryForm();
    });
}

// 加载分类列表
function loadCategories() {
    $('#categoryList').html('<tr><td colspan="5" class="text-center py-3"><div class="spinner-border spinner-border-sm" role="status"></div> 加载中...</td></tr>');

    $.ajax({
        url: API_BASE_URL + '/api/admin/categories',
        type: 'GET',
        data: { shop_id: categoryShopId },
        headers: getAuthHeaders(),
        success: function(res) {
            if (res.code === 0) {
                categories = res.data || [];
                renderCategoryList(categories);
                
                // 更新设计表单中的分类选择
                updateDesignCategorySelect();
            } else {
                showToast('错误', res.msg || '加载分类失败', 'error');
                $('#categoryList').html('<tr><td colspan="5" class="text-center text-danger">加载失败</td></tr>');
            }
        },
        error: function() {
            showToast('错误', '网络错误，请稍后重试', 'error');
            $('#categoryList').html('<tr><td colspan="5" class="text-center text-danger">加载失败</td></tr>');
        }
    });
}

// 渲染分类列表
function renderCategoryList(categories) {
    if (!categories || categories.length === 0) {
        $('#categoryList').html('<tr><td colspan="5" class="text-center">暂无分类，请添加</td></tr>');
        return;
    }

    let html = '';
    categories.forEach(function(category) {
        html += `
        <tr>
            <td>${category.id}</td>
            <td>${category.name}</td>
            <td><code>${category.code}</code></td>
            <td>${formatDateTime(category.created_at)}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary me-1" onclick="editCategory(${category.id})">
                    <i class="bi bi-pencil-square"></i>
                </button>
                <button class="btn btn-sm btn-outline-danger" onclick="deleteCategory(${category.id})">
                    <i class="bi bi-trash"></i>
                </button>
            </td>
        </tr>`;
    });

    $('#categoryList').html(html);
}

// 更新设计表单中的分类选择
function updateDesignCategorySelect() {
    const $designCategory = $('#designCategory');
    const $designCategoryFilter = $('#designCategoryFilter');
    
    // 保存当前选中的值
    const selectedCategory = $designCategory.val();
    const selectedFilter = $designCategoryFilter.val();
    
    // 清空选项
    $designCategory.empty();
    $designCategoryFilter.empty();
    
    // 添加筛选的空选项
    $designCategoryFilter.append('<option value="">所有分类</option>');
    
    // 如果没有分类，添加默认选项
    if (!categories || categories.length === 0) {
        $designCategory.append('<option value="other">其他图案</option>');
        $designCategoryFilter.append('<option value="other">其他图案</option>');
    } else {
        // 添加所有分类
        categories.forEach(function(category) {
            $designCategory.append(`<option value="${category.code}">${category.name}</option>`);
            $designCategoryFilter.append(`<option value="${category.code}">${category.name}</option>`);
        });
    }
    
    // 还原选中的值
    if (selectedCategory) $designCategory.val(selectedCategory);
    if (selectedFilter) $designCategoryFilter.val(selectedFilter);
}

// 保存分类
function saveCategory() {
    const categoryId = $('#categoryId').val();
    const isEdit = !!categoryId;
    
    // 表单验证
    if (!$('#categoryName').val()) {
        showToast('提示', '请输入分类名称', 'warning');
        return;
    }
    
    const categoryCode = $('#categoryCode').val();
    if (!categoryCode) {
        showToast('提示', '请输入分类代码', 'warning');
        return;
    }
    
    // 验证代码格式
    if (!/^[a-zA-Z_]+$/.test(categoryCode)) {
        showToast('提示', '分类代码只能包含字母和下划线', 'warning');
        return;
    }
    
    // 验证代码是否重复（仅在添加时）
    if (!isEdit && categories.some(c => c.code === categoryCode)) {
        showToast('提示', '分类代码已存在', 'warning');
        return;
    }
    
    // 准备数据
    const data = {
        name: $('#categoryName').val(),
        code: categoryCode,
        sort_order: parseInt($('#categorySort').val()) || 0,
        shop_id: categoryShopId
    };
    
    // 显示保存中
    $('#saveCategoryBtn').prop('disabled', true).html('<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 保存中...');
    
    // 发送请求
    $.ajax({
        url: API_BASE_URL + `/api/admin/categories${isEdit ? '/' + categoryId : ''}`,
        type: isEdit ? 'PUT' : 'POST',
        data: JSON.stringify(data),
        headers: getAuthHeaders(),
        contentType: 'application/json',
        success: function(res) {
            $('#saveCategoryBtn').prop('disabled', false).text('保存');
            if (res.code === 0) {
                $('#categoryModal').modal('hide');
                showToast('成功', isEdit ? '分类更新成功' : '分类添加成功', 'success');
                loadCategories();
            } else {
                showToast('错误', res.msg || '操作失败', 'error');
            }
        },
        error: function() {
            $('#saveCategoryBtn').prop('disabled', false).text('保存');
            showToast('错误', '网络错误，请稍后重试', 'error');
        }
    });
}

// 编辑分类
function editCategory(id) {
    // 重置表单
    resetCategoryForm();
    
    // 查找分类
    const category = categories.find(c => c.id === id);
    if (!category) {
        showToast('错误', '分类不存在', 'error');
        return;
    }
    
    // 设置模态框标题
    $('#categoryModalLabel').text('编辑分类');
    
    // 填充表单
    $('#categoryId').val(category.id);
    $('#categoryName').val(category.name);
    $('#categoryCode').val(category.code);
    $('#categorySort').val(category.sort_order || 0);
    
    // 分类代码在编辑时不可修改
    $('#categoryCode').prop('readonly', true);
    
    // 显示模态框
    $('#categoryModal').modal('show');
}

// 删除分类
function deleteCategory(id) {
    const category = categories.find(c => c.id === id);
    if (!category) {
        showToast('错误', '分类不存在', 'error');
        return;
    }
    
    if (confirm(`确定要删除 "${category.name}" 分类吗？此操作不可恢复，该分类下的图案将被设为"其他"分类。`)) {
        $.ajax({
            url: API_BASE_URL + `/api/admin/categories/${id}`,
            type: 'DELETE',
            headers: getAuthHeaders(),
            success: function(res) {
                if (res.code === 0) {
                    showToast('成功', '分类已删除', 'success');
                    loadCategories();
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

// 重置分类表单
function resetCategoryForm() {
    $('#categoryForm')[0].reset();
    $('#categoryId').val('');
    $('#categoryModalLabel').text('添加分类');
    $('#categoryCode').prop('readonly', false);
}

// 获取分类文本（按代码）
function getCategoryTextByCode(code) {
    const category = categories.find(c => c.code === code);
    return category ? category.name : '其他图案';
}

// 页面加载完成后初始化
$(function() {
    // 延迟执行初始化，确保商店ID已经加载
    setTimeout(initCategories, 500);
}); 