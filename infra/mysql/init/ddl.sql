CREATE DATABASE IF NOT EXISTS `goTechDojoDB` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `goTechDojoDB`;

DROP TABLE IF EXISTS Users CASCADE;
DROP TABLE IF EXISTS Collections CASCADE;
DROP TABLE IF EXISTS Scores CASCADE;
DROP TABLE IF EXISTS User_Collections CASCADE;

-- Users Table
CREATE TABLE Users (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    coins INT DEFAULT 0,
    high_score INT DEFAULT 0
);

-- Collections Table
CREATE TABLE Collections (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    rarity INT NOT NULL,
    weight INT NOT NULL
);

-- Scores Table
CREATE TABLE Scores (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36),
    value INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- UserCollections Table
CREATE TABLE User_Collections (
    user_id CHAR(36),
    collection_id CHAR(36),
    UNIQUE(user_id, collection_id),
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (collection_id) REFERENCES Collections(id)
);