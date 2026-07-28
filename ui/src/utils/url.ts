const appBaseUrl = new URL(import.meta.env.BASE_URL, document.baseURI);

// PocketBase 会把自身的 /api 路径拼接到这里，因此基础路径必须以斜杠结尾。
export const APP_BASE_PATH = appBaseUrl.pathname;

export const resolveAppPath = (path: string) => {
  return new URL(path.replace(/^\/+/, ""), appBaseUrl).pathname;
};
