-- 创建知识库表
CREATE TABLE `rag_dataset`
(
    `id`                       varchar(64) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT (uuid()) COMMENT '主键id',
    `subject_id`               varchar(64) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '租户id',
    `name`                     varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '知识库名称',
    `description`              varchar(500) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '知识库描述',
    `data_source_type`         varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '来源类型，文档上传/null',
    `indexing_technique`       varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '索引类型，高质量/经济',
    `created_at`               timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
    `created_by`               varchar(64) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '创建者',
    `updated_at`               timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
    `updated_by`               varchar(64) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '更新者',
    `embedding_model`          varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '向量模型',
    `embedding_model_provider` varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '向量模型供应商',
    `retrieval_model`          varchar(1000) COLLATE utf8mb4_general_ci         DEFAULT NULL COMMENT '检索模型',
    `built_in_field_enabled`   tinyint(1)                              NOT NULL DEFAULT '0' COMMENT '是否启用内置字段',
    `ext_info`                 varchar(1000) COLLATE utf8mb4_general_ci         DEFAULT NULL COMMENT '扩展字段',
    `process_rule`             varchar(2000) COLLATE utf8mb4_general_ci         DEFAULT NULL COMMENT '分段规则',
    `scope`                    varchar(32) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '作用域类型',
    `env`                      varchar(20) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '环境标',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='知识库';

-- 创建文档表
CREATE TABLE `rag_document`
(
    `id`                 varchar(64) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT (uuid()) COMMENT '主键id',
    `subject_id`         varchar(64) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '租户id',
    `dataset_id`         varchar(64) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '知识库id',
    `position`           int                                     NOT NULL COMMENT '序号',
    `name`               varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '文档名称',
    `doc_form`           varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'text_model' COMMENT '文档来源',
    `data_source_type`   varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '数据源类型',
    `data_source_info`   varchar(1000) COLLATE utf8mb4_general_ci         DEFAULT NULL COMMENT '数据源信息',
    `batch`              varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '批次',
    `word_count`         int                                     NOT NULL COMMENT '字符数',
    `tokens`             int                                     NOT NULL COMMENT 'token数',
    `status`             varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'waiting' COMMENT '索引状态',
    `error`              text COLLATE utf8mb4_general_ci COMMENT '异常信息',
    `status_change_time` timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '状态变更时间',
    `enabled`            tinyint(1)                              NOT NULL DEFAULT '1' COMMENT '是否启用',
    `doc_language`       varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '文档语言类型',
    `doc_metadata`       text COLLATE utf8mb4_general_ci COMMENT '文档元数据',
    `indexing_log`       varchar(2000) COLLATE utf8mb4_general_ci         DEFAULT NULL COMMENT '索引留痕',
    `created_at`         timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by`         varchar(64) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '创建者',
    `updated_at`         timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by`         varchar(64) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '更新者',
    `disabled_at`        timestamp                               NULL     DEFAULT NULL COMMENT '禁用时间',
    `disabled_by`        varchar(64) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '禁用者',
    `ext_info`           varchar(1000) COLLATE utf8mb4_general_ci         DEFAULT NULL COMMENT '扩展字段',
    `scope`              varchar(32) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '作用域类型',
    `env`                varchar(20) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '环境标',
    `status_when_error`  varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '失败前状态',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='文档表';

-- 创建文本分段表
CREATE TABLE `rag_document_segment`
(
    `id`                 varchar(64) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT (uuid()) COMMENT '主键id',
    `subject_id`         varchar(64) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '租户id',
    `dataset_id`         varchar(64) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '知识库id',
    `document_id`        varchar(64) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '文档id',
    `parent_position`    int                                     NOT NULL DEFAULT '0' COMMENT '父级序号，同一父级下子段都一样',
    `parent_group`       varchar(64) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '父级标识',
    `segment_position`   int                                     NOT NULL COMMENT '序号',
    `content`            text COLLATE utf8mb4_general_ci         NOT NULL COMMENT '分段内容',
    `answer`             text COLLATE utf8mb4_general_ci COMMENT '答案',
    `image_url`          varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '图片地址',
    `bbox_position`      text COLLATE utf8mb4_general_ci,
    `word_count`         int                                     NOT NULL COMMENT '字符数',
    `index_node_id`      varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '索引节点id',
    `index_node_hash`    varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '索引节点哈希',
    `hit_count`          int                                     NOT NULL COMMENT '命中次数',
    `enabled`            tinyint(1)                              NOT NULL DEFAULT '1' COMMENT '是否启用',
    `status`             varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'waiting' COMMENT '导入状态',
    `status_change_time` timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '状态变更时间',
    `error`              text COLLATE utf8mb4_general_ci COMMENT '错误信息',
    `created_at`         timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by`         varchar(64) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '创建者',
    `updated_at`         timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by`         varchar(64) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '更新者',
    `disabled_at`        timestamp                               NULL     DEFAULT NULL COMMENT '禁用时间',
    `disabled_by`        varchar(64) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '禁用者',
    `indexing_log`       varchar(2000) COLLATE utf8mb4_general_ci         DEFAULT NULL COMMENT '索引留痕',
    `ext_info`           varchar(1000) COLLATE utf8mb4_general_ci         DEFAULT NULL COMMENT '扩展字段',
    `scope`              varchar(32) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '作用域类型',
    `env`                varchar(20) COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '环境标',
    `status_when_error`  varchar(255) COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '失败前状态',
    PRIMARY KEY (`id`),
    KEY `idx_rag_document_segment_document_id` (`document_id`),
    KEY `idx_rag_document_segment_dataset_id` (`dataset_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='文档分段表'
