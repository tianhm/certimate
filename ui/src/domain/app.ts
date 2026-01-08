import i18next from "i18next";

export const APP_VERSION = "v" + (__APP_VERSION__ || "0.0.0-dev").replace(/^v/, "");

export const APP_REPO_URL = "https://github.com/certimate-go/certimate";

export const APP_DOWNLOAD_URL = APP_REPO_URL + "/releases";

const APP_DOCUMENT_URLBASE = "https://docs.certimate.me";
export let APP_DOCUMENT_URL = APP_DOCUMENT_URLBASE;

i18next.on("languageChanged", (language) => {
  if (language.startsWith("en")) {
    APP_DOCUMENT_URL = APP_DOCUMENT_URLBASE + "/en-US";
  } else if (language.startsWith("zh")) {
    APP_DOCUMENT_URL = APP_DOCUMENT_URLBASE + "/zh-CN";
  } else {
    APP_DOCUMENT_URL = APP_DOCUMENT_URLBASE;
  }
});
