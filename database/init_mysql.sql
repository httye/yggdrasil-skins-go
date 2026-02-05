-- YggdrasilGo Minecraft皮肤站 - MySQL数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS yggdrasil_skin_server DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE yggdrasil_skin_server;

-- 用户表（扩展YggdrasilGo原有结构）
CREATE TABLE IF NOT EXISTS users (
    uuid VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    primary_player_name VARCHAR(255) UNIQUE COMMENT '主游戏名（Minecraft角色名）',
    player_uuid VARCHAR(36) COMMENT '游戏名对应UUID',
    qq_number VARCHAR(20) COMMENT 'QQ号码',
    email_verified BOOLEAN DEFAULT FALSE COMMENT '邮箱验证状态',
    email_verification_token VARCHAR(255) COMMENT '邮箱验证令牌',
    agreed_to_terms BOOLEAN DEFAULT FALSE COMMENT '用户协议同意状态',
    registration_ip VARCHAR(45) COMMENT '注册IP地址',
    last_login_ip VARCHAR(45) COMMENT '最后登录IP地址',
    last_login_at TIMESTAMP NULL COMMENT '最后登录时间',
    max_profiles INT DEFAULT 5 COMMENT '角色数量限制',
    is_banned BOOLEAN DEFAULT FALSE COMMENT '封禁状态',
    banned_reason TEXT COMMENT '封禁原因',
    banned_at TIMESTAMP NULL COMMENT '封禁时间',
    banned_by VARCHAR(36) NULL COMMENT '封禁管理员UUID',
    is_admin BOOLEAN DEFAULT FALSE COMMENT '管理员标识',
    permission_group_id INT DEFAULT 1 COMMENT '权限组ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_primary_player_name (primary_player_name),
    INDEX idx_player_uuid (player_uuid),
    INDEX idx_is_banned (is_banned),
    INDEX idx_is_admin (is_admin),
    INDEX idx_permission_group (permission_group_id),
    INDEX idx_email_verified (email_verified)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户主表';

-- 角色表（Minecraft游戏角色）
CREATE TABLE IF NOT EXISTS profiles (
    uuid VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL COMMENT '角色名',
    user_uuid VARCHAR(36) NOT NULL,
    skin_id INT NULL COMMENT '当前皮肤ID',
    cape_id INT NULL COMMENT '当前披风ID',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
    FOREIGN KEY (skin_id) REFERENCES skins(id) ON DELETE SET NULL,
    FOREIGN KEY (cape_id) REFERENCES capes(id) ON DELETE SET NULL,
    INDEX idx_user_uuid (user_uuid),
    INDEX idx_name (name),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Minecraft角色表';

-- 皮肤表
CREATE TABLE IF NOT EXISTS skins (
    id INT PRIMARY KEY AUTO_INCREMENT,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL COMMENT '皮肤名称',
    hash VARCHAR(64) UNIQUE NOT NULL COMMENT '文件哈希值',
    type ENUM('steve', 'alex') DEFAULT 'steve' COMMENT '皮肤类型',
    model_type ENUM('default', 'slim') DEFAULT 'default' COMMENT '模型类型',
    file_path VARCHAR(500) NOT NULL COMMENT '文件存储路径',
    file_size INT NOT NULL COMMENT '文件大小(字节)',
    upload_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    uploader_uuid VARCHAR(36) NOT NULL,
    is_public BOOLEAN DEFAULT FALSE COMMENT '是否公开',
    download_count INT DEFAULT 0 COMMENT '下载次数',
    likes_count INT DEFAULT 0 COMMENT '点赞数',
    is_verified BOOLEAN DEFAULT FALSE COMMENT '是否审核通过',
    verified_by VARCHAR(36) NULL COMMENT '审核管理员UUID',
    verified_at TIMESTAMP NULL COMMENT '审核时间',
    FOREIGN KEY (uploader_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
    FOREIGN KEY (verified_by) REFERENCES users(uuid) ON DELETE SET NULL,
    INDEX idx_uploader (uploader_uuid),
    INDEX idx_hash (hash),
    INDEX idx_type (type),
    INDEX idx_is_public (is_public),
    INDEX idx_is_verified (is_verified),
    INDEX idx_upload_time (upload_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='皮肤资源表';

-- 披风表
CREATE TABLE IF NOT EXISTS capes (
    id INT PRIMARY KEY AUTO_INCREMENT,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL COMMENT '披风名称',
    hash VARCHAR(64) UNIQUE NOT NULL COMMENT '文件哈希值',
    file_path VARCHAR(500) NOT NULL COMMENT '文件存储路径',
    file_size INT NOT NULL COMMENT '文件大小(字节)',
    upload_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    uploader_uuid VARCHAR(36) NOT NULL,
    is_public BOOLEAN DEFAULT FALSE COMMENT '是否公开',
    download_count INT DEFAULT 0 COMMENT '下载次数',
    is_verified BOOLEAN DEFAULT FALSE COMMENT '是否审核通过',
    verified_by VARCHAR(36) NULL COMMENT '审核管理员UUID',
    verified_at TIMESTAMP NULL COMMENT '审核时间',
    FOREIGN KEY (uploader_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
    FOREIGN KEY (verified_by) REFERENCES users(uuid) ON DELETE SET NULL,
    INDEX idx_uploader (uploader_uuid),
    INDEX idx_hash (hash),
    INDEX idx_is_public (is_public),
    INDEX idx_is_verified (is_verified)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='披风资源表';

-- 公告表
CREATE TABLE IF NOT EXISTS announcements (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL COMMENT '公告标题',
    content TEXT NOT NULL COMMENT '公告内容',
    type ENUM('info', 'warning', 'error', 'success', 'update', 'maintenance') DEFAULT 'info' COMMENT '公告类型',
    priority INT DEFAULT 0 COMMENT '优先级(越高越重要)',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    target_group ENUM('all', 'users', 'admins', 'banned') DEFAULT 'all' COMMENT '目标用户组',
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    end_time TIMESTAMP NULL COMMENT '结束时间',
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(uuid) ON DELETE CASCADE,
    INDEX idx_active (is_active, start_time, end_time),
    INDEX idx_priority (priority),
    INDEX idx_type (type),
    INDEX idx_target_group (target_group)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统公告表';

-- 管理员操作日志表
CREATE TABLE IF NOT EXISTS admin_logs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    admin_uuid VARCHAR(36) NOT NULL COMMENT '管理员UUID',
    action VARCHAR(100) NOT NULL COMMENT '操作类型',
    target_user_uuid VARCHAR(36) NULL COMMENT '目标用户UUID',
    details JSON COMMENT '操作详情',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
    FOREIGN KEY (target_user_uuid) REFERENCES users(uuid) ON DELETE SET NULL,
    INDEX idx_admin (admin_uuid, created_at),
    INDEX idx_target (target_user_uuid),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员操作日志';

-- 用户操作日志表
CREATE TABLE IF NOT EXISTS user_logs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_uuid VARCHAR(36) NOT NULL COMMENT '用户UUID',
    action VARCHAR(100) NOT NULL COMMENT '操作类型',
    details JSON COMMENT '操作详情',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
    INDEX idx_user (user_uuid, created_at),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户操作日志';

-- 服务器状态表
CREATE TABLE IF NOT EXISTS server_status (
    id INT PRIMARY KEY AUTO_INCREMENT,
    server_name VARCHAR(255) NOT NULL COMMENT '服务器名称',
    server_type ENUM('survival', 'creative', 'minigames', 'auth', 'lobby', 'bedwars', 'skywars') NOT NULL COMMENT '服务器类型',
    server_address VARCHAR(255) NOT NULL COMMENT '服务器地址',
    server_port INT DEFAULT 25565 COMMENT '服务器端口',
    status ENUM('online', 'offline', 'maintenance', 'starting', 'stopping') DEFAULT 'offline' COMMENT '服务器状态',
    player_count INT DEFAULT 0 COMMENT '当前玩家数',
    max_players INT DEFAULT 0 COMMENT '最大玩家数',
    motd TEXT COMMENT '服务器描述(MOTD)',
    version VARCHAR(100) COMMENT '服务器版本',
    tps FLOAT DEFAULT 20.0 COMMENT '每秒刻数(Ticks Per Second)',
    uptime_seconds BIGINT DEFAULT 0 COMMENT '运行时间(秒)',
    memory_used BIGINT DEFAULT 0 COMMENT '已用内存(字节)',
    memory_max BIGINT DEFAULT 0 COMMENT '最大内存(字节)',
    cpu_usage FLOAT DEFAULT 0.0 COMMENT 'CPU使用率(%)',
    last_ping TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后检查时间',
    next_ping TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '下次检查时间',
    is_monitoring BOOLEAN DEFAULT TRUE COMMENT '是否监控中',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_last_ping (last_ping),
    INDEX idx_server_type (server_type),
    INDEX idx_is_monitoring (is_monitoring)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务器状态监控表';

-- 权限组表
CREATE TABLE IF NOT EXISTS permission_groups (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) UNIQUE NOT NULL COMMENT '组名称',
    description TEXT COMMENT '组描述',
    permissions JSON NOT NULL COMMENT '权限配置',
    is_default BOOLEAN DEFAULT FALSE COMMENT '是否默认组',
    is_system BOOLEAN DEFAULT FALSE COMMENT '是否系统组(不可删除)',
    priority INT DEFAULT 0 COMMENT '优先级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_default (is_default),
    INDEX idx_system (is_system),
    INDEX idx_priority (priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限组表';

-- 用户点赞皮肤关联表
CREATE TABLE IF NOT EXISTS skin_likes (
    skin_id INT NOT NULL,
    user_uuid VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (skin_id, user_uuid),
    FOREIGN KEY (skin_id) REFERENCES skins(id) ON DELETE CASCADE,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
    INDEX idx_user_likes (user_uuid, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='皮肤点赞表';

-- 皮肤标签关联表
CREATE TABLE IF NOT EXISTS skin_tags (
    id INT PRIMARY KEY AUTO_INCREMENT,
    skin_id INT NOT NULL,
    tag_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (skin_id) REFERENCES skins(id) ON DELETE CASCADE,
    INDEX idx_skin_tags (skin_id),
    INDEX idx_tag_name (tag_name),
    UNIQUE KEY unique_skin_tag (skin_id, tag_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='皮肤标签表';

-- 插入默认权限组
INSERT INTO permission_groups (id, name, description, permissions, is_default, is_system, priority) VALUES
(1, 'default_user', '普通用户', '{"upload_skin": true, "max_skins": 10, "create_profile": true, "max_profiles": 5, "download_skin": true, "rate_skin": true, "comment_skin": false}', TRUE, TRUE, 1),
(2, 'premium_user', '高级用户', '{"upload_skin": true, "max_skins": 50, "create_profile": true, "max_profiles": 20, "download_skin": true, "rate_skin": true, "comment_skin": true}', FALSE, FALSE, 10),
(3, 'moderator', '版主', '{"upload_skin": true, "max_skins": 100, "create_profile": true, "max_profiles": 50, "download_skin": true, "rate_skin": true, "comment_skin": true, "verify_skin": true, "delete_skin": true, "ban_user": true}', FALSE, TRUE, 50),
(4, 'admin', '管理员', '{"upload_skin": true, "max_skins": -1, "create_profile": true, "max_profiles": -1, "download_skin": true, "rate_skin": true, "comment_skin": true, "verify_skin": true, "delete_skin": true, "ban_user": true, "manage_users": true, "manage_announcements": true, "view_logs": true, "system_config": true}', FALSE, TRUE, 100);

-- 插入默认管理员用户 (密码: admin123)
INSERT INTO users (uuid, email, username, password, is_admin, permission_group_id, max_profiles) VALUES
('00000000-0000-0000-0000-000000000001', 'admin@yggdrasil.com', 'admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', TRUE, 4, -1);

-- 插入默认公告
INSERT INTO announcements (title, content, type, priority, is_active, target_group, start_time, end_time, created_by) VALUES
('欢迎来到YggdrasilGo皮肤站！', '这是全新的Minecraft皮肤站系统，支持皮肤上传、角色管理等功能。如有问题请联系管理员。', 'success', 10, TRUE, 'all', NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), '00000000-0000-0000-0000-000000000001'),
('服务器维护通知', '系统将于每周三凌晨2:00-4:00进行例行维护，期间可能影响正常使用。', 'warning', 5, TRUE, 'all', NOW(), DATE_ADD(NOW(), INTERVAL 365 DAY), '00000000-0000-0000-0000-000000000001');

-- 创建必要的存储过程
DELIMITER //

-- 获取用户角色数量存储过程
CREATE PROCEDURE GetUserProfileCount(IN user_uuid VARCHAR(36))
BEGIN
    SELECT COUNT(*) as profile_count 
    FROM profiles 
    WHERE user_uuid = user_uuid AND is_active = TRUE;
END //

-- 检查用户是否可以创建新角色存储过程
CREATE PROCEDURE CanUserCreateProfile(IN user_uuid VARCHAR(36), OUT can_create BOOLEAN, OUT max_allowed INT)
BEGIN
    DECLARE current_count INT;
    DECLARE user_max_profiles INT;
    DECLARE user_banned BOOLEAN;
    
    -- 获取用户限制和状态
    SELECT max_profiles, is_banned INTO user_max_profiles, user_banned
    FROM users 
    WHERE uuid = user_uuid;
    
    -- 获取当前角色数量
    SELECT COUNT(*) INTO current_count
    FROM profiles 
    WHERE user_uuid = user_uuid AND is_active = TRUE;
    
    -- 检查是否可以创建
    IF user_banned THEN
        SET can_create = FALSE;
    ELSEIF user_max_profiles = -1 THEN
        SET can_create = TRUE; -- 无限制
    ELSEIF current_count < user_max_profiles THEN
        SET can_create = TRUE;
    ELSE
        SET can_create = FALSE;
    END IF;
    
    SET max_allowed = user_max_profiles;
END //

-- 记录管理员操作日志存储过程
CREATE PROCEDURE LogAdminAction(
    IN admin_uuid VARCHAR(36),
    IN action_type VARCHAR(100),
    IN target_uuid VARCHAR(36),
    IN action_details JSON,
    IN ip_addr VARCHAR(45),
    IN user_agent TEXT
)
BEGIN
    INSERT INTO admin_logs (admin_uuid, action, target_user_uuid, details, ip_address, user_agent, created_at)
    VALUES (admin_uuid, action_type, target_uuid, action_details, ip_addr, user_agent, NOW());
END //

DELIMITER ;

-- 创建触发器：用户封禁时自动禁用其所有角色
DELIMITER //

CREATE TRIGGER after_user_ban
    AFTER UPDATE ON users
    FOR EACH ROW
BEGIN
    IF NEW.is_banned = TRUE AND OLD.is_banned = FALSE THEN
        -- 禁用用户的所有角色
        UPDATE profiles SET is_active = FALSE WHERE user_uuid = NEW.uuid;
        
        -- 记录封禁日志
        INSERT INTO admin_logs (admin_uuid, action, target_user_uuid, details, created_at)
        VALUES (NEW.banned_by, 'ban_user', NEW.uuid, 
                JSON_OBJECT('reason', NEW.banned_reason, 'previous_status', 'active'), NOW());
    END IF;
END //

DELIMITER ;

-- 创建视图：用户完整信息视图
CREATE OR REPLACE VIEW user_full_info AS
SELECT 
    u.uuid,
    u.email,
    u.username,
    u.primary_player_name,
    u.player_uuid,
    u.qq_number,
    u.email_verified,
    u.agreed_to_terms,
    u.registration_ip,
    u.last_login_ip,
    u.last_login_at,
    u.max_profiles,
    u.is_banned,
    u.banned_reason,
    u.banned_at,
    u.is_admin,
    u.permission_group_id,
    pg.name as permission_group_name,
    u.created_at,
    u.updated_at,
    (SELECT COUNT(*) FROM profiles p WHERE p.user_uuid = u.uuid AND p.is_active = TRUE) as current_profiles,
    (SELECT COUNT(*) FROM skins s WHERE s.uploader_uuid = u.uuid) as uploaded_skins,
    (SELECT COUNT(*) FROM capes c WHERE c.uploader_uuid = u.uuid) as uploaded_capes
FROM users u
LEFT JOIN permission_groups pg ON u.permission_group_id = pg.id;

-- 创建视图：皮肤统计视图
CREATE OR REPLACE VIEW skin_statistics AS
SELECT 
    s.id,
    s.uuid,
    s.name,
    s.type,
    s.model_type,
    s.upload_time,
    s.uploader_uuid,
    u.username as uploader_name,
    s.is_public,
    s.download_count,
    s.likes_count,
    s.is_verified,
    (SELECT COUNT(*) FROM user_skins us WHERE us.skin_uuid = s.uuid AND us.is_active = TRUE) as active_users,
    (SELECT GROUP_CONCAT(tag_name SEPARATOR ', ') FROM skin_tags st WHERE st.skin_id = s.id) as tags
FROM skins s
LEFT JOIN users u ON s.uploader_uuid = u.uuid;

-- 创建视图：公告有效视图
CREATE OR REPLACE VIEW active_announcements AS
SELECT 
    a.id,
    a.title,
    a.content,
    a.type,
    a.priority,
    a.target_group,
    a.start_time,
    a.end_time,
    u.username as created_by_name,
    a.created_at
FROM announcements a
LEFT JOIN users u ON a.created_by = u.uuid
WHERE a.is_active = TRUE 
  AND a.start_time <= NOW() 
  AND (a.end_time IS NULL OR a.end_time > NOW())
ORDER BY a.priority DESC, a.created_at DESC;

-- 记录用户操作日志存储过程
DELIMITER //

CREATE PROCEDURE LogUserAction(
    IN p_user_uuid VARCHAR(36),
    IN p_action VARCHAR(100),
    IN p_details JSON,
    IN p_ip_address VARCHAR(45),
    IN p_user_agent TEXT
)
BEGIN
    INSERT INTO user_logs (user_uuid, action, details, ip_address, user_agent, created_at)
    VALUES (p_user_uuid, p_action, p_details, p_ip_address, p_user_agent, NOW());
END //

DELIMITER ;

-- 权限说明
-- 普通用户(default_user): 上传皮肤(10个), 创建角色(5个), 下载皮肤, 点赞
-- 高级用户(premium_user): 上传皮肤(50个), 创建角色(20个), 下载皮肤, 点赞, 评论
-- 版主(moderator): 所有用户权限 + 审核皮肤, 删除皮肤, 封禁用户
-- 管理员(admin): 所有权限无限制

-- 数据库维护建议：
-- 1. 定期清理过期公告 (end_time < NOW())
-- 2. 定期归档旧的操作日志 (created_at < DATE_SUB(NOW(), INTERVAL 1 YEAR))
-- 3. 定期清理未使用的皮肤文件
-- 4. 定期备份数据库
-- 5. 监控数据库性能，必要时添加索引