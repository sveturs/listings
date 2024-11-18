import os

def collect_files_content(file_paths, output_file):
    """
    Собирает содержимое указанных файлов в один текстовый файл.

    :param file_paths: Список путей к файлам для анализа.
    :param output_file: Путь к файлу, куда записать собранные данные.
    """
    try:
        with open(output_file, 'w', encoding='utf-8') as output:
            for file_path in file_paths:
                if os.path.exists(file_path):
                    output.write(f"--- Содержимое файла: {file_path} ---\n\n")
                    with open(file_path, 'r', encoding='utf-8') as input_file:
                        output.write(input_file.read())
                    output.write("\n\n")
                else:
                    output.write(f"--- Файл не найден: {file_path} ---\n\n")
        print(f"Содержимое файлов успешно собрано в '{output_file}'.")
    except Exception as e:
        print(f"Ошибка при сборе содержимого файлов: {e}")

# Пример использования
file_paths = [
    "backend/main.go",
    "frontend/hostel-frontend/src/components/RoomList.js",
    "frontend/hostel-frontend/src/components/AddRoom.js",
    "frontend/hostel-frontend/src/api/axios.js",
    "frontend/hostel-frontend/src/components/AddUser.js",
    "frontend/hostel-frontend/src/components/AdminPanel.js",
    "frontend/hostel-frontend/src/components/RoomList.js",
    "frontend/hostel-frontend/src/pages/AddRoomPage.js",
    "frontend/hostel-frontend/src/pages/AddUserPage.js",
    "frontend/hostel-frontend/src/pages/HomePage.js",
    "frontend/hostel-frontend/src/App.css",
    "frontend/hostel-frontend/src/App.js",
    "frontend/hostel-frontend/src/index.js",
    "deploy/docker-compose.yml",
    "backend/.env",
    "docker-compose.yml"
]
output_file = "collected_files_content.txt"

collect_files_content(file_paths, output_file)

