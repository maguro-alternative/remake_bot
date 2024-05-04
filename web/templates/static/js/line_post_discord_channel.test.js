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
        //document.forms['form'].elements = 
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

        const mockInput = document.createElement('input',{
            id: 'ng_111',
            type: 'checkbox',
            name: 'ng_111',
            value: 'false',
            checked: false
        });

        jest.spyOn(document, 'getElementById').mockImplementation(() =>
            mockInput
        );

        console.log(document.forms['form'].elements)

        const jsonData = await createJsonData(document.forms['form'].elements, formData)

        expect(jsonData).toEqual('{"default_channel_id":"111","debug_mode":"false"}');
    });
});