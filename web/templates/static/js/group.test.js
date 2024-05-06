const { JSDOM } = require('jsdom');

const dom = new JSDOM('<!doctype html><html><body><form id="form"><input id="id" type="text" name="name" value="pas"/></form></body></html>');
global.document = dom.window.document;
global.window = dom.window;

global.fetch = jest.fn();
container = global.document.getElementById('form');

const { createJsonData } = require('./group');

describe('fetchGroupData', () => {
    afterEach(() => {
        jest.restoreAllMocks();
    });

    // ここにテストケースを書く
    it('formからjsonに変換できること', async () => {
        const formData = new FormData();
        formData.append("defaultChannelId", "111");
        formData.append("debugMode", false);

        const jsonData = await createJsonData(formData)

        expect(jsonData).toEqual('{"defaultChannelId":"111","debugMode":"false"}');
    });
});