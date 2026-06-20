import { getCurrentWindow } from "@tauri-apps/api/window";
import { MenuManager, ICONS } from "./tools/ui";
import { NetClient } from "./tools/netclient";
import { ModelManager } from "./tools/liveModel";
import { eventBus } from "./tools/plugins";

// 变量定义 ----------------------------------------------------
declare const PIXI: any;
// const { Live2DModel } = PIXI.live2d;
const appWindow = getCurrentWindow();
const RootPath: string = "http://localhost:9527";

let wsclient: any = null;
let app: any = null;

let slots: any[] = [];
let modelIndex = 0;
let isLoading: boolean = false;

interface Command {
  name?: string;
  args?: any;
}

// let client: WebSocketClient | null = null;

// 主函数 ------------------------------------------------------
async function init() {
  // 创建 PIXI 应用
  app = new PIXI.Application({
    view: document.getElementById("canvas") as HTMLCanvasElement,
    autoStart: true,
    backgroundAlpha: 0,
    width: 450,
    height: 600,
  });

  // ws客户端 ---------------------------------------------------------
  wsclient = new NetClient(methons);
  // 注意写法 =()=>{}
  wsclient.onOpen = () => {
    console.log("This is Client");
    wsclient.SendMessage("hello server");

    wsclient.SendJson({
      name: "client_get_live2d_model",
      args: {},
    });
  };
  wsclient.StartClient();

  // 模型控制器

  // 插件---------------------------------------------------------
  async function methons(comm: Command) {
    let f: any;
    switch (comm.name) {
      // 关闭软件
      case "server_set_close":
        app.destroy();
        break;

      // 切换表情
      case "server_set_emotion":
        slots[modelIndex].setExpression(comm.args["emotion"]);
        wsclient.SendJson({
          name: "client_set_cur_emo",
          args: { emotion: slots[modelIndex].API_GetCurEmo() },
        });
        break;

      // 获得所有表情
      case "server_get_emotions":
        f = eventBus.get("send_all_emotion");
        if (f) {
          f(slots[modelIndex], wsclient);
        }
        break;

      // 传入语音音量列表,模型适配嘴型
      case "server_set_volume_list":
        f = eventBus.get("add_volume_list");
        if (f) {
          f(slots[modelIndex], comm.args["volist"]);
        }
        break;

      // 模型变换
      case "server_set_live2d_trans":
        console.log(comm);
        f = eventBus.get("trans_live2d_model");
        if (f) {
          f(slots[modelIndex], app, comm);
        }
        break;

      // 模型旋转
      case "server_set_live2d_rotate":
        slots[modelIndex].api_setRotationWeek(comm.args["week"]);
        break;

      // 模型摇头
      case "server_set_live2d_shake":
        slots[modelIndex].api_shakingHead();
        break;

      // 设置/加载 模型初始化
      case "server_set_live2d_model":
        console.log("收到模型配置，开始加载...");
        if (isLoading) break;
        modelFolder = RootPath + "/live2d/" + comm.args["FolderPath"];
        modelModel3 = modelFolder + "/" + comm.args["JsonName"];
        console.log(modelModel3);

        isLoading = true;
        if (!slots[0]) {
          slots[0] = new ModelManager(app);
          await slots[0].LoadModel(modelModel3);
          await slots[0].loadExpressionEyeConfig(modelFolder);
          modelIndex = 0;

          if (slots[1]) {
            slots[1].destory(app);
            app.ticker.remove(slots[1].update);
            slots[1] = null;
          }
        } else {
          slots[1] = new ModelManager(app);
          await slots[1].LoadModel(modelModel3);
          await slots[1].loadExpressionEyeConfig(modelFolder);
          modelIndex = 1;

          if (slots[0]) {
            slots[0].destory(app);
            app.ticker.remove(slots[0].update);
            slots[0] = null;
          }
        }

        isLoading = false;
        wsclient.SendJson({
          name: "client_set_all_emo",
          args: { emotion: slots[modelIndex].getAllExp() },
        });
        wsclient.SendJson({
          name: "client_set_cur_emo",
          args: { emotion: slots[modelIndex].getCurExp() },
        });
        wsclient.SendJson({
          name: "client_get_live2d_trans",
          args: {},
        });
        slots[modelIndex].API_OnPat = () => {
          wsclient.SendJson({
            name: "client_set_onpat",
            args: {},
          });
        };

        break;

      default:
        break;
    }
  }

  // 初始化 模型 -----------------------------------------------------
  let modelFolder: string = "";
  let modelModel3: string = "";

  // 更新 ---------------------------------------------------------
  const fps: number = 40;
  const dt: number = 1 / fps;
  app.ticker.maxFPS = fps;
  app.ticker.add(() => {
    if (slots[modelIndex]) {
      slots[modelIndex].update(appWindow, dt);
    }
  });

  // 界面ui相关 ---------------------------------------------------
  // 获取或创建包装器元素, 创建一个"画框", 后续的内容在这个画框上绘制
  let wrapper = document.getElementById("app");
  if (!wrapper) {
    wrapper = document.createElement("div");
    wrapper.id = "app";
    Object.assign(wrapper.style, {
      position: "relative",
      width: "450px",
      height: "600px",
    });
    document.body.appendChild(wrapper);
  }

  // 使用 MenuManager 创建菜单
  const menu = new MenuManager(wrapper);
  menu
    .addButton("Close", ICONS.CLOSE, () => {
      wsclient.SendJson({
        name: "client_close",
        args: {},
      });
      appWindow.close();
    })
    .addButton("Next Expression", ICONS.NEXT, () => {
      slots[modelIndex].setNextExpression();
    })
    .addButton("Connect WS Server", ICONS.WIFIERROR, () => {
      wsclient.reConnect();
    })
    .addButton("Chat", ICONS.CHAT, () => {
      wsclient.SendJson({
        name: "client_set_chat",
        args: {},
      });
    });
  const wifiBtn = wrapper.querySelector(
    '.live2d-menu-item[title="Connect WS Server"]',
  );
  if (wifiBtn) {
    wifiBtn.classList.add("live2d-menu-wifi");
  }

  // 根据连接状态动态更新WiFi图标
  const updateWifiIcon = () => {
    if (wifiBtn) {
      wifiBtn.innerHTML = wsclient.IsConnect() ? ICONS.WIFI : ICONS.WIFIERROR;
    }
  };

  // 包装原有 onOpen，追加图标更新
  const originalOnOpen = wsclient.onOpen;
  wsclient.onOpen = () => {
    updateWifiIcon();
    originalOnOpen();
  };
  wsclient.onClose = () => updateWifiIcon();
  wsclient.onError = () => updateWifiIcon();

  // 初始化时同步一次图标状态
  updateWifiIcon();

  // 左键拖拽窗口（排除菜单按钮区域） ----------------------------------
  window.addEventListener("mousedown", (e) => {
    const target = e.target as HTMLElement;
    const isMenuButton = target.closest(".live2d-menu-item");
    if (e.buttons === 1 && !isMenuButton) {
      appWindow.startDragging();
    }
  });
}

// 启动应用
init().catch(console.error);
