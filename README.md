# Пакетный менеджер (pm)

CLI-инструмент на Go для упаковки файлов в ZIP, загрузки на сервер по SFTP, скачивания и распаковки с учётом версий.

## Требования
- Go 1.21+
- SSH-сервер (например, OpenSSH) с пользователем и паролем
- Зависимости (в `go.mod`):
    - `golang.org/x/crypto v0.17.0`
    - `github.com/pkg/sftp v1.13.6`
    - `gopkg.in/yaml.v2 v2.4.0`

## Установка
```bash
git clone https://github.com/yourusername/pm.git
cd pm
go build -o pm ./cmd/pm
```

## Проверка работоспособности

### 1. Настройка
1. Настройте локальный SSH-сервер:
   ```bash
   sudo apt-get install openssh-server
   sudo systemctl start ssh
   sudo mkdir /packages
   sudo chown testuser:testuser /packages
   ```
   Пользователь: `testuser`, пароль: `testpass`.
2. Создайте тестовые файлы:
   ```bash
   mkdir -p archive_this1 archive_this2
   echo "Test 1" > archive_this1/file1.txt
   echo "Test 2" > archive_this2/file2.txt
   echo "Temp" > archive_this2/temp.tmp
   ```
3. Создайте `packet.json`:
   ```json
   {
       "name": "packet-1",
       "ver": "1.10",
       "targets": [
           "./archive_this1/*.txt",
           {"path": "./archive_this2/*", "exclude": "*.tmp"}
       ]
   }
   ```
4. Создайте `packages.json`:
   ```json
   {
       "packages": [
           {"name": "packet-1", "ver": ">=1.10"}
       ]
   }
   ```

### 2. Команда `create`
```bash
./pm create packet.json -ssh-host=localhost:22 -ssh-user=testuser -ssh-pass=testpass
```

**Ожидаемый вывод**:
```
Пакет packet-1-1.10.zip загружен на /packages/packet-1-1.10.zip
```

**Проверка**:
```bash
sftp testuser@localhost
cd /packages
ls
```
Ожидается: `packet-1-1.10.zip`.

### 3. Команда `update`
```bash
./pm update packages.json -ssh-host=localhost:22 -ssh-user=testuser -ssh-pass=testpass
```

**Ожидаемый вывод**:
```
Обновлено packet-1 до 1.10
```

**Проверка**:
```bash
ls archive_this1 archive_this2
```
Ожидается: `archive_this1/file1.txt`, `archive_this2/file2.txt`.

## Флаги
- `-ssh-host`: Хост:порт SSH (default: `localhost:22`)
- `-ssh-user`: Пользователь SSH (default: `user`)
- `-ssh-pass`: Пароль SSH
- `-remote-dir`: Удалённая директория (default: `/packages`)

## Ограничения
- Парольная SSH-аутентификация (небезопасно, для продакшена нужны ключи).
- Используется `panic` вместо обработки ошибок.
- Нет юнит-тестов.
- Версии только в формате `major.minor`.