/*
Navicat MySQL Data Transfer

Source Server         : mysql8.1.0
Source Server Version : 80100
Source Host           : localhost:3308
Source Database       : _t

Target Server Type    : MYSQL
Target Server Version : 80100
File Encoding         : 65001

Date: 2024-11-24 23:48:15
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for tab
-- ----------------------------
DROP TABLE IF EXISTS `tab`;
CREATE TABLE `tab` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of tab
-- ----------------------------
