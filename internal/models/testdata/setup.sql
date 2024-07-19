-- Creación de la tabla 'snippets'
CREATE TABLE snippets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- Creación de índice para la columna 'created' de la tabla 'snippets'
CREATE INDEX idx_snippets_created ON snippets(created);

-- Creación de la tabla 'users'
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

-- Creación de una restricción de unicidad para la columna 'email' de la tabla 'users'
CREATE UNIQUE INDEX users_uc_email ON users(email);

-- Inserción de datos en la tabla 'users'
INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 10:00:00'
);
