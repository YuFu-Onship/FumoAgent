import * as PIXI from "pixi.js";
import { Live2DModel } from "pixi-live2d-display";
import { invoke } from "@tauri-apps/api/core";

// 定义不同表情下的参数结构
interface ExpressionEyeMap {
  [expressionName: string]: {
    [paramId: string]: number;
  };
}

export class ModelManager {
  private model: any = null;
  private app: PIXI.Application;

  // 表情
  private modelExpNames: string[] = [];
  private currentExpIdx: number = -1;
  private curEmoName: string = "Default";
  private emoTimer: number = 0;
  private emoReset: boolean = false;

  // 嘴部运动
  private mouthQueue: number[] = [];
  private isTickerInit = false;

  // 眨眼相关
  private blinkTimer = 0;
  private nextBlinkTime = 4.0;
  private blinkState: "idle" | "closing" | "closed" | "opening" = "idle";
  private stateTimer = 0;

  // 眨眼
  public expressionEyeData: ExpressionEyeMap = {};
  private eyeParamIds = ["ParamEyeLOpen", "ParamEyeROpen"];

  // 尺寸
  public baseS: number = 1.0;
  public baseX: number = 0;
  public baseY: number = 0;

  //摸头检测
  private patTracker = new HeadPatTracker();
  public API_OnPat: (() => void) | null = null;
  private isPatCooldown: boolean = false;

  // 视线检测
  private idleGazeTracker = new IdleGazeTracker({
    idleThresholdTime: 2.0, // 2秒不动就发呆
    mouseMoveThreshold: 10, // 10px 误差
    randomGazeInterval: 3.0, // 每3秒换个地方看
  });

  // 旋转
  private rotation_cur_week: number = 0;
  private rotation_tar_week: number = 0;
  // private is_rotation: boolean = false;

  constructor(app: PIXI.Application) {
    this.app = app;
  }

  // 异步加载模型
  async LoadModel(modelPath: string) {
    if (this.model) {
      this.app.stage.removeChild(this.model);
      this.model.destroy();
    }
    try {
      // 初始化应用
      this.model = await Live2DModel.from(modelPath);
      this.model.eventMode = "none";
      this.app.stage.addChild(this.model);
      this.modelLayoyt(this.app, this.model);

      // 获取到当前模型的表情
      this.modelExpNames = this.getExpression(this.model);
      this.setExpression("Default");
    } catch (error) {
      console.error(error);
    }
  }

  // 更新 --------------------------------------------------------------------
  public update(appWindow: any, dt: number) {
    if (!this.model || !this.model.internalModel) return;
    this.modelSpeak();
    this.eyeBlink(dt);
    this.gazeTrackerGlobal(appWindow, dt);
    this.updateRotate(dt);

    if (this.emoReset) {
      this.emoTimer += dt;
      if (this.emoTimer >= 10) {
        this.emoTimer = 0;
        this.emoReset = false;
        this.curEmoName = "Default";
        this.currentExpIdx = 0;
        this.model.internalModel.motionManager.expressionManager.resetExpression();
      }
    }
  }

  // 销毁 ---------------------------------------------------------------------
  public destory(app: PIXI.Application) {
    if (this.model) {
      app.stage.removeChild(this.model);
      this.model.destroy();
      this.model = null;
    }
  }

  // 变换 ---------------------------------------------------------------------
  // 计算模型缩放与位置
  public modelLayoyt(app: any, model: any) {
    const ratio = model.width / model.height;
    model.height = app.screen.height * 0.8 * this.baseS;
    model.width = model.height * ratio * this.baseS;
    model.x = app.screen.width / 2;
    model.y = app.screen.height / 2;
    model.anchor.set(0.5, 0.5);
  }

  public api_initModelLayout(
    app: PIXI.Application,
    model: any,
    scale: number,
    x: number,
    y: number,
  ) {
    this.baseS = scale;
    this.baseX = x;
    this.baseY = y;
    this.modelLayoyt(app, model);
  }

  public api_modelLayoyt(
    app: any,
    model: any,
    scale: number = this.baseS,
    transX: number = 0,
    transY: number = 0,
  ) {
    if (!model) return;
    const ratio = model.width / model.height;
    const targetHeight = app.screen.height * 0.8 * scale;
    model.height = targetHeight;
    model.width = targetHeight * ratio;

    model.anchor.set(0.5, 0.5);
    model.x = app.screen.width / 2 + transX;
    model.y = app.screen.height / 2 + transY;
  }

  // 旋转更新
  public updateRotate(dt: number) {
    if (this.rotation_tar_week == 0) return;
    this.model.rotation += 4 * dt;
    if (this.model.rotation >= Math.PI * 2) {
      this.model.rotation = 0;
      this.rotation_cur_week += 1;
    }
    if (this.rotation_cur_week >= this.rotation_tar_week) {
      this.rotation_cur_week = 0;
      this.rotation_tar_week = 0;
      this.model.rotation = 0;
    }
  }
  public api_setRotationWeek(value: number) {
    this.rotation_tar_week = value;
  }

  // 跟随鼠标 -----------------------------------------------------------------
  // 模型视线追踪
  public gazeTracker(relX: number, relY: number) {
    this.model.focus(relX, relY);
  }
  public async gazeTrackerGlobal(appWindow: any, dt: number) {
    try {
      if (!this.model || this.model.destroyed) return;
      // 视线追踪 ----------------------------
      const globalPos = await invoke<{ x: number; y: number }>(
        "get_global_mouse_pos",
      );
      const windowPos = await appWindow.innerPosition();
      const windowSize = await appWindow.innerSize();

      const relX = globalPos.x - windowPos.x;
      const relY = globalPos.y - windowPos.y;

      const gazePos = this.idleGazeTracker.update(
        relX,
        relY,
        windowSize.width,
        windowSize.height,
        dt,
      );
      this.gazeTracker(gazePos.x, gazePos.y);

      // 摸头 --------------------------------
      const isPatted = this.patTracker.update(
        relX,
        relY,
        windowSize.width,
        windowSize.height,
      );

      if (isPatted) {
        if (this.isPatCooldown) {
          return;
        }

        if (typeof this.API_OnPat === "function") {
          this.API_OnPat();
          this.isPatCooldown = true;
          setTimeout(() => {
            this.isPatCooldown = false;
          }, 10000);
        }
      }
    } catch (err) {}
  }

  public api_shakingHead() {
    this.idleGazeTracker.triggerShakeHead();
  }

  // 嘴巴闭合 ---------------------------------------------------------------------
  public mouthInit(app: any, fps: number) {
    if (this.isTickerInit) return;
    app.ticker.maxFPS = fps;
    app.ticker.add(() => {
      if (!this.model || !this.model.internalModel) return;

      // 如果队列里有数据，就取出一个（先进先出）
      if (this.mouthQueue.length > 0) {
        const currentVolume = this.mouthQueue.shift()!; // 弹出数组第一个元素

        // 映射到你的模型：后端传来 0~1，模型需要 0~10
        const mouthOpenValue = currentVolume * 10;

        this.model.internalModel.coreModel.setParameterValueById(
          "ParamMouthOpenY",
          mouthOpenValue,
        );
      } else {
        this.model.internalModel.coreModel.setParameterValueById(
          "ParamMouthOpenY",
          0,
        );
      }
    });

    this.isTickerInit = true;
  }
  public mouthAddData(volumes: number[]) {
    this.mouthQueue = [];
    this.mouthQueue.push(...volumes);
  }
  // 说话, 嘴巴动起来
  public modelSpeak() {
    if (this.mouthQueue.length > 0) {
      let mouthOpenValue = 0;
      if (this.mouthQueue.length > 0) {
        mouthOpenValue = this.mouthQueue.shift()!;
      }
      this.model.internalModel.coreModel.setParameterValueById(
        "ParamMouthOpenY",
        mouthOpenValue * 10,
      );
    }
  }

  // 眨眼 ----------------------------------------------------------------------
  public eyeBlink(dt: number) {
    if (!this.model || !this.model.internalModel) return;
    if (!this.expressionEyeData?.["Default"]?.["ParamEyeLOpen"]) return;
    const coreModel = this.model.internalModel.coreModel;

    this.blinkTimer += dt;

    if (this.blinkState === "idle") {
      if (this.blinkTimer >= this.nextBlinkTime) {
        this.blinkState = "closing";
        this.blinkTimer = 0;
      }
    } else if (this.blinkState === "closing") {
      if (this.blinkTimer >= 0.05) {
        this.blinkState = "closed";
        this.blinkTimer = 0;
      }
    } else if (this.blinkState === "closed") {
      if (this.blinkTimer >= 0.05) {
        this.blinkState = "opening";
        this.blinkTimer = 0;
      }
    } else if (this.blinkState === "opening") {
      if (this.blinkTimer >= 0.05) {
        this.blinkState = "idle";
        this.blinkTimer = 0;
      }
    }

    const baseOpenValue = this.expressionEyeData["Default"]["ParamEyeLOpen"];
    let maxOpenValue: number;
    if (this.curEmoName == "Default") {
      maxOpenValue = baseOpenValue;
    } else {
      maxOpenValue =
        baseOpenValue +
        this.expressionEyeData[this.curEmoName]["ParamEyeLOpen"];
    }
    if (maxOpenValue > baseOpenValue) {
      return;
    }

    let eyeOpenValue = baseOpenValue; // 默认睁眼

    switch (this.blinkState) {
      case "idle":
        eyeOpenValue = baseOpenValue;
        break;
      case "closing":
        eyeOpenValue = maxOpenValue - this.stateTimer / 0.1;
        break;
      case "closed":
        eyeOpenValue = 0;
        break;
      case "opening":
        eyeOpenValue = this.stateTimer / 0.05;
        break;
    }

    eyeOpenValue = Math.max(0, Math.min(maxOpenValue, eyeOpenValue));
    coreModel.setParameterValueById("ParamEyeROpen", eyeOpenValue);
    coreModel.setParameterValueById("ParamEyeLOpen", eyeOpenValue);
    coreModel.setParameterValueByIndex(1, 0, 0.3);
  }

  // 表情相关 ----------------------------------------------------------------------
  // 获得表情列表
  public getExpression(model: any) {
    // model.internalModel.blinkState = false;
    // return model.internalModel.settings.expressions.map((exp: any) => exp.Name);
    return (
      model?.internalModel?.settings?.expressions?.map(
        (exp: any) => exp.Name,
      ) ?? []
    );
  }

  // 设置表情
  public setExpression(name: string) {
    if (name !== "Default" && !this.modelExpNames.includes(name)) {
      return;
    }
    this.curEmoName = name;
    if (name === "Default") {
      this.model.internalModel?.motionManager?.expressionManager?.resetExpression();
    } else {
      this.model.expression(name);
      this.emoReset = true;
      this.emoTimer = 0;
    }
  }

  // 设置下一个表情
  public setNextExpression() {
    if (!this.model) {
      return;
    }
    const names = this.modelExpNames;
    if (names.length === 0) {
      return;
    }
    this.currentExpIdx = (this.currentExpIdx + 1) % (names.length + 1);
    let name: string;
    if (this.currentExpIdx === names.length) {
      name = "Default";
    } else {
      name = names[this.currentExpIdx];
    }
    this.setExpression(name);
    this.curEmoName = name;
  }

  // 获取到当前的表情名称
  public getCurExp(): string {
    if (!this.model || !this.model.internalModel) {
      return "Default";
    }
    if (this.modelExpNames.length == 0) {
      return "Default";
    }
    // 获取表达式管理器
    const expManager = this.model.internalModel.motionManager.expressionManager;
    const currentIdx = expManager.currentExpressionIndex;

    if (
      currentIdx !== undefined &&
      currentIdx >= 0 &&
      currentIdx < this.modelExpNames.length
    ) {
      return this.modelExpNames[currentIdx];
    }
    return "Default";
  }

  // 获取到表情列表
  public getAllExp(): string {
    let l = "";
    for (let i = 0; i < this.modelExpNames.length; i++) {
      l += this.modelExpNames[i] + ",";
    }
    l += "Default";
    return l;
  }

  // 获取到当前的表情
  public getCurrentExpression() {
    return;
  }

  // 文件io --------------------------------------------------

  /**
   * 加载并解析所有表情 JSON，提取眼部参数目标值
   * @param modelUrl model3.json 的完整路径
   * @param model3Json 已经加载好的 model3.json 对象数据
   */
  public async loadExpressionEyeConfig(floderPath: string) {
    if (!this.model || !this.model.internalModel) return;
    const model3Json = this.model.internalModel.settings.json;
    const coreModel = this.model.internalModel.coreModel;
    const expressions = model3Json.FileReferences?.Expressions;
    if (!expressions || expressions.length === 0) {
      console.log("该模型没有配置表情文件。");
      return;
    }
    // 1. 获取模型所在的文件夹基础路径
    const baseFolder = floderPath + "/";
    console.log(baseFolder);
    // 2. 遍历所有表情
    for (const expr of expressions) {
      const exprName = expr.Name; // 例如 "Sorpresa" 或 "Happy"
      const exprFile = expr.File; // 例如 "Cirno 01.exp3.json"
      const exprUrl = `${baseFolder}${exprFile}`;
      // 初始化该表情的字典对象
      this.expressionEyeData[exprName] = {};
      try {
        // 发送网络请求读取具体的 .exp3.json
        const response = await fetch(exprUrl);
        const exprJson = await response.json();
        // 3. 检查并提取我们关心的眼部参数
        this.eyeParamIds.forEach((paramId) => {
          // 在 Parameters 数组中查找对应的 ID
          const targetParam = exprJson.Parameters?.find(
            (p: any) => p.Id === paramId,
          );
          if (targetParam !== undefined && targetParam.Value !== undefined) {
            // 如果 JSON 里有定义，直接读取（例如 0.25）
            this.expressionEyeData[exprName][paramId] = targetParam.Value;
          } else {
            // 如果 JSON 里没有定义，使用你提到的官方底层核心 API 获取默认值兜底
            // const defaultValue = coreModel.getParameterValueById(paramId);
            this.expressionEyeData[exprName][paramId] = 0;
          }
        });
      } catch (error) {
        console.error(`解析表情文件失败: ${exprUrl}`, error);
        // 发生错误时，全部用模型默认值兜底，防止程序崩溃
        this.expressionEyeData[exprName] = {};
        this.eyeParamIds.forEach((paramId) => {
          this.expressionEyeData[exprName][paramId] =
            coreModel.getParameterValueById(paramId);
        });
      }
    }
    this.expressionEyeData["Default"] = {};
    this.eyeParamIds.forEach((paramId) => {
      const defaultValue = coreModel.getParameterValueById(paramId);
      this.expressionEyeData["Default"][paramId] = defaultValue;
    });
    // 打印最终生成的字典验证结果
    console.log("成功建立表情眼部参数字典:", this.expressionEyeData);
  }

  // api -----------------------------------------------
  public API_GetCurEmo(): string {
    return this.curEmoName;
  }
}

// 摸头检测 --------------------------------------------
class HeadPatTracker {
  private lastX: number | null = null;
  private currentDirection: number = 0; // 0: 未知, 1: 向右, -1: 向左
  private currentDistance: number = 0; // 当前单向移动的累计距离
  private strokeCount: number = 0; // 成功来回的次数

  // 配置项
  private readonly minDistance = 100; // 单向最小距离
  private readonly requiredStrokes = 8; // 需要来回的次数
  private readonly extendX = 100; // 左右延伸距离

  /**
   * 每一帧调用此方法
   * @param relX 鼠标相对于窗口的 X 坐标
   * @param relY 鼠标相对于窗口的 Y 坐标
   * @param windowWidth 窗口的当前宽度
   * @param windowHeight 窗口的当前高度
   * @returns 是否触发摸头
   */
  public update(
    relX: number,
    relY: number,
    windowWidth: number,
    windowHeight: number,
  ): boolean {
    // 1. 区域判定：窗口上半部分 (0 ~ windowHeight/2)，左右各延伸 100px
    const inXRange =
      relX >= -this.extendX && relX <= windowWidth + this.extendX;
    const inYRange = relY >= 0 && relY <= windowHeight / 2;

    if (!inXRange || !inYRange) {
      this.reset(); // 离开区域，清空状态
      return false;
    }

    // 第一次记录坐标
    if (this.lastX === null) {
      this.lastX = relX;
      return false;
    }

    const deltaX = relX - this.lastX;
    this.lastX = relX; // 更新上一帧坐标

    if (deltaX === 0) return false; // 鼠标没动

    const dir = Math.sign(deltaX); // 1 表示向右, -1 表示向左

    // 2. 方向判定
    if (this.currentDirection === 0) {
      // 刚开始移动，初始化方向
      this.currentDirection = dir;
      this.currentDistance = Math.abs(deltaX);
    } else if (dir === this.currentDirection) {
      // 保持同方向移动，累加距离
      this.currentDistance += Math.abs(deltaX);
    } else {
      // 方向发生改变（来回拐弯了！）
      // 检查刚刚结束的那个方向，移动距离是否达标
      if (this.currentDistance >= this.minDistance) {
        this.strokeCount++;
        console.log(
          `[摸头中] 成功完成 1 次单向移动，当前计数: ${this.strokeCount}`,
        );

        // 3. 触发判定
        if (this.strokeCount >= this.requiredStrokes) {
          this.reset(); // 触发后重置
          return true; // 成功触发摸头
        }
      } else {
        // 如果拐弯了但是距离不够，重置之前的计数（或者只重置当前单次，看你对“连续性”的要求）
        // 这里采用严格模式：一旦中间有一次短距离调头，就重新算
        this.strokeCount = 0;
      }

      // 切换到新方向，并重新开始计算新方向的距离
      this.currentDirection = dir;
      this.currentDistance = Math.abs(deltaX);
    }

    return false;
  }

  public reset() {
    this.lastX = null;
    this.currentDirection = 0;
    this.currentDistance = 0;
    this.strokeCount = 0;
  }
}

// 视线注视与更改 ------------------------------------------------------
export interface GazeConfig {
  idleThresholdTime?: number; // 鼠标静止多少秒后开始随机看 (默认2秒)
  mouseMoveThreshold?: number; // 鼠标触发移动的阈值 (默认10px)
  randomGazeInterval?: number; // 每隔多少秒换一个随机注视点 (默认3秒)
  lerpSpeed?: number; // 眼睛/头部转动的平滑速度 (默认5)
}

// 定义状态机状态
export enum GazeState {
  FollowMouse, // 追随鼠标
  IdleWait, // 鼠标刚停下的过渡期
  RandomGaze, // 随机注视漫游
  ShakingHead, // 摇头状态
}

export class IdleGazeTracker {
  // 配置项
  private config: Required<GazeConfig>;

  // 内部状态机
  private currentState: GazeState = GazeState.FollowMouse;

  // 基础数据
  private lastMousePos = { x: 0, y: 0 };
  private mouseIdleTime = 0;
  private currentPos = { x: 0, y: 0 }; // 统一维护当前输出的平滑位置

  // 随机注视相关
  private randomGazeTimer = 0;
  private targetRandomX = 0;
  private targetRandomY = 0;

  // 摇头（ShakeHead）相关变量
  private shakeDuration = 0; // 当前摇头的总持续时间（秒）
  private shakeTimer = 0; // 摇头已消耗的时间
  private shakeSpeed = 0; // 摇头正弦波频率速度
  private shakeBaseY = 0; // 摇头时的固定Y轴高度

  constructor(config?: GazeConfig) {
    this.config = {
      idleThresholdTime: config?.idleThresholdTime ?? 2.0,
      mouseMoveThreshold: config?.mouseMoveThreshold ?? 10,
      randomGazeInterval: config?.randomGazeInterval ?? 3.0,
      lerpSpeed: config?.lerpSpeed ?? 5.0,
    };
  }

  /**
   * 外部触发接口：命令模型开始摇头
   * 摇头结束后会自动进入随机注视状态
   */
  public triggerShakeHead(): void {
    this.currentState = GazeState.ShakingHead;
    this.shakeTimer = 0;

    // 随机生成摇头持续时间 (5 ~ 10 秒)
    this.shakeDuration = 5 + Math.random() * 5;

    // 随机生成摇头速度（正弦波角速度，可以根据具体表现微调）
    this.shakeSpeed = 4 + Math.random() * 3;

    // 重置时间戳，确保每次摇头都从中间或特定振幅开始
    this.randomGazeTimer = 0;
  }

  /**
   * 核心更新方法
   * @param relX 当前鼠标相对窗口X
   * @param relY 当前鼠标相对窗口Y
   * @param width 窗口宽度
   * @param height 窗口高度
   * @param dt 增量时间 (秒)
   * @returns 最终模型应该注视的 {x, y} 坐标
   */
  public update(
    relX: number,
    relY: number,
    width: number,
    height: number,
    dt: number,
  ): { x: number; y: number } {
    // 1. 如果当前正在【摇头】，则完全无视鼠标移动，直到摇头自然结束
    if (this.currentState === GazeState.ShakingHead) {
      this.updateShakeHeadLogic(width, height, dt);

      // 关键：在摇头的过程中，我们需要在后台静默同步鼠标位置
      // 这样摇头一旦结束切回其他状态时，不会因为“积累了很久的鼠标位移”而产生画面瞬间跳变
      this.lastMousePos = { x: relX, y: relY };
      return this.currentPos;
    }

    // 2. 非摇头状态下，再进行普通的鼠标移动距离判断
    const dx = relX - this.lastMousePos.x;
    const dy = relY - this.lastMousePos.y;
    const distance = Math.sqrt(dx * dx + dy * dy);

    if (distance > this.config.mouseMoveThreshold) {
      // 鼠标移动了，切回跟随状态
      this.currentState = GazeState.FollowMouse;
      this.mouseIdleTime = 0;
      this.lastMousePos = { x: relX, y: relY };
    }

    // 3. 状态机逻辑分支处理（此时已被排除 ShakingHead 状态）
    switch (this.currentState) {
      case GazeState.FollowMouse:
        this.currentPos = { x: relX, y: relY };
        this.lastMousePos = { x: relX, y: relY };
        // 这一帧跟随完了，下一帧默认进入静止等待，除非下一帧鼠标继续动
        this.currentState = GazeState.IdleWait;
        break;

      case GazeState.IdleWait:
        this.mouseIdleTime += dt;
        this.currentPos = { x: this.lastMousePos.x, y: this.lastMousePos.y };

        if (this.mouseIdleTime >= this.config.idleThresholdTime) {
          this.currentState = GazeState.RandomGaze;
          this.mouseIdleTime = 0;
          this.randomGazeTimer = this.config.randomGazeInterval;
        }
        break;

      case GazeState.RandomGaze:
        this.updateRandomTarget(width, height, dt);
        break;
    }
    return this.currentPos;
  }

  // 内部方法：更新随机点并平滑插值
  private updateRandomTarget(width: number, height: number, dt: number) {
    this.randomGazeTimer += dt;

    if (
      this.randomGazeTimer >= this.config.randomGazeInterval ||
      (this.targetRandomX === 0 && this.targetRandomY === 0)
    ) {
      this.randomGazeTimer = 0;

      // 限制在屏幕中心 60% 区域内随机
      const paddingX = width * 0.2;
      const paddingY = height * 0.2;
      this.targetRandomX = paddingX + Math.random() * (width - paddingX * 2);
      this.targetRandomY = paddingY + Math.random() * (height - paddingY * 2);
    }

    // 平滑插值 (Lerp)
    this.lerpToTarget(this.targetRandomX, this.targetRandomY, dt);
  }

  // 内部方法：处理摇头逻辑
  private updateShakeHeadLogic(width: number, height: number, dt: number) {
    this.shakeTimer += dt;

    // 检查摇头是否结束
    if (this.shakeTimer >= this.shakeDuration) {
      // 摇头结束 -> 自动进入随机注视状态
      this.currentState = GazeState.RandomGaze;
      this.randomGazeTimer = this.config.randomGazeInterval; // 确保立刻计算新的随机注视点
      return;
    }

    // 摇头计算：高度固定在 0.25 左右位置 (这里取 25% 屏幕高度)
    this.shakeBaseY = height * 0.25;

    // 左右正弦运动：
    // Math.sin 的范围是 [-1, 1]
    // 映射到屏幕上：以屏幕中心 (width * 0.5) 为基准，向左右各摆动 30% 屏幕宽度 (width * 0.3)
    const targetX =
      width * 0.5 + Math.sin(this.shakeTimer * this.shakeSpeed) * (width * 0.3);
    const targetY = this.shakeBaseY;

    // 摇头的插值速度可以稍微快点或者保持一致，这里用原本的平滑速度
    this.lerpToTarget(targetX, targetY, dt);
  }

  // 公共平滑插值工具
  private lerpToTarget(targetX: number, targetY: number, dt: number) {
    const speed = this.config.lerpSpeed * dt;
    this.currentPos.x += (targetX - this.currentPos.x) * Math.min(speed, 1);
    this.currentPos.y += (targetY - this.currentPos.y) * Math.min(speed, 1);
  }
}
