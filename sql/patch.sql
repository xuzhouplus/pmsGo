ALTER TABLE `pms_files`
    ADD COLUMN `status` int NOT NULL COMMENT '状态，0上传成功，1文件处理中，2可用' AFTER `preview`;

ALTER TABLE `pms_files`
    ADD COLUMN `poster` varchar(255) NULL COMMENT '封面图' AFTER `name`;