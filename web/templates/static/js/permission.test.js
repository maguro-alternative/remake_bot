const { JSDOM } = require('jsdom');

const dom = new JSDOM('<!doctype html><html><body><form id="form"><input id="id" type="text" name="name" value="pas"/></form></body></html>');

global.document = dom.window.document;
global.window = dom.window;

global.fetch = jest.fn();
container = global.document.getElementById('form');

const { createJsonData } = require('./permission');

describe('fetchPermissionData', () => {
    beforeAll(() => {
        const testUser1Option = document.createElement('option');
        testUser1Option.name = 'testuser1';
        testUser1Option.value = '111';
        testUser1Option.selected = false;
        testUser1Option.defaultSelected = false;

        const testUser2Option = document.createElement('option');
        testUser2Option.name = 'testuser2';
        testUser2Option.value = '222';
        testUser2Option.selected = true;
        testUser2Option.defaultSelected = true;

        const testUser3Option = document.createElement('option');
        testUser3Option.name = 'testuser3';
        testUser3Option.value = '333';
        testUser3Option.selected = false;
        testUser3Option.defaultSelected = false;

        const testUser4Option = document.createElement('option');
        testUser4Option.name = 'testuser4';
        testUser4Option.value = '444';
        testUser4Option.selected = true;
        testUser4Option.defaultSelected = true;

        const testRole1Option = document.createElement('option');
        testRole1Option.name = 'testrole1';
        testRole1Option.value = '111';
        testRole1Option.selected = false;
        testRole1Option.defaultSelected = false;

        const testRole2Option = document.createElement('option');
        testRole2Option.name = 'testrole2';
        testRole2Option.value = '222';
        testRole2Option.selected = true;
        testRole2Option.defaultSelected = true;

        const testRole3Option = document.createElement('option');
        testRole3Option.name = 'testrole3';
        testRole3Option.value = '333';
        testRole3Option.selected = false;
        testRole3Option.defaultSelected = false;

        const testRole4Option = document.createElement('option');
        testRole4Option.name = 'testrole4';
        testRole4Option.value = '444';
        testRole4Option.selected = true;
        testRole4Option.defaultSelected = true;

        const guildIdInput = document.createElement('input');
        guildIdInput.id = 'guild_id';
        guildIdInput.type = 'text';
        guildIdInput.name = 'guild_id';
        guildIdInput.value = '111';
        document.forms['form'].appendChild(guildIdInput);

        const lineBotPermissionCodeInput = document.createElement('input');
        lineBotPermissionCodeInput.type = 'number';
        lineBotPermissionCodeInput.name = 'lineBotPermissionCode';
        lineBotPermissionCodeInput.value = 8;
        document.forms['form'].appendChild(lineBotPermissionCodeInput);

        const linePostDiscordChannelPermissionCodeInput = document.createElement('input');
        linePostDiscordChannelPermissionCodeInput.type = 'number';
        linePostDiscordChannelPermissionCodeInput.name = 'linePostDiscordChannelPermissionCode';
        linePostDiscordChannelPermissionCodeInput.value = 8;
        document.forms['form'].appendChild(linePostDiscordChannelPermissionCodeInput);

        const vcSignalPermissionCodeInput = document.createElement('input');
        vcSignalPermissionCodeInput.type = 'number';
        vcSignalPermissionCodeInput.name = 'vcSignalPermissionCode';
        vcSignalPermissionCodeInput.value = 8;
        document.forms['form'].appendChild(vcSignalPermissionCodeInput);

        const webhookPermissionCodeInput = document.createElement('input');
        webhookPermissionCodeInput.type = 'number';
        webhookPermissionCodeInput.name = 'webhookPermissionCode';
        webhookPermissionCodeInput.value = 8;
        document.forms['form'].appendChild(webhookPermissionCodeInput);

        const lineBotMemberPermissionIdSelect = document.createElement('select');
        lineBotMemberPermissionIdSelect.name = 'lineBotMemberPermissionId';
        lineBotMemberPermissionIdSelect.multiple = true;
        lineBotMemberPermissionIdSelect.appendChild(testUser1Option);
        lineBotMemberPermissionIdSelect.appendChild(testUser2Option);
        lineBotMemberPermissionIdSelect.appendChild(testUser3Option);
        lineBotMemberPermissionIdSelect.appendChild(testUser4Option);
        document.forms['form'].appendChild(lineBotMemberPermissionIdSelect);

        const linePostDiscordChannelMemberPermissionIdSelect = document.createElement('select');
        linePostDiscordChannelMemberPermissionIdSelect.name = 'linePostDiscordChannelMemberPermissionId';
        linePostDiscordChannelMemberPermissionIdSelect.multiple = true;
        linePostDiscordChannelMemberPermissionIdSelect.appendChild(testUser1Option);
        linePostDiscordChannelMemberPermissionIdSelect.appendChild(testUser2Option);
        linePostDiscordChannelMemberPermissionIdSelect.appendChild(testUser3Option);
        linePostDiscordChannelMemberPermissionIdSelect.appendChild(testUser4Option);
        document.forms['form'].appendChild(linePostDiscordChannelMemberPermissionIdSelect);

        const vcSignalMemberPermissionIdSelect = document.createElement('select');
        vcSignalMemberPermissionIdSelect.name = 'vcSignalMemberPermissionId';
        vcSignalMemberPermissionIdSelect.multiple = true;
        vcSignalMemberPermissionIdSelect.appendChild(testUser1Option);
        vcSignalMemberPermissionIdSelect.appendChild(testUser2Option);
        vcSignalMemberPermissionIdSelect.appendChild(testUser3Option);
        vcSignalMemberPermissionIdSelect.appendChild(testUser4Option);
        document.forms['form'].appendChild(vcSignalMemberPermissionIdSelect);

        const webhookMemberPermissionIdSelect = document.createElement('select');
        webhookMemberPermissionIdSelect.name = 'webhookMemberPermissionId';
        webhookMemberPermissionIdSelect.multiple = true;
        webhookMemberPermissionIdSelect.appendChild(testUser1Option);
        webhookMemberPermissionIdSelect.appendChild(testUser2Option);
        webhookMemberPermissionIdSelect.appendChild(testUser3Option);
        webhookMemberPermissionIdSelect.appendChild(testUser4Option);
        document.forms['form'].appendChild(webhookMemberPermissionIdSelect);

        const lineBotRolePermissionIdSelect = document.createElement('select');
        lineBotRolePermissionIdSelect.name = 'lineBotRolePermissionId';
        lineBotRolePermissionIdSelect.multiple = true;
        lineBotRolePermissionIdSelect.appendChild(testRole1Option);
        lineBotRolePermissionIdSelect.appendChild(testRole2Option);
        lineBotRolePermissionIdSelect.appendChild(testRole3Option);
        lineBotRolePermissionIdSelect.appendChild(testRole4Option);
        document.forms['form'].appendChild(lineBotRolePermissionIdSelect);

        const linePostDiscordChannelRolePermissionIdSelect = document.createElement('select');
        linePostDiscordChannelRolePermissionIdSelect.name = 'linePostDiscordChannelRolePermissionId';
        linePostDiscordChannelRolePermissionIdSelect.multiple = true;
        linePostDiscordChannelRolePermissionIdSelect.appendChild(testRole1Option);
        linePostDiscordChannelRolePermissionIdSelect.appendChild(testRole2Option);
        linePostDiscordChannelRolePermissionIdSelect.appendChild(testRole3Option);
        linePostDiscordChannelRolePermissionIdSelect.appendChild(testRole4Option);
        document.forms['form'].appendChild(linePostDiscordChannelRolePermissionIdSelect);

        const vcSignalRolePermissionIdSelect = document.createElement('select');
        vcSignalRolePermissionIdSelect.name = 'vcSignalRolePermissionId';
        vcSignalRolePermissionIdSelect.multiple = true;
        vcSignalRolePermissionIdSelect.appendChild(testRole1Option);
        vcSignalRolePermissionIdSelect.appendChild(testRole2Option);
        vcSignalRolePermissionIdSelect.appendChild(testRole3Option);
        vcSignalRolePermissionIdSelect.appendChild(testRole4Option);
        document.forms['form'].appendChild(vcSignalRolePermissionIdSelect);

        const webhookRolePermissionIdSelect = document.createElement('select');
        webhookRolePermissionIdSelect.name = 'webhookRolePermissionId';
        webhookRolePermissionIdSelect.multiple = true;
        webhookRolePermissionIdSelect.appendChild(testRole1Option);
        webhookRolePermissionIdSelect.appendChild(testRole2Option);
        webhookRolePermissionIdSelect.appendChild(testRole3Option);
        webhookRolePermissionIdSelect.appendChild(testRole4Option);
        document.forms['form'].appendChild(webhookRolePermissionIdSelect);
    });

    test('createJsonData', async () => {
        const formData = new FormData();
        formData.append('guildId', '111');
        formData.append('lineBotPermissionCode', '8');
        formData.append('linePostDiscordChannelPermissionCode', '8');
        formData.append('vcSignalPermissionCode', '8');
        formData.append('webhookPermissionCode', '8');
        formData.append('lineBotMemberPermissionId', '111');
        formData.append('lineBotMemberPermissionId', '222');
        formData.append('linePostDiscordChannelMemberPermissionId', '111');
        formData.append('linePostDiscordChannelMemberPermissionId', '222');
        formData.append('vcSignalMemberPermissionId', '111');
        formData.append('vcSignalMemberPermissionId', '222');
        formData.append('webhookMemberPermissionId', '111');
        formData.append('webhookMemberPermissionId', '222');
        formData.append('lineBotRolePermissionId', '111');
        formData.append('lineBotRolePermissionId', '222');
        formData.append('linePostDiscordChannelRolePermissionId', '111');
        formData.append('linePostDiscordChannelRolePermissionId', '222');
        formData.append('vcSignalRolePermissionId', '111');
        formData.append('vcSignalRolePermissionId', '222');
        formData.append('webhookRolePermissionId', '111');
        formData.append('webhookRolePermissionId', '222');

        const jsonData = await createJsonData(111, formData, document.forms['form'].elements);

        // jsonDataをオブジェクトにパース
        const parsedJsonData = JSON.parse(jsonData);

        // 期待値のオブジェクトを定義
        const expectedObject = {
            permissionCodes: [
                { guildId: 111, type: "lineBot", code: 8 },
                { guildId: 111, type: "linePostDiscordChannel", code: 8 },
                { guildId: 111, type: "vcSignal", code: 8 },
                { guildId: 111, type: "webhook", code: 8 }
            ],
            permissionUserIds: [
                { guildId: 111, type: "lineBot", userId: "111", permission: "all" },
                { guildId: 111, type: "lineBot", userId: "222", permission: "all" },
                { guildId: 111, type: "linePostDiscordChannel", userId: "111", permission: "all" },
                { guildId: 111, type: "linePostDiscordChannel", userId: "222", permission: "all" },
                { guildId: 111, type: "vcSignal", userId: "111", permission: "all" },
                { guildId: 111, type: "vcSignal", userId: "222", permission: "all" },
                { guildId: 111, type: "webhook", userId: "111", permission: "all" },
                { guildId: 111, type: "webhook", userId: "222", permission: "all" }
            ],
            permissionRoleIds: [
                { guildId: 111, type: "lineBot", roleId: "111", permission: "all" },
                { guildId: 111, type: "lineBot", roleId: "222", permission: "all" },
                { guildId: 111, type: "linePostDiscordChannel", roleId: "111", permission: "all" },
                { guildId: 111, type: "linePostDiscordChannel", roleId: "222", permission: "all" },
                { guildId: 111, type: "vcSignal", roleId: "111", permission: "all" },
                { guildId: 111, type: "vcSignal", roleId: "222", permission: "all" },
                { guildId: 111, type: "webhook", roleId: "111", permission: "all" },
                { guildId: 111, type: "webhook", roleId: "222", permission: "all" }
            ]
        };

        // パースしたオブジェクトと期待値のオブジェクトを比較
        expect(parsedJsonData).toEqual(expectedObject);
    });
});
