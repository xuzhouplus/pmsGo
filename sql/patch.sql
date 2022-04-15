ALTER TABLE `pms_files`
    ADD COLUMN `status` int NOT NULL COMMENT '状态，0上传成功，1文件处理中，2可用' AFTER `preview`;

ALTER TABLE `pms_files`
    ADD COLUMN `poster` varchar(255) NULL COMMENT '封面图' AFTER `name`;

ALTER TABLE `pms_carousels`
    ADD COLUMN `timeout` int NULL COMMENT '展示时长，单位s' AFTER `switch_type`;

ALTER TABLE `pms_files`
    ADD COLUMN `extension` varchar(255) NULL COMMENT '文件后缀' AFTER `uuid`;

ALTER TABLE `pms_carousels`
    ADD COLUMN `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态，1启用，0禁用' AFTER `timeout`;

ALTER TABLE `pms_carousels`
    ADD COLUMN `title_style` text NULL COMMENT '标题文字样式' AFTER `status`,
ADD COLUMN `description_style` text NULL COMMENT '描述文字样式' AFTER `title_style`;