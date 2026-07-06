-- 帖子全文检索、热度排序、敏感词拦截
-- 依赖 MySQL InnoDB FULLTEXT；中文搜索使用内置 ngram parser。

SET @ddl := (
    SELECT IF(
        EXISTS (
            SELECT 1
            FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'posts'
              AND COLUMN_NAME = 'reply_num'
        ),
        'SELECT ''posts.reply_num already exists'' AS result',
        'ALTER TABLE `posts` ADD COLUMN `reply_num` int(11) NOT NULL DEFAULT 0 AFTER `like_num`'
    )
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @ddl := (
    SELECT IF(
        EXISTS (
            SELECT 1
            FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'posts'
              AND COLUMN_NAME = 'hot_score'
        ),
        'SELECT ''posts.hot_score already exists'' AS result',
        'ALTER TABLE `posts` ADD COLUMN `hot_score` decimal(18,4) NOT NULL DEFAULT 0 AFTER `score`'
    )
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @ddl := (
    SELECT IF(
        EXISTS (
            SELECT 1
            FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'posts'
              AND INDEX_NAME = 'idx_ft_posts_title_content'
        ),
        'SELECT ''idx_ft_posts_title_content already exists'' AS result',
        'ALTER TABLE `posts` ADD FULLTEXT KEY `idx_ft_posts_title_content` (`title`, `content`) WITH PARSER ngram'
    )
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @ddl := (
    SELECT IF(
        EXISTS (
            SELECT 1
            FROM information_schema.STATISTICS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'posts'
              AND INDEX_NAME = 'idx_posts_hot_score'
        ),
        'SELECT ''idx_posts_hot_score already exists'' AS result',
        'ALTER TABLE `posts` ADD KEY `idx_posts_hot_score` (`hot_score`, `id`)'
    )
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `sensitive_words`
(
    `id`         int(11) AUTO_INCREMENT PRIMARY KEY,
    `word`       varchar(100) NOT NULL,
    `enabled`    tinyint(1) NOT NULL DEFAULT 1,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `uniq_sensitive_word` (`word`),
    KEY `idx_sensitive_words_enabled` (`enabled`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

DROP VIEW IF EXISTS `posts_hotness_view`;
DROP TRIGGER IF EXISTS `trg_posts_bi_sensitive_hotness`;
DROP TRIGGER IF EXISTS `trg_posts_bu_sensitive_hotness`;
DROP TRIGGER IF EXISTS `trg_comments_bi_sensitive`;
DROP TRIGGER IF EXISTS `trg_comments_bu_sensitive`;
DROP TRIGGER IF EXISTS `trg_comments_ai_refresh_post_hotness`;
DROP TRIGGER IF EXISTS `trg_comments_au_refresh_post_hotness`;
DROP TRIGGER IF EXISTS `trg_comments_ad_refresh_post_hotness`;
DROP PROCEDURE IF EXISTS `refresh_post_hotness`;
DROP FUNCTION IF EXISTS `calc_post_hotness`;
DROP FUNCTION IF EXISTS `contains_sensitive_word`;

DELIMITER $$

CREATE FUNCTION `calc_post_hotness`(
    p_create_time varchar(30),
    p_reply_count int,
    p_like_count int
) RETURNS decimal(18,4)
    NOT DETERMINISTIC
    NO SQL
BEGIN
    DECLARE v_created_at datetime;
    DECLARE v_age_hours int DEFAULT 0;
    DECLARE v_reply_count int DEFAULT 0;
    DECLARE v_like_count int DEFAULT 0;

    SET v_created_at = CASE
        WHEN p_create_time IS NULL OR p_create_time = '' THEN NOW()
        WHEN p_create_time LIKE '%T%+__:__' THEN STR_TO_DATE(SUBSTRING(p_create_time, 1, 19), '%Y-%m-%dT%H:%i:%s')
        WHEN p_create_time LIKE '%T%Z' THEN STR_TO_DATE(SUBSTRING(p_create_time, 1, 19), '%Y-%m-%dT%H:%i:%s')
        ELSE COALESCE(STR_TO_DATE(p_create_time, '%Y-%m-%d %H:%i:%s'), NOW())
    END;
    SET v_age_hours = GREATEST(TIMESTAMPDIFF(HOUR, v_created_at, NOW()), 0);
    SET v_reply_count = GREATEST(COALESCE(p_reply_count, 0), 0);
    SET v_like_count = GREATEST(COALESCE(p_like_count, 0), 0);

    RETURN ROUND(((v_like_count * 2) + (v_reply_count * 3) + 10) / POW(v_age_hours + 2, 1.2), 4);
END$$

CREATE FUNCTION `contains_sensitive_word`(
    p_content longtext
) RETURNS tinyint(1)
    NOT DETERMINISTIC
    READS SQL DATA
BEGIN
    DECLARE v_hit int DEFAULT 0;

    IF p_content IS NULL OR p_content = '' THEN
        RETURN 0;
    END IF;

    SELECT EXISTS (
        SELECT 1
          FROM `sensitive_words`
         WHERE `enabled` = 1
           AND `word` <> ''
           AND INSTR(
               CONVERT(p_content USING utf8mb4) COLLATE utf8mb4_unicode_ci,
               CONVERT(`word` USING utf8mb4) COLLATE utf8mb4_unicode_ci
           ) > 0
         LIMIT 1
    )
      INTO v_hit
    ;

    RETURN IF(v_hit > 0, 1, 0);
END$$

CREATE PROCEDURE `refresh_post_hotness`(IN p_post_id int)
BEGIN
    UPDATE `posts` p
    LEFT JOIN (
        SELECT `target_id`, COUNT(*) AS reply_num
          FROM `comments`
         WHERE `target_type` = 'post'
           AND `deleted_at` = 0
         GROUP BY `target_id`
    ) c ON c.`target_id` = p.`id`
       SET p.`reply_num` = COALESCE(c.`reply_num`, 0),
           p.`hot_score` = `calc_post_hotness`(p.`create_time`, COALESCE(c.`reply_num`, 0), p.`like_num`)
     WHERE p.`id` = p_post_id;
END$$

CREATE TRIGGER `trg_posts_bi_sensitive_hotness`
BEFORE INSERT ON `posts`
FOR EACH ROW
BEGIN
    IF `contains_sensitive_word`(CONCAT_WS(' ', NEW.`title`, NEW.`content`, NEW.`summary`)) = 1 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '含敏感词无法发表';
    END IF;

    SET NEW.`reply_num` = COALESCE(NEW.`reply_num`, 0);
    SET NEW.`like_num` = COALESCE(NEW.`like_num`, 0);
    SET NEW.`hot_score` = `calc_post_hotness`(NEW.`create_time`, NEW.`reply_num`, NEW.`like_num`);
END$$

CREATE TRIGGER `trg_posts_bu_sensitive_hotness`
BEFORE UPDATE ON `posts`
FOR EACH ROW
BEGIN
    IF (NOT (NEW.`title` <=> OLD.`title`)
        OR NOT (NEW.`content` <=> OLD.`content`)
        OR NOT (NEW.`summary` <=> OLD.`summary`))
       AND `contains_sensitive_word`(CONCAT_WS(' ', NEW.`title`, NEW.`content`, NEW.`summary`)) = 1 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '含敏感词无法发表';
    END IF;

    SET NEW.`reply_num` = COALESCE(NEW.`reply_num`, 0);
    SET NEW.`like_num` = COALESCE(NEW.`like_num`, 0);
    SET NEW.`hot_score` = `calc_post_hotness`(NEW.`create_time`, NEW.`reply_num`, NEW.`like_num`);
END$$

CREATE TRIGGER `trg_comments_bi_sensitive`
BEFORE INSERT ON `comments`
FOR EACH ROW
BEGIN
    IF `contains_sensitive_word`(NEW.`content`) = 1 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '含敏感词无法发表';
    END IF;
END$$

CREATE TRIGGER `trg_comments_bu_sensitive`
BEFORE UPDATE ON `comments`
FOR EACH ROW
BEGIN
    IF NOT (NEW.`content` <=> OLD.`content`)
       AND `contains_sensitive_word`(NEW.`content`) = 1 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '含敏感词无法发表';
    END IF;
END$$

CREATE TRIGGER `trg_comments_ai_refresh_post_hotness`
AFTER INSERT ON `comments`
FOR EACH ROW
BEGIN
    IF NEW.`target_type` = 'post' AND NEW.`deleted_at` = 0 THEN
        UPDATE `posts`
           SET `hot_score` = `calc_post_hotness`(`create_time`, `reply_num` + 1, `like_num`),
               `reply_num` = `reply_num` + 1
         WHERE `id` = NEW.`target_id`;
    END IF;
END$$

CREATE TRIGGER `trg_comments_au_refresh_post_hotness`
AFTER UPDATE ON `comments`
FOR EACH ROW
BEGIN
    IF OLD.`target_type` = 'post' AND OLD.`deleted_at` = 0 THEN
        UPDATE `posts`
           SET `hot_score` = `calc_post_hotness`(`create_time`, GREATEST(`reply_num` - 1, 0), `like_num`),
               `reply_num` = GREATEST(`reply_num` - 1, 0)
         WHERE `id` = OLD.`target_id`;
    END IF;

    IF NEW.`target_type` = 'post' AND NEW.`deleted_at` = 0 THEN
        UPDATE `posts`
           SET `hot_score` = `calc_post_hotness`(`create_time`, `reply_num` + 1, `like_num`),
               `reply_num` = `reply_num` + 1
         WHERE `id` = NEW.`target_id`;
    END IF;
END$$

CREATE TRIGGER `trg_comments_ad_refresh_post_hotness`
AFTER DELETE ON `comments`
FOR EACH ROW
BEGIN
    IF OLD.`target_type` = 'post' AND OLD.`deleted_at` = 0 THEN
        UPDATE `posts`
           SET `hot_score` = `calc_post_hotness`(`create_time`, GREATEST(`reply_num` - 1, 0), `like_num`),
               `reply_num` = GREATEST(`reply_num` - 1, 0)
         WHERE `id` = OLD.`target_id`;
    END IF;
END$$

DELIMITER ;

UPDATE `posts` p
LEFT JOIN (
    SELECT `target_id`, COUNT(*) AS reply_num
      FROM `comments`
     WHERE `target_type` = 'post'
       AND `deleted_at` = 0
     GROUP BY `target_id`
) c ON c.`target_id` = p.`id`
   SET p.`reply_num` = COALESCE(c.`reply_num`, 0),
       p.`hot_score` = `calc_post_hotness`(p.`create_time`, COALESCE(c.`reply_num`, 0), p.`like_num`);

CREATE OR REPLACE VIEW `posts_hotness_view` AS
SELECT p.`id` AS `post_id`,
       p.`create_time`,
       p.`reply_num`,
       p.`like_num`,
       p.`hot_score`
  FROM `posts` p
 WHERE p.`re` = 0
   AND p.`is_report` = 0;
