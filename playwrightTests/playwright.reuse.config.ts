import base from './playwright.config';
export default {
  ...base,
  webServer: base.webServer ? { ...base.webServer, reuseExistingServer: true } : undefined,
};
