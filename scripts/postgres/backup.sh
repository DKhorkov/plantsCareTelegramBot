#!/bin/bash

# === Настройки ===
BACKUP_DIR="/backups"
LOG_FILE="$BACKUP_DIR/backup.log"
MIN_FREE_SPACE_MB=500  # Минимум 500 МБ свободно
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")

# === Проверка переменных ===
if [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_DB" ] || [ -z "$POSTGRES_PORT" ] || [ -z "$POSTGRES_PASSWORD" ]; then
    echo "❌ Не все переменные окружения заданы: POSTGRES_USER, POSTGRES_DB, POSTGRES_PORT, POSTGRES_PASSWORD" | tee -a "$LOG_FILE"
    exit 1
fi

# === Проверка свободного места ===
FREE_SPACE_KB=$(df --output=avail "$BACKUP_DIR" | tail -n1)
if [ -z "$FREE_SPACE_KB" ]; then
    echo "❌ Не удалось определить свободное место в $BACKUP_DIR" | tee -a "$LOG_FILE"
    exit 1
fi

FREE_SPACE_MB=$((FREE_SPACE_KB / 1024))

if [ $FREE_SPACE_MB -lt $MIN_FREE_SPACE_MB ]; then
    echo "❌ Недостаточно места: $FREE_SPACE_MB MB. Требуется минимум $MIN_FREE_SPACE_MB MB." | tee -a "$LOG_FILE"
    exit 1
fi

# === Подготовка ===
mkdir -p "$BACKUP_DIR"
export PGPASSWORD="$POSTGRES_PASSWORD"
BACKUP_FILE="$BACKUP_DIR/${POSTGRES_DB}_$TIMESTAMP.sql.gz"

echo "$(date): Начинаем резервное копирование базы '$POSTGRES_DB'..." >> "$LOG_FILE"

# === Создание сжатого бэкапа ===
echo "📦 Создаём и сжимаем дамп базы: $POSTGRES_DB"
pg_dump \
    --username="$POSTGRES_USER" \
    --host=localhost \
    --port="$POSTGRES_PORT" \
    --no-password \
    --verbose \
    "$POSTGRES_DB" | gzip > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Бэкап успешно создан: $BACKUP_FILE"
    echo "$(date): Успешно создан бэкап: $BACKUP_FILE" >> "$LOG_FILE"

    # Симлинк на последний бэкап
    ln -sf "$(basename "$BACKUP_FILE")" "$BACKUP_DIR/latest.sql.gz"
    echo "🔗 Актуальный бэкап: $BACKUP_DIR/latest.sql.gz"
else
    echo "❌ Ошибка при создании бэкапа!"
    echo "$(date): ОШИБКА при создании бэкапа" >> "$LOG_FILE"
    exit 1
fi

