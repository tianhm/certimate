const appBaseUrl = new URL(import.meta.env.BASE_URL, document.baseURI);
const urlSchemePattern = /^[a-z][a-z\d+.-]*:/i;

// PocketBase 会把自身的 /api 路径拼接到这里，因此基础路径必须以斜杠结尾。
export const APP_BASE_PATH = appBaseUrl.pathname;

export const withBasePath = (path: string) => {
  // 完整 URL 和协议相对 URL 不属于应用内资源，保持调用方传入的地址不变。
  if (urlSchemePattern.test(path) || path.startsWith("//")) {
    return path;
  }

  const url = new URL(path.replace(/^\/+/, ""), appBaseUrl);
  return `${url.pathname}${url.search}${url.hash}`;
};
