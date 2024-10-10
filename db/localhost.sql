-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Oct 08, 2024 at 05:32 PM
-- Server version: 10.6.19-MariaDB-cll-lve
-- PHP Version: 8.3.9

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `negrikui_ypgmerchant`
--
CREATE DATABASE IF NOT EXISTS `negrikui_ypgmerchant` DEFAULT CHARACTER SET latin1 COLLATE latin1_swedish_ci;
USE `negrikui_ypgmerchant`;

-- --------------------------------------------------------

--
-- Table structure for table `UserOrder`
--

CREATE TABLE `UserOrder` (
  `OrderID` varchar(32) NOT NULL,
  `Usertoken` text NOT NULL,
  `ReceivedTime` timestamp NOT NULL DEFAULT current_timestamp(),
  `PaymentValidateTime` timestamp NULL DEFAULT NULL,
  `PaymentReceivedTime` timestamp NULL DEFAULT NULL,
  `PaymentValidatePayload` mediumtext DEFAULT NULL,
  `PaymentReceivedPayload` mediumtext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `UserOrder`
--
ALTER TABLE `UserOrder`
  ADD PRIMARY KEY (`OrderID`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
