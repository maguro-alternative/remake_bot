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
        ngInput.id = 'ng111';
        ngInput.type = 'checkbox';
        ngInput.name = 'ng111';
        ngInput.value = 'false';
        ngInput.checked = false;
        document.forms['form'].appendChild(ngInput);

        const botMessageInput = document.createElement('input');
        botMessageInput.id = 'botMessage111';
        botMessageInput.type = 'checkbox';
        botMessageInput.name = 'botMessage111';
        botMessageInput.value = 'false';
        botMessageInput.checked = false;
        document.forms['form'].appendChild(botMessageInput);

        const ngTypeInput1 = document.createElement('input');
        ngTypeInput1.id = 'ngTypes111';
        ngTypeInput1.type = 'checkbox';
        ngTypeInput1.name = 'ngTypes111[]';
        ngTypeInput1.value = '1';
        ngTypeInput1.checked = false;
        document.forms['form'].appendChild(ngTypeInput1);

        const ngTypeInput2 = document.createElement('input');
        ngTypeInput2.id = 'ngTypes111';
        ngTypeInput2.type = 'checkbox';
        ngTypeInput2.name = 'ngTypes111[]';
        ngTypeInput2.value = '2';
        ngTypeInput2.checked = false;
        document.forms['form'].appendChild(ngTypeInput2);

        const ngUserInput1 = document.createElement('input');
        ngUserInput1.id = 'ngUsers111';
        ngUserInput1.type = 'checkbox';
        ngUserInput1.name = 'ngUsers111[]';
        ngUserInput1.value = '1111';
        ngUserInput1.checked = false;
        document.forms['form'].appendChild(ngUserInput1);

        const ngUserInput2 = document.createElement('input');
        ngUserInput2.id = 'ngUsers111';
        ngUserInput2.type = 'checkbox';
        ngUserInput2.name = 'ngUsers111[]';
        ngUserInput2.value = '2222';
        ngUserInput2.checked = false;
        document.forms['form'].appendChild(ngUserInput2);

        const ngRoleInput1 = document.createElement('input');
        ngRoleInput1.id = 'ngRoles111';
        ngRoleInput1.type = 'checkbox';
        ngRoleInput1.name = 'ngRoles111[]';
        ngRoleInput1.value = '3333';
        ngRoleInput1.checked = false;
        document.forms['form'].appendChild(ngRoleInput1);

        const ngRoleInput2 = document.createElement('input');
        ngRoleInput2.id = 'ngRoles111';
        ngRoleInput2.type = 'checkbox';
        ngRoleInput2.name = 'ngRoles111[]';
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
        formData.append("ng111", false);
        formData.append("botMessage111", false);
        formData.append("ngTypes111[]", 1);
        formData.append("ngTypes111[]", 2);
        formData.append("ngUsers111[]", "1111");
        formData.append("ngUsers111[]", "2222");
        formData.append("ngRoles111[]", "3333");
        formData.append("ngRoles111[]", "4444");

        const jsonData = await createJsonData(document.forms['form'].elements, formData)

        // jsonDataをオブジェクトにパース
        const parsedJsonData = JSON.parse(jsonData);

        // 期待値のオブジェクトを定義
        const expectedObject = {
            channels: [
                {
                    channelId: "111",
                    ng: false,
                    botMessage: false,
                    ngTypes: [1, 2],
                    ngUsers: ["1111", "2222"],
                    ngRoles: ["3333", "4444"]
                }
            ]
        };

        // パースしたオブジェクトと期待値のオブジェクトを比較
        expect(parsedJsonData).toEqual(expectedObject);
    });
});