const { JSDOM } = require('jsdom');

const dom = new JSDOM('<!doctype html><html><body><form id="form"><input id="id" type="text" name="name" value="pas"/></form></body></html>');
global.document = dom.window.document;
global.window = dom.window;

global.fetch = jest.fn();
container = global.document.getElementById('form');

const { createJsonData } = require('./vc_signal');

describe('fetchGroupData', () => {
    beforeAll(() => {
        const testUser1Option = document.createElement('option');
        testUser1Option.name = 'testuser1';
        testUser1Option.value = '1111';
        testUser1Option.selected = true;
        testUser1Option.defaultSelected = false;

        const testUser2Option = document.createElement('option');
        testUser2Option.name = 'testuser2';
        testUser2Option.value = '2222';
        testUser2Option.selected = true;
        testUser2Option.defaultSelected = false;

        const testUser3Option = document.createElement('option');
        testUser3Option.name = 'testuser3';
        testUser3Option.value = '1112';
        testUser3Option.selected = true;
        testUser3Option.defaultSelected = false;

        const testRole1Option = document.createElement('option');
        testRole1Option.name = 'testrole1';
        testRole1Option.value = '11111';
        testRole1Option.selected = true;
        testRole1Option.defaultSelected = false;

        const testRole2Option = document.createElement('option');
        testRole2Option.name = 'testrole2';
        testRole2Option.value = '11112';
        testRole2Option.selected = true;
        testRole2Option.defaultSelected = false;

        const vcSignalMentionUserIdsSelect = document.createElement('select');
        vcSignalMentionUserIdsSelect.name = 'vcSignalMentionUserIds111[]';
        vcSignalMentionUserIdsSelect.multiple = true;
        vcSignalMentionUserIdsSelect.appendChild(testUser1Option);
        vcSignalMentionUserIdsSelect.appendChild(testUser2Option);
        document.forms['form'].appendChild(vcSignalMentionUserIdsSelect);

        const vcSignalMentionRoleIdsSelect = document.createElement('select');
        vcSignalMentionRoleIdsSelect.name = 'vcSignalMentionRoleIds111[]';
        vcSignalMentionRoleIdsSelect.multiple = true;
        vcSignalMentionRoleIdsSelect.appendChild(testRole2Option);
        document.forms['form'].appendChild(vcSignalMentionRoleIdsSelect);

        const vcSignalNgUserIdsSelect = document.createElement('select');
        vcSignalNgUserIdsSelect.name = 'vcSignalNgUserIds111[]';
        vcSignalNgUserIdsSelect.multiple = true;
        vcSignalNgUserIdsSelect.appendChild(testUser3Option);
        document.forms['form'].appendChild(vcSignalNgUserIdsSelect);

        const vcSignalNgRoleIdsSelect = document.createElement('select');
        vcSignalNgRoleIdsSelect.name = 'vcSignalNgRoleIds111[]';
        vcSignalNgRoleIdsSelect.multiple = true;
        vcSignalNgRoleIdsSelect.appendChild(testRole2Option);
        document.forms['form'].appendChild(vcSignalNgRoleIdsSelect);

        const sendSignalInput = document.createElement('input');
        sendSignalInput.id = 'sendSignal111';
        sendSignalInput.type = 'checkbox';
        sendSignalInput.name = 'sendSignal111';
        sendSignalInput.value = 'on';
        sendSignalInput.checked = true;
        document.forms['form'].appendChild(sendSignalInput);

        const joinBotInput = document.createElement('input');
        joinBotInput.id = 'joinBot111';
        joinBotInput.type = 'checkbox';
        joinBotInput.name = 'joinBot111';
        joinBotInput.value = 'on';
        joinBotInput.checked = true;
        document.forms['form'].appendChild(joinBotInput);

        const everyoneMentionInput = document.createElement('input');
        everyoneMentionInput.id = 'everyoneMention111';
        everyoneMentionInput.type = 'checkbox';
        everyoneMentionInput.name = 'everyoneMention111';
        everyoneMentionInput.value = 'on';
        everyoneMentionInput.checked = true;
        document.forms['form'].appendChild(everyoneMentionInput);

        const sendChannelIdInput = document.createElement('input');
        sendChannelIdInput.id = 'sendChannelId111';
        sendChannelIdInput.type = 'text';
        sendChannelIdInput.name = 'sendChannelId111';
        sendChannelIdInput.value = '111111';
        document.forms['form'].appendChild(sendChannelIdInput);
    })
    afterEach(() => {
        jest.restoreAllMocks();
    });

    // ここにテストケースを書く
    it('formからjsonに変換できること', async () => {
        const formData = new FormData();
        formData.append('sendSignal111', 'on');
        formData.append('joinBot111', 'on');
        formData.append('everyoneMention111', 'on')
        formData.append('vcSignalMentionUserIds111[]', '1111')
        formData.append('vcSignalMentionUserIds111[]', '2222')
        formData.append('vcSignalMentionRoleIds111[]', '11112')
        formData.append('vcSignalNgUserIds111[]', '1112')
        formData.append('vcSignalNgRoleIds111[]', '11112')
        formData.append('defaultChannelId111', '111111')

        const jsonData = await createJsonData(document.forms['form'].elements, formData)

        // jsonDataをオブジェクトにパース
        const parsedJsonData = JSON.parse(jsonData);

        // 期待値のオブジェクトを定義
        const expectedObject = {
            vcSignals: [
                {
                    vcChannelId: "111",
                    sendSignal: true,
                    sendChannelId: "111111",
                    joinBot: true,
                    everyoneMention: true,
                    vcSignalMentionUserIds: ["1111", "2222"],
                    vcSignalMentionRoleIds: ["11112"],
                    vcSignalNgUserIds: ["1112"],
                    vcSignalNgRoleIds: ["11112"]
                }
            ]
        }
        // パースしたオブジェクトと期待値のオブジェクトを比較
        expect(parsedJsonData).toEqual(expectedObject);
    });
});
