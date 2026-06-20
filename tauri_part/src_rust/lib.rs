use device_query::{DeviceQuery, DeviceState};
use serde::Serialize;

// 原有的 greet 函数保留...
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}
#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        // 注册指令
        .invoke_handler(tauri::generate_handler![greet, get_global_mouse_pos])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

#[derive(Serialize)]
struct MousePos {
    x: i32,
    y: i32,
}

// 定义一个获取全局坐标的命令
#[tauri::command]
fn get_global_mouse_pos() -> MousePos {
    let device_state = DeviceState::new();
    let mouse = device_state.get_mouse();
    MousePos {
        x: mouse.coords.0,
        y: mouse.coords.1,
    }
}
