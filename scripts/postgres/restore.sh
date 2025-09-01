#!/bin/bash

# === Настройки ===
DUMP_FILE="${1:-/backups/latest.sql.gz}"

# === Проверка переменных ===
if [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_DB" ] || [ -z "$POSTGRES_PORT" ] || [ -z "$POSTGRES_PASSWORD" ]; then
    echo "❌ Не все переменные окружения заданы: POSTGRES_USER, POSTGRES_DB, POSTGRES_PORT, POSTGRES_PASSWORD" | tee -a "$LOG_FILE"
    exit 1
fi

# === Проверка существования файла ===
if [ ! -f "$DUMP_FILE" ]; then
    echo "❌ Файл бэкапа не найден: $DUMP_FILE"
    exit 1
fi

export PGPASSWORD="$POSTGRES_PASSWORD"

echo "🔄 Начинаем восстановление базы '$POSTGRES_DB' из: $DUMP_FILE"

echo "🔌 Завершаем активные подключения к базе '$POSTGRES_DB'..."
psql -U "$POSTGRES_USER" -d postgres -c "
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = '$POSTGRES_DB'
  AND pid <> pg_backend_pid();
" > /dev/null 2>&1

echo "🗑️ Удаляем старую базу данных..."
psql -U "$POSTGRES_USER" -d postgres -c "DROP DATABASE IF EXISTS $POSTGRES_DB;" > /dev/null 2>&1

echo "🆕 Создаём новую базу данных..."
psql -U "$POSTGRES_USER" -d postgres -c "CREATE DATABASE $POSTGRES_DB OWNER $POSTGRES_USER;" > /dev/null 2>&1

echo "📥 Распаковываем и восстанавливаем данные..."
gunzip -c "$DUMP_FILE" | psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" > /dev/null 2>&1

if [ $? -ne 0 ]; then
 echo "❌ Ошибка при восстановлении данных из $DUMP_FILE!"
 exit 1
fi

echo "✅ Данные восстановлены. Проверяем целостность..."

# Количество таблиц в схеме public
TABLE_COUNT=$(psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -t -c "SELECT count(*) FROM pg_tables WHERE schemaname = 'public';" 2>/dev/null | xargs)

if [ -z "$TABLE_COUNT" ] || ! [[ "$TABLE_COUNT" =~ ^[0-9]+$ ]]; then
 echo "❌ Не удалось определить количество таблиц. Возможна ошибка подключения или повреждение БД."
 exit 1
fi

echo "📊 В схеме public обнаружено таблиц: $TABLE_COUNT"

if [ -z "$TABLE_COUNT" ] || [[ $TABLE_COUNT -eq 0 ]]; then
 echo "❌ Восстановленная база пуста — возможно, дамп был повреждён или пуст."
 exit 1
fi

# Дополнительная проверка: простой запрос
HEALTH_CHECK=$(psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -t -c "SELECT 1;" 2>/dev/null | xargs)

if [ -z "$HEALTH_CHECK" ] || [[ "$HEALTH_CHECK" != "1" ]]; then
 echo "❌ База не прошла health-check: SELECT 1 вернул пустой результат."
 exit 1
fi

# === Финал ===
echo "✅ Восстановление успешно завершено и проверено:"
echo "   - Файл: $DUMP_FILE"
echo "   - База: $POSTGRES_DB"
echo "   - Таблиц в public: $TABLE_COUNT"
echo "   - Health-check: пройден"
