// FILEPATH: /c:/Users/bi_wa/gopro/remake_bot/web/templates/static/js/group.test.js
const { JSDOM } = require('jsdom');

const dom = new JSDOM('<!doctype html><html><body><form id="form"><input id="id" type="text" name="name" value="pas"/></form></body></html>');
global.document = dom.window.document;
global.window = dom.window;

global.fetch = jest.fn();
container = global.document.getElementById('form');

const { fetchGroupData, createJsonData } = require('./group');

describe('fetchGroupData', () => {
    afterEach(() => {
        jest.restoreAllMocks();
    });

    // ここにテストケースを書く
    it('should send a POST request with correct headers and body', async () => {
        const guildId = '123';
        const mockResponse = { status: 'success' };
        global.fetch.mockResolvedValue({
            ok: true,
            json: () => Promise.resolve(mockResponse),
        });
        const formData = new FormData();
        formData.append("username", "Groucho");
        formData.append("accountnum", 123456);

        const jsonData = await createJsonData(formData)
        const data = await fetchGroupData(guildId, jsonData);

        expect(fetch).toHaveBeenCalledWith(`/api/${guildId}/group`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonData
        });
        expect(data).toEqual(mockResponse);
    });

    /*it('should throw an error if the response is not ok', async () => {
        const guildId = '123';
        global.fetch.mockResolvedValue({
            ok: false,
        });
        const jsonData = await createJsonData(formData)

        await expect(fetchGroupData(guildId, jsonData)).rejects.toThrow('Error');
    });*/
});