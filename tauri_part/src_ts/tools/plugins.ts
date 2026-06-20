// live2dPlugin.ts

export const eventBus = new Map<string, (...args: any[]) => void>();

// -> 发送 当前表情
eventBus.set(
  "send_cur_emotion",
  function (modelManager: any, wsclient: any, emotionName: string) {
    modelManager.setExpression(emotionName);
    wsclient.SendJson({
      name: "client_set_cur_emo",
      args: { emotion: modelManager.API_GetCurEmo() },
    });
  },
);

// 发送 所有表情,当前表情
eventBus.set("send_all_emotion", function (modelManager: any, wsclient: any) {
  const emotion = modelManager.getAllExp();
  wsclient.SendJson({
    name: "client_set_all_emo",
    args: { emotion: emotion },
  });
  wsclient.SendJson({
    name: "client_set_cur_emo",
    args: { emotion: modelManager.API_GetCurEmo() },
  });
});

// 接受 音量序列并处理
eventBus.set("add_volume_list", function (modelManager: any, volist: number[]) {
  modelManager.mouthAddData(volist);
});

// 接受 模型变换
eventBus.set(
  "trans_live2d_model",
  function (modelManager: any, app: any, comm: any) {
    modelManager.api_modelLayoyt(
      app,
      modelManager.model,
      comm.args["scale"],
      comm.args["transx"],
      comm.args["transy"],
    );
  },
);
