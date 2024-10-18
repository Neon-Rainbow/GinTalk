USE `GinTalk`;

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user`
(
    `id`          bigint(20)                             NOT NULL AUTO_INCREMENT COMMENT '自增主键，唯一标识用户记录',
    `user_id`     bigint(20)                             NOT NULL COMMENT '用户ID，用于业务中的用户唯一标识',
    `username`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名，唯一且不区分大小写',
    `password`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户密码，存储的是哈希值',
    `email`       varchar(64) COLLATE utf8mb4_general_ci COMMENT '用户邮箱，可为空',
    `gender`      tinyint(4)                             NOT NULL DEFAULT '0' COMMENT '用户性别：0-未知，1-男，2-女',
    `create_time` timestamp                              NULL     DEFAULT CURRENT_TIMESTAMP COMMENT '记录的创建时间',
    `update_time` timestamp                              NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录的最后更新时间',
    `delete_time` bigint                           NULL DEFAULT 0 COMMENT '逻辑删除时间，NULL表示未删除',

    PRIMARY KEY (`id`) COMMENT '主键索引',

    -- 联合唯一索引：确保未删除的用户名唯一
    UNIQUE KEY `idx_username_delete_time` (`username`, `delete_time`) USING BTREE COMMENT '联合索引：用户名和删除时间确保未删除的用户名唯一',

    -- 联合唯一索引：确保未删除的用户ID唯一
    UNIQUE KEY `idx_user_id_delete_time` (`user_id`, `delete_time`) USING BTREE COMMENT '联合索引：用户ID和删除时间确保未删除的用户ID唯一'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci
    COMMENT = '用户信息表：存储用户基本信息及状态';


DROP TABLE IF EXISTS `community`;
CREATE TABLE `community`
(
    `id`             int(11)                                 NOT NULL AUTO_INCREMENT,
    `community_id`   int(10) unsigned                        NOT NULL,
    `community_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
    `introduction`   varchar(256) COLLATE utf8mb4_general_ci NOT NULL,
    `create_time`    timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`    bigint                            NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_community_id_delete_time` (`community_id`, `delete_time`),
    UNIQUE KEY `idx_community_name_delete_time` (`community_name`, `delete_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;
INSERT INTO `community`
VALUES ('1', '1', 'Go', 'Golang', '2016-11-01 08:10:10', '2016-11-01 08:10:10', NULL);
INSERT INTO `community`
VALUES ('2', '2', 'leetcode', '刷题刷题刷题', '2020-01-01 08:00:00', '2020-01-01 08:00:00', NULL);
INSERT INTO `community`
VALUES ('3', '3', 'PUBG', '大吉大利，今晚吃鸡。', '2018-08-07 08:30:00', '2018-08-07 08:30:00', NULL);
INSERT INTO `community`
VALUES ('4', '4', 'LOL', '欢迎来到英雄联盟!', '2016-01-01 08:00:00', '2016-01-01 08:00:00', NULL);

DROP TABLE IF EXISTS `post`;

CREATE TABLE `post`
(
    `id`           bigint(20)                               NOT NULL AUTO_INCREMENT COMMENT '自增主键，唯一标识每条帖子记录',
    `post_id`      bigint(20)                               NOT NULL COMMENT '帖子ID，用于业务中的帖子唯一标识',
    `title`        varchar(128) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '帖子标题',
    `content`      varchar(8192) COLLATE utf8mb4_general_ci NOT NULL COMMENT '帖子内容，最大支持8192字符',
    `author_id`    bigint(20)                               NOT NULL COMMENT '作者的用户ID，用于关联用户表',
    `community_id` bigint(20)                               NOT NULL COMMENT '所属社区ID，用于关联社区表',
    `status`       tinyint(4)                               NOT NULL DEFAULT '1' COMMENT '帖子状态：1-正常，0-隐藏或删除',
    `create_time`  timestamp                                NULL     DEFAULT CURRENT_TIMESTAMP COMMENT '帖子创建时间，默认当前时间',
    `update_time`  timestamp                                NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '帖子更新时间，每次更新时自动修改',
    `delete_time`  bigint                               NULL DEFAULT 0 COMMENT '逻辑删除时间，NULL表示未删除',

    PRIMARY KEY (`id`) COMMENT '主键索引',

    UNIQUE KEY `idx_post_id_delete_time` (`post_id`, `delete_time`) COMMENT '联合索引：帖子ID和删除时间确保未删除的帖子ID唯一',

    KEY `idx_author_id` (`author_id`) COMMENT '普通索引：按作者ID查询帖子',

    KEY `idx_community_id` (`community_id`) COMMENT '普通索引：按社区ID查询帖子'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci
    COMMENT = '帖子表：存储用户发布的帖子及其状态';


DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`
(
    `id`          bigint(20)                      NOT NULL AUTO_INCREMENT COMMENT '自增主键，唯一标识每条评论记录',
    `comment_id`  bigint(20) unsigned             NOT NULL COMMENT '评论ID，用于业务中的评论唯一标识',
    `content`     text COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论内容',
    `post_id`     bigint(20)                      NOT NULL COMMENT '评论所属的帖子ID',
    `author_id`   bigint(20)                      NOT NULL COMMENT '评论作者的用户ID',
    `author_name` varchar(64)                     NOT NULL COMMENT '评论时的用户的名字',
    `parent_id`   bigint(20)                      NOT NULL DEFAULT '0' COMMENT '该评论回复的评论ID，为0表示原生评论,即第一层的评论，不为0表示回复评论',
    `reply_id`    bigint(20)                      NOT NULL COMMENT '父评论ID, 为0表示原生评论，不为0表示回复评论',
    `status`      tinyint(3) unsigned             NOT NULL DEFAULT '1' COMMENT '评论状态：1-正常，0-删除',
    `create_time` timestamp                       NULL     DEFAULT CURRENT_TIMESTAMP COMMENT '评论创建时间，默认当前时间',
    `update_time` timestamp                       NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '评论更新时间，每次更新时自动修改',
    `delete_time` bigint                      NULL DEFAULT 0 COMMENT '逻辑删除时间，NULL表示未删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_comment_id_delete_time` (`comment_id`, `delete_time`) COMMENT '联合索引：评论ID和删除时间确保未删除的评论ID唯一',
    KEY `idx_author_Id` (`author_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

DROP TABLE IF EXISTS `vote`;
CREATE TABLE `vote`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键，唯一标识每条投票记录',
    `post_id`     bigint(20) NOT NULL COMMENT '投票所属的帖子ID',
    `user_id`     bigint(20) NOT NULL COMMENT '投票用户的用户ID',
    `vote`        tinyint(4) NOT NULL COMMENT '投票类型：1-赞，-1-踩',
    `create_time` timestamp  NULL DEFAULT CURRENT_TIMESTAMP COMMENT '投票创建时间，默认当前时间',
    `update_time` timestamp  NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '投票更新时间，每次更新时自动修改',
    `delete_time` bigint  NULL DEFAULT 0 COMMENT '逻辑删除时间，NULL表示未删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_post_id_user_id_delete_time` (`post_id`, `user_id`, `delete_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

DROP TABLE IF EXISTS `content_votes`;
CREATE TABLE `content_votes`
(
    `post_id`     bigint(20) NOT NULL COMMENT '投票所属的帖子ID',
    `count`       int(11)    NOT NULL DEFAULT '0' COMMENT '投票总数',
    `up`          int(11)    NOT NULL DEFAULT '0' COMMENT '赞数',
    `down`        int(11)    NOT NULL DEFAULT '0' COMMENT '踩数',
    `create_time` timestamp  NULL     DEFAULT CURRENT_TIMESTAMP COMMENT '投票创建时间，默认当前时间',
    `update_time` timestamp  NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '投票更新时间，每次更新时自动修改',
    `delete_time` bigint  NULL DEFAULT 0 COMMENT '逻辑删除时间，NULL表示未删除',
    UNIQUE KEY `idx_post_id_delete_time` (`post_id`, `delete_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;