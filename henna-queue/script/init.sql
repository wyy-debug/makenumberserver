-- 创建数据库
CREATE DATABASE IF NOT EXISTS henna_queue DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE henna_queue;

-- 创建店铺表
CREATE TABLE IF NOT EXISTS shops (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    phone VARCHAR(20),
    business_hours VARCHAR(100),
    description TEXT,
    cover_image VARCHAR(255),
    rating DECIMAL(2,1) DEFAULT 5.0,
    status TINYINT DEFAULT 1 COMMENT '1: 营业中, 0: 休息中',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建管理员表
CREATE TABLE IF NOT EXISTS admins (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    shop_id INT UNSIGNED,
    role TINYINT NOT NULL DEFAULT 1 COMMENT '1: 普通管理员, 2: 超级管理员',
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(50) PRIMARY KEY COMMENT 'OpenID',
    union_id VARCHAR(50),
    nickname VARCHAR(50),
    avatar_url VARCHAR(255),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建服务表
CREATE TABLE IF NOT EXISTS services (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    shop_id INT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    duration INT NOT NULL COMMENT '服务时长(分钟)',
    description TEXT,
    status TINYINT DEFAULT 1 COMMENT '1: 可用, 0: 不可用',
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建排队表
CREATE TABLE IF NOT EXISTS queues (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    shop_id INT UNSIGNED NOT NULL,
    queue_number VARCHAR(10) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    service_id INT UNSIGNED NOT NULL,
    status TINYINT DEFAULT 0 COMMENT '0: 等待中, 1: 就位中, 2: 服务中, 3: 已完成, 4: 已取消',
    estimated_wait_time INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE,
    INDEX (shop_id, status),
    INDEX (user_id, shop_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建统计表
CREATE TABLE IF NOT EXISTS statistics (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    shop_id INT UNSIGNED NOT NULL,
    date DATE NOT NULL,
    served_count INT DEFAULT 0,
    queue_count INT DEFAULT 0,
    cancel_count INT DEFAULT 0,
    avg_wait_time INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
    UNIQUE KEY (shop_id, date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建分类表
CREATE TABLE IF NOT EXISTS categories (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    shop_id INT UNSIGNED NOT NULL,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(50) NOT NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
    UNIQUE KEY idx_shop_code (shop_id, code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建图案表
CREATE TABLE IF NOT EXISTS tattoo_designs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    shop_id INT UNSIGNED NOT NULL,
    title VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    description TEXT,
    likes INT DEFAULT 0,
    status TINYINT DEFAULT 1 COMMENT '1: 可用, 0: 不可用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
    INDEX (shop_id, category),
    INDEX (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建收藏表
CREATE TABLE IF NOT EXISTS favorites (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    design_id INT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (design_id) REFERENCES tattoo_designs(id) ON DELETE CASCADE,
    UNIQUE KEY (user_id, design_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入初始超级管理员账号
INSERT INTO admins (username, password_hash, role, created_at, updated_at)
VALUES ('admin', '$2a$10$GGKQzVzP9hsN9VBOuB1eaOiVrIDyxFEfKsHWK7xTc87Nkgi/gEre6', 2, NOW(), NOW());
-- 密码为 'admin123'，使用bcrypt加密