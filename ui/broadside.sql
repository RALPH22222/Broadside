-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Jun 27, 2025 at 02:00 PM
-- Server version: 11.4.5-MariaDB
-- PHP Version: 8.0.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `broadside`
--

-- --------------------------------------------------------

--
-- Table structure for table `leaderboard`
--

CREATE TABLE `leaderboard` (
  `id` int(11) NOT NULL,
  `user_id` int(11) DEFAULT NULL,
  `score` int(11) DEFAULT NULL,
  `quests_completed` int(11) DEFAULT NULL,
  `weapon_boosts` int(11) DEFAULT NULL,
  `accuracy` float DEFAULT NULL,
  `bonus_success` float DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `leaderboard`
--

INSERT INTO `leaderboard` (`id`, `user_id`, `score`, `quests_completed`, `weapon_boosts`, `accuracy`, `bonus_success`, `created_at`) VALUES
(1, 1, 110, 1, 1, 0.5, 1, '2025-06-25 12:00:04'),
(2, 2, 0, 0, 1, 0, 1, '2025-06-25 12:03:35'),
(3, 3, 110, 1, 1, 0.5, 1, '2025-06-25 12:56:38'),
(4, 4, 192, 1, 1, 0.5, 1, '2025-06-25 13:00:52'),
(5, 5, 0, 0, 0, 0, 0, '2025-06-25 13:03:45'),
(6, 6, 0, 0, 1, 0, 1, '2025-06-25 13:04:09'),
(7, 7, 0, 0, 0, 0, 0, '2025-06-25 14:36:02'),
(8, 8, 110, 1, 1, 0.5, 1, '2025-06-25 14:46:27'),
(9, 9, 196, 1, 1, 1, 1, '2025-06-25 14:47:32'),
(10, 10, 110, 1, 1, 1, 1, '2025-06-25 14:48:46'),
(11, 11, 159, 1, 1, 1, 1, '2025-06-25 14:49:22'),
(12, 12, 40, 0, 0, 0.25, 0, '2025-06-25 14:50:09'),
(13, 13, 110, 1, 1, 1, 1, '2025-06-25 15:35:00'),
(14, 14, 110, 1, 1, 1, 1, '2025-06-25 15:38:12'),
(15, 15, 110, 1, 1, 0.5, 1, '2025-06-25 15:39:38'),
(16, 16, 110, 1, 1, 1, 1, '2025-06-25 15:43:54'),
(17, 17, 110, 1, 1, 1, 1, '2025-06-25 15:47:58'),
(18, 18, 110, 1, 1, 0.5, 1, '2025-06-25 15:50:39'),
(19, 19, 196, 1, 1, 1, 1, '2025-06-25 15:51:25'),
(20, 20, 0, 0, 0, 0, 0, '2025-06-26 05:28:46'),
(21, 21, 110, 1, 1, 1, 1, '2025-06-26 06:57:06'),
(22, 22, 80, 1, 0, 0.727273, 0, '2025-06-26 09:02:04'),
(23, 23, 42, 0, 1, 0.428571, 1, '2025-06-26 09:03:53');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `name`, `created_at`) VALUES
(1, 'chex', '2025-06-25 12:00:04'),
(2, 'ralph', '2025-06-25 12:03:35'),
(3, 'chex', '2025-06-25 12:56:38'),
(4, 'ralph', '2025-06-25 13:00:52'),
(5, 'chester\'', '2025-06-25 13:03:45'),
(6, 'r', '2025-06-25 13:04:09'),
(7, 'c', '2025-06-25 14:36:02'),
(8, 'chestir', '2025-06-25 14:46:27'),
(9, 'chestir', '2025-06-25 14:47:32'),
(10, 'chestir', '2025-06-25 14:48:46'),
(11, 'chestir', '2025-06-25 14:49:22'),
(12, 'chestir', '2025-06-25 14:50:09'),
(13, 'ral', '2025-06-25 15:35:00'),
(14, 'chex', '2025-06-25 15:38:12'),
(15, 'ches', '2025-06-25 15:39:38'),
(16, 'q', '2025-06-25 15:43:54'),
(17, 'w', '2025-06-25 15:47:58'),
(18, 'e', '2025-06-25 15:50:39'),
(19, 'e', '2025-06-25 15:51:25'),
(20, 'r', '2025-06-26 05:28:46'),
(21, 'y', '2025-06-26 06:57:06'),
(22, 's', '2025-06-26 09:02:04'),
(23, 's', '2025-06-26 09:03:53'),
(24, 'z', '2025-06-26 16:38:47'),
(25, 'v', '2025-06-26 16:39:26'),
(26, 'ww', '2025-06-26 16:41:44'),
(27, 'b', '2025-06-26 16:55:49'),
(28, 'm', '2025-06-27 05:04:46'),
(29, 'mmm', '2025-06-27 05:48:51'),
(30, 'yyy', '2025-06-27 05:49:15'),
(31, 'll', '2025-06-27 05:53:48'),
(32, 'll', '2025-06-27 06:03:31'),
(33, 'gg', '2025-06-27 06:14:55'),
(34, 'ffff', '2025-06-27 10:33:08'),
(35, 'gggg', '2025-06-27 10:35:07');

-- --------------------------------------------------------

--
-- Table structure for table `user_progress`
--

CREATE TABLE `user_progress` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `subject` varchar(64) NOT NULL,
  `difficulty` varchar(16) NOT NULL,
  `completed_at` datetime DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `leaderboard`
--
ALTER TABLE `leaderboard`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `user_progress`
--
ALTER TABLE `user_progress`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `unique_progress` (`user_id`,`subject`,`difficulty`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `leaderboard`
--
ALTER TABLE `leaderboard`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=24;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=36;

--
-- AUTO_INCREMENT for table `user_progress`
--
ALTER TABLE `user_progress`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `leaderboard`
--
ALTER TABLE `leaderboard`
  ADD CONSTRAINT `leaderboard_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Constraints for table `user_progress`
--
ALTER TABLE `user_progress`
  ADD CONSTRAINT `user_progress_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
