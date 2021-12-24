ALTER TABLE `pms_files`
    ADD COLUMN `status` int NOT NULL COMMENT '状态，0上传成功，1文件处理中，2可用' AFTER `preview`;