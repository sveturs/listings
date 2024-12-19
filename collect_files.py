import os

def collect_content(paths, output_file):
    """
    Собирает содержимое указанных файлов и файлов в папках в один текстовый файл.

    :param paths: Список путей к файлам и папкам.
    :param output_file: Путь к файлу, куда записать собранные данные.
    """
    try:
        with open(output_file, 'w', encoding='utf-8') as output:
            for path in paths:
                if os.path.isfile(path):
                    # Если это файл, записываем его содержимое
                    try:
                        output.write(f"--- Содержимое файла: {path} ---\n\n")
                        with open(path, 'r', encoding='utf-8') as input_file:
                            output.write(input_file.read())
                        output.write("\n\n")
                    except Exception as e:
                        output.write(f"--- Ошибка чтения файла: {path} ({e}) ---\n\n")
                elif os.path.isdir(path):
                    # Если это папка, обходим её содержимое
                    for root, _, files in os.walk(path):
                        for file in files:
                            file_path = os.path.join(root, file)
                            try:
                                output.write(f"--- Содержимое файла: {file_path} ---\n\n")
                                with open(file_path, 'r', encoding='utf-8') as input_file:
                                    output.write(input_file.read())
                                output.write("\n\n")
                            except Exception as e:
                                output.write(f"--- Ошибка чтения файла: {file_path} ({e}) ---\n\n")
                else:
                    # Если путь не существует
                    output.write(f"--- Путь не найден: {path} ---\n\n")
        print(f"Содержимое успешно собрано в '{output_file}'.")
    except Exception as e:
        print(f"Ошибка при сборе содержимого: {e}")

paths = [
    "backend/main.go",
    "backend/Dockerfile",
    "backend/.env",
    "backend/.env.local",  
    "backend/.dockerignore",      
#    "backend/auth/auth.go",
    "backend/database/db.go",
    "backend/migrations",
    "frontend/hostel-frontend/src/index.js",
    "frontend/hostel-frontend/src/App.js",
    "frontend/hostel-frontend/src/api/axios.js",
    "frontend/hostel-frontend/src/components/accommodation/RoomList.js",            
    "frontend/hostel-frontend/src/components/accommodation/BookingsList.js",
    "deploy/docker-compose.yml",
#    "frontend/hostel-frontend/public",
#    "frontend/hostel-frontend/src",
#    "frontend/hostel-frontend/src/api/axios.js",        
#    "frontend/hostel-frontend/src/components",
#    "frontend/hostel-frontend/src/contexts",
#    "frontend/hostel-frontend/src/pages",
    "frontend/hostel-frontend/.env",
    "frontend/hostel-frontend/.env.local",
    "frontend/hostel-frontend/.gitignore",
#    "frontend/hostel-frontend/package.json",        
    ".gitignore",
    "deploy.sh",
    "docker-compose.prod.yml",
    "docker-compose.yml",
#    "init-ssl.sh",
    "nginx.conf",
#    "package.json"

]
output_file = "collected_content.txt"

collect_content(paths, output_file)

