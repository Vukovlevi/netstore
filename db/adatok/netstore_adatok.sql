-- phpMyAdmin SQL Dump
-- version 5.1.1deb5ubuntu1
-- https://www.phpmyadmin.net/
--
-- Gép: localhost:3306
-- Létrehozás ideje: 2026. Már 29. 12:16
-- Kiszolgáló verziója: 10.6.23-MariaDB-0ubuntu0.22.04.1
-- PHP verzió: 8.1.2-1ubuntu2.23

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Adatbázis: `netstore`
--
CREATE DATABASE IF NOT EXISTS `netstore` DEFAULT CHARACTER SET utf8mb3 COLLATE utf8mb3_hungarian_ci;
USE `netstore`;

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `brand`
--

CREATE TABLE `brand` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `is_own` tinyint(1) NOT NULL,
  `is_temporary` tinyint(1) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `brand`
--

INSERT INTO `brand` (`id`, `name`, `is_own`, `is_temporary`, `deleted_at`) VALUES
(1, 'FreshFarm', 1, 0, NULL),
(2, 'Coca-Cola', 0, 0, NULL),
(3, 'Nestlé', 0, 0, NULL),
(4, 'LocalBakery', 1, 1, NULL),
(5, 'Rauch', 0, 0, NULL),
(6, 'Parkside', 1, 0, NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `category`
--

CREATE TABLE `category` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `category`
--

INSERT INTO `category` (`id`, `name`, `deleted_at`) VALUES
(1, 'Italok', NULL),
(2, 'Élelmiszer', NULL),
(3, 'Háztartás', NULL),
(4, 'Ruhák', NULL),
(5, 'Gépjármű', NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `contract`
--

CREATE TABLE `contract` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `contract_type_id` int(11) NOT NULL,
  `salary` int(11) NOT NULL,
  `file` varchar(255) DEFAULT NULL,
  `starts_at` date NOT NULL,
  `ends_at` date DEFAULT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `contract`
--

INSERT INTO `contract` (`id`, `user_id`, `contract_type_id`, `salary`, `file`, `starts_at`, `ends_at`, `deleted_at`) VALUES
(1, 15, 6, 4, NULL, '2025-11-24', '2025-11-29', '2025-11-24'),
(2, 15, 1, 610000, NULL, '2025-12-19', NULL, NULL),
(3, 8, 6, 300000, NULL, '2026-01-10', NULL, '2026-01-10'),
(4, 8, 1, 600000, 'Órarend-13C-V3.pdf', '2026-01-31', '2026-02-08', '2026-01-25'),
(5, 8, 6, 120000, NULL, '2026-01-27', NULL, NULL),
(6, 16, 1, 450000, 'Osztályok20252026-1.pdf', '2026-01-26', NULL, '2026-01-25'),
(7, 16, 6, 230000, 'Osztályok20252026-1.pdf', '2026-01-26', NULL, '2026-01-25'),
(8, 16, 6, 350000, NULL, '2026-03-27', '2026-03-28', NULL),
(9, 16, 6, 510000, NULL, '2026-03-23', '2026-04-05', '2026-03-29');

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `contract_day`
--

CREATE TABLE `contract_day` (
  `id` int(11) NOT NULL,
  `starting_time` time NOT NULL,
  `ending_time` time NOT NULL,
  `contract_id` int(11) NOT NULL,
  `week_day_id` int(11) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `contract_day`
--

INSERT INTO `contract_day` (`id`, `starting_time`, `ending_time`, `contract_id`, `week_day_id`, `deleted_at`) VALUES
(3, '20:38:00', '23:38:00', 1, 1, '2025-11-24'),
(4, '22:38:00', '20:44:00', 1, 5, '2025-11-24'),
(41, '12:39:00', '18:45:00', 3, 1, '2026-01-10'),
(48, '10:36:00', '15:41:00', 4, 1, '2026-01-25'),
(50, '11:05:00', '14:09:00', 5, 1, NULL),
(51, '12:08:00', '17:14:00', 6, 1, '2026-01-25'),
(55, '11:09:00', '17:15:00', 7, 1, '2026-01-25'),
(58, '12:36:00', '18:42:00', 2, 1, NULL),
(59, '14:48:00', '19:54:00', 2, 2, NULL),
(61, '11:27:00', '11:30:00', 8, 1, NULL),
(62, '11:27:00', '11:28:00', 9, 1, '2026-03-29');

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `contract_type`
--

CREATE TABLE `contract_type` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `weekly_hours` tinyint(4) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `contract_type`
--

INSERT INTO `contract_type` (`id`, `name`, `weekly_hours`, `deleted_at`) VALUES
(1, 'Teljes állás', 40, NULL),
(2, 'torolheto', 120, '2025-10-26'),
(3, 'Fel allas', 20, '2025-11-09'),
(5, 'Felallas baby', 20, '2025-11-09'),
(6, 'Félállás', 20, NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `ingredient`
--

CREATE TABLE `ingredient` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `is_allergen` tinyint(1) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `ingredient`
--

INSERT INTO `ingredient` (`id`, `name`, `is_allergen`, `deleted_at`) VALUES
(1, 'Cukor', 0, NULL),
(2, 'Szénsav', 0, NULL),
(3, 'Narancs', 0, NULL),
(4, 'Kávébab', 0, NULL),
(5, 'Búzaliszt', 1, NULL),
(6, 'Élesztő', 0, NULL),
(7, 'Tej', 1, NULL),
(8, 'Mosószer-alapanyag', 0, NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `open_day`
--

CREATE TABLE `open_day` (
  `id` int(11) NOT NULL,
  `week_day_id` int(11) NOT NULL,
  `open_hour_id` int(11) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `open_day`
--

INSERT INTO `open_day` (`id`, `week_day_id`, `open_hour_id`, `deleted_at`) VALUES
(6, 1, 1, '2025-12-06'),
(7, 2, 1, '2025-12-06'),
(8, 3, 1, '2025-12-06'),
(9, 5, 1, '2025-12-06'),
(15, 1, 2, NULL),
(16, 2, 2, NULL),
(17, 3, 2, NULL),
(18, 4, 2, NULL),
(19, 5, 2, NULL),
(20, 6, 4, NULL),
(21, 7, 5, NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `open_hour`
--

CREATE TABLE `open_hour` (
  `id` int(11) NOT NULL,
  `opens_at` time NOT NULL,
  `closes_at` time NOT NULL,
  `special_date` date DEFAULT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `open_hour`
--

INSERT INTO `open_hour` (`id`, `opens_at`, `closes_at`, `special_date`, `deleted_at`) VALUES
(1, '13:45:00', '17:45:00', NULL, '2025-12-06'),
(2, '06:00:00', '22:00:00', NULL, NULL),
(3, '00:00:00', '00:00:00', '2026-03-15', NULL),
(4, '08:00:00', '20:00:00', NULL, NULL),
(5, '09:00:00', '19:00:00', NULL, NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `product`
--

CREATE TABLE `product` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `amount` int(11) NOT NULL,
  `expires_at` date DEFAULT NULL,
  `size` double NOT NULL,
  `size_type` varchar(10) NOT NULL,
  `price` int(11) NOT NULL,
  `warranty` date DEFAULT NULL,
  `discount` double NOT NULL,
  `other_properties` varchar(255) DEFAULT NULL,
  `type_id` int(11) NOT NULL,
  `brand_id` int(11) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `product`
--

INSERT INTO `product` (`id`, `name`, `description`, `amount`, `expires_at`, `size`, `size_type`, `price`, `warranty`, `discount`, `other_properties`, `type_id`, `brand_id`, `deleted_at`) VALUES
(1, 'Coca-Cola 0.5L', 'Szénsavas üdítőital', 100, '2026-01-01', 0.5, 'L', 350, NULL, 0, NULL, 1, 2, NULL),
(2, 'Narancslé 1L', '100% narancslé', 80, '2025-12-01', 1, 'L', 500, NULL, 0.1, NULL, 2, 3, NULL),
(3, 'Arabica Kávé', 'Őrölt arabica kávé', 50, '2026-05-01', 0.25, 'kg', 1500, NULL, 0, NULL, 3, 3, NULL),
(4, 'Teljes kiőrlésű kenyér', 'Friss pékáru jajaj', 30, '2025-10-02', 0.5, 'kg', 600, NULL, 0, NULL, 4, 4, NULL),
(5, 'Trappista sajt', 'Klasszikus magyar sajt', 40, '2025-12-15', 1, 'kg', 2200, NULL, 0.05, NULL, 5, 1, NULL),
(6, 'Persil mosószer', 'Mosószer fehér ruhákhoz vagy nem', 20, '2028-01-01', 2, 'L', 3200, NULL, 0, NULL, 6, 3, NULL),
(7, 'Vázsonyi kenyér', 'A legfinomabb kenyér', 10, '2026-02-08', 1, 'kg', 400, NULL, 0.1, NULL, 4, 4, NULL),
(8, '12V Akkumulátortöltő szgk-hoz', 'A legjobb akkutöltő csak neked, csak most, csak itt (EV-t ne töltcséé)', 8, NULL, 1, 'doboz', 12000, '2029-03-04', 0, NULL, 7, 6, NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `product_ingredient`
--

CREATE TABLE `product_ingredient` (
  `id` int(11) NOT NULL,
  `ingredient_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `amount` double NOT NULL,
  `amount_type` varchar(10) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `product_ingredient`
--

INSERT INTO `product_ingredient` (`id`, `ingredient_id`, `product_id`, `amount`, `amount_type`, `deleted_at`) VALUES
(1, 1, 1, 50, 'g', NULL),
(2, 2, 1, 2, 'g', NULL),
(3, 3, 2, 100, 'ml', NULL),
(4, 4, 3, 250, 'g', NULL),
(5, 5, 4, 300, 'g', NULL),
(6, 6, 4, 20, 'g', NULL),
(7, 7, 5, 500, 'g', NULL),
(8, 8, 6, 200, 'ml', NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `product_type`
--

CREATE TABLE `product_type` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `sub_id` int(11) NOT NULL,
  `storing_condition_id` int(11) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `product_type`
--

INSERT INTO `product_type` (`id`, `name`, `description`, `sub_id`, `storing_condition_id`, `deleted_at`) VALUES
(1, 'Szénsavas üdítő', 'Cukros, szénsavas italok', 1, 1, NULL),
(2, 'Gyümölcslé', '100%-os vagy sűrítményből készült gyümölcslé', 1, 1, NULL),
(3, 'Őrölt kávé', 'Kávéfőzőhöz őrölt kávé', 2, 3, NULL),
(4, 'Kenyér', 'Friss pékáru', 3, 1, NULL),
(5, 'Sajt', 'Tejből készült sajtok', 4, 2, NULL),
(6, 'Mosószer', 'Ruhatisztításhoz használt mosószerek', 5, 4, NULL),
(7, 'Akkumulátortöltő', 'A piacon elérhető legjobb akkumulátortöltő (Cseh Balács részére, formázni nem tud)', 6, 1, NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `role`
--

CREATE TABLE `role` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `role`
--

INSERT INTO `role` (`id`, `name`) VALUES
(6, 'Egyéb dolgozó'),
(2, 'HR'),
(5, 'Pénztáros'),
(4, 'Raktárkezelő'),
(3, 'Raktárvezető'),
(1, 'Üzletvezető');

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `session`
--

CREATE TABLE `session` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `token` varchar(64) NOT NULL,
  `expires_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `store_detail`
--

CREATE TABLE `store_detail` (
  `address` varchar(255) NOT NULL,
  `central_server_address` varchar(255) NOT NULL,
  `central_server_port` int(11) NOT NULL,
  `store_type_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `store_detail`
--

INSERT INTO `store_detail` (`address`, `central_server_address`, `central_server_port`, `store_type_id`) VALUES
('8200 Veszprém, Iskola utca 4.', '', 0, 2);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `store_type`
--

CREATE TABLE `store_type` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `store_type`
--

INSERT INTO `store_type` (`id`, `name`) VALUES
(3, 'Hipermarket'),
(1, 'Kisbolt'),
(2, 'Szupermarket');

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `storing_condition`
--

CREATE TABLE `storing_condition` (
  `id` int(11) NOT NULL,
  `description` varchar(255) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `storing_condition`
--

INSERT INTO `storing_condition` (`id`, `description`, `deleted_at`) VALUES
(1, 'Szobahőmérsékleten tárolandó (15-25°C)', NULL),
(2, 'Hűtve tárolandó (0-5°C)', NULL),
(3, 'Száraz, hűvös helyen tárolandó', NULL),
(4, 'Vegyszerek számára kijelölt tárolóhelyen', NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `sub_category`
--

CREATE TABLE `sub_category` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `category_id` int(11) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `sub_category`
--

INSERT INTO `sub_category` (`id`, `name`, `category_id`, `deleted_at`) VALUES
(1, 'Üdítők', 1, NULL),
(2, 'Kávé & Tea', 1, NULL),
(3, 'Pékáruk', 2, NULL),
(4, 'Tejtermékek', 2, NULL),
(5, 'Tisztítószerek', 3, NULL),
(6, 'Személygépjármű', 5, NULL);

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `user`
--

CREATE TABLE `user` (
  `id` int(11) NOT NULL,
  `firstname` varchar(255) NOT NULL,
  `lastname` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(72) NOT NULL,
  `password_changed` tinyint(1) NOT NULL DEFAULT 0,
  `phone_number` varchar(20) NOT NULL,
  `email` varchar(255) NOT NULL,
  `role_id` int(11) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `user`
--

INSERT INTO `user` (`id`, `firstname`, `lastname`, `username`, `password`, `password_changed`, `phone_number`, `email`, `role_id`, `deleted_at`) VALUES
(1, '-', '-', 'admin', '$2a$12$EKzHWgDilgOwYFjMm7A7GO1uUe0LA5ZFKwTsqlTVBTPeqdQUQpILO', 1, '-', '-', 1, NULL),
(4, 'admin2', 'adminovics2', 'admin22', '$2a$12$P.s0LsonRQujn5QZMA4tGe0ddLKQbEs5753WXFKP60DtqBMu9oCmG', 1, '+36205123687', 'admin@admin.com', 2, '2025-11-08'),
(5, 'admin2', 'adminovics2', 'admin2', '$2a$12$qpppHwGQTE1AHKEfWPEfke7PNG55dz0mW78FyalN6QejwXuRnrXsa', 0, '+36205123688', 'admin2@admin.com', 6, '2025-10-26'),
(6, 'admin2', 'adminovics2', 'admin3', '$2a$12$OBb9I9j4PAu1awfVHe2/2.RhWunEFRFWJJUxCPFpwmydtRGNCZxDS', 1, '+36205123689', 'admin3@admin.com', 1, '2025-10-26'),
(7, 'zsigmont', 'lakatos', 'lakatos.zsigmont', '$2a$12$Q/TEGML18UEgb.cCckH5T.GpnHuJ3e7idd0Il2hlM8kBZrQHlqw5i', 1, '', '', 6, '2025-11-08'),
(8, 'Pista', 'Kiss', 'kiss.pista', '$2a$12$8q2y6/kCWfDqaqj/by..q.s.WxApOdOxvX/WHElkvATsR18sNroqe', 1, '+36208976512', '', 6, NULL),
(9, 'asd', 'asd', 'asd', '$2a$12$3lu92Ovw1QrBBZlO1hKtkOIA90UWxUppGTs3SJluH1ERX0JooL1mS', 1, '', '', 1, '2025-11-08'),
(10, 'zsigmont', 'janos', 'janika', '$2a$12$EOjVWaC2CIvAFuaM3GxIzOfoRZVS2zyHlyRvhemtMFOB6642aZwm6', 0, '', '', 2, '2025-11-08'),
(12, 'asd', 'asd', 'asdasd', '$2a$12$3iOPSp3g4N5TEakIcCLtYedtfcOlw73wqXGgGI0.nOdxydGhuVf4m', 1, '', '', 5, '2025-11-08'),
(13, 'asd', 'asd', 'wasdasdas', '$2a$12$ccZl8ZzDe9aIfdeSQHR8pO.4/VUmGkOitrx38OI7qHrj2bSpwL8lS', 1, '', '', 5, '2025-11-08'),
(14, 'laci', 'nagy', 'nagy.laci', '$2a$12$2mcDL.ubWALynIsRBizr2ODAdgHPtgEbxhEnm4HVYU5BpopvL3mb2', 1, '', '', 1, '2025-11-10'),
(15, 'Lajos', 'Nagy', 'nagy.lajos', '$2a$12$mGCFox1M06pkjImJ46aq2OX3bvGS72zvT/DVGYiZe8LZeX0T62vaa', 1, '', 'nagy.lajos@gmail.com', 2, NULL),
(16, 'Dzsenifer', 'Lakatos', 'lakatos.dzsenifer', '$2a$12$9SagNbjY583DB2rzQJyB5.7cN4JsiT6NIznmcFjGQSRtWF6rRwW/q', 1, '', '', 5, NULL),
(17, 'wasd', 'asd', 'wasd', '$2a$12$qPtbJ29WLhRoehRKMof/POUpXvbdSqSfSOMSSVMMCMRec7aPXPYw6', 1, '', '', 2, '2025-11-24');

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `week_day`
--

CREATE TABLE `week_day` (
  `id` int(11) NOT NULL,
  `name` varchar(9) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

--
-- A tábla adatainak kiíratása `week_day`
--

INSERT INTO `week_day` (`id`, `name`) VALUES
(4, 'Csütörtök'),
(1, 'Hétfő'),
(2, 'Kedd'),
(5, 'Péntek'),
(3, 'Szerda'),
(6, 'Szombat'),
(7, 'Vasárnap');

--
-- Indexek a kiírt táblákhoz
--

--
-- A tábla indexei `brand`
--
ALTER TABLE `brand`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- A tábla indexei `category`
--
ALTER TABLE `category`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- A tábla indexei `contract`
--
ALTER TABLE `contract`
  ADD PRIMARY KEY (`id`),
  ADD KEY `contract_to_user` (`user_id`),
  ADD KEY `contract_to_type` (`contract_type_id`);

--
-- A tábla indexei `contract_day`
--
ALTER TABLE `contract_day`
  ADD PRIMARY KEY (`id`),
  ADD KEY `day_to_contract` (`contract_id`),
  ADD KEY `day_to_week_day` (`week_day_id`);

--
-- A tábla indexei `contract_type`
--
ALTER TABLE `contract_type`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- A tábla indexei `ingredient`
--
ALTER TABLE `ingredient`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- A tábla indexei `open_day`
--
ALTER TABLE `open_day`
  ADD PRIMARY KEY (`id`),
  ADD KEY `open_to_day` (`week_day_id`),
  ADD KEY `open_to_open_hour` (`open_hour_id`);

--
-- A tábla indexei `open_hour`
--
ALTER TABLE `open_hour`
  ADD PRIMARY KEY (`id`);

--
-- A tábla indexei `product`
--
ALTER TABLE `product`
  ADD PRIMARY KEY (`id`),
  ADD KEY `product_to_type` (`type_id`),
  ADD KEY `product_to_brand` (`brand_id`);

--
-- A tábla indexei `product_ingredient`
--
ALTER TABLE `product_ingredient`
  ADD PRIMARY KEY (`id`),
  ADD KEY `con_to_ingredient` (`ingredient_id`),
  ADD KEY `product_id` (`product_id`);

--
-- A tábla indexei `product_type`
--
ALTER TABLE `product_type`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`),
  ADD KEY `type_to_sub` (`sub_id`),
  ADD KEY `type_to_condition` (`storing_condition_id`);

--
-- A tábla indexei `role`
--
ALTER TABLE `role`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- A tábla indexei `session`
--
ALTER TABLE `session`
  ADD PRIMARY KEY (`id`),
  ADD KEY `session_to_user` (`user_id`);

--
-- A tábla indexei `store_detail`
--
ALTER TABLE `store_detail`
  ADD KEY `store_to_type` (`store_type_id`);

--
-- A tábla indexei `store_type`
--
ALTER TABLE `store_type`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`),
  ADD UNIQUE KEY `name_2` (`name`);

--
-- A tábla indexei `storing_condition`
--
ALTER TABLE `storing_condition`
  ADD PRIMARY KEY (`id`);

--
-- A tábla indexei `sub_category`
--
ALTER TABLE `sub_category`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`),
  ADD KEY `sub_to_category` (`category_id`);

--
-- A tábla indexei `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `username` (`username`),
  ADD KEY `user_to_role` (`role_id`);

--
-- A tábla indexei `week_day`
--
ALTER TABLE `week_day`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- A kiírt táblák AUTO_INCREMENT értéke
--

--
-- AUTO_INCREMENT a táblához `brand`
--
ALTER TABLE `brand`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT a táblához `category`
--
ALTER TABLE `category`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT a táblához `contract`
--
ALTER TABLE `contract`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

--
-- AUTO_INCREMENT a táblához `contract_day`
--
ALTER TABLE `contract_day`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=63;

--
-- AUTO_INCREMENT a táblához `contract_type`
--
ALTER TABLE `contract_type`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT a táblához `ingredient`
--
ALTER TABLE `ingredient`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT a táblához `open_day`
--
ALTER TABLE `open_day`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=22;

--
-- AUTO_INCREMENT a táblához `open_hour`
--
ALTER TABLE `open_hour`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT a táblához `product`
--
ALTER TABLE `product`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT a táblához `product_ingredient`
--
ALTER TABLE `product_ingredient`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT a táblához `product_type`
--
ALTER TABLE `product_type`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT a táblához `role`
--
ALTER TABLE `role`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT a táblához `session`
--
ALTER TABLE `session`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `store_type`
--
ALTER TABLE `store_type`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT a táblához `storing_condition`
--
ALTER TABLE `storing_condition`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT a táblához `sub_category`
--
ALTER TABLE `sub_category`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT a táblához `user`
--
ALTER TABLE `user`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=18;

--
-- AUTO_INCREMENT a táblához `week_day`
--
ALTER TABLE `week_day`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- Megkötések a kiírt táblákhoz
--

--
-- Megkötések a táblához `contract`
--
ALTER TABLE `contract`
  ADD CONSTRAINT `contract_to_type` FOREIGN KEY (`contract_type_id`) REFERENCES `contract_type` (`id`),
  ADD CONSTRAINT `contract_to_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

--
-- Megkötések a táblához `contract_day`
--
ALTER TABLE `contract_day`
  ADD CONSTRAINT `day_to_contract` FOREIGN KEY (`contract_id`) REFERENCES `contract` (`id`),
  ADD CONSTRAINT `day_to_week_day` FOREIGN KEY (`week_day_id`) REFERENCES `week_day` (`id`);

--
-- Megkötések a táblához `open_day`
--
ALTER TABLE `open_day`
  ADD CONSTRAINT `open_to_day` FOREIGN KEY (`week_day_id`) REFERENCES `week_day` (`id`),
  ADD CONSTRAINT `open_to_open_hour` FOREIGN KEY (`open_hour_id`) REFERENCES `open_hour` (`id`);

--
-- Megkötések a táblához `product`
--
ALTER TABLE `product`
  ADD CONSTRAINT `product_to_brand` FOREIGN KEY (`brand_id`) REFERENCES `brand` (`id`),
  ADD CONSTRAINT `product_to_type` FOREIGN KEY (`type_id`) REFERENCES `product_type` (`id`);

--
-- Megkötések a táblához `product_ingredient`
--
ALTER TABLE `product_ingredient`
  ADD CONSTRAINT `con_to_ingredient` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredient` (`id`),
  ADD CONSTRAINT `product_ingredient_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`);

--
-- Megkötések a táblához `product_type`
--
ALTER TABLE `product_type`
  ADD CONSTRAINT `type_to_condition` FOREIGN KEY (`storing_condition_id`) REFERENCES `storing_condition` (`id`),
  ADD CONSTRAINT `type_to_sub` FOREIGN KEY (`sub_id`) REFERENCES `sub_category` (`id`);

--
-- Megkötések a táblához `session`
--
ALTER TABLE `session`
  ADD CONSTRAINT `session_to_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

--
-- Megkötések a táblához `store_detail`
--
ALTER TABLE `store_detail`
  ADD CONSTRAINT `store_to_type` FOREIGN KEY (`store_type_id`) REFERENCES `store_type` (`id`);

--
-- Megkötések a táblához `sub_category`
--
ALTER TABLE `sub_category`
  ADD CONSTRAINT `sub_to_category` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`);

--
-- Megkötések a táblához `user`
--
ALTER TABLE `user`
  ADD CONSTRAINT `user_to_role` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
