const { JSDOM } = require('jsdom');

const dom = new JSDOM('<!doctype html><html><body><form id="form"><input id="id" type="text" name="name" value="pas"/></form></body></html>');

global.document = dom.window.document;
global.window = dom.window;

global.fetch = jest.fn();
container = global.document.getElementById('form');

const { createJsonData } = require('./line_post_discord_channel');

describe('fetchLinePostDiscordChannelData', () => {
    beforeAll(() => {
        const ngInput = document.createElement('input');
        ngInput.id = 'ng_111';
        ngInput.type = 'checkbox';
        ngInput.name = 'ng_111';
        ngInput.value = 'false';
        ngInput.checked = false;
        document.forms['form'].appendChild(ngInput);

        const botMessageInput = document.createElement('input');
        botMessageInput.id = 'bot_message_111';
        botMessageInput.type = 'checkbox';
        botMessageInput.name = 'bot_message_111';
        botMessageInput.value = 'false';
        botMessageInput.checked = false;
        document.forms['form'].appendChild(botMessageInput);

        const ngTypeInput1 = document.createElement('input');
        ngTypeInput1.id = 'ng_types_111';
        ngTypeInput1.type = 'checkbox';
        ngTypeInput1.name = 'ng_types_111[]';
        ngTypeInput1.value = '1';
        ngTypeInput1.checked = false;
        document.forms['form'].appendChild(ngTypeInput1);

        const ngTypeInput2 = document.createElement('input');
        ngTypeInput2.id = 'ng_types_111';
        ngTypeInput2.type = 'checkbox';
        ngTypeInput2.name = 'ng_types_111[]';
        ngTypeInput2.value = '2';
        ngTypeInput2.checked = false;
        document.forms['form'].appendChild(ngTypeInput2);

        const ngUserInput1 = document.createElement('input');
        ngUserInput1.id = 'ng_users_111';
        ngUserInput1.type = 'checkbox';
        ngUserInput1.name = 'ng_users_111[]';
        ngUserInput1.value = '1111';
        ngUserInput1.checked = false;
        document.forms['form'].appendChild(ngUserInput1);

        const ngUserInput2 = document.createElement('input');
        ngUserInput2.id = 'ng_users_111';
        ngUserInput2.type = 'checkbox';
        ngUserInput2.name = 'ng_users_111[]';
        ngUserInput2.value = '2222';
        ngUserInput2.checked = false;
        document.forms['form'].appendChild(ngUserInput2);

        const ngRoleInput1 = document.createElement('input');
        ngRoleInput1.id = 'ng_roles_111';
        ngRoleInput1.type = 'checkbox';
        ngRoleInput1.name = 'ng_roles_111[]';
        ngRoleInput1.value = '3333';
        ngRoleInput1.checked = false;
        document.forms['form'].appendChild(ngRoleInput1);

        const ngRoleInput2 = document.createElement('input');
        ngRoleInput2.id = 'ng_roles_111';
        ngRoleInput2.type = 'checkbox';
        ngRoleInput2.name = 'ng_roles_111[]';
        ngRoleInput2.value = '4444';
        ngRoleInput2.checked = false;
        document.forms['form'].appendChild(ngRoleInput2);
    });
    afterEach(() => {
        jest.restoreAllMocks();
    });

    // ここにテストケースを書く
    it('formからjsonに変換できること', async () => {
        const formData = new FormData();
        formData.append("ng_111", false);
        formData.append("bot_message_111", false);
        formData.append("ng_types_111[]", 1);
        formData.append("ng_types_111[]", 2);
        formData.append("ng_users_111[]", "1111");
        formData.append("ng_users_111[]", "2222");
        formData.append("ng_roles_111[]", "3333");
        formData.append("ng_roles_111[]", "4444");

        const jsonData = await createJsonData(document.forms['form'].elements, formData)

        expect(jsonData).toEqual('{"channels":[{"channel_id":"111","ng":false,"bot_message":false,"ng_types":[1,2],"ng_users":["1111","2222"],"ng_roles":["3333","4444"]}]}');
    });
});