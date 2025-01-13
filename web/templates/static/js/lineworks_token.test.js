const { JSDOM } = require('jsdom');

const dom = new JSDOM(`
  <!doctype html>
  <html>
      <body>
          <form id="form">
            <label for="lineWorksClientID">LINE Works Client ID</label>
            <input id="lineWorksClientID" type="password" name="lineWorksClientID">
            <br/>
            <label for="lineWorksClientIdDelete">LINE Works Client ID削除</label>
            <input id="lineWorksClientIdDelete" type="checkbox" name="lineWorksClientIdDelete">
            <br/>
            <label for="lineWorksClientSecret">LINE Works Client Secret</label>
            <input id="lineWorksClientSecret" type="password" name="lineWorksClientSecret">
            <br/>
            <label for="lineWorksClientSecretDelete">LINE Works Client Secret削除</label>
            <input id="lineWorksClientSecretDelete" type="checkbox" name="lineWorksClientSecretDelete">
            <br/>
            <label for="lineWorksServiceAccount">LINE Works Service Account</label>
            <input id="lineWorksServiceAccount" type="password" name="lineWorksServiceAccount">
            <br/>
            <label for="lineWorksServiceAccountDelete">LINE Works Service Account削除</label>
            <input id="lineWorksServiceAccountDelete" type="checkbox" name="lineWorksServiceAccountDelete">
            <br/>
            <label for="lineWorksPrivateKey">LINE Works Private Key</label>
            <input id="lineWorksPrivateKey" type="password" name="lineWorksPrivateKey">
            <br/>
            <label for="lineWorksPrivateKeyDelete">LINE Works Private Key削除</label>
            <input id="lineWorksPrivateKeyDelete" type="checkbox" name="lineWorksPrivateKeyDelete">
            <br/>
            <label for="lineWorksDomainId">LINE Works Domain ID</label>
            <input id="lineWorksDomainId" type="password" name="lineWorksDomainId">
            <br/>
            <label for="lineWorksDomainIdDelete">LINE Works Domain ID削除</label>
            <input id="lineWorksDomainIdDelete" type="checkbox" name="lineWorksDomainIdDelete">
            <br/>
            <label for="lineWorksAdminId">LINE Works Admin ID</label>
            <input id="lineWorksAdminId" type="password" name="lineWorksAdminId">
            <br/>
            <label for="lineWorksAdminIdDelete">LINE Works Admin ID削除</label>
            <input id="lineWorksAdminIdDelete" type="checkbox" name="lineWorksAdminIdDelete">
            <br/>
            <label for="lineWorksBotID">LINE Works Bot ID</label>
            <input id="lineWorksBotID" type="password" name="lineWorksBotID">
            <br/>
            <label for="lineWorksBotIdDelete">LINE Works Bot ID削除</label>
            <input id="lineWorksBotIdDelete" type="checkbox" name="lineWorksBotIdDelete">
            <br/>
            <label for="lineWorksBotSecret">LINE Works Bot Secret</label>
            <input id="lineWorksBotSecret" type="password" name="lineWorksBotSecret">
            <br/>
            <label for="lineWorksBotSecretDelete">LINE Works Bot Secret削除</label>
            <input id="lineWorksBotSecretDelete" type="checkbox" name="lineWorksBotSecretDelete">
            <br/>
            <label for="lineWorksGroupID">LINE Works Group ID</label>
            <input id="lineWorksGroupID" type="password" name="lineWorksGroupID">
            <br/>
            <label for="lineWorksGroupIdDelete">LINE Works Group ID削除</label>
            <input id="lineWorksGroupIdDelete" type="checkbox" name="lineWorksGroupIdDelete">
            <br/>
            <label for="defaultChannel">Discordに送信するデフォルトのチャンネル</label>
            <select id="defaultChannel" name="defaultChannelId" >
            </select>
            <br/>
            <label for="debugMode">デバッグモード</label>
            <input id="debugMode" type="checkbox" name="debugMode">
            <br/>
            <input type="submit" value="送信">
          </form>
      </body>
  </html>
`);
global.document = dom.window.document;
global.window = dom.window;

global.fetch = jest.fn();
container = global.document.getElementById('form');

const localStorageMock = (() => {
  let store = {};

  return {
      getItem(key) {
          return store[key] || null;
      },
      setItem(key, value) {
          store[key] = value.toString();
      },
      removeItem(key) {
          delete store[key];
      },
      clear() {
          store = {};
      },
  };
})();

Object.defineProperty(window, 'sessionStorage', {
  value: localStorageMock,
});

const { createJsonData } = require('./lineworks_token');

describe('createJsonData function', () => {
  it('createJsonData function should return JSON string', async() => {
    // モックデータの準備
    const formElements = [
      { name: 'lineWorksClientID', value: 'testClientID' },
      { name: 'lineWorksClientSecret', value: 'testClientSecret' },
      { name: 'lineWorksServiceAccount', value: 'testServiceAccount' },
      { name: 'lineWorksPrivateKey', value: 'testPrivateKey' },
      { name: 'lineWorksDomainId', value: 'testDomainId' },
      { name: 'lineWorksAdminId', value: 'testAdminId' },
      { name: 'lineWorksBotID', value: 'testBotID' },
      { name: 'lineWorksBotSecret', value: 'testBotSecret' },
      { name: 'lineWorksGroupID', value: 'testGroupID' },
      { name: 'lineWorksClientIdDelete', value: 'on' },
      { name: 'lineWorksClientSecretDelete', value: 'on' },
      { name: 'lineWorksServiceAccountDelete', value: 'on' },
      { name: 'lineWorksPrivateKeyDelete', value: 'on' },
      { name: 'lineWorksDomainIdDelete', value: 'on' },
      { name: 'lineWorksAdminIdDelete', value: 'on' },
      { name: 'lineWorksBotIdDelete', value: 'on' },
      { name: 'lineWorksBotSecretDelete', value: 'on' },
      { name: 'lineWorksGroupIdDelete', value: 'on' },
      { name: 'defaultChannel', value: 'testChannel' },
      { name: 'debugMode', value: 'on' },
    ];
    const form = new FormData();
    formElements.forEach((element) => {
      form.append(element.name, element.value);
    });

    console.log(form);

    // 実行
    const jsonData = await createJsonData(form);

    // 期待値の作成
    const expectedJsonData = {
      lineWorksClientID: 'testClientID',
      lineWorksClientSecret: 'testClientSecret',
      lineWorksServiceAccount: 'testServiceAccount',
      lineWorksPrivateKey: 'testPrivateKey',
      lineWorksDomainId: 'testDomainId',
      lineWorksAdminId: 'testAdminId',
      lineWorksBotID: 'testBotID',
      lineWorksBotSecret: 'testBotSecret',
      lineWorksGroupID: 'testGroupID',
      lineWorksClientIdDelete: true,
      lineWorksClientSecretDelete: true,
      lineWorksServiceAccountDelete: true,
      lineWorksPrivateKeyDelete: true,
      lineWorksDomainIdDelete: true,
      lineWorksAdminIdDelete: true,
      lineWorksBotIdDelete: true,
      lineWorksBotSecretDelete: true,
      lineWorksGroupIdDelete: true,
      defaultChannel: 'testChannel',
      debugMode: true,
    };

    // 検証
    expect(JSON.parse(jsonData)).toEqual(expectedJsonData);
  });
});
