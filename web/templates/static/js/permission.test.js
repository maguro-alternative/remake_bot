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
        lineBotPermissionCodeInput.name = 'line_bot_permission_code';
        lineBotPermissionCodeInput.value = 8;
        document.forms['form'].appendChild(lineBotPermissionCodeInput);

        const linePostDiscordChannelPermissionCodeInput = document.createElement('input');
        linePostDiscordChannelPermissionCodeInput.type = 'number';
        linePostDiscordChannelPermissionCodeInput.name = 'line_post_discord_channel_permission_code';
        linePostDiscordChannelPermissionCodeInput.value = 8;
        document.forms['form'].appendChild(linePostDiscordChannelPermissionCodeInput);

        const vcSignalPermissionCodeInput = document.createElement('input');
        vcSignalPermissionCodeInput.type = 'number';
        vcSignalPermissionCodeInput.name = 'vc_signal_permission_code';
        vcSignalPermissionCodeInput.value = 8;
        document.forms['form'].appendChild(vcSignalPermissionCodeInput);

        const webhookPermissionCodeInput = document.createElement('input');
        webhookPermissionCodeInput.type = 'number';
        webhookPermissionCodeInput.name = 'webhook_permission_code';
        webhookPermissionCodeInput.value = 8;
        document.forms['form'].appendChild(webhookPermissionCodeInput);

        const lineBotMemberPermissionIdSelect = document.createElement('select');
        lineBotMemberPermissionIdSelect.name = 'line_bot_member_permission_id';
        lineBotMemberPermissionIdSelect.multiple = true;
        lineBotMemberPermissionIdSelect.appendChild(testUser1Option);
        lineBotMemberPermissionIdSelect.appendChild(testUser2Option);
        lineBotMemberPermissionIdSelect.appendChild(testUser3Option);
        lineBotMemberPermissionIdSelect.appendChild(testUser4Option);
        document.forms['form'].appendChild(lineBotMemberPermissionIdSelect);

        const linePostDiscordChannelMemberPermissionIdSelect = document.createElement('select');
        linePostDiscordChannelMemberPermissionIdSelect.name = 'line_post_discord_channel_member_permission_id';
        linePostDiscordChannelMemberPermissionIdSelect.multiple = true;
        linePostDiscordChannelMemberPermissionIdSelect.appendChild(testUser1Option);
        linePostDiscordChannelMemberPermissionIdSelect.appendChild(testUser2Option);
        linePostDiscordChannelMemberPermissionIdSelect.appendChild(testUser3Option);
        linePostDiscordChannelMemberPermissionIdSelect.appendChild(testUser4Option);
        document.forms['form'].appendChild(linePostDiscordChannelMemberPermissionIdSelect);

        const vcSignalMemberPermissionIdSelect = document.createElement('select');
        vcSignalMemberPermissionIdSelect.name = 'vc_signal_member_permission_id';
        vcSignalMemberPermissionIdSelect.multiple = true;
        vcSignalMemberPermissionIdSelect.appendChild(testUser1Option);
        vcSignalMemberPermissionIdSelect.appendChild(testUser2Option);
        vcSignalMemberPermissionIdSelect.appendChild(testUser3Option);
        vcSignalMemberPermissionIdSelect.appendChild(testUser4Option);
        document.forms['form'].appendChild(vcSignalMemberPermissionIdSelect);

        const webhookMemberPermissionIdSelect = document.createElement('select');
        webhookMemberPermissionIdSelect.name = 'webhook_member_permission_id';
        webhookMemberPermissionIdSelect.multiple = true;
        webhookMemberPermissionIdSelect.appendChild(testUser1Option);
        webhookMemberPermissionIdSelect.appendChild(testUser2Option);
        webhookMemberPermissionIdSelect.appendChild(testUser3Option);
        webhookMemberPermissionIdSelect.appendChild(testUser4Option);
        document.forms['form'].appendChild(webhookMemberPermissionIdSelect);

        const lineBotRolePermissionIdSelect = document.createElement('select');
        lineBotRolePermissionIdSelect.name = 'line_bot_role_permission_id';
        lineBotRolePermissionIdSelect.multiple = true;
        lineBotRolePermissionIdSelect.appendChild(testRole1Option);
        lineBotRolePermissionIdSelect.appendChild(testRole2Option);
        lineBotRolePermissionIdSelect.appendChild(testRole3Option);
        lineBotRolePermissionIdSelect.appendChild(testRole4Option);
        document.forms['form'].appendChild(lineBotRolePermissionIdSelect);

        const linePostDiscordChannelRolePermissionIdSelect = document.createElement('select');
        linePostDiscordChannelRolePermissionIdSelect.name = 'line_post_discord_channel_role_permission_id';
        linePostDiscordChannelRolePermissionIdSelect.multiple = true;
        linePostDiscordChannelRolePermissionIdSelect.appendChild(testRole1Option);
        linePostDiscordChannelRolePermissionIdSelect.appendChild(testRole2Option);
        linePostDiscordChannelRolePermissionIdSelect.appendChild(testRole3Option);
        linePostDiscordChannelRolePermissionIdSelect.appendChild(testRole4Option);
        document.forms['form'].appendChild(linePostDiscordChannelRolePermissionIdSelect);

        const vcSignalRolePermissionIdSelect = document.createElement('select');
        vcSignalRolePermissionIdSelect.name = 'vc_signal_role_permission_id';
        vcSignalRolePermissionIdSelect.multiple = true;
        vcSignalRolePermissionIdSelect.appendChild(testRole1Option);
        vcSignalRolePermissionIdSelect.appendChild(testRole2Option);
        vcSignalRolePermissionIdSelect.appendChild(testRole3Option);
        vcSignalRolePermissionIdSelect.appendChild(testRole4Option);
        document.forms['form'].appendChild(vcSignalRolePermissionIdSelect);

        const webhookRolePermissionIdSelect = document.createElement('select');
        webhookRolePermissionIdSelect.name = 'webhook_role_permission_id';
        webhookRolePermissionIdSelect.multiple = true;
        webhookRolePermissionIdSelect.appendChild(testRole1Option);
        webhookRolePermissionIdSelect.appendChild(testRole2Option);
        webhookRolePermissionIdSelect.appendChild(testRole3Option);
        webhookRolePermissionIdSelect.appendChild(testRole4Option);
        document.forms['form'].appendChild(webhookRolePermissionIdSelect);
    });

    test('createJsonData', async () => {
        const formData = new FormData();
        formData.append('guild_id', '111');
        formData.append('line_bot_permission_code', '8');
        formData.append('line_post_discord_channel_permission_code', '8');
        formData.append('vc_signal_permission_code', '8');
        formData.append('webhook_permission_code', '8');
        formData.append('line_bot_member_permission_id', '111');
        formData.append('line_bot_member_permission_id', '222');
        formData.append('line_post_discord_channel_member_permission_id', '111');
        formData.append('line_post_discord_channel_member_permission_id', '222');
        formData.append('vc_signal_member_permission_id', '111');
        formData.append('vc_signal_member_permission_id', '222');
        formData.append('webhook_member_permission_id', '111');
        formData.append('webhook_member_permission_id', '222');
        formData.append('line_bot_role_permission_id', '111');
        formData.append('line_bot_role_permission_id', '222');
        formData.append('line_post_discord_channel_role_permission_id', '111');
        formData.append('line_post_discord_channel_role_permission_id', '222');
        formData.append('vc_signal_role_permission_id', '111');
        formData.append('vc_signal_role_permission_id', '222');
        formData.append('webhook_role_permission_id', '111');
        formData.append('webhook_role_permission_id', '222');

        const jsonData = await createJsonData(111, formData, document.forms['form'].elements);

        expect(jsonData).toBe('{"permission_codes":[{"guild_id":111,"type":"line_bot","code":8},{"guild_id":111,"type":"line_post_discord_channel","code":8},{"guild_id":111,"type":"vc_signal","code":8},{"guild_id":111,"type":"webhook","code":8}],"permission_user_ids":[{"guild_id":111,"type":"line_bot","user_id":"111","permission":"all"},{"guild_id":111,"type":"line_bot","user_id":"222","permission":"all"},{"guild_id":111,"type":"line_post_discord_channel","user_id":"111","permission":"all"},{"guild_id":111,"type":"line_post_discord_channel","user_id":"222","permission":"all"},{"guild_id":111,"type":"vc_signal","user_id":"111","permission":"all"},{"guild_id":111,"type":"vc_signal","user_id":"222","permission":"all"},{"guild_id":111,"type":"webhook","user_id":"111","permission":"all"},{"guild_id":111,"type":"webhook","user_id":"222","permission":"all"}],"permission_role_ids":[{"guild_id":111,"type":"line_bot","role_id":"111","permission":"all"},{"guild_id":111,"type":"line_bot","role_id":"222","permission":"all"},{"guild_id":111,"type":"line_post_discord_channel","role_id":"111","permission":"all"},{"guild_id":111,"type":"line_post_discord_channel","role_id":"222","permission":"all"},{"guild_id":111,"type":"vc_signal","role_id":"111","permission":"all"},{"guild_id":111,"type":"vc_signal","role_id":"222","permission":"all"},{"guild_id":111,"type":"webhook","role_id":"111","permission":"all"},{"guild_id":111,"type":"webhook","role_id":"222","permission":"all"}]}');
    });
});
