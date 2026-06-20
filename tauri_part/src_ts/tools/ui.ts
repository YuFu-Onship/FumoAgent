import { getCurrentWindow } from "@tauri-apps/api/window";

const appWindow = getCurrentWindow();

export class MenuManager {
  private menuContainer: HTMLElement;
  private wrapper: HTMLElement;

  constructor(wrapper: HTMLElement) {
    this.wrapper = wrapper;
    // 初始化容器
    this.menuContainer = document.createElement("div");
    this.menuContainer.className =
      "live2d-menus live2d-transition-all live2d-opacity-0 live2d-hidden";
    this.menuContainer.style.setProperty("--live2d-duration", "300ms");

    this.wrapper.appendChild(this.menuContainer);
    this.setupEvents();
  }

  // 核心：为按钮添加动态点击波纹效果的方法
  private createRipple(button: HTMLElement, e: MouseEvent) {
    // 1. 创建波纹 span 元素
    const ripple = document.createElement("span");
    ripple.className = "live2d-ripple-span";

    // 2. 计算点击位置和波纹大小
    const rect = button.getBoundingClientRect();
    const size = Math.max(rect.width, rect.height);

    ripple.style.width = `${size}px`;
    ripple.style.height = `${size}px`;
    ripple.style.left = `${e.clientX - rect.left - size / 2}px`;
    ripple.style.top = `${e.clientY - rect.top - size / 2}px`;

    // 3. 移除旧的波纹，防止连续狂点产生堆积
    const oldRipple = button.querySelector(".live2d-ripple-span");
    if (oldRipple) oldRipple.remove();

    // 4. 插入并定时销毁
    button.appendChild(ripple);
    setTimeout(() => ripple.remove(), 1500); // 与 CSS 动画时间保持一致
  }

  // 内部辅助方法：创建按钮
  // private createButton(
  //   title: string,
  //   svgIcon: string,
  //   onClick: () => void,
  // ): HTMLElement {
  //   const button = document.createElement("div");
  //   button.className = "live2d-menu-item live2d-flex-center";
  //   button.title = title;
  //   button.innerHTML = svgIcon;
  //   button.addEventListener("click", (e) => {
  //     e.stopPropagation();
  //     this.createRipple(button, e); // 触发波纹
  //     onClick();
  //   });
  //   return button;
  // }

  // 公开方法：添加自定义按钮
  public addButton(title: string, svgIcon: string, onClick: () => void) {
    const btn = document.createElement("div");
    btn.className = "live2d-menu-item live2d-flex-center";
    btn.title = title;
    btn.innerHTML = svgIcon;
    btn.addEventListener("click", (e) => {
      e.stopPropagation(); // 防止触发窗口拖拽或其它事件
      this.createRipple(btn, e); // 触发波纹

      // 让点击逻辑稍微延迟一丁点（可选），给波纹动画展现的时间
      setTimeout(() => onClick(), 50);
    });
    this.menuContainer.appendChild(btn);
    return this; // 允许链式调用
  }

  // 动画逻辑
  // private async fadeIn() {
  //   this.menuContainer.classList.remove("live2d-hidden");
  //   await new Promise((resolve) => requestAnimationFrame(resolve));
  //   this.menuContainer.classList.remove("live2d-opacity-0");
  //   this.menuContainer.classList.add("live2d-opacity-1");
  // }

  // private async fadeOut() {
  //   this.menuContainer.classList.remove("live2d-opacity-1");
  //   this.menuContainer.classList.add("live2d-opacity-0");
  //   await new Promise((resolve) => setTimeout(resolve, 300));
  //   this.menuContainer.classList.add("live2d-hidden");
  // }

  // 绑定事件
  private setupEvents() {
    this.wrapper.addEventListener("mouseover", () => this.toggle(true));
    this.wrapper.addEventListener("mouseleave", () => this.toggle(false));
  }
  private async toggle(show: boolean) {
    if (show) {
      this.menuContainer.classList.remove("live2d-hidden");
      await new Promise((r) => requestAnimationFrame(r));
      this.menuContainer.classList.replace(
        "live2d-opacity-0",
        "live2d-opacity-1",
      );
    } else {
      this.menuContainer.classList.replace(
        "live2d-opacity-1",
        "live2d-opacity-0",
      );
      await new Promise((r) => setTimeout(r, 300));
      this.menuContainer.classList.add("live2d-hidden");
    }
  }

  // 预设：添加关闭按钮
  public addCloseButton() {
    return this.addButton("关闭", `<svg>...</svg>`, () => appWindow.close());
  }
}
export const ICONS = {
  SETTINGS: `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#e3e3e3"><path d="m370-80-16-128q-13-5-24.5-12T307-235l-119 50L78-375l103-78q-1-7-1-13.5v-27q0-6.5 1-13.5L78-585l110-190 119 50q11-8 23-15t24-12l16-128h220l16 128q13 5 24.5 12t22.5 15l119-50 110 190-103 78q1 7 1 13.5v27q0 6.5-2 13.5l103 78-110 190-118-50q-11 8-23 15t-24 12L590-80H370Zm70-80h79l14-106q31-8 57.5-23.5T639-327l99 41 39-68-86-65q5-14 7-29.5t2-31.5q0-16-2-31.5t-7-29.5l86-65-39-68-99 42q-22-23-48.5-38.5T533-694l-13-106h-79l-14 106q-31 8-57.5 23.5T321-633l-99-41-39 68 86 64q-5 15-7 30t-2 32q0 16 2 31t7 30l-86 65 39 68 99-42q22 23 48.5 38.5T427-266l13 106Zm42-180q58 0 99-41t41-99q0-58-41-99t-99-41q-59 0-99.5 41T342-480q0 58 40.5 99t99.5 41Zm-2-140Z"/></svg>`,
  INFO: `
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
        <path d="M256 512c141.4 0 256-114.6 256-256S397.4 0 256 0S0 114.6 0 256S114.6 512 256 512zM216 336h24V272H216c-13.3 0-24-10.7-24-24s10.7-24 24-24h48c13.3 0 24 10.7 24 24v88h8c13.3 0 24 10.7 24 24s-10.7 24-24 24H216c-13.3 0-24-10.7-24-24s10.7-24 24-24zm40-144c-17.7 0-32-14.3-32-32s14.3-32 32-32s32 14.3 32 32s-14.3 32-32 32z"/>
      </svg>
    `,
  NEXT: `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#e3e3e3"><path d="M504-480 320-664l56-56 240 240-240 240-56-56 184-184Z"/></svg>`,
  CLOSE: `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#e3e3e3"><path d="m256-200-56-56 224-224-224-224 56-56 224 224 224-224 56 56-224 224 224 224-56 56-224-224-224 224Z"/></svg>`,
  WIFI: `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#e3e3e3"><path d="M409-149q-29-29-29-71t29-71q29-29 71-29t71 29q29 29 29 71t-29 71q-29 29-71 29t-71-29ZM254-346l-84-86q59-59 138.5-93.5T480-560q92 0 171.5 35T790-430l-84 84q-44-44-102-69t-124-25q-66 0-124 25t-102 69ZM84-516 0-600q92-94 215-147t265-53q142 0 265 53t215 147l-84 84q-77-77-178.5-120.5T480-680q-116 0-217.5 43.5T84-516Z"/></svg>`,
  WIFIERROR: `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#e3e3e3"><path d="M790-56 414-434q-47 11-87.5 33T254-346l-84-86q32-32 69-56t79-42l-90-90q-41 21-76.5 46.5T84-516L0-602q32-32 66.5-57.5T140-708l-84-84 56-56 736 736-58 56Zm-381-93.5Q380-179 380-220q0-42 29-71t71-29q42 0 71 29t29 71q0 41-29 70.5T480-120q-42 0-71-29.5ZM716-358l-29-29-29-29-144-144q81 8 151.5 41T790-432l-74 74Zm160-158q-77-77-178.5-120.5T480-680q-21 0-40.5 1.5T400-674L298-776q44-12 89.5-18t92.5-6q142 0 265 53t215 145l-84 86Z"/></svg>`,
  CHAT: `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#e3e3e3"><path d="M240-400h320v-80H240v80Zm0-120h480v-80H240v80Zm0-120h480v-80H240v80ZM80-80v-720q0-33 23.5-56.5T160-880h640q33 0 56.5 23.5T880-800v480q0 33-23.5 56.5T800-240H240L80-80Zm126-240h594v-480H160v525l46-45Zm-46 0v-480 480Z"/></svg>`,
};
