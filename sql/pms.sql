SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for pms_admins
-- ----------------------------
DROP TABLE IF EXISTS `pms_admins`;
CREATE TABLE `pms_admins`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'UUID',
  `type` tinyint(1) NULL DEFAULT 2 COMMENT '类型，1超管，2普通',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '头像',
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账号',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `salt` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'salt',
  `status` tinyint(1) NULL DEFAULT 1 COMMENT '状态，1启用，2禁用',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `auth_key` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uuid`(`uuid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '管理账号' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pms_admins
-- ----------------------------
INSERT INTO `pms_admins` VALUES (1, 'a48c92341c1140b48f6cdc2ee77c1935', 1, '', 'admin', 'e82de84cb69a0b769c39fb0c8117420b', '6740cdf6a2fd4a59a0e9a210595afadb', 1, '2021-04-09 06:01:21', '2021-04-15 02:48:42', 'zvgfya87878');

-- ----------------------------
-- Table structure for pms_carousels
-- ----------------------------
DROP TABLE IF EXISTS `pms_carousels`;
CREATE TABLE `pms_carousels`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'uuid',
  `file_id` int NULL DEFAULT NULL COMMENT '使用的文件id',
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类型，image图片，video视频，ad广告，html网页',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '地址',
  `width` int NULL DEFAULT NULL COMMENT '幅面宽',
  `height` int NULL DEFAULT NULL COMMENT '幅面高',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '描述',
  `order` int NOT NULL DEFAULT 99 COMMENT '顺序',
  `thumb` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '预览图',
  `link` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '预览图',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uuid`(`uuid`) USING BTREE,
  INDEX `OrderIndex`(`order`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '首页幻灯片' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pms_carousels
-- ----------------------------

-- ----------------------------
-- Table structure for pms_connects
-- ----------------------------
DROP TABLE IF EXISTS `pms_connects`;
CREATE TABLE `pms_connects`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `admin_id` int NOT NULL COMMENT '所属账号',
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'wechat' COMMENT '对接类型，wechat微信，weibo微博，qq QQ',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '头像',
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账号',
  `union_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '三方授权唯一标识',
  `status` tinyint(1) NULL DEFAULT 1 COMMENT '状态，1启用，2禁用',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `union_id`(`union_id`) USING BTREE,
  INDEX `admin_id`(`admin_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '第三方账号互联' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pms_connects
-- ----------------------------

-- ----------------------------
-- Table structure for pms_files
-- ----------------------------
DROP TABLE IF EXISTS `pms_files`;
CREATE TABLE `pms_files`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
  `thumb` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '缩略图',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '文件路径',
  `width` int NULL DEFAULT NULL COMMENT '幅面宽',
  `height` int NULL DEFAULT NULL COMMENT '幅面高',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '说明',
  `preview` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '预览图',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文件管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pms_files
-- ----------------------------

-- ----------------------------
-- Table structure for pms_migration
-- ----------------------------
DROP TABLE IF EXISTS `pms_migration`;
CREATE TABLE `pms_migration`  (
  `version` varchar(180) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `apply_time` int NULL DEFAULT NULL,
  PRIMARY KEY (`version`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pms_migration
-- ----------------------------

-- ----------------------------
-- Table structure for pms_posts
-- ----------------------------
DROP TABLE IF EXISTS `pms_posts`;
CREATE TABLE `pms_posts`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'uuid',
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'html' COMMENT '类型，rt富文本，md Markdown',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `sub_title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '二级标题',
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '封面',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '内容',
  `status` int NOT NULL DEFAULT 2 COMMENT '是否启用，1启用，2禁用',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uuid`(`uuid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '稿件' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pms_posts
-- ----------------------------

-- ----------------------------
-- Table structure for pms_settings
-- ----------------------------
DROP TABLE IF EXISTS `pms_settings`;
CREATE TABLE `pms_settings`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `key` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置标识',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置名称',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'input' COMMENT '配置类型，input输入框，radio单选框，checkbox复选框，select下拉选择，multiSelect多选下拉选择，textarea文本域',
  `private` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否私有，1是，2否',
  `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '配置值',
  `options` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置选项',
  `required` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否必填，1是，2否',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置说明',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `key`(`key`) USING BTREE,
  UNIQUE INDEX `name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 32 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pms_settings
-- ----------------------------
INSERT INTO `pms_settings` VALUES (1, 'title', '站点名', 'input', 2, 'Cool', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (2, 'icp', '备案号', 'input', 2, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (3, 'version', '版本', 'input', 2, 'v1', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (4, 'maintain', '维护状态', 'radio', 2, 'false', '{\"false\":\"\\u5426\",\"true\":\"\\u662f\"}', 1, NULL);
INSERT INTO `pms_settings` VALUES (5, 'icon', 'ICON', 'input', 2, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (6, 'logo', 'LOGO', 'input', 2, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (7, 'login_duration', '登录有效时长', 'input', 1, '1800', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (8, 'carousel_type', '轮播类型', 'select', 2, 'webgl', '{\"webgl\":\"webGL\",\"bootstrap\":\"Bootstrap\"}', 1, NULL);
INSERT INTO `pms_settings` VALUES (9, 'carousel_limit', '轮播数量限制', 'input', 1, '3', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (10, 'alipay_app_id', '支付宝应用APPID', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (11, 'alipay_app_primary_key', '支付宝应用私钥', 'textarea', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (12, 'alipay_public_key', '支付宝公钥', 'textarea', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (13, 'baidu_api_key', '百度网盘ApiKey', 'textarea', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (14, 'baidu_secret_key', '百度网盘SecretKey', 'textarea', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (15, 'wechat_app_id', '微信APP_ID', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (16, 'wechat_app_secret', '微信APP_SECRET', 'textarea', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (17, 'qq_app_id', 'QQ APP_ID', 'textarea', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (18, 'qq_app_secret', 'QQ APP_SECRET', 'textarea', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (19, 'weibo_app_id', '微博APP_ID', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (20, 'weibo_app_secret', '微博APP_SECRET', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (21, 'github_application_name', 'GitHub应用名称', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (22, 'github_app_id', 'GitHub APP_ID', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (23, 'github_app_secret', 'GitHub APP_SECRET', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (24, 'facebook_app_id', 'Facebook APP_ID', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (25, 'facebook_app_secret', 'Facebook APP_SECRET', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (26, 'twitter_app_id', 'Twitter APP_ID', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (27, 'twitter_app_secret', 'Twitter APP_SECRET', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (28, 'line_app_id', 'Line APP_ID', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (29, 'line_app_secret', 'Line APP_SECRET', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (30, 'google_app_id', 'Google APP_ID', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (31, 'google_app_secret', 'Google APP_SECRET', 'input', 1, '', '', 1, NULL);
INSERT INTO `pms_settings` VALUES (32, 'carousel_interval', '轮播图间隔时间', 'input', 2, '3', '', 1, '单位秒');
INSERT INTO `pms_settings` VALUES (34, 'gitee_application_name', '码云应用名称', 'input', 1, '', NULL, 1, NULL);
INSERT INTO `pms_settings` VALUES (35, 'gitee_app_id', '码云 Client ID', 'input', 1, '', NULL, 1, NULL);
INSERT INTO `pms_settings` VALUES (36, 'gitee_app_secret', '码云 Client Secret', 'input', 1, '', NULL, 1, NULL);

SET FOREIGN_KEY_CHECKS = 1;
