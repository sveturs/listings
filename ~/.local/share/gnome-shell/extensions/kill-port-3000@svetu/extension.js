const { St, GLib, Gio, Clutter } = imports.gi;
const Main = imports.ui.main;
const PanelMenu = imports.ui.panelMenu;
const PopupMenu = imports.ui.popupMenu;

let button;

function init() {
}

function enable() {
    // Создаем кнопку в панели
    button = new PanelMenu.Button(0.0, 'Kill Port 3000', false);
    
    // Создаем иконку с символом остановки
    let icon = new St.Icon({
        icon_name: 'process-stop-symbolic',
        style_class: 'system-status-icon'
    });
    
    button.add_child(icon);
    
    // Добавляем обработчик клика
    button.connect('button-press-event', () => {
        killPort3000();
    });
    
    // Добавляем кнопку в панель рядом с системными индикаторами
    Main.panel.addToStatusArea('kill-port-3000', button, 1, 'right');
}

function disable() {
    if (button) {
        button.destroy();
        button = null;
    }
}

function killPort3000() {
    try {
        // Сначала находим PID процесса на порту 3000
        let [success, stdout, stderr] = GLib.spawn_command_line_sync('netstat -tlnp 2>/dev/null | grep :3000');
        
        if (success && stdout) {
            let output = stdout.toString();
            // Парсим вывод netstat для получения PID
            // Формат: tcp  0  0 0.0.0.0:3000  0.0.0.0:*  LISTEN  1861063/main
            let match = output.match(/(\d+)\/\w+/);
            
            if (match && match[1]) {
                let pid = match[1];
                
                // Убиваем процесс
                let [killSuccess] = GLib.spawn_command_line_sync(`kill -9 ${pid}`);
                
                if (killSuccess) {
                    Main.notify('Kill Port 3000', `Процесс ${pid} на порту 3000 успешно остановлен`);
                } else {
                    Main.notify('Kill Port 3000', 'Ошибка при остановке процесса');
                }
            } else {
                Main.notify('Kill Port 3000', 'Процесс на порту 3000 не найден');
            }
        } else {
            Main.notify('Kill Port 3000', 'Процесс на порту 3000 не найден');
        }
    } catch (e) {
        Main.notify('Kill Port 3000', `Ошибка: ${e.message}`);
    }
}