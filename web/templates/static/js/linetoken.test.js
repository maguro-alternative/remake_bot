const { JSDOM } = require('jsdom');

const dom = new JSDOM('<!doctype html><html><body><form id="form"><input id="id" type="text" name="name" value="pas"/></form></body></html>');

global.document = dom.window.document;
global.window = dom.window;

global.fetch = jest.fn();
container = global.document.getElementById('form');

const { createJsonData } = require('./linetoken');

describe('fetchLineTokenData', () => {
    beforeAll(() => {
        const guildIdInput = document.createElement('input');
        guildIdInput.id = 'guildId';
        guildIdInput.type = 'text';
        guildIdInput.name = 'guildId';
        guildIdInput.value = '111';
        document.body.appendChild(guildIdInput);

        const lineNotifyTokenInput = document.createElement('input');
        lineNotifyTokenInput.id = 'lineNotifyToken';
        lineNotifyTokenInput.type = 'password';
        lineNotifyTokenInput.name = 'lineNotifyToken';
        lineNotifyTokenInput.value = 'lineNotifyToken';
        document.body.appendChild(lineNotifyTokenInput);

        const lineBotTokenInput = document.createElement('input');
        lineBotTokenInput.id = 'lineBotToken';
        lineBotTokenInput.type = 'password';
        lineBotTokenInput.name = 'lineBotToken';
        lineBotTokenInput.value = 'lineBotToken';
        document.body.appendChild(lineBotTokenInput);

        const lineBotSecretInput = document.createElement('input');
        lineBotSecretInput.id = 'lineBotSecret';
        lineBotSecretInput.type = 'password';
        lineBotSecretInput.name = 'lineBotSecret';
        lineBotSecretInput.value = 'lineBotSecret';
        document.body.appendChild(lineBotSecretInput);

        const lineGroupIDInput = document.createElement('input');
        lineGroupIDInput.id = 'lineGroupId';
        lineGroupIDInput.type = 'password';
        lineGroupIDInput.name = 'lineGroupId';
        lineGroupIDInput.value = 'lineGroupId';
        document.body.appendChild(lineGroupIDInput);

        const lineClientIDInput = document.createElement('input');
        lineClientIDInput.id = 'lineClientId';
        lineClientIDInput.type = 'password';
        lineClientIDInput.name = 'lineClientId';
        lineClientIDInput.value = 'lineClientId';
        document.body.appendChild(lineClientIDInput);

        const lineClientSecretInput = document.createElement('input');
        lineClientSecretInput.id = 'lineClientSecret';
        lineClientSecretInput.type = 'password';
        lineClientSecretInput.name = 'lineClientSecret';
        lineClientSecretInput.value = 'lineClientSecret';
        document.body.appendChild(lineClientSecretInput);

        const debugModeInput = document.createElement('input');
        debugModeInput.id = 'debugMode';
        debugModeInput.type = 'checkbox';
        debugModeInput.name = 'debugMode';
        debugModeInput.value = 'on';
        debugModeInput.checked = true;
        document.body.appendChild(debugModeInput);
    });

    afterEach(() => {
        jest.restoreAllMocks();
    });

    it('formからjsonに変換できること', async () => {
        const formData = new FormData();
        formData.append('guildId', '111');
        formData.append('lineNotifyToken', 'line_notify_token');
        formData.append('lineBotToken', 'line_bot_token');
        formData.append('lineBotSecret', 'line_bot_secret');
        formData.append('lineGroupId', 'line_group_id');
        formData.append('lineClientId', 'line_client_id');
        formData.append('lineClientSecret', 'line_client_secret');
        formData.append('debugMode', 'on');

        const jsonData = await createJsonData(formData)

        // jsonDataをオブジェクトにパース
        const parsedJsonData = JSON.parse(jsonData);

        // 期待値のオブジェクトを定義
        const expectedObject = {
            guildId: "111",
            lineNotifyToken: "line_notify_token",
            lineBotToken: "line_bot_token",
            lineBotSecret: "line_bot_secret",
            lineGroupId: "line_group_id",
            lineClientId: "line_client_id",
            lineClientSecret: "line_client_secret",
            debugMode: true,
            lineNotifyTokenDelete: false,
            lineBotTokenDelete: false,
            lineBotSecretDelete: false,
            lineGroupIdDelete: false,
            lineClientIdDelete: false,
            lineClientSecretDelete: false
        };

        // パースしたオブジェクトと期待値のオブジェクトを比較
        expect(parsedJsonData).toEqual(expectedObject);
    })
});

describe('fetchLineTokenData', () => {
    beforeAll(() => {
        const guildIdInput = document.createElement('input');
        guildIdInput.id = 'guildId';
        guildIdInput.type = 'text';
        guildIdInput.name = 'guildId';
        guildIdInput.value = '111';
        document.body.appendChild(guildIdInput);

        const debugModeInput = document.createElement('input');
        debugModeInput.id = 'lineBotTokenDelete';
        debugModeInput.type = 'checkbox';
        debugModeInput.name = 'lineBotTokenDelete';
        debugModeInput.value = 'on';
        debugModeInput.checked = true;
        document.body.appendChild(debugModeInput);
    });

    it('formからjsonに変換できること(deleteフラグ)', async () => {
        const formData = new FormData();
        formData.append('guildId', '111');
        formData.append('lineBotTokenDelete', 'on')

        const jsonData = await createJsonData(formData)

        // jsonDataをオブジェクトにパース
        const parsedJsonData = JSON.parse(jsonData);

        // 期待値のオブジェクトを定義
        const expectedObject = {
            guildId: "111",
            lineBotTokenDelete: true,
            debugMode: false,
            lineNotifyTokenDelete: false,
            lineBotSecretDelete: false,
            lineGroupIdDelete: false,
            lineClientIdDelete: false,
            lineClientSecretDelete: false
        };

        // パースしたオブジェクトと期待値のオブジェクトを比較
        expect(parsedJsonData).toEqual(expectedObject);
    })
});
