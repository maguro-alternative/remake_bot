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
        guildIdInput.id = 'guild_id';
        guildIdInput.type = 'text';
        guildIdInput.name = 'guild_id';
        guildIdInput.value = '111';
        document.body.appendChild(guildIdInput);

        const lineNotifyTokenInput = document.createElement('input');
        lineNotifyTokenInput.id = 'line_notify_token';
        lineNotifyTokenInput.type = 'password';
        lineNotifyTokenInput.name = 'line_notify_token';
        lineNotifyTokenInput.value = 'line_notify_token';
        document.body.appendChild(lineNotifyTokenInput);

        const lineBotTokenInput = document.createElement('input');
        lineBotTokenInput.id = 'line_bot_token';
        lineBotTokenInput.type = 'password';
        lineBotTokenInput.name = 'line_bot_token';
        lineBotTokenInput.value = 'line_bot_token';
        document.body.appendChild(lineBotTokenInput);

        const lineBotSecretInput = document.createElement('input');
        lineBotSecretInput.id = 'line_bot_secret';
        lineBotSecretInput.type = 'password';
        lineBotSecretInput.name = 'line_bot_secret';
        lineBotSecretInput.value = 'line_bot_secret';
        document.body.appendChild(lineBotSecretInput);

        const lineGroupIDInput = document.createElement('input');
        lineGroupIDInput.id = 'line_group_id';
        lineGroupIDInput.type = 'password';
        lineGroupIDInput.name = 'line_group_id';
        lineGroupIDInput.value = 'line_group_id';
        document.body.appendChild(lineGroupIDInput);

        const lineClientIDInput = document.createElement('input');
        lineClientIDInput.id = 'line_client_id';
        lineClientIDInput.type = 'password';
        lineClientIDInput.name = 'line_client_id';
        lineClientIDInput.value = 'line_client_id';
        document.body.appendChild(lineClientIDInput);

        const lineClientSecretInput = document.createElement('input');
        lineClientSecretInput.id = 'line_client_secret';
        lineClientSecretInput.type = 'password';
        lineClientSecretInput.name = 'line_client_secret';
        lineClientSecretInput.value = 'line_client_secret';
        document.body.appendChild(lineClientSecretInput);

        const debugModeInput = document.createElement('input');
        debugModeInput.id = 'debug_mode';
        debugModeInput.type = 'checkbox';
        debugModeInput.name = 'debug_mode';
        debugModeInput.value = 'on';
        debugModeInput.checked = true;
        document.body.appendChild(debugModeInput);
    });

    afterEach(() => {
        jest.restoreAllMocks();
    });

    it('formからjsonに変換できること', async () => {
        const formData = new FormData();
        formData.append('guild_id', '111');
        formData.append('line_notify_token', 'line_notify_token');
        formData.append('line_bot_token', 'line_bot_token');
        formData.append('line_bot_secret', 'line_bot_secret');
        formData.append('line_group_id', 'line_group_id');
        formData.append('line_client_id', 'line_client_id');
        formData.append('line_client_secret', 'line_client_secret');
        formData.append('debug_mode', 'on');

        const jsonData = await createJsonData(formData)

        expect(jsonData).toEqual('{"guild_id":"111","line_notify_token":"line_notify_token","line_bot_token":"line_bot_token","line_bot_secret":"line_bot_secret","line_group_id":"line_group_id","line_client_id":"line_client_id","line_client_secret":"line_client_secret","debug_mode":true}');
    })
});