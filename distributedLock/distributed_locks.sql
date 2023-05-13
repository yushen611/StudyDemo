/*
 Navicat Premium Data Transfer

 Source Server         : wintest
 Source Server Type    : MySQL
 Source Server Version : 80011
 Source Host           : localhost:3306
 Source Schema         : offermaker

 Target Server Type    : MySQL
 Target Server Version : 80011
 File Encoding         : 65001

 Date: 13/05/2023 22:50:48
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for distributed_locks
-- ----------------------------
DROP TABLE IF EXISTS `distributed_locks`;
CREATE TABLE `distributed_locks`  (
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '锁的名称',
  `expire_time` bigint(20) NOT NULL COMMENT '过期时间',
  PRIMARY KEY (`name`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '分布式锁表' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
