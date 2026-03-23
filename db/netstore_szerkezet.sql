-- phpMyAdmin SQL Dump
-- version 5.1.1deb5ubuntu1
-- https://www.phpmyadmin.net/
--
-- Gép: localhost:3306
-- Létrehozás ideje: 2026. Már 23. 16:55
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

-- --------------------------------------------------------

--
-- Tábla szerkezet ehhez a táblához `category`
--

CREATE TABLE `category` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

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
  `other_properties` varchar(255) NOT NULL,
  `type_id` int(11) NOT NULL,
  `brand_id` int(11) NOT NULL,
  `deleted_at` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_hungarian_ci;

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
('-', '-', 0, 1);

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
(1, '-', '-', 'admin', '$2a$12$EKzHWgDilgOwYFjMm7A7GO1uUe0LA5ZFKwTsqlTVBTPeqdQUQpILO', 1, '-', '-', 1, NULL);

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
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `category`
--
ALTER TABLE `category`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `contract`
--
ALTER TABLE `contract`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `contract_day`
--
ALTER TABLE `contract_day`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `contract_type`
--
ALTER TABLE `contract_type`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `ingredient`
--
ALTER TABLE `ingredient`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `open_day`
--
ALTER TABLE `open_day`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `open_hour`
--
ALTER TABLE `open_hour`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `product`
--
ALTER TABLE `product`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `product_ingredient`
--
ALTER TABLE `product_ingredient`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `product_type`
--
ALTER TABLE `product_type`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

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
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT a táblához `sub_category`
--
ALTER TABLE `sub_category`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

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
