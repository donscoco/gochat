-- MySQL dump 10.13  Distrib 8.0.32, for Linux (x86_64)
--
-- Host: localhost    Database: gochat
-- ------------------------------------------------------
-- Server version	8.0.32-0ubuntu0.20.04.2

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `chat_contacts`
--

DROP TABLE IF EXISTS `chat_contacts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chat_contacts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `friend_id` bigint NOT NULL DEFAULT '0' COMMENT '好友用户id',
  `friend_nick_name` varchar(20) NOT NULL DEFAULT '0' COMMENT '好友昵称备注',
  `create_time` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '和好友的关系类型',
  PRIMARY KEY (`id`),
  KEY `ind_id` (`id`),
  KEY `ind_uf` (`user_id`,`friend_id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chat_conversation`
--

DROP TABLE IF EXISTS `chat_conversation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chat_conversation` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `recv_id` bigint NOT NULL DEFAULT '0' COMMENT '会话id，',
  `group_id` bigint NOT NULL DEFAULT '0' COMMENT '群会话id，',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '会话类型',
  `read_offset` bigint NOT NULL DEFAULT '0' COMMENT '读offset',
  `write_offset` bigint NOT NULL DEFAULT '0' COMMENT '写offset',
  `create_time` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `ind_uid_fid` (`user_id`,`recv_id`),
  KEY `ind_uid_gid` (`user_id`,`group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chat_group`
--

DROP TABLE IF EXISTS `chat_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chat_group` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '群id',
  `group_name` varchar(20) NOT NULL DEFAULT '0' COMMENT '群名称',
  `owner_id` bigint NOT NULL DEFAULT '0' COMMENT '群主id',
  `head_image` varchar(200) NOT NULL DEFAULT '0' COMMENT '群头像',
  `notice` varchar(20) NOT NULL DEFAULT '0' COMMENT '群公告',
  `create_time` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chat_group_member`
--

DROP TABLE IF EXISTS `chat_group_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chat_group_member` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
  `group_id` bigint NOT NULL DEFAULT '0' COMMENT '群id',
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `alias_name` varchar(20) NOT NULL DEFAULT '0' COMMENT '群成员别称',
  `remark` varchar(20) NOT NULL DEFAULT '0' COMMENT '群成员备注',
  `create_time` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `ind_uid` (`user_id`),
  KEY `ind_group_id` (`group_id`),
  KEY `ind_gid_uid` (`group_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chat_msg`
--

DROP TABLE IF EXISTS `chat_msg`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chat_msg` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '群id',
  `source_id` bigint NOT NULL DEFAULT '0' COMMENT '消息发送方id',
  `target_id` bigint NOT NULL DEFAULT '0' COMMENT '消息接收方id',
  `sequence` bigint NOT NULL DEFAULT '0' COMMENT '消息序列号',
  `is_withdraw` tinyint NOT NULL DEFAULT '0' COMMENT '消息是否撤回',
  `create_time` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `ind_source_id` (`source_id`,`is_withdraw`),
  KEY `ind_target_id` (`target_id`,`is_withdraw`),
  KEY `ind_sequence` (`sequence`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chat_user`
--

DROP TABLE IF EXISTS `chat_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chat_user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名称',
  `nick_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `head_image` varchar(200) NOT NULL DEFAULT '' COMMENT '用户头像',
  `sex` tinyint NOT NULL DEFAULT '0',
  `signature` varchar(20) NOT NULL DEFAULT '',
  `password` varchar(20) NOT NULL DEFAULT '' COMMENT '用户密码',
  `last_login_time` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '最后登陆时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `ind_name` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-02-21  8:57:49
