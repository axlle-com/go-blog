-- Создаем вторую базу данных
CREATE DATABASE cms_test;

-- Подключаемся к базе данных cms_test
\connect cms_test;

-- Создаем пользователя для второй базы данных
CREATE USER postgres_test WITH PASSWORD 'secret_test';

-- Даем пользователю права на базу данных
GRANT ALL PRIVILEGES ON DATABASE cms_test TO postgres_test;

-- Даем пользователю права на схему public
GRANT ALL PRIVILEGES ON SCHEMA public TO postgres_test;

-- Если нужно, можно также предоставить права на все существующие таблицы в схеме public
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres_test;

-- Для обеспечения, что пользователь имеет права на все таблицы, создаваемые в будущем
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO postgres_test;
